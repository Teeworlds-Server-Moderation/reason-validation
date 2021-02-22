package main

import (
	"encoding/json"
	"log"

	"github.com/Teeworlds-Server-Moderation/common/amqp"
	"github.com/Teeworlds-Server-Moderation/common/dto"
	"github.com/Teeworlds-Server-Moderation/common/events"
)

// ReasonEvent are events that have a Reason
// and a player that started a vote
type ReasonEvent struct {
	events.BaseEvent
	Reason string
	Source dto.Player
}

func processMessage(message string, publisher *amqp.Publisher, cfg *Config, cs *CSet) error {
	event := ReasonEvent{}
	err := json.Unmarshal([]byte(message), &event)
	if err != nil {
		return err
	}

	reaction, ok := cs.Get(event.Reason, event.Type)
	if !ok {
		// unknown reason
		if err := publisher.Publish("", unknownReasonQueue, message); err != nil {
			return err
		}
		reaction = "abort"
	}

	switch reaction {
	case "voteban":
		return requestCommandExecForPlayer(
			publisher,
			cfg.DefaultVotebanDuration,
			event.Source,
			cfg.DefaultVotebanCommand,
			"funvote",
			event.EventSource,
			cfg.BroadcastNonAbortActions,
		)
	case "ignore":
		log.Println("Ignoring: ", message)
		return nil
	default:
		// abort
		err := requestCommandExecForPlayer(
			publisher,
			0,
			event.Source,
			"vote no",
			"",
			event.EventSource,
			false,
		)
		if err != nil {
			return err
		}
		// send info message
		err = requestCommandExecForPlayer(
			publisher,
			0,
			event.Source,
			"say Votes are usually aborted because it might make sense to move the player to the spectators instead.",
			"",
			event.EventSource,
			false,
		)
		if err != nil {
			return err
		}
	}
	return nil
}
