package parse

import (
	"errors"

	"github.com/Teeworlds-Server-Moderation/common/amqp"
	"github.com/Teeworlds-Server-Moderation/common/concurrent"
)

var (

	// ServerState is modified by individual calls to Parsing functions
	// This state can be retrieved concurrently and will represent the current playerlist
	ServerState = concurrent.NewServerState()

	// ErrInvalidLineFormat is returned by parsing functions that cannot parse a given line
	ErrInvalidLineFormat = errors.New("invalid line format")
)

// Handler is the basic function signature of parser functions
type Handler []func(string, string, string) ([]amqp.Message, error)

// toMsgList creates a list from a single amqp.Message
func toMsgList(msg amqp.Message, err error) ([]amqp.Message, error) {
	return []amqp.Message{msg}, err
}
