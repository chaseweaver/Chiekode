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
	ctx.session.ChannelMessageSend(ctx.channel.ID, "🏓 Pong!")
}

// Eval is the bot's evaluate command
func Eval(ctx Context) {
	out, err := golpal.New().Execute(strings.Join(ctx.args, ""))
	if err != nil {
		FormatString("**RESULT**\n"+err.Error(), "golang")
	}
	ctx.session.ChannelMessageSend(ctx.channel.ID, FormatString("**RESULT**\n"+out, "golang"))
}
