package main

import (
	"context"
	"log"

	mqtt "github.com/Teeworlds-Server-Moderation/common/mqtt"
)

func eventProducerRoutine(ctx context.Context, source string, lineChan chan string, publisher *mqtt.Publisher) {
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
			publisher.Publish(msg)
			log.Printf("Processed: %s\n", line)
		}
	}
}
