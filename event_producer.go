package main

import (
	"context"
	"log"

	"github.com/Teeworlds-Server-Moderation/common/amqp"
)

func eventProducerRoutine(ctx context.Context, source string, lineChan chan string, publisher *amqp.Publisher) {
	for {
		select {
		case <-ctx.Done():
			log.Println("Closing event parser routine...")
			return
		case line, ok := <-lineChan:
			if !ok {
				continue
			}

			msg, err := parseEvent(source, line)
			if err != nil {
				log.Printf("Skipped: %s\n", line)
				continue
			}
			if err := publisher.Publish(msg.Queue, msg.Payload); err != nil {
				log.Printf("Error: %s\nError: %s\n", line, err)
				continue
			}

			log.Printf("Processed: %s\n", line)
		}
	}
}
