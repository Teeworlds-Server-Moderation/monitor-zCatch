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
	teamChatRegex = regexp.MustCompile(`([\d]+):[\d]+:(.{1,16}): (.*)$`)
)

// ChatTeam returns event messages when the logLine contains the proper line.
func ChatTeam(source, timestamp, logLine string) ([]amqp.Message, error) {
	match := teamChatRegex.FindStringSubmatch(logLine)
	if len(match) == 0 {
		return nil, fmt.Errorf("%w: ChatTeam: %s", ErrInvalidLineFormat, logLine)
	}

	id, _ := strconv.Atoi(match[1])
	text := match[3]
	player := ServerState.GetPlayer(id)

	// chat team event
	event := events.NewChatTeamEvent()
	event.SetEventSource(source)
	event.Source = player
	event.Text = text

	msg := amqp.Message{
		Exchange: event.Type,
		Payload:  event.Marshal(),
	}

	return toMsgList(msg, nil)
}
