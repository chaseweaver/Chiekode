package main

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

/**
 * events.go
 * Chase Weaver
 *
 * This package bundles event commands when they are triggered.
 */

// MessageCreate triggers on a message that is visible to the bot
func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Checks if message content begins with prefix
	if !strings.HasPrefix(m.Content, conf.Prefix) {
		return
	}

	// Fetches channel object
	channel, err := s.State.Channel(m.ChannelID)
	if err != nil {
		return
	}

	// Give context for command pass-in
	ctx := Context{
		Session: s,
		Event:   m,
		Channel: channel,
		Name:    strings.Split(strings.TrimPrefix(m.Content, conf.Prefix), " ")[0],
	}

	// Fetches guild object if text channel is NOT a DM
	if ctx.Channel.Type == discordgo.ChannelTypeGuildText {
		guild, err := s.State.Guild(ctx.Channel.GuildID)
		if err != nil {
			return
		}
		ctx.Guild = guild
	}

	// Returns a valid command using a name/alias
	ctx.Command = FetchCommand(ctx.Name)

	// Splits command arguments
	ctx.Args = strings.Split(ctx.Event.Content, ctx.Command.ArgsDelim)[1:]

	// Checks if the config for the command passes all checks
	if !CommandIsValid(ctx) {
		return
	}

	// Fetch command funcs from command properties init()
	funcs := map[string]interface{}{
		ctx.Command.Name: ctx.Command.Func,
	}

	// Log commands to console
	LogCommands(ctx)

	// Call command with args pass-in
	Call(funcs, FetchCommandName(ctx.Name), ctx)
}
