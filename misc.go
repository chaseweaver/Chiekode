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
func Help(ctx *Context) {
	if len(ctx.args) > 0 {
		c := FetchCommand(ctx.args[0])
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
		ctx.session.ChannelMessageSend(ctx.channel.ID, FormatString(str, "asciidoc"))
	} else {
		str := "== Help ==\n\n"
		for k := range commands {
			str += fmt.Sprintf("%s:\n\t%s\n", commands[k].Name, commands[k].Description)
		}
		ctx.session.ChannelMessageSend(ctx.channel.ID, FormatString(str, "asciidoc"))
	}
	return
}

// Avatar command will return member's avatar
func Avatar(ctx *Context) {

	// Author avatar if no members are mentioned
	if len(ctx.event.Message.Mentions) == 0 {
		ctx.session.ChannelMessageSendEmbed(ctx.channel.ID, &discordgo.MessageEmbed{
			Title: fmt.Sprintf("%s's Avatar", ctx.event.Author.Username),
			Color: rand.Intn(32767-(-32768)+1) - 32768,
			Image: &discordgo.MessageEmbedImage{
				URL:    ctx.event.Author.AvatarURL("2048"),
				Width:  2048,
				Height: 2048,
			},
			URL: ctx.event.Author.AvatarURL("2048"),
		})
	}

	// Returns every mentioned member's avatar
	for u := range ctx.event.Message.Mentions {
		ctx.session.ChannelMessageSendEmbed(ctx.channel.ID, &discordgo.MessageEmbed{
			Title: fmt.Sprintf("%s's Avatar", ctx.event.Message.Mentions[u].Username),
			Color: rand.Intn(32767-(-32768)+1) - 32768,
			Image: &discordgo.MessageEmbedImage{
				URL:    ctx.event.Message.Mentions[u].AvatarURL("2048"),
				Width:  2048,
				Height: 2048,
			},
			URL: ctx.event.Message.Mentions[u].AvatarURL("2048"),
		})
	}
	return
}
