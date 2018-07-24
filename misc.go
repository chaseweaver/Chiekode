package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/gomodule/redigo/redis"
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
		RunIn:           []string{"Text"},
		Aliases:         []string{},
		UserPermissions: []string{},
		ArgsDelim:       " ",
		Usage:           []string{"[Command Name]"},
		Description:     "Displays a helpful help menu for all commands, or just one.",
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
		Usage:           []string{"[@Member(s)|ID(s)|Name(s)]"},
		Description:     "Fetches the avatar for the requested member, or command author.",
	})
}

// Help :
// Returns help per all-basis or per command-basis
func Help(ctx Context) {

	if len(ctx.Args) == 0 {

		help := "== Helpful Help Menu ==\n\n"
		tmp := []string{}

		// Fetch commands the user has permissions to run
		for _, v := range commands {
			if len(v.UserPermissions) == 0 {
				tmp = append(tmp, fmt.Sprintf("%-10s::  %s", v.Name, v.Description))
			} else {
				isValid := false
				for _, k := range v.UserPermissions {
					if MemberHasPermission(ctx, k) {
						isValid = true
					}
				}
				if isValid {
					tmp = append(tmp, fmt.Sprintf("%-10s::  %s", v.Name, v.Description))
				}
			}
		}

		// Sort commands by alphabetical order
		sort.Strings(tmp)

		// Index command numbers
		index := 1
		for _, v := range tmp {
			help += fmt.Sprintf("%2d. %s\n", index, v)
			index++
		}

		// Creates DM channel between bot and message author
		channel, err := ctx.Session.UserChannelCreate(ctx.Event.Author.ID)

		// Return if the user has DMs blocked
		if err != nil {
			ctx.Session.ChannelMessageSend(ctx.Channel.ID, fmt.Sprintf("<@%s>, I cannot DM you the commands because either I am blocked or you do not accept messages!", ctx.Event.Author.ID))
			log.Println(err)
			return
		}

		// Send a message to the channel indicating the command list has been sent to the command author
		msg, err := ctx.Session.ChannelMessageSend(ctx.Channel.ID, fmt.Sprintf("<@%s>, 📥 I have sent the command list to your inbox!", ctx.Event.Author.ID))

		// DM the command author
		_, err = ctx.Session.ChannelMessageSend(channel.ID, FormatString(help, "asciidoc"))

		if err != nil {
			log.Println(err)
		}

		// Delete the bot response
		DeleteMessageWithTime(ctx, msg.ID, 7500)

		if err != nil {
			log.Println(err)
			return
		}
	} else {

		// Fetch command from message args
		cmd := FetchCommand(strings.Join(ctx.Args, ctx.Command.ArgsDelim))

		// Return if the args cannot find the requested command
		if cmd.isEmpty() {
			msg, err := ctx.Session.ChannelMessageSend(ctx.Channel.ID, fmt.Sprintf("`%s` is not a valid command!", strings.Join(ctx.Args, ctx.Command.ArgsDelim)))

			if err != nil {
				log.Println(err)
			}

			DeleteMessageWithTime(ctx, msg.ID, 3000)
			return
		}

		// Return if the user doesn't have permissions to run said command
		isValid := false
		if len(cmd.UserPermissions) == 0 {
			isValid = true
		} else {
			for _, v := range cmd.UserPermissions {
				if MemberHasPermission(ctx, v) {
					isValid = true
				}
			}
		}

		if !isValid {
			return
		}

		// Get guild information from database
		data, err := redis.Bytes(p.Do("GET", ctx.Guild.ID))
		if err != nil {
			log.Println(err)
		}

		var g Guild
		err = json.Unmarshal(data, &g)

		if err != nil {
			log.Println(err)
		}

		runIn := strings.Join(cmd.RunIn, ", ")
		aliases := strings.Join(cmd.Aliases, ", ")
		permissions := strings.Join(cmd.UserPermissions, ", ")
		usage := g.GuildPrefix + cmd.Name + " " + strings.Join(cmd.Usage, cmd.ArgsDelim)

		if len(cmd.Aliases) == 0 {
			aliases = "N/A"
		}

		if len(cmd.UserPermissions) == 0 {
			permissions = "N/A"
		}

		help := fmt.Sprintf("== %s Help ==\n\nName        :: %s\nRuns In     :: %s\nAliases     :: %s\nPermissions :: %s\nUsage       :: %s\nDescription :: %s", cmd.Name, cmd.Name, runIn, aliases, permissions, usage, cmd.Description)

		// Creates DM channel between bot and message author
		channel, err := ctx.Session.UserChannelCreate(ctx.Event.Author.ID)

		// Return if the user has DMs blocked
		if err != nil {
			ctx.Session.ChannelMessageSend(ctx.Channel.ID, fmt.Sprintf("<@%s>, I cannot DM you the commands because either I am blocked or you do not accept messages!", ctx.Event.Author.ID))
			log.Println(err)
			return
		}

		// Send a message to the channel indicating the command list has been sent to the command author
		msg, err := ctx.Session.ChannelMessageSend(ctx.Channel.ID, fmt.Sprintf("<@%s>, 📥 I have sent the command list to your inbox!", ctx.Event.Author.ID))

		// DM the command author
		_, err = ctx.Session.ChannelMessageSend(channel.ID, FormatString(help, "asciidoc"))

		if err != nil {
			log.Println(err)
		}

		// Delete the bot response
		DeleteMessageWithTime(ctx, msg.ID, 7500)

		if err != nil {
			log.Println(err)
			return
		}

	}
}

// Avatar :
// Returns User's Avatar by ID / Name / Mention.
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
