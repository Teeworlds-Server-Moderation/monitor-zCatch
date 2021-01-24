package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Teeworlds-Server-Moderation/common/mqtt"
	"github.com/joho/godotenv"
	configo "github.com/jxsl13/simple-configo"
	"github.com/jxsl13/twapi/econ"
)

var (
	config = &Config{}
)

func init() {
	var env map[string]string
	env, err := godotenv.Read()
	if err != nil {
		log.Fatalf("Failed to get environment variables: %s", err)
	}

	err = configo.Parse(config, env)
	if err != nil {
		log.Fatalf("Invalid configutaion parameters provided:\n%s", err)
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	conn, err := econ.DialTo(config.EconAddress, config.Password)
	if err != nil {
		log.Fatalf("Failed to connect to or authenticate at %s: %s", config.EconAddress, err)
	}
	defer conn.Close()
	log.Println("Successfully connected to Econ: ", config.EconAddress)

	err = conn.WriteLine("ec_output_level 2")
	if err != nil {
		log.Fatalln("Failed to set: ec_output_level 2")
	}

	publisher, err := mqtt.NewPublisher(config.BrokerAddress, config.EconAddress+"monitor", "")
	if err != nil {
		log.Fatalf("Failed to connect to broker %s: %s", config.BrokerAddress, err)
	}

	lineChan := make(chan string, 64)

	go econLineReader(ctx, conn, lineChan)
	go eventProducerRoutine(ctx, config.EconAddress, lineChan, publisher)

	// Messages will be delivered asynchronously so we just need to wait for a signal to shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	fmt.Println("Startup finished, running...")
	<-sig
	fmt.Println("signal caught - exiting")
	fmt.Println("shutdown complete")
}
