package parse

import (
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
