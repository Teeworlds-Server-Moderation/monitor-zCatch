package main

import (
	"fmt"
	"strings"

	"github.com/Teeworlds-Server-Moderation/common/events"
	"github.com/jxsl13/twapi/econ"
)

var (
	econEscape = strings.NewReplacer(
		"#",
		"_",
	)
)

// econExecAndLog executes and logs the executed command with their corresponding author
// in the econ ia the echo command. This is necessary in order to have proper logging of every command executed
// and in order to see abuse or potential bugs where people may bypass access control.
func econExecAndLog(conn *econ.Conn, event events.RequestCommandExecEvent) error {
	user := econEscape.Replace(event.Requestor)
	err := conn.WriteLine(fmt.Sprintf("echo User '%s' executed command '%s'", user, event.Command))
	if err != nil {
		return err
	}
	err = conn.WriteLine(event.Command)
	if err != nil {
		return err
	}
	return nil
}
