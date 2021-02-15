package main

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/jxsl13/twapi/econ"
)

func econLineReader(ctx context.Context, cancel context.CancelFunc, conn *econ.Conn, lChan chan string) {
	retries := 0
	for {
		line, err := conn.ReadLine()
		select {
		case <-ctx.Done():
			log.Println("Closing econ reader...")
			return
		default:
			if errors.Is(err, econ.ErrNetwork) {
				time.Sleep(time.Second)
				retries++
				if retries >= 60 {
					log.Fatalln("Failed to reestablish econ connection.")
				}
				continue
			} else if err != nil {
				log.Println("Failed to read line: ", err)
				log.Panicln("Closing application...")
				cancel()
				return
			}
			// reset retries
			retries = 0
			// push read line int channel
			lChan <- line
		}
	}
}
