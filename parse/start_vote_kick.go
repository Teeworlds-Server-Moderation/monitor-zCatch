package parse

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/Teeworlds-Server-Moderation/common/dto"
	"github.com/Teeworlds-Server-Moderation/common/events"
	"github.com/Teeworlds-Server-Moderation/common/mqtt"
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
	startVoteKickRegex = regexp.MustCompile(`'([\d]{1,2}):(.*)' voted kick '([\d]{1,2}):(.*)' reason='(.{1,20})' cmd='(.*)' force=([\d])`)
)

func StartVoteKick(source, timestamp, logLine string) (mqtt.Message, error) {
	match := startVoteKickRegex.FindStringSubmatch(logLine)
	if len(match) == 0 {
		return emptyMsg, fmt.Errorf("Invalid StartVoteKick line format: %s", logLine)
	}

	idVoter, _ := strconv.Atoi(match[1])
	idVictim, _ := strconv.Atoi(match[3])
	forced, _ := strconv.Atoi(match[7])

	voteKickStartEvent := events.NewVoteKickStartedEvent()
	voteKickStartEvent.Timestamp = timestamp
	voteKickStartEvent.EventSource = source
	voteKickStartEvent.Source = dto.Player{
		ID:   idVoter,
		Name: match[2],
	}
	voteKickStartEvent.Target = dto.Player{
		ID:   idVictim,
		Name: match[4],
	}
	voteKickStartEvent.Forced = forced != 0

	msg := mqtt.Message{
		Topic:   events.TypeVoteKickStarted,
		Payload: voteKickStartEvent.Marshal(),
	}
	return msg, nil
}
