package main

import (
	"log"

	"github.com/Teeworlds-Server-Moderation/common/amqp"
	"github.com/Teeworlds-Server-Moderation/common/events"
	"github.com/jxsl13/twapi/econ"
)

func requestProcessor(cfg *Config, subscriber *amqp.Subscriber, publisher *amqp.Publisher, conn *econ.Conn) {

	base := events.BaseEvent{}

	next, err := subscriber.Consume(cfg.EconAddress)
	if err != nil {
		log.Fatalf("Failed to consume from queue %s, closing.", cfg.EconAddress)
	}

	for msg := range next {
		payload := string(msg.Body)
		err := base.Unmarshal(payload)
		if err != nil {
			log.Printf("Failed to unmarshal BaseEvent(%s): %s\n", cfg.EconAddress, payload)
			continue
		}

		switch base.Type {
		case events.TypeRequestCommandExec:
			event := events.NewRequestCommandExecEvent()
			err = event.Unmarshal(payload)
			if err != nil {
				log.Printf("Failed to unmarshal expected request event type: %s:%s\n", base.Type, err)
				continue
			}
			err = econExecAndLog(conn, event)
			if err != nil {
				log.Printf("Failed to econExecAndLog command: %s\n", err)
				continue
			}
		default:
			log.Printf("Unknown request event type received: %s\n", base.Type)
		}

	}
}
