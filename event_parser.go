package main

import (
	"fmt"
	"regexp"

	"github.com/Teeworlds-Server-Moderation/common/amqp"
	"github.com/Teeworlds-Server-Moderation/monitor-zcatch/parse"
)

var (
	// [2020-05-22 23:01:09][client_enter]: id=0 addr=192.168.178.25:64139 version=1796 name='MisterFister:(' clan='FistingTea`' country=-1
	// 0: full 1: timestamp 2: log level 3: log line
	initialLoglevelRegex = regexp.MustCompile(`^\[([\d\s-:]+)\]\[([^:]+)\]: (.+)$`)

	// dummy used as empty return value
	emptyMsg = amqp.Message{}
)

// different handler functions that handle specific
var serverLogLevelHandlers = []func(string, string, string) ([]amqp.Message, error){
	parse.MapChange,
	parse.StartVoteKick,
	parse.StartVoteSpec,
	parse.StartVoteOption,
}

var gameLogLevelHandlers = []func(string, string, string) (amqp.Message, error){}

// handle allows the homogenous handling of the above defined paarser function lists
func handle(source, timestamp, logLine string, parserList []func(string, string, string) ([]amqp.Message, error)) ([]amqp.Message, error) {
	var err error
	for _, handler := range parserList {
		msg, err := handler(source, timestamp, logLine)
		if err == nil {
			return msg, nil
		}
	}
	return nil, err
}

// returns a message or an error in case something went wrong
func parseEvent(source, line string) ([]amqp.Message, error) {
	matches := initialLoglevelRegex.FindStringSubmatch(line)
	if len(matches) == 0 {
		return nil, fmt.Errorf("%s: %s", parse.ErrInvalidLineFormat, line)
	}

	timestamp := matches[1]
	logLevel := matches[2]
	logLine := matches[3]

	switch logLevel {
	case "client_enter":
		return parse.PlayerJoined(source, timestamp, logLine)
	case "client_drop":
		return parse.PlayerLeft(source, timestamp, logLine)
	case "server":
		return handle(source, timestamp, logLine, serverLogLevelHandlers)
	}
	return nil, fmt.Errorf("Unknown log level: %s", logLevel)
}
