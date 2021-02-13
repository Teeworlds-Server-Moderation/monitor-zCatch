package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Teeworlds-Server-Moderation/common/amqp"
	env "github.com/Teeworlds-Server-Moderation/common/env"
	"github.com/jxsl13/twapi/econ"
)

var (
	cfg = &Config{}
)

func init() {
	err := env.Parse(cfg)
	if err != nil {
		log.Fatalf("Failed to get environment variables: %s\n", err)
	}
}

func brokerCredentials(c *Config) (address, username, password string) {
	return cfg.BrokerAddress, cfg.BrokerUsername, cfg.BrokerPassword
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conn, err := econ.New(cfg.EconAddress, cfg.EconPassword)
	if err != nil {
		log.Fatalf("Failed to connect to or authenticate at %s: %s", cfg.EconAddress, err)
	}
	defer conn.Close()
	log.Println("Successfully connected to Econ: ", cfg.EconAddress)

	// Make logging more verbose in order to get more information
	err = conn.WriteLine("ec_output_level 2")
	if err != nil {
		log.Fatalln("Failed to set: ec_output_level 2")
	}

	publisher, err := amqp.NewPublisher(brokerCredentials(cfg))
	if err != nil {
		log.Fatalf("Failed to connect to broker %s: %s", cfg.BrokerAddress, err)
	}

	subscriber, err := amqp.NewSubscriber(brokerCredentials(cfg))
	if err != nil {
		log.Fatalf("Failed to connect to broker %s: %s", cfg.BrokerAddress, err)
	}

	// buffered channel, blocks when full.
	lineChan := make(chan string, 64)

	// Connects to the Teeworlds server and reads lines an dpushes them into
	// the lineChan channel (buffered channel)
	// cancel is able to close the whole application, as it might make no sense to
	// continue rrunning when there is no input data provided
	go econLineReader(ctx, cancel, conn, lineChan)
	// receives those lines in the lineChan channel and parses them in order to create events
	// that are then pushed to their corresponding broker topics
	go eventProducerRoutine(ctx, cfg.EconAddress, lineChan, publisher)

	// There are two topics that the monitor currently listens on, the first one being the
	// "IP:Port" topic, where individual messages can be received and the second topic being the
	// broadcast topic that multicasts one message to all subscribing monitors.
	go requestProcessor(cfg, subscriber, publisher, conn)

	// Messages will be delivered asynchronously so we just need to wait for a signal to shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	fmt.Println("Startup finished, running...")
	<-sig
	fmt.Println("signal caught - exiting")
	fmt.Println("shutdown complete")
}
