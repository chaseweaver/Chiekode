package main

import (
	"fmt"
	"math/rand"

	"github.com/bwmarrin/discordgo"
)

/**
 * misc.go
 * Chase Weaver
 *
 * This package bundles miscellaneous commands for guilds.
 */

// Avatar command will return member's avatar
func Avatar(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {

	// Author avatar if no members are mentioned
	if len(m.Message.Mentions) == 0 {
		s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
			Title: fmt.Sprintf("%s's Avatar", m.Author.Username),
			Color: rand.Intn(32767-(-32768)+1) - 32768,
			Image: &discordgo.MessageEmbedImage{
				URL:    m.Author.AvatarURL("2048"),
				Width:  2048,
				Height: 2048,
			},
			URL: m.Author.AvatarURL("2048"),
		})
	}

	// Returns every mentioned member's avatar
	for u := range m.Message.Mentions {
		s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
			Title: fmt.Sprintf("%s's Avatar", m.Message.Mentions[u].Username),
			Color: rand.Intn(32767-(-32768)+1) - 32768,
			Image: &discordgo.MessageEmbedImage{
				URL:    m.Message.Mentions[u].AvatarURL("2048"),
				Width:  2048,
				Height: 2048,
			},
			URL: m.Message.Mentions[u].AvatarURL("2048"),
		})
	}

	return
}

