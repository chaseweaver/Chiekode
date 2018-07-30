package main

import (
	"fmt"
	"log"
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
		Cooldown:        30,
		RunIn:           []string{"Text", "DM"},
		Aliases:         []string{},
		UserPermissions: []string{},
		ArgsDelim:       "",
		Usage:           []string{},
		Description:     "Pong! Responds with the heartbeat.",
	})

	RegisterNewCommand(Command{
		Name:            "resetguilddatabase",
		Func:            ResetGuildDatabase,
		Enabled:         true,
		NSFWOnly:        false,
		IgnoreSelf:      true,
		IgnoreBots:      true,
		Cooldown:        0,
		RunIn:           []string{"DM", "Text"},
		Aliases:         []string{},
		UserPermissions: []string{"Bot Owner"},
		ArgsDelim:       " ",
		Usage:           []string{},
		Description:     "CAUTION! Flushes the database and reinitializes guild settings!",
	})

	RegisterNewCommand(Command{
		Name:            "test",
		Func:            Test,
		Enabled:         true,
		NSFWOnly:        false,
		IgnoreSelf:      true,
		IgnoreBots:      true,
		Cooldown:        0,
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
	ctx.Session.ChannelMessageSend(ctx.Channel.ID, fmt.Sprintf("🏓 | Pong! Heatbeat: `%s`", t))
}

// ResetGuildDatabase :
// Deletes all redis keys, reinitializes guilds
func ResetGuildDatabase(ctx Context) {

	_, err := p.Do("FLUSHDB")
	if err != nil {
		log.Println(err)
	}

	for _, v := range ctx.Session.State.Guilds {
		RegisterNewGuild(v)
	}

	ctx.Session.ChannelMessageSend(ctx.Channel.ID, "✅ | All guilds purged from database. Guild settings have been reset.")
}

// Test :
// Bot owner's test command
func Test(ctx Context) {
	return
}
