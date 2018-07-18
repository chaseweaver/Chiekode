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
		Name:            "lock",
		Func:            Lock,
		Enabled:         true,
		NSFWOnly:        false,
		IgnoreSelf:      true,
		IgnoreBots:      true,
		RunIn:           []string{"Text"},
		Aliases:         []string{"lockdown"},
		UserPermissions: []string{"Bot Owner", "Administrator", "Manage Roles", "Manage Channels"},
		ArgsDelim:       "",
		ArgsUsage:       "",
		Description:     "Locks a channel (prevent message sending) for the default @everyone permission.",
	})

	RegisterNewCommand(Command{
		Name:            "unlock",
		Func:            Unlock,
		Enabled:         true,
		NSFWOnly:        false,
		IgnoreSelf:      true,
		IgnoreBots:      true,
		RunIn:           []string{"Text"},
		Aliases:         []string{},
		UserPermissions: []string{"Bot Owner", "Administrator", "Manage Roles", "Manage Channels"},
		ArgsDelim:       "",
		ArgsUsage:       "",
		Description:     "Unlocks a channel (allows message sending) for the default @everyone permission.",
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
		ArgsUsage:       "<@member(s)|ID(s)> [warnings|mutes|kicks|bans]",
		Description:     "Checks the warnings, mutes, kicks, and bans of a mentioned user.",
	})
}

// Warn command will warn a mentioned user and log it to the redis database
func Warn(ctx Context) {

	mem := FetchMessageContentUsers(ctx)
	reason := RemoveMessageIDs(strings.Join(ctx.Args[0:], " "))

	if len(reason) == 0 {
		reason = "N/A"
	}

	if len(mem) == 0 {
		msg, err := ctx.Session.ChannelMessageSend(ctx.Channel.ID, "I cannot find that user!")

		if err != nil {
			log.Println(err)
		}

		// Delete author message
		DeleteMessageWithTime(ctx, ctx.Event.Message.ID, 0)

		// Delete bot message
		DeleteMessageWithTime(ctx, msg.ID, 7500)
	}

	DeleteMessageWithTime(ctx, ctx.Event.Message.ID, 0)

	for _, member := range mem {
		n1 := member.Username + "#" + member.Discriminator
		n2 := ctx.Event.Message.Author.Username + "#" + ctx.Event.Message.Author.Discriminator
		channel, err := ctx.Session.UserChannelCreate(member.ID)

		ctx.Session.ChannelMessageSend(ctx.Channel.ID, fmt.Sprintf("`%s` has been warned by `%s`", n1, n2))
		LogWarning(ctx, member, reason)

		if err != nil {
			log.Println(err)
			break
		}

		if reason != "N/A" {
			_, err = ctx.Session.ChannelMessageSend(channel.ID, fmt.Sprintf("You have been warned by `%s` with reason `%s`", ctx.Event.Message.Author.Username+"#"+ctx.Event.Message.Author.Discriminator, reason))

			if err != nil {
				log.Println(err)
				break
			}
		} else {
			_, err = ctx.Session.ChannelMessageSend(channel.ID, fmt.Sprintf("You have been warned by `%s` without a reason.", ctx.Event.Message.Author.Username+"#"+ctx.Event.Message.Author.Discriminator))

			if err != nil {
				log.Println(err)
				break
			}
		}

	}
}

// Kick command will kick a mentioned user and log it to the redis database
func Kick(ctx Context) {

	mem := FetchMessageContentUsers(ctx)
	reason := RemoveMessageIDs(strings.Join(ctx.Args[0:], " "))

	if len(reason) == 0 {
		reason = "N/A"
	}

	if len(mem) == 0 {
		msg, err := ctx.Session.ChannelMessageSend(ctx.Channel.ID, "I cannot find that user!")

		if err != nil {
			log.Println(err)
		}

		// Delete author message
		DeleteMessageWithTime(ctx, ctx.Event.Message.ID, 0)

		// Delete bot message
		DeleteMessageWithTime(ctx, msg.ID, 7500)
	}

	DeleteMessageWithTime(ctx, ctx.Event.Message.ID, 0)

	for _, member := range mem {
		err = ctx.Session.GuildMemberDeleteWithReason(ctx.Guild.ID, member.ID, reason)

		if err != nil {
			msg, _ := ctx.Session.ChannelMessageSend(ctx.Channel.ID, "I cannot kick this user!")
			DeleteMessageWithTime(ctx, ctx.Event.Message.ID, 0)
			DeleteMessageWithTime(ctx, msg.ID, 7500)
			break
		}

		n1 := member.Username + "#" + member.Discriminator
		n2 := ctx.Event.Message.Author.Username + "#" + ctx.Event.Message.Author.Discriminator
		channel, err := ctx.Session.UserChannelCreate(member.ID)

		ctx.Session.ChannelMessageSend(ctx.Channel.ID, fmt.Sprintf("`%s` has been kicked by `%s`", n1, n2))
		LogKick(ctx, member, reason)

		if err != nil {
			log.Println(err)
			break
		}

		if reason != "N/A" {
			_, err = ctx.Session.ChannelMessageSend(channel.ID, fmt.Sprintf("You have been kicked by `%s` with reason `%s`", ctx.Event.Message.Author.Username+"#"+ctx.Event.Message.Author.Discriminator, reason))

			if err != nil {
				log.Println(err)
				break
			}
		} else {
			_, err = ctx.Session.ChannelMessageSend(channel.ID, fmt.Sprintf("You have been kicked by `%s` without a reason.", ctx.Event.Message.Author.Username+"#"+ctx.Event.Message.Author.Discriminator))

			if err != nil {
				log.Println(err)
				break
			}
		}

	}
}

