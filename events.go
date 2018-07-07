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
		session: s,
		event:   m,
		guild:   guild,
		channel: channel,
		name:    strings.Split(strings.TrimPrefix(m.Content, conf.Prefix), " ")[0],
	}

	ctx.command = FetchCommand(ctx.name)
	ctx.args = strings.Split(ctx.event.Content, ctx.command.ArgsDelim)[1:]

	if !CheckValidPrereq(ctx.session, ctx.event, ctx.command) {
		return
	}

	// Fetch command funcs from command properties init()
	funcs := map[string]interface{}{
		ctx.command.Name: ctx.command.Func,
	}

	LogCommands(ctx)
	Call(funcs, ctx.name, ctx)
}
