package parse

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/Teeworlds-Server-Moderation/common/amqp"
	"github.com/Teeworlds-Server-Moderation/common/dto"
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

func StartVoteOption(source, timestamp, logLine string) (amqp.Message, error) {
	match := startVoteOptionRegex.FindStringSubmatch(logLine)
	if len(match) == 0 {
		return emptyMsg, fmt.Errorf("Invalid StartVoteOption line format: %s", logLine)
	}

	id, _ := strconv.Atoi(match[1])
	forced, _ := strconv.Atoi(match[6])

	voteSpecStartEvent := events.NewVoteSpecStartedEvent()
	voteSpecStartEvent.Timestamp = formatedTimestamp()
	voteSpecStartEvent.EventSource = source
	voteSpecStartEvent.Source = dto.Player{
		ID:   id,
		Name: match[2],
	}

	voteSpecStartEvent.Forced = forced != 0

	msg := amqp.Message{
		Queue:   events.TypeVoteKickStarted,
		Payload: voteSpecStartEvent.Marshal(),
	}
	return msg, nil
}
