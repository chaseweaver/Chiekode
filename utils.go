package main

import (
	"fmt"
	"log"
)

/**
 * utils.go
 * Chase Weaver
 *
 * This package handles various utilities for shorthands and logging.
 */

// Reply shorthand
func Reply(ctx Context, s string) {
	ctx.session.ChannelMessageSend(ctx.channel.ID, fmt.Sprintf("<@!%s>, %s", ctx.event.Author.ID, s))
}

// FormatString adds string formatting (i.e. asciidoc)
func FormatString(s string, t string) string {
	return fmt.Sprintf("```%s\n"+s+"```", t)
}

// LogCommands logs commands being run
func LogCommands(ctx Context) {
	log.Printf(
		"\n"+
			"Guild:     %s / %s\n"+
			"User:      %s / %s\n"+
			"Command:   %s\n"+
			"Args:      %s"+
			"\n\n",
		ctx.guild.Name, ctx.guild.ID,
		ctx.event.Author.Username+ctx.event.Author.Discriminator,
		ctx.event.Author.ID, ctx.name, ctx.args)
	return
}
