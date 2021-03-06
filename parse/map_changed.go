package parse

import (
	"fmt"
	"regexp"

	"github.com/Teeworlds-Server-Moderation/common/amqp"
	"github.com/Teeworlds-Server-Moderation/common/events"
)

var (

	// 0: Full
	// 1: map/
	// 2: ctf5_carribean.map
	// [2021-03-06 14:50:16][server]: maps/ctf5_carribbean.map crc is 814ce0a4
	mapChangedRegexp = regexp.MustCompile(`(.*\/)([^/]+\.map) crc is .+`)
)

// MapChange returns event messages when the logLine contains the proper line.
func MapChange(source, timestamp, logLine string) ([]amqp.Message, error) {
	match := mapChangedRegexp.FindStringSubmatch(logLine)
	if len(match) == 0 {
		return nil, fmt.Errorf("%w: MapChange: %s", ErrInvalidLineFormat, logLine)
	}

	// we get multiple events from the map change
	oldMap := ServerState.GetMap()
	newMap := match[2]
	leftPlayers := ServerState.PlayerLeaveAll()

	// n left player events + map change event
	eventMessages := make([]amqp.Message, 0, len(leftPlayers)+1)

	for _, player := range leftPlayers {
		// player left event for every player
		event := events.NewPlayerLeftEvent()
		event.SetEventSource(source)
		event.Player = player
		event.Reason = "map change"

		eventMessages = append(eventMessages, amqp.Message{
			Exchange: event.Type,
			Payload:  event.Marshal(),
		})
	}

	// map change event
	event := events.NewMapChanedEvent()
	event.SetEventSource(source)
	event.OldMap = oldMap
	event.NewMap = newMap

	eventMessages = append(eventMessages, amqp.Message{
		Exchange: event.Type,
		Payload:  event.Marshal(),
	})

	return eventMessages, nil
}
