package main

import (
	"context"
	"log"
	"time"

	"github.com/jxsl13/twapi/econ"
)

func econLineReader(ctx context.Context, conn *econ.Conn, lChan chan string) {
	for {
		line, err := conn.ReadLine()
		select {
		case <-ctx.Done():
			log.Println("Closing econ reader...")
			return
		default:
			if err != nil {
				log.Println("Failed to read line: ", err)
				time.Sleep(time.Second)
				return
			}
			// push read line int channel
			lChan <- line
		}
	}
}
