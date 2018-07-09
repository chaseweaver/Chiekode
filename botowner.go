package main

import (
	"strings"

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
		UserPermissions: []string{"BotOwner"},
		ArgsDelim:       " ",
		ArgsUsage:       "<golang expression>",
		Description:     "Evaluation command for bot-owner only.",
	})
}

// Ping command will return Pong!
func Ping(ctx Context) {
	ctx.Session.ChannelMessageSend(ctx.Channel.ID, "🏓 Pong!")
}

// Eval is the bot's evaluate command for complex functions
func Eval(ctx Context) {
	out, err := golpal.New().AddLibs("strings", "runtime", "fmt", "github.com/bwmarrin/discordgo").ExecuteRaw(strings.Join(ctx.Args, " "))
	if err != nil {
		ctx.Session.ChannelMessageSend(ctx.Channel.ID, "**ERROR**\n"+FormatString(err.Error(), "go"))
	} else {
		ctx.Session.ChannelMessageSend(ctx.Channel.ID, "**RESULT**\n"+FormatString(out, "go"))
	}
}
