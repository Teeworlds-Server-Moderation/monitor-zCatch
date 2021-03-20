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
	"github.com/Teeworlds-Server-Moderation/common/events"
	"github.com/Teeworlds-Server-Moderation/common/topics"
	"github.com/jxsl13/twapi/econ"
)

var (
	cfg        = &Config{}
	ctx        context.Context
	cancel     func()
	conn       *econ.Conn
	publisher  *amqp.Publisher
	subscriber *amqp.Subscriber
)

func brokerCredentials(c *Config) (address, username, password string) {
	return c.BrokerAddress, c.BrokerUsername, c.BrokerPassword
}

// ExchangeCreator can be publisher or subscriber
type ExchangeCreator interface {
	CreateExchange(string) error
}

// QueueCreateBinder creates queues and binds them to exchanges
type QueueCreateBinder interface {
	CreateQueue(queue string) error
	BindQueue(queue, exchange string) error
}

func createExchanges(ec ExchangeCreator, exchanges ...string) {
	for _, exchange := range exchanges {
		if err := ec.CreateExchange(exchange); err != nil {
			log.Fatalf("Failed to create exchange '%s': %v\n", exchange, err)
		}
	}
}

func createQueueAndBindToExchanges(qcb QueueCreateBinder, queue string, exchanges ...string) {
	if err := qcb.CreateQueue(queue); err != nil {
		log.Fatalf("Failed to create queue '%s'\n", queue)
	}

	for _, exchange := range exchanges {
		if err := qcb.BindQueue(queue, exchange); err != nil {
			log.Fatalf("Failed to bind queue '%s' to exchange '%s'\n", queue, exchange)
		}

	}

}

func init() {
	err := env.Parse(cfg)
	if err != nil {
		log.Fatalf("Failed to get environment variables: %s\n", err)
	}

	// context
	ctx, cancel = context.WithCancel(context.Background())

	// econ connection
	conn, err = econ.New(cfg.EconAddress, cfg.EconPassword)
	if err != nil {
		log.Fatalf("Failed to connect to or authenticate at %s: %s", cfg.EconAddress, err)
	}

	// Make logging more verbose in order to get more information
	err = conn.WriteLine("ec_output_level 2")
	if err != nil {
		log.Fatalln("Failed to set: ec_output_level 2")
	}
	log.Println("Successfully connected to Econ: ", cfg.EconAddress)

	publisher, err = amqp.NewPublisher(brokerCredentials(cfg))
	if err != nil {
		log.Fatalf("Failed to connect to broker %s: %s", cfg.BrokerAddress, err)
	}

	subscriber, err = amqp.NewSubscriber(brokerCredentials(cfg))
	if err != nil {
		log.Fatalf("Failed to connect to broker %s: %s", cfg.BrokerAddress, err)
	}

	// exchanges that the publisher uses and publishes messages to
	createExchanges(
		publisher,
		events.TypeChat,
		events.TypeChatTeam,
		events.TypeChatWhisper,
		events.TypeMapChanged,
		events.TypePlayerJoined,
		events.TypePlayerLeft,
		events.TypeVoteKickStarted,
		events.TypeVoteSpecStarted,
		events.TypeVoteOptionStarted,
	)

	// exchanges that the subscriber uses
	createExchanges(
		subscriber,
		topics.Broadcast,
	)

	// get all messages from broadcast
	createQueueAndBindToExchanges(subscriber,
		cfg.EconAddress,
		topics.Broadcast,
	)

	// persistent queues that contain some debugging an dlogging data
	createQueueAndBindToExchanges(subscriber,
		"join-log",
		events.TypePlayerJoined,
	)

	// leave log for testing
	createQueueAndBindToExchanges(subscriber,
		"leave-log",
		events.TypePlayerLeft,
	)

	createQueueAndBindToExchanges(subscriber,
		"vote-log",
		events.TypeVoteKickStarted,
		events.TypeVoteSpecStarted,
	)

	// log map changes for testing
	createQueueAndBindToExchanges(subscriber,
		"map-change-log",
		events.TypeMapChanged,
	)

	// chat logs for testing
	createQueueAndBindToExchanges(subscriber,
		"chat-log",
		events.TypeChat,
		events.TypeChatTeam,
		events.TypeChatWhisper,
	)

}

func main() {
	defer cancel()
	defer conn.Close()
	defer publisher.Close()
	defer subscriber.Close()

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
