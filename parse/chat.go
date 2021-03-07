package parse

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/Teeworlds-Server-Moderation/common/amqp"
	"github.com/Teeworlds-Server-Moderation/common/events"
)

var (

	// 0: full
	// 1: ID
	// 2: nick
	// 3: text
	chatRegex = regexp.MustCompile(`([\d]+):[\d]+:(.{1,16}): (.*)$`)
)

// Chat returns event messages when the logLine contains the proper line.
func Chat(source, timestamp, logLine string) ([]amqp.Message, error) {
	match := chatRegex.FindStringSubmatch(logLine)
	if len(match) == 0 {
		return nil, fmt.Errorf("%w: Chat: %s", ErrInvalidLineFormat, logLine)
	}

	id, _ := strconv.Atoi(match[1])
	text := match[3]
	player := ServerState.GetPlayer(id)

	// chat event
	event := events.NewChatEvent()
	event.SetEventSource(source)
	event.Source = player
	event.Text = text

	msg := amqp.Message{
		Exchange: event.Type,
		Payload:  event.Marshal(),
	}

	return toMsgList(msg, nil)
}
