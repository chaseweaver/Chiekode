package main

import (
	"fmt"
	"time"
)

/**
 * owner.go
 * Chase Weaver
 *
 * This package bundles commands for the owner of the bot.
 */

func init() {
	RegisterNewCommand(Command{
		Name:            "ping",
		Func:            Ping,
		Enabled:         true,
		NSFWOnly:        false,
		IgnoreSelf:      true,
		IgnoreBots:      true,
		RunIn:           []string{"Text", "DM"},
		Aliases:         []string{},
		UserPermissions: []string{},
		ArgsDelim:       "",
		Usage:           []string{},
		Description:     "Pong! Responds with the heartbeat.",
	})

	RegisterNewCommand(Command{
		Name:            "test",
		Func:            Test,
		Enabled:         true,
		NSFWOnly:        false,
		IgnoreSelf:      true,
		IgnoreBots:      true,
		RunIn:           []string{"DM", "Text"},
		Aliases:         []string{},
		UserPermissions: []string{"Bot Owner"},
		ArgsDelim:       " ",
		Usage:           []string{},
		Description:     "Bot-owner testing function.",
	})
}

// Ping :
// Command will return Pong! with the last heartbeat.
func Ping(ctx Context) {
	t := time.Now().Sub(ctx.Session.LastHeartbeatAck) / 1000
	ctx.Session.ChannelMessageSend(ctx.Channel.ID, fmt.Sprintf("🏓 Pong! Heatbeat: `%s`", t))
}

// Test :
// Bot owner's test command
func Test(ctx Context) {
	return
}
