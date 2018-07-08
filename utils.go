package main

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
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
	if ctx.Channel.Type == discordgo.ChannelTypeGuildText {
		log.Printf(
			"\n"+
				"User:      %s / %s\n"+
				"Guild:     %s / %s\n"+
				"Channel:   %s / %s\n"+
				"Command:   %s\n"+
				"Args:      %s"+
				"\n\n",
			ctx.Event.Author.Username+"#"+ctx.Event.Author.Discriminator, ctx.Event.Author.ID,
			ctx.Guild.Name, ctx.Guild.ID, ctx.Channel.Name, ctx.Channel.ID,
			ctx.Name, ctx.Args)
	} else {
		log.Printf(
			"\n"+
				"User:      %s / %s\n"+
				"DM:        %s\n"+
				"Command:   %s\n"+
				"Args:      %s"+
				"\n\n",
			ctx.Event.Author.Username+"#"+ctx.Event.Author.Discriminator,
			ctx.Event.Author.ID, ctx.Channel.ID, ctx.Name, ctx.Args)
	}
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
