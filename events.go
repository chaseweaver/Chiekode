package main

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/gomodule/redigo/redis"
)

/**
 * events.go
 * Chase Weaver
 *
 * This package bundles event commands when they are triggered.
 */

// MessageCreate :
// Triggers on a message that is visible to the bot
// Handles message and command responses
func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Default bot prefix
	prefix := conf.Prefix

	// Fetches channel object
	channel, err := s.State.Channel(m.ChannelID)
	if err != nil {
		return
	}

	// Gets the guild prefix from database
	if channel.Type == discordgo.ChannelTypeGuildText {
		guild, err := s.State.Guild(channel.GuildID)

		if err != nil {
			log.Println(err)
		}

		data, err := redis.Bytes(p.Do("GET", guild.ID))

		if err != nil {
			log.Println(err)
		}

		var g Guild
		err = json.Unmarshal(data, &g)
		prefix = g.GuildPrefix
	}

	// Checks if message content begins with prefix
	if !strings.HasPrefix(m.Content, prefix) {
		return
	}

	// Give context for command pass-in
	ctx := Context{
		Session: s,
		Event:   m,
		Channel: channel,
		Name:    strings.Split(strings.TrimPrefix(m.Content, prefix), " ")[0],
	}

	// Fetches guild object if text channel is NOT a DM
	if ctx.Channel.Type == discordgo.ChannelTypeGuildText {
		guild, err := s.State.Guild(ctx.Channel.GuildID)
		if err != nil {
			return
		}

		ctx.Guild = guild

		// Registers a new guild if not done already
		RegisterNewGuild(ctx.Guild)
	}

	// Returns a valid command using a name/alias
	ctx.Command = FetchCommand(ctx.Name)

	// Splits command arguments
	tmp := strings.TrimPrefix(m.Content, prefix)

	// Returns all Members, Channels, and Roles by Mention, ID, and Name, and removes them from the string
	ctx.Args = strings.Split(tmp, ctx.Command.ArgsDelim)[1:]

	// Checks if the config for the command passes all checks and is part of a text channel in a guild
	if ctx.Channel.Type == discordgo.ChannelTypeGuildText {
		if !CommandIsValid(ctx) {
			return
		}
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

// GuildCreate :
// Initializes a new guild when the bot is first added
func GuildCreate(s *discordgo.Session, m *discordgo.GuildCreate) {

	// Register new guild
	_, err := RegisterNewGuild(m.Guild)

	if err != nil {
		log.Panicln(err)
	}

	log.Println("== New Guild Added ==\nGuild Name: %s\nGuild ID:   %s\n", m.Guild.Name, m.Guild.ID)
}

// GuildDelete :
// Initializes a new guild when the bot is removed from a guild
func GuildDelete(s *discordgo.Session, m *discordgo.GuildDelete) {

	// Delete guild key
	_, err := DeleteGuild(m.Guild)

	if err != nil {
		log.Panicln(err)
	}

	log.Println("== Guild Removed ==\nGuild Name: %s\nGuild ID:   %s\n", m.Guild.Name, m.Guild.ID)
}
