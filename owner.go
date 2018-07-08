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

// Ping command will return Pong!
func Ping(ctx Context) {
	ctx.Session.ChannelMessageSend(ctx.Channel.ID, "🏓 Pong!")
}

// Eval is the bot's evaluate command
func Eval(ctx Context) {
	out, err := golpal.New().Execute(strings.Join(ctx.Args, " "))
	if err != nil {
		ctx.Session.ChannelMessageSend(ctx.Channel.ID, "**ERROR**\n"+FormatString(err.Error(), "go"))
	} else {
		ctx.Session.ChannelMessageSend(ctx.Channel.ID, "**RESULT**\n"+FormatString(out, "go"))
	}
}
