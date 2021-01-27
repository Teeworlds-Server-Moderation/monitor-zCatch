package parse

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/Teeworlds-Server-Moderation/common/events"
	"github.com/Teeworlds-Server-Moderation/common/mqtt"
)

var (
	// 0: full 1: ID 2: IP 3: reason
	playerLeftRegex = regexp.MustCompile(`id=([\d]+) addr=([a-fA-F0-9\.\:\[\]]+) reason='(.*)'$`)
)

// PlayerLeft parses potential leaving players with as much information as possible.
// Any empty struct field will be set to the default empty value.
func PlayerLeft(source, timestamp, logLine string) (mqtt.Message, error) {
	match := playerLeftRegex.FindStringSubmatch(logLine)
	if len(match) == 0 {
		return emptyMsg, fmt.Errorf("Invalid PlayerLeft line format: %s", logLine)
	}
	id, _ := strconv.Atoi(match[1])
	reason := match[3]

	playerLeftEvent := events.NewPlayerLeftEvent()
	playerLeftEvent.EventSource = source
	playerLeftEvent.Timestamp = timestamp

	player := ServerState.PlayerLeave(id)
	playerLeftEvent.Player = player
	playerLeftEvent.Reason = reason

	msg := mqtt.Message{
		Topic:   events.TypePlayerLeft,
		Payload: playerLeftEvent.Marshal(),
	}
	return msg, nil
}
