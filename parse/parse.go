package parse

import (
	"github.com/Teeworlds-Server-Moderation/common/concurrent"
	"github.com/Teeworlds-Server-Moderation/common/mqtt"
)

var (
	// dummy used as empty return value
	emptyMsg = mqtt.Message{}

	// ServerState is modified by individual calls to Parsing functions
	// This state can be retrieved concurrently and will represent the current playerlist
	ServerState = concurrent.NewServerState()
)
