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
func Reply(s *discordgo.Session, m *discordgo.MessageCreate, msg string) {
	r := fmt.Sprintf("<@!%s>, %s", m.Author.ID, msg)
	s.ChannelMessageSend(m.ChannelID, r)
	return
}

// FormatString adds string formatting (i.e. asciidoc)
func FormatString(s string, t string) string {
	return fmt.Sprintf("```%s\n" + s + "```", t) 
}

// LogCommands logs commands being run
func LogCommands(s *discordgo.Session, m *discordgo.MessageCreate, cmd string, args []string) {
	guild, err := s.State.Guild(m.GuildID)
	if err != nil {
		guild, err = s.Guild(m.GuildID)
		if err != nil {
			return
		}
	}
	log.Printf(
		"\n"+
		"Guild:     %s / %s\n"+
		"User:      %s / %s\n"+
		"Command:   %s\n"+
		"Args:      %s"+
		"\n\n",
		guild.Name, m.GuildID, m.Author.Username+m.Author.Discriminator, m.Author.ID, cmd, args)
	return
}