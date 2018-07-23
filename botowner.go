package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/novalagung/golpal"
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
		ArgsUsage:       "",
		Description:     "Pong!",
	})

	RegisterNewCommand(Command{
		Name:            "eval",
		Func:            Eval,
		Enabled:         true,
		NSFWOnly:        false,
		IgnoreSelf:      true,
		IgnoreBots:      true,
		RunIn:           []string{"Text", "DM"},
		Aliases:         []string{"e"},
		UserPermissions: []string{"Bot Owner"},
		ArgsDelim:       " ",
		ArgsUsage:       "<golang expression>",
		Description:     "Evaluation command for bot-owner only.",
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
		ArgsUsage:       "",
		Description:     "Bot owner testing function.",
	})
}

// Ping :
// Command will return Pong! with the last heartbeat.
func Ping(ctx Context) {
	t := time.Now().Sub(ctx.Session.LastHeartbeatAck) / 1000
	ctx.Session.ChannelMessageSend(ctx.Channel.ID, fmt.Sprintf("🏓 Pong! Heatbeat: `%s`", t))
}

// Eval :
// Bot's evaluate command for complex functions.
func Eval(ctx Context) {
	out, err := golpal.New().AddLibs("strings", "runtime", "fmt", "github.com/bwmarrin/discordgo").ExecuteRaw(strings.Join(ctx.Args, " "))
	if err != nil {
		ctx.Session.ChannelMessageSend(ctx.Channel.ID, "**ERROR**\n"+FormatString(err.Error(), "go"))
	} else {
		ctx.Session.ChannelMessageSend(ctx.Channel.ID, "**RESULT**\n"+FormatString(out, "go"))
	}
}

// Test :
// Bot owner's test command
func Test(ctx Context) {
	return
}
