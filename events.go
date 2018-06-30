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

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if !strings.HasPrefix(m.Content, conf.Prefix) {
		return
	}

	msg := strings.TrimPrefix(m.Content, "+")
	nme := strings.Split(msg, " ")[0]
	cmd := FetchCommand(nme)

	if !CheckValidPrereq(s, m, FetchCommand(nme)) {
		return
	}

	args := strings.Split(msg, cmd.ArgsDelim)[1:]

	// Fetch command funcs from command properties init()
	funcs := map[string]interface{}{
		FetchCommand(nme).Name: FetchCommand(nme).Func,
	}

	LogCommands(s, m, FetchCommandName(nme), args)
	Call(funcs, FetchCommandName(nme), s, m, args)
}
