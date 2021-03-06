package parse

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/Teeworlds-Server-Moderation/common/amqp"
	"github.com/Teeworlds-Server-Moderation/common/events"
)

var (

	// 0: Full
	// 1: ID
	// 2: Name
	// 3: Option
	// 4: Reason
	// 5: Command
	// 6: Forced
	startVoteOptionRegex = regexp.MustCompile(`'([\d]{1,2}):(.*)' voted option '(.+)' reason='(.{1,20})' cmd='(.+)' force=([\d])`)
)

// StartVoteOption returns event messages when the logLine contains the proper line.
func StartVoteOption(source, timestamp, logLine string) ([]amqp.Message, error) {
	match := startVoteOptionRegex.FindStringSubmatch(logLine)
	if len(match) == 0 {
		return nil, fmt.Errorf("%w: StartVoteOption: %s", ErrInvalidLineFormat, logLine)
	}

	id, _ := strconv.Atoi(match[1])
	reason := match[4]
	option := match[5]
	forced, _ := strconv.Atoi(match[6])

	event := events.NewVoteOptionStartedEvent()
	event.SetEventSource(source)

	event.Source = ServerState.GetPlayer(id)
	event.Reason = reason
	event.Option = option
	event.Forced = forced != 0

	msg := amqp.Message{
		Exchange: event.Type,
		Payload:  event.Marshal(),
	}
	return toMsgList(msg, nil)
}
