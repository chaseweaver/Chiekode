package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/gomodule/redigo/redis"
)

/**
 * mod.go
 * Chase Weaver
 *
 * This package bundles commands for the moderation commands of the bot.
 */

func init() {
	RegisterNewCommand(Command{
		Name:            "warn",
		Func:            Warn,
		Enabled:         true,
		NSFWOnly:        false,
		IgnoreSelf:      true,
		IgnoreBots:      true,
		RunIn:           []string{"Text"},
		Aliases:         []string{},
		UserPermissions: []string{"Bot Owner", "Administrator", "Kick Members"},
		ArgsDelim:       " ",
		ArgsUsage:       "<@member(s)|ID(s)>",
		Description:     "Warns a member via mention or ID.",
	})

	RegisterNewCommand(Command{
		Name:            "kick",
		Func:            Kick,
		Enabled:         true,
		NSFWOnly:        false,
		IgnoreSelf:      true,
		IgnoreBots:      true,
		RunIn:           []string{"Text"},
		Aliases:         []string{},
		UserPermissions: []string{"Bot Owner", "Administrator", "Kick Members"},
		ArgsDelim:       " ",
		ArgsUsage:       "<@member(s)|ID(s)>",
		Description:     "Kicks a member via mention or ID.",
	})

	RegisterNewCommand(Command{
		Name:            "ban",
		Func:            Ban,
		Enabled:         true,
		NSFWOnly:        false,
		IgnoreSelf:      true,
		IgnoreBots:      true,
		RunIn:           []string{"Text"},
		Aliases:         []string{},
		UserPermissions: []string{"Bot Owner", "Administrator", "Ban Members"},
		ArgsDelim:       " ",
		ArgsUsage:       "<@member(s)|ID(s)>",
		Description:     "Bans a member via mention or ID.",
	})

	RegisterNewCommand(Command{
		Name:            "check",
		Func:            Check,
		Enabled:         true,
		NSFWOnly:        false,
		IgnoreSelf:      true,
		IgnoreBots:      true,
		RunIn:           []string{"Text"},
		Aliases:         []string{},
		UserPermissions: []string{"Bot Owner", "Administrator", "Kick Members"},
		ArgsDelim:       " ",
		ArgsUsage:       "<@member(s)|ID(s)>",
		Description:     "Checks the warnings, mutes, kicks, and bans of a mentioned user.",
	})
}

// Warn command will warn a mentioned user and log it to the redis database
func Warn(ctx Context) {
	reason := "N/A"

	if len(ctx.Args) != 0 {
		reason = strings.Join(ctx.Args[0:], " ")
	}

	mem := FetchMessageContentUsers(ctx)
	if len(mem) == 0 {
		ctx.Session.ChannelMessageSend(ctx.Channel.ID, "I cannot find that user!")
	}

	for _, usr := range mem {
		LogWarning(ctx, usr, reason)
	}
}

// Kick command will kick a mentioned user and log it to the redis database
func Kick(ctx Context) {
	reason := "N/A"

	if len(ctx.Args) != 0 {
		reason = strings.Join(ctx.Args[0:], " ")
	}

	mem := FetchMessageContentUsers(ctx)
	if len(mem) == 0 {
		ctx.Session.ChannelMessageSend(ctx.Channel.ID, "I cannot find that user!")
	}

	for _, member := range mem {
		ctx.Session.GuildMemberDeleteWithReason(ctx.Guild.ID, member.ID, reason)
		LogKick(ctx, member, reason)
	}
}

// Ban command will ban a mentioned user and log it to the redis database
func Ban(ctx Context) {
	reason := "N/A"

	if len(ctx.Args) != 0 {
		reason = strings.Join(ctx.Args[0:], " ")
	}

	mem := FetchMessageContentUsers(ctx)
	if len(mem) == 0 {
		ctx.Session.ChannelMessageSend(ctx.Channel.ID, "I cannot find that user!")
	}

	for _, member := range mem {
		ctx.Session.GuildBanCreateWithReason(ctx.Guild.ID, member.ID, reason, 0)
		LogBan(ctx, member, reason)
	}
}

// Check command will check the user's warnings, mutes, kicks, bans, nicknames, and usernames, and return it from the redis database
func Check(ctx Context) {

	if len(ctx.Args) == 0 {
		ctx.Session.ChannelMessageSend(ctx.Channel.ID, "Invalid usage! Please `@mention` (a) member(s)!")
	}

	data, err := redis.Bytes(p.Do("GET", ctx.Guild.ID))
	if err != nil {
		log.Println(err)
	}

	var g Guild
	err = json.Unmarshal(data, &g)

	if err != nil {
		log.Println(err)
	}

	for key := range ctx.Event.Message.Mentions {
		mem := ctx.Event.Message.Mentions[key]

		for k := range g.Users {
			if g.Users[k].ID == mem.ID {

				embed := &discordgo.MessageEmbed{
					Title:       fmt.Sprintf("%s#%s / %s", mem.Username, mem.Discriminator, mem.ID),
					Color:       RandomInt(0, 16777215),
					Description: fmt.Sprintf("Run `%scheck @member [warnings/mutes/kicks/bans]` for a complete list of information.", g.GuildPrefix),
					Thumbnail: &discordgo.MessageEmbedThumbnail{
						URL:    mem.AvatarURL("2048"),
						Width:  2048,
						Height: 2048,
					},
					Fields: []*discordgo.MessageEmbedField{
						&discordgo.MessageEmbedField{
							Name:   fmt.Sprintf("⮞ Total Warnings"),
							Value:  fmt.Sprintf("%d", len(g.Users[k].Warnings)),
							Inline: false,
						},
						&discordgo.MessageEmbedField{
							Name:   fmt.Sprintf("⮞ Total Mutes"),
							Value:  fmt.Sprintf("%d", len(g.Users[k].Mutes)),
							Inline: false,
						},
						&discordgo.MessageEmbedField{
							Name:   fmt.Sprintf("⮞ Total Kicks"),
							Value:  fmt.Sprintf("%d", len(g.Users[k].Kicks)),
							Inline: false,
						},
						&discordgo.MessageEmbedField{
							Name:   fmt.Sprintf("⮞ Total Bans"),
							Value:  fmt.Sprintf("%d", len(g.Users[k].Bans)),
							Inline: false,
						},
					},
					Timestamp: time.Now().Format(time.RFC3339),
				}

				ctx.Session.ChannelMessageSendEmbed(ctx.Channel.ID, embed)
			}
		}

	}
}
