package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

/**
 * misc.go
 * Chase Weaver
 *
 * This package bundles miscellaneous commands for guilds.
 */

func init() {
	RegisterNewCommand(Command{
		Name:            "help",
		Func:            Help,
		Enabled:         true,
		NSFWOnly:        false,
		IgnoreSelf:      true,
		IgnoreBots:      true,
		RunIn:           []string{"Text", "DM"},
		Aliases:         []string{},
		UserPermissions: []string{},
		ArgsDelim:       " ",
		ArgsUsage:       "[command]",
		Description:     "Displays a helpful help menu.",
	})

	RegisterNewCommand(Command{
		Name:            "avatar",
		Func:            Avatar,
		Enabled:         true,
		NSFWOnly:        false,
		IgnoreSelf:      true,
		IgnoreBots:      true,
		RunIn:           []string{"Text", "DM"},
		Aliases:         []string{"pfp", "icon"},
		UserPermissions: []string{},
		ArgsDelim:       " ",
		ArgsUsage:       "[@member]",
		Description:     "Fetches the avatar/pfp for the requested member.",
	})
}

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
}

// Avatar command will return member's avatar
func Avatar(ctx Context) {

	mem := FetchMessageContentUsers(ctx)

	// Author avatar if no members are mentioned
	if len(ctx.Args) == 0 {
		ctx.Session.ChannelMessageSendEmbed(ctx.Channel.ID, &discordgo.MessageEmbed{
			Title: fmt.Sprintf("%s's Avatar", ctx.Event.Author.Username),
			Color: RandomInt(0, 16777215),
			Image: &discordgo.MessageEmbedImage{
				URL:    ctx.Event.Author.AvatarURL("2048"),
				Width:  2048,
				Height: 2048,
			},
			URL: ctx.Event.Author.AvatarURL("2048"),
		})
	}

	// Returns every mentioned member's avatar
	for u := range mem {
		ctx.Session.ChannelMessageSendEmbed(ctx.Channel.ID, &discordgo.MessageEmbed{
			Title: fmt.Sprintf("%s's Avatar", mem[u].Username),
			Color: RandomInt(0, 16777215),
			Image: &discordgo.MessageEmbedImage{
				URL:    mem[u].AvatarURL("2048"),
				Width:  2048,
				Height: 2048,
			},
			URL: mem[u].AvatarURL("2048"),
		})
	}

	if len(mem) == 0 && len(ctx.Args) != 0 {
		ctx.Session.ChannelMessageSend(ctx.Channel.ID, "I cannot find that user!")
	}
}
