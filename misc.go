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
func Help(ctx Context) {
	if len(ctx.Args) > 0 {
		c := FetchCommand(ctx.Args[0])
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
		ctx.Session.ChannelMessageSend(ctx.Channel.ID, FormatString(str, "asciidoc"))
	} else {
		str := "== Help ==\n\n"
		for k := range commands {
			str += fmt.Sprintf("%s:\n\t%s\n", commands[k].Name, commands[k].Description)
		}
		ctx.Session.ChannelMessageSend(ctx.Channel.ID, FormatString(str, "asciidoc"))
	}
	return
}

// Avatar command will return member's avatar
func Avatar(ctx Context) {

	// Author avatar if no members are mentioned
	if len(ctx.Event.Message.Mentions) == 0 {
		ctx.Session.ChannelMessageSendEmbed(ctx.Channel.ID, &discordgo.MessageEmbed{
			Title: fmt.Sprintf("%s's Avatar", ctx.Event.Author.Username),
			Color: rand.Intn(32767-(-32768)+1) - 32768,
			Image: &discordgo.MessageEmbedImage{
				URL:    ctx.Event.Author.AvatarURL("2048"),
				Width:  2048,
				Height: 2048,
			},
			URL: ctx.Event.Author.AvatarURL("2048"),
		})
	}

	// Returns every mentioned member's avatar
	for u := range ctx.Event.Message.Mentions {
		ctx.Session.ChannelMessageSendEmbed(ctx.Channel.ID, &discordgo.MessageEmbed{
			Title: fmt.Sprintf("%s's Avatar", ctx.Event.Message.Mentions[u].Username),
			Color: rand.Intn(32767-(-32768)+1) - 32768,
			Image: &discordgo.MessageEmbedImage{
				URL:    ctx.Event.Message.Mentions[u].AvatarURL("2048"),
				Width:  2048,
				Height: 2048,
			},
			URL: ctx.Event.Message.Mentions[u].AvatarURL("2048"),
		})
	}
	return
}
