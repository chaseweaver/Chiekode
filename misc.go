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
		ArgsUsage:       "[@member|ID]",
		Description:     "Fetches the avatar/pfp for the requested member.",
	})
}

// Help command returns help per all-basis or per command-basis
func Help(ctx Context) {
	if len(ctx.Args) > 0 {

		c := FetchCommand(ctx.Command.Name)
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

// Avatar :
// Returns User's Avatar by ID / Name / Mention
func Avatar(ctx Context) {

	// Returns the command author's avatar if no arguments are given, or the command is within a DM
	if ctx.Channel.Type != discordgo.ChannelTypeGuildText || len(ctx.Args) == 0 {
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
		return
	}

	// Fetch users from message content
	members := FetchMessageContentUsers(ctx, strings.Join(ctx.Args, ctx.Command.ArgsDelim))

	// Returns every mentioned member's avatar as seperate messages
	for _, m := range members {
		ctx.Session.ChannelMessageSendEmbed(ctx.Channel.ID, &discordgo.MessageEmbed{
			Title: fmt.Sprintf("%s's Avatar", m.Username),
			Color: RandomInt(0, 16777215),
			Image: &discordgo.MessageEmbedImage{
				URL:    m.AvatarURL("2048"),
				Width:  2048,
				Height: 2048,
			},
			URL: m.AvatarURL("2048"),
		})
	}

	// Return if the message is not empty and no users can be found
	if len(members) == 0 && len(ctx.Args) != 0 {
		ctx.Session.ChannelMessageSend(ctx.Channel.ID, "I cannot find that user!")
	}
}
