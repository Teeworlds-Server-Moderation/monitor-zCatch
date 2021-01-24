package main

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/Teeworlds-Server-Moderation/common/events"
	"github.com/Teeworlds-Server-Moderation/common/mqtt"
)

var (
	playerEnteredRegex = regexp.MustCompile(`id=([\d]+) addr=([a-fA-F0-9\.\:\[\]]+):([\d]+) version=(\d+) name='(.{0,20})' clan='(.{0,16})' country=([-\d]+)$`)
)

func parseClientEnter(source, timestamp, logLine string) (mqtt.Message, error) {
	match := playerEnteredRegex.FindStringSubmatch(logLine)
	if len(match) != 8 {
		return emptyMsg, fmt.Errorf("Invalid ClientEnter line format: %s", logLine)
	}
	port, _ := strconv.Atoi(match[3])
	id, _ := strconv.Atoi(match[1])
	country, _ := strconv.Atoi(match[7])
	version, _ := strconv.Atoi(match[4])

	playerJoinEvent := events.NewEventPlayerJoin(
		source,
		timestamp,
		match[5],
		match[6],
		match[2],
		port,
		id,
		country,
		version,
	)

	msg := mqtt.Message{
		Topic:   events.TypePlayerJoin,
		Payload: playerJoinEvent.Marshal(),
	}
	return msg, nil
}
