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
	// 1: ID Voter
	// 2: Name Voter
	// 3: ID Victim
	// 4: Name Victim
	// 5: Reason
	// 6: Command
	// 7: Forced
	startVoteSpecRegex = regexp.MustCompile(`'([\d]{1,2}):(.*)' voted spectate '([\d]{1,2}):(.*)' reason='(.{1,20})' cmd='(.*)' force=([\d])`)
)

func StartVoteSpec(source, timestamp, logLine string) (amqp.Message, error) {
	match := startVoteSpecRegex.FindStringSubmatch(logLine)
	if len(match) == 0 {
		return emptyMsg, fmt.Errorf("Invalid StartVoteSpec line format: %s", logLine)
	}

	idVoter, _ := strconv.Atoi(match[1])
	idVictim, _ := strconv.Atoi(match[3])
	forced, _ := strconv.Atoi(match[7])

	voteSpecStartEvent := events.NewVoteSpecStartedEvent()
	voteSpecStartEvent.Timestamp = timestamp
	voteSpecStartEvent.EventSource = source
	voteSpecStartEvent.Source = dto.Player{
		ID:   idVoter,
		Name: match[2],
	}
	voteSpecStartEvent.Target = dto.Player{
		ID:   idVictim,
		Name: match[4],
	}
	voteSpecStartEvent.Forced = forced != 0

	msg := amqp.Message{
		Queue:   events.TypeVoteKickStarted,
		Payload: voteSpecStartEvent.Marshal(),
	}
	return msg, nil
}
