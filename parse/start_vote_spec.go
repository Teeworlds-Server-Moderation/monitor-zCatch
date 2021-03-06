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
	// 1: ID Voter
	// 2: Name Voter
	// 3: ID Victim
	// 4: Name Victim
	// 5: Reason
	// 6: Command
	// 7: Forced
	startVoteSpecRegex = regexp.MustCompile(`'([\d]{1,2}):(.*)' voted spectate '([\d]{1,2}):(.*)' reason='(.{1,20})' cmd='(.*)' force=([\d])`)
)

// StartVoteSpec returns event messages when the logLine contains the proper line.
func StartVoteSpec(source, timestamp, logLine string) ([]amqp.Message, error) {
	match := startVoteSpecRegex.FindStringSubmatch(logLine)
	if len(match) == 0 {
		return nil, fmt.Errorf("%w: StartVoteSpec: %s", ErrInvalidLineFormat, logLine)
	}

	idVoter, _ := strconv.Atoi(match[1])
	idVictim, _ := strconv.Atoi(match[3])
	reason := match[5]
	forced, _ := strconv.Atoi(match[7])

	event := events.NewVoteSpecStartedEvent()
	event.SetEventSource(source)

	event.Source = ServerState.GetPlayer(idVoter)
	event.Target = ServerState.GetPlayer(idVictim)
	event.Reason = reason
	event.Forced = forced != 0

	msg := amqp.Message{
		Exchange: event.Type,
		Payload:  event.Marshal(),
	}
	return toMsgList(msg, nil)
}
