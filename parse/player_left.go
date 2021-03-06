package parse

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/Teeworlds-Server-Moderation/common/amqp"
	"github.com/Teeworlds-Server-Moderation/common/events"
)

var (
	// 0: full 1: ID 2: IP 3: reason
	playerLeftRegex = regexp.MustCompile(`id=([\d]+) addr=([a-fA-F0-9\.\:\[\]]+) reason='(.*)'$`)
)

// PlayerLeft parses potential leaving players with as much information as possible.
// Any empty struct field will be set to the default empty value.
func PlayerLeft(source, timestamp, logLine string) ([]amqp.Message, error) {
	match := playerLeftRegex.FindStringSubmatch(logLine)
	if len(match) == 0 {
		return nil, fmt.Errorf("%w: PlayerLeft: %s", ErrInvalidLineFormat, logLine)
	}
	id, _ := strconv.Atoi(match[1])
	reason := match[3]

	event := events.NewPlayerLeftEvent()
	event.SetEventSource(source)

	player := ServerState.PlayerLeave(id)
	event.Player = player
	event.Reason = reason

	msg := amqp.Message{
		Exchange: event.Type,
		Payload:  event.Marshal(),
	}
	return toMsgList(msg, nil)
}
