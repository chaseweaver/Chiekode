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
	ctx.Session.ChannelMessageSend(ctx.Channel.ID, fmt.Sprintf("<@!%s>, %s", ctx.Event.Author.ID, s))
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
		ctx.Guild.Name, ctx.Guild.ID,
		ctx.Event.Author.Username+"#"+ctx.Event.Author.Discriminator,
		ctx.Event.Author.ID, ctx.Name, ctx.Args)
	return
}

// Contains checks if element is in array
func Contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}
