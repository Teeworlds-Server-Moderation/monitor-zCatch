package parse

import (
	"time"

	"github.com/Teeworlds-Server-Moderation/common/amqp"
	"github.com/Teeworlds-Server-Moderation/common/concurrent"
)

var (
	// dummy used as empty return value
	emptyMsg = amqp.Message{}

	// ServerState is modified by individual calls to Parsing functions
	// This state can be retrieved concurrently and will represent the current playerlist
	ServerState = concurrent.NewServerState()
)

func formatedTimestamp() string {
	return time.Now().Format("2006-01-02T15:04:05.999999-07:00")
}
