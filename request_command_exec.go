package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Teeworlds-Server-Moderation/common/amqp"
	"github.com/Teeworlds-Server-Moderation/common/dto"
	"github.com/Teeworlds-Server-Moderation/common/events"
	"github.com/Teeworlds-Server-Moderation/common/topics"
)

// serverTopic may be either the server's ip:port address or the broadcast topic
func requestCommandExecForPlayer(publisher *amqp.Publisher, defaultDuration time.Duration, player dto.Player, command, reason, targetServer string, broadcast bool) error {
	event := events.NewRequestCommandExecEvent()
	event.Timestamp = time.Now().Format("2006-01-02 15:04:05")
	event.Requestor = applicationID
	event.EventSource = applicationID

	// construct command and replace
	replacer := strings.NewReplacer(
		"{IP}",
		player.IP,
		"{ID}",
		fmt.Sprintf("%d", player.ID),
		"{DURATION:MINUTES}",
		fmt.Sprintf("%d", int64(defaultDuration/time.Minute)),
		"{DURATION:SECONDS}",
		fmt.Sprintf("%d", int64(defaultDuration/time.Second)),
		"{REASON}",
		reason,
	)

	broadcastFeasible := true
	if strings.Contains(command, "{ID}") {
		broadcastFeasible = false
	}

	banCommand := replacer.Replace(command)
	event.Command = banCommand

	payload := event.Marshal()
	log.Println("Requesting: ", payload)
	if broadcast && broadcastFeasible {
		// ban on all servers
		// if broadcasting makes sense
		// if the ban command contains an ID,
		// it makes no sense to broadcast it
		publisher.Publish(topics.Broadcast, "", payload)
	} else {
		// only ban on the server where the player joined
		// do not publish to exchange, but directly to the queue
		publisher.Publish("", targetServer, payload)
	}
	return nil
}
