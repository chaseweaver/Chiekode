package main

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/bwmarrin/discordgo"
)

/**
 * misc.go
 * Chase Weaver
 *
 * This package bundles miscellaneous commands for guilds.
 */

// Help command returns help per all-basis or per command-basis
func Help(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) > 0 {
		c := FetchCommand(args[0])
		str := fmt.Sprintf("== %s ==\n", c.Name)
		tmp := c.ArgsDelim
		na := c.Name

		if c.ArgsDelim == " " {
			tmp = "[SPACE]"
		}

		if len(c.Aliases) > 0 {
			na = c.Name + "|" + strings.Join(c.Aliases, "|")
		}

		str += fmt.Sprintf(
			"Command     ::  %s\n"+
				"Description ::  %s\n"+
				"Usage       ::  %s\n"+
				"Run In      ::  %s\n"+
				"Arg Delim   ::  %s\n",
			c.Name,
			c.Description,
			fmt.Sprintf("%s<%s> %s", conf.Prefix, na, c.ArgsUsage),
			strings.Join(c.RunIn, ", "), tmp)
		s.ChannelMessageSend(m.ChannelID, FormatString(str, "asciidoc"))
	} else {
		str := "== Help ==\n\n"
		for k := range commands {
			str += fmt.Sprintf("%s:\n\t%s\n", commands[k].Name, commands[k].Description)
		}
		s.ChannelMessageSend(m.ChannelID, FormatString(str, "asciidoc"))
	}
	return
}

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
