package main

import (
	"log"
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
	if !strings.HasPrefix(m.Content, conf.Prefix) {
		return
	}

	channel, err := s.State.Channel(m.ChannelID)
	if err != nil {
		log.Println("Could not get source channel,", err)
		return
	}

	guild, err := s.State.Guild(channel.GuildID)
	if err != nil {
		log.Println("Could not get source guild,", err)
		return
	}

	// Give context for command pass-in
	ctx := Context{
		Session: s,
		Event:   m,
		Guild:   guild,
		Channel: channel,
		Name:    strings.Split(strings.TrimPrefix(m.Content, conf.Prefix), " ")[0],
	}

	ctx.Command = FetchCommand(ctx.Name)
	ctx.Args = strings.Split(ctx.Event.Content, ctx.Command.ArgsDelim)[1:]

	if !CheckValidPrereq(ctx.Session, ctx.Event, ctx.Command) {
		return
	}

	// Fetch command funcs from command properties init()
	funcs := map[string]interface{}{
		ctx.Command.Name: ctx.Command.Func,
	}

	LogCommands(ctx)
	Call(funcs, ctx.Name, ctx)
}
