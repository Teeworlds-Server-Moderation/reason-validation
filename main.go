package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"github.com/Teeworlds-Server-Moderation/common/amqp"
	"github.com/Teeworlds-Server-Moderation/common/env"
	"github.com/Teeworlds-Server-Moderation/common/events"
	"github.com/Teeworlds-Server-Moderation/common/topics"
)

var (
	cfg                = &Config{}
	store              = NewCSet()
	applicationID      = "reason-validation"
	unknownReasonQueue = "unknown-vote-reason"
	subscriber         *amqp.Subscriber
	publisher          *amqp.Publisher
	profile            *os.File
	startupTime        = time.Now()
)

func brokerCredentials(c *Config) (address, username, password string) {
	return c.BrokerAddress, c.BrokerUsername, c.BrokerPassword
}

// ExchangeCreator can be publisher or subscriber
type ExchangeCreator interface {
	CreateExchange(string) error
}

// QueueCreateBinder creates queues and binds them to exchanges
type QueueCreateBinder interface {
	CreateQueue(queue string) error
	BindQueue(queue, exchange string) error
}

func createExchanges(ec ExchangeCreator, exchanges ...string) {
	for _, exchange := range exchanges {
		if err := ec.CreateExchange(exchange); err != nil {
			log.Fatalf("Failed to create exchange '%s': %v\n", exchange, err)
		}
	}
}

func createQueueAndBindToExchanges(qcb QueueCreateBinder, queue string, exchanges ...string) {
	if err := qcb.CreateQueue(queue); err != nil {
		log.Fatalf("Failed to create queue '%s'\n", queue)
	}

	for _, exchange := range exchanges {
		if err := qcb.BindQueue(queue, exchange); err != nil {
			log.Fatalf("Failed to bind queue '%s' to exchange '%s'\n", queue, exchange)
		}

	}
}

func init() {

	err := env.Parse(cfg)
	if err != nil {
		log.Fatalln(err)
	}

	subscriber, err = amqp.NewSubscriber(brokerCredentials(cfg))
	if err != nil {
		log.Fatalln("Could not establish subscriber connection: ", err)
	}

	publisher, err = amqp.NewPublisher(brokerCredentials(cfg))
	if err != nil {
		log.Fatalln("Could not establish publisher connection: ", err)
	}

	createExchanges(
		publisher,
		topics.Broadcast,
	)

	createExchanges(
		subscriber,
		events.TypeVoteKickStarted,
		events.TypeVoteSpecStarted,
	)

	createQueueAndBindToExchanges(
		subscriber,
		applicationID,
		events.TypeVoteKickStarted,
		events.TypeVoteSpecStarted,
	)

	err = initializeKeyValueStore(store, cfg.DataPath)
	if err != nil {
		log.Fatalln(err)
	}

	var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
	flag.Parse()
	if *cpuprofile != "" {
		profile, err = os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer publisher.Close()
	defer subscriber.Close()
	defer cancel()
	if profile != nil {
		defer profile.Close()
	}

	// message processing
	go func() {
		next, err := subscriber.Consume(applicationID)
		if err != nil {
			log.Fatalln(err)
		}
		for msg := range next {
			if err := processMessage(string(msg.Body), publisher, cfg, store); err != nil {
				log.Printf("Error processing message: %s\n", err)
			}
		}
	}()

	// database backups
	go func() {
		ticker := time.NewTicker(cfg.BackupInterval)

		for {
			select {
			case <-ctx.Done():
				log.Println("Closing backup routine...")

				if time.Now().Before(startupTime.Add(cfg.DurationBeforeFirstBackup)) {
					return
				}

				log.Println("Creating backup...")
				filename := time.Now().Format(time.RFC3339) + ".csv"
				filename = path.Join(cfg.DataPath, filename)
				data, err := store.DumpCSV()
				if err != nil {
					log.Printf("Failed to retrieve data for backup: %v\n", err)
					continue
				}
				err = ioutil.WriteFile(filename, data, 0660)
				if err != nil {
					log.Printf("Failed to write data to file '%s': %v", filename, err)
					continue
				}

				return
			case now := <-ticker.C:
				filename := now.Format(time.RFC3339) + ".csv"
				filename = path.Join(cfg.DataPath, filename)
				data, err := store.DumpCSV()
				if err != nil {
					log.Printf("Failed to retrieve data for backup: %v\n", err)
					continue
				}
				err = ioutil.WriteFile(filename, data, 0660)
				if err != nil {
					log.Printf("Failed to write data to file '%s': %v", filename, err)
					continue
				}

			}
		}
	}()

	// Messages will be delivered asynchronously so we just need to wait for a signal to shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	fmt.Println("Connection is up, press Ctrl-C to shutdown")
	<-sig
	fmt.Println("Signal caught - exiting")
	fmt.Println("Shutdown complete")
}
