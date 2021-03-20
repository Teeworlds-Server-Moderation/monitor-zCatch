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
	// 1: source ID
	// 2: target ID
	// 3: text
	whisperChatRegex = regexp.MustCompile(`([\d]+):([\d]+):.{1,16}: (.*)$`)
)

// ChatWhisper returns event messages when the logLine contains the proper line.
func ChatWhisper(source, timestamp, logLine string) ([]amqp.Message, error) {
	match := whisperChatRegex.FindStringSubmatch(logLine)
	if len(match) == 0 {
		return nil, fmt.Errorf("%w: ChatWhisper: %s", ErrInvalidLineFormat, logLine)
	}

	sourceID, _ := strconv.Atoi(match[1])
	targetID, _ := strconv.Atoi(match[2])
	text := match[3]
	sourcePlayer := ServerState.GetPlayer(sourceID)
	targetPlayer := ServerState.GetPlayer(targetID)

	// chat team event
	event := events.NewChatWhisperEvent()
	event.SetEventSource(source)
	event.Source = sourcePlayer
	event.Target = targetPlayer
	event.Text = text

	msg := amqp.Message{
		Exchange: event.Type,
		Payload:  event.Marshal(),
	}

	return toMsgList(msg, nil)
}
