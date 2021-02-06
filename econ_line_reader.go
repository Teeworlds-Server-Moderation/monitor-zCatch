package main

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/jxsl13/twapi/econ"
)

func econLineReader(ctx context.Context, cancel context.CancelFunc, conn *econ.Conn, lChan chan string) {
	for {
		line, err := conn.ReadLine()
		select {
		case <-ctx.Done():
			log.Println("Closing econ reader...")
			return
		default:
			if errors.Is(err, econ.ErrNetwork) {
				time.Sleep(time.Second)
				continue
			} else if err != nil {
				log.Println("Failed to read line: ", err)
				log.Panicln("Closing application...")
				cancel()
				return
			}
			// push read line int channel
			lChan <- line
		}
	}
}