// Ban command will ban a mentioned user and log it to the redis database
func Ban(ctx Context) {

	mem := FetchMessageContentUsers(ctx)
	reason := RemoveMessageIDs(strings.Join(ctx.Args[0:], " "))

	if len(reason) == 0 {
		reason = "N/A"
	}

	if len(mem) == 0 {
		msg, err := ctx.Session.ChannelMessageSend(ctx.Channel.ID, "I cannot find that user!")

		if err != nil {
			log.Println(err)
		}

		// Delete author message
		DeleteMessageWithTime(ctx, ctx.Event.Message.ID, 0)

		// Delete bot message
		DeleteMessageWithTime(ctx, msg.ID, 7500)
	}

	DeleteMessageWithTime(ctx, ctx.Event.Message.ID, 0)

	for _, member := range mem {
		err = ctx.Session.GuildBanCreateWithReason(ctx.Guild.ID, member.ID, reason, 0)

		if err != nil {
			msg, _ := ctx.Session.ChannelMessageSend(ctx.Channel.ID, "I cannot ban this user!")
			DeleteMessageWithTime(ctx, ctx.Event.Message.ID, 0)
			DeleteMessageWithTime(ctx, msg.ID, 7500)
			break
		}

		n1 := member.Username + "#" + member.Discriminator
		n2 := ctx.Event.Message.Author.Username + "#" + ctx.Event.Message.Author.Discriminator
		channel, err := ctx.Session.UserChannelCreate(member.ID)

		ctx.Session.ChannelMessageSend(ctx.Channel.ID, fmt.Sprintf("`%s` has been banned by `%s`", n1, n2))
		LogBan(ctx, member, reason)

		if err != nil {
			log.Println(err)
			break
		}

		if reason != "N/A" {
			_, err = ctx.Session.ChannelMessageSend(channel.ID, fmt.Sprintf("You have been banned by `%s` with reason `%s`", ctx.Event.Message.Author.Username+"#"+ctx.Event.Message.Author.Discriminator, reason))

			if err != nil {
				log.Println(err)
				break
			}
		} else {
			_, err = ctx.Session.ChannelMessageSend(channel.ID, fmt.Sprintf("You have been banned by `%s` without a reason.", ctx.Event.Message.Author.Username+"#"+ctx.Event.Message.Author.Discriminator))

			if err != nil {
				log.Println(err)
				break
			}
		}

	}
}

// Lock will lock a channel (prevent message sending) for the default @everyone permission
func Lock(ctx Context) {

	// This is the default @everyone role that cannot change index, hence the hardcode
	everyone := ctx.Guild.Roles[0]

	err := ctx.Session.ChannelPermissionSet(ctx.Channel.ID, everyone.ID, "1", 0, discordgo.PermissionSendMessages)
	DeleteMessageWithTime(ctx, ctx.Event.Message.ID, 0)

	if err != nil {
		log.Println(err)
	}

	_, err = ctx.Session.ChannelMessageSend(ctx.Channel.ID, fmt.Sprintf("This channel is now under lockdown!"))

	if err != nil {
		log.Println(err)
	}
}

// Unlock will unlock a channel (allows message sending) for the default @everyone permission
func Unlock(ctx Context) {

	// This is the default @everyone role that cannot change index, hence the hardcode
	everyone := ctx.Guild.Roles[0]

	err := ctx.Session.ChannelPermissionSet(ctx.Channel.ID, everyone.ID, "1", discordgo.PermissionSendMessages, 0)
	DeleteMessageWithTime(ctx, ctx.Event.Message.ID, 0)

	if err != nil {
		log.Println(err)
	}

	_, err = ctx.Session.ChannelMessageSend(ctx.Channel.ID, fmt.Sprintf("This channel is no longer under lockdown!"))

	if err != nil {
		log.Println(err)
	}
}

