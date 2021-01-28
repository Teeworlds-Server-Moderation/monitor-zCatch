package main

import (
	"log"

	"github.com/Teeworlds-Server-Moderation/common/events"
	"github.com/Teeworlds-Server-Moderation/common/mqtt"
	"github.com/jxsl13/twapi/econ"
)

func requestProcessor(subscriber *mqtt.Subscriber, publisher *mqtt.Publisher, conn *econ.Conn) {
	base := events.BaseEvent{}
	for msg := range subscriber.Next() {
		err := base.Unmarshal(msg.Payload)
		if err != nil {
			log.Printf("Failed to unmarshal BaseEvent(%s): %s\n", msg.Topic, msg.Payload)
			continue
		}

		switch base.Type {
		case events.TypeCommandExec:
			event := events.NewCommandExecEvent()
			err = event.Unmarshal(msg.Payload)
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
