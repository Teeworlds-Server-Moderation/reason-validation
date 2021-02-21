package main

import (
	"log"

	"github.com/Teeworlds-Server-Moderation/common/amqp"
)

func processMessage(message string, publisher *amqp.Publisher, cfg *Config) error {
	log.Printf("Recived message: %s\n", message)
	return nil
}