// Check command will check the user's warnings, mutes, kicks, bans, nicknames, and usernames, and return it from the redis database
func Check(ctx Context) {

	mem := FetchMessageContentUsers(ctx)
	gtype := strings.ToUpper(RemoveMessageIDs(strings.Join(ctx.Args, " ")))

	if len(mem) == 0 {
		msg, err := ctx.Session.ChannelMessageSend(ctx.Channel.ID, "I cannot find that user!")

		if err != nil {
			log.Println(err)
		}

		// Delete bot message
		DeleteMessageWithTime(ctx, msg.ID, 7500)
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

	switch gtype {
	case "WARN":
		fallthrough
	case "WARNS":
		fallthrough
	case "WARNING":
		fallthrough
	case "WARNINGS":

		for _, member := range mem {
			for _, usr := range g.Users {
				if usr.ID == member.ID {

					embed := &discordgo.MessageEmbed{
						Title: fmt.Sprintf("Warning Stats [%d]", len(usr.Warnings)),
						Color: warningColor,
						Author: &discordgo.MessageEmbedAuthor{
							URL:     member.AvatarURL("2048"),
							IconURL: member.AvatarURL("256"),
							Name:    fmt.Sprintf("%s#%s / %s", member.Username, member.Discriminator, member.ID),
						},
						Thumbnail: &discordgo.MessageEmbedThumbnail{
							URL:    member.AvatarURL("2048"),
							Width:  2048,
							Height: 2048,
						},
						Description: FormatWarning(usr.Warnings),
						Timestamp:   time.Now().Format(time.RFC3339),
					}
					_, err := ctx.Session.ChannelMessageSendEmbed(ctx.Channel.ID, embed)

					if err != nil {
						return
					}

				}
			}
		}

	case "KICK":
		fallthrough
	case "KICKS":
		for _, member := range mem {
			for _, usr := range g.Users {
				if usr.ID == member.ID {

					embed := &discordgo.MessageEmbed{
						Title: fmt.Sprintf("Kick Stats [%d]", len(usr.Kicks)),
						Color: kickColor,
						Author: &discordgo.MessageEmbedAuthor{
							URL:     member.AvatarURL("2048"),
							IconURL: member.AvatarURL("256"),
							Name:    fmt.Sprintf("%s#%s / %s", member.Username, member.Discriminator, member.ID),
						},
						Thumbnail: &discordgo.MessageEmbedThumbnail{
							URL:    member.AvatarURL("2048"),
							Width:  2048,
							Height: 2048,
						},
						Description: FormatKick(usr.Kicks),
						Timestamp:   time.Now().Format(time.RFC3339),
					}
					_, err := ctx.Session.ChannelMessageSendEmbed(ctx.Channel.ID, embed)

					if err != nil {
						return
					}

				}
			}
		}

	case "BAN":
		fallthrough
	case "BANS":
		for _, member := range mem {
			for _, usr := range g.Users {
				if usr.ID == member.ID {

					embed := &discordgo.MessageEmbed{
						Title: fmt.Sprintf("Ban Stats [%d]", len(usr.Bans)),
						Color: banColor,
						Author: &discordgo.MessageEmbedAuthor{
							URL:     member.AvatarURL("2048"),
							IconURL: member.AvatarURL("256"),
							Name:    fmt.Sprintf("%s#%s / %s", member.Username, member.Discriminator, member.ID),
						},
						Thumbnail: &discordgo.MessageEmbedThumbnail{
							URL:    member.AvatarURL("2048"),
							Width:  2048,
							Height: 2048,
						},
						Description: FormatBan(usr.Bans),
						Timestamp:   time.Now().Format(time.RFC3339),
					}
					_, err := ctx.Session.ChannelMessageSendEmbed(ctx.Channel.ID, embed)

					if err != nil {
						return
					}

				}
			}
		}

	default:
		for _, member := range mem {
			for _, usr := range g.Users {
				if usr.ID == member.ID {

					embed := &discordgo.MessageEmbed{
						Title:       fmt.Sprintf("%s#%s / %s", member.Username, member.Discriminator, member.ID),
						Color:       RandomInt(0, 16777215),
						Description: fmt.Sprintf("Run `%scheck @member [warnings/mutes/kicks/bans]` for a complete list of information.", g.GuildPrefix),
						Thumbnail: &discordgo.MessageEmbedThumbnail{
							URL:    member.AvatarURL("2048"),
							Width:  2048,
							Height: 2048,
						},
						Fields: []*discordgo.MessageEmbedField{
							&discordgo.MessageEmbedField{
								Name:   fmt.Sprintf("⮞ Total Warnings"),
								Value:  fmt.Sprintf("%d", len(usr.Warnings)),
								Inline: false,
							},
							&discordgo.MessageEmbedField{
								Name:   fmt.Sprintf("⮞ Total Mutes"),
								Value:  fmt.Sprintf("%d", len(usr.Mutes)),
								Inline: false,
							},
							&discordgo.MessageEmbedField{
								Name:   fmt.Sprintf("⮞ Total Kicks"),
								Value:  fmt.Sprintf("%d", len(usr.Kicks)),
								Inline: false,
							},
							&discordgo.MessageEmbedField{
								Name:   fmt.Sprintf("⮞ Total Bans"),
								Value:  fmt.Sprintf("%d", len(usr.Bans)),
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
}
