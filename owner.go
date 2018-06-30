package main

import (
	"github.com/bwmarrin/discordgo"
)

/**
 * owner.go
 * Chase Weaver
 *
 * This package bundles commands for the owner of the bot.
 */

// Ping command will return Pong!
func Ping(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	s.ChannelMessageSend(m.ChannelID, "🏓 Pong!")
	return
}