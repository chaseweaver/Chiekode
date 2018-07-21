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
		ArgsUsage:       "<@member(s)|ID(s)|Name(s)>",
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
		ArgsUsage:       "<@member(s)|ID(s)|Name(s)>",
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
		ArgsUsage:       "<@member(s)|ID(s)|Name(s)>",
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
		ArgsUsage:       "<@member(s)|ID(s)|Name(s)> [warnings|mutes|kicks|bans]",
		Description:     "Checks the warnings, mutes, kicks, and bans of a mentioned user.",
	})
}

// Warn :
// Warn a user by ID / Name#xxxx / Mention, logs it to the redis database.
func Warn(ctx Context) {

	// Fetch users from message content, returns list of members and the remaining string with the member removed
	members, reason := FetchMessageContentUsersString(ctx, strings.Join(ctx.Args, ctx.Command.ArgsDelim))

	log.Println(ctx.Args)
	log.Println(strings.Join(ctx.Args, ctx.Command.ArgsDelim))
	log.Println(members)

	// Returns if a user cannot be found in the message, deletes command message, then deletes delayed response
	if len(members) == 0 {
		msg, err := ctx.Session.ChannelMessageSend(ctx.Channel.ID, "I cannot find that user!")

		if err != nil {
			log.Println(err)
		}

		DeleteMessageWithTime(ctx, ctx.Event.Message.ID, 0)
		DeleteMessageWithTime(ctx, msg.ID, 7500)
		return
	}

	// Delete command message
	DeleteMessageWithTime(ctx, ctx.Event.Message.ID, 0)

	// If no reason is specified, set one for the database logger
	if len(reason) == 0 {
		reason = "N/A"
	}

	// Warns all members found within the message, logs warning to redis database
	for _, member := range members {

		// Target username
		target := member.Username + "#" + member.Discriminator

		// Author username
		author := ctx.Event.Message.Author.Username + "#" + ctx.Event.Message.Author.Discriminator

		// Creates DM channel between bot and target
		channel, err := ctx.Session.UserChannelCreate(member.ID)

		// Sends warn message to channel the command was instantiated in
		ctx.Session.ChannelMessageSend(ctx.Channel.ID, fmt.Sprintf("`%s` has been warned by `%s`", target, author))

		// Logs warning to redis database
		LogWarning(ctx, member, reason)

		// Exits the loop if the user has DMs blocked
		if err != nil {
			break
		}

		// Sends a DM to the user with the warning information if the user can accept DMs
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

// Kick :
// Kicks a user by ID / Name#xxxx / Mention, logs it to the redis database.
func Kick(ctx Context) {

	// Fetch users from message content, returns list of members and the remaining string with the member removed
	members, reason := FetchMessageContentUsersString(ctx, strings.Join(ctx.Args, ctx.Command.ArgsDelim))

	// Returns if a user cannot be found in the message, deletes command message, then deletes delayed response
	if len(members) == 0 {
		msg, err := ctx.Session.ChannelMessageSend(ctx.Channel.ID, "I cannot find that user!")

		if err != nil {
			log.Println(err)
		}

		DeleteMessageWithTime(ctx, ctx.Event.Message.ID, 0)
		DeleteMessageWithTime(ctx, msg.ID, 7500)
		return
	}

	// Delete command message
	DeleteMessageWithTime(ctx, ctx.Event.Message.ID, 0)

	// If no reason is specified, set one for the database logger
	if len(reason) == 0 {
		reason = "N/A"
	}

	// Kicks all members found within the message, logs warning to redis database
	for _, member := range members {

		// Kicks the guild member with given reason
		err = ctx.Session.GuildMemberDeleteWithReason(ctx.Guild.ID, member.ID, reason)

		if err != nil {
			msg, _ := ctx.Session.ChannelMessageSend(ctx.Channel.ID, "I cannot kick this user!")
			DeleteMessageWithTime(ctx, msg.ID, 7500)
			break
		}

		// Target username
		target := member.Username + "#" + member.Discriminator

		// Author username
		author := ctx.Event.Message.Author.Username + "#" + ctx.Event.Message.Author.Discriminator

		// Creates DM channel between bot and target
		channel, err := ctx.Session.UserChannelCreate(member.ID)

		// Sends kick message to channel the command was instantiated in
		ctx.Session.ChannelMessageSend(ctx.Channel.ID, fmt.Sprintf("`%s` has been kicked by `%s`", target, author))

		// Logs kick to redis database
		LogKick(ctx, member, reason)

		// Exits the loop if the user has DMs blocked
		if err != nil {
			break
		}

		// Sends a DM to the user with the kick information if the user can accept DMs
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

// Ban :
// Bans a user by ID / Name#xxxx / Mention, logs it to the redis database.
func Ban(ctx Context) {

	// Fetch users from message content, returns list of members and the remaining string with the member removed
	members, reason := FetchMessageContentUsersString(ctx, strings.Join(ctx.Args, ctx.Command.ArgsDelim))

	// Returns if a user cannot be found in the message, deletes command message, then deletes delayed response
	if len(members) == 0 {
		msg, err := ctx.Session.ChannelMessageSend(ctx.Channel.ID, "I cannot find that user!")

		if err != nil {
			log.Println(err)
		}

		DeleteMessageWithTime(ctx, ctx.Event.Message.ID, 0)
		DeleteMessageWithTime(ctx, msg.ID, 7500)
		return
	}

	// Delete command message
	DeleteMessageWithTime(ctx, ctx.Event.Message.ID, 0)

	// If no reason is specified, set one for the database logger
	if len(reason) == 0 {
		reason = "N/A"
	}

	// Bans all members found within the message, logs warning to redis database
	for _, member := range members {

		// Bans the guild member with given reason, deletes 0 messages
		err = ctx.Session.GuildBanCreateWithReason(ctx.Guild.ID, member.ID, reason, 0)

		if err != nil {
			msg, _ := ctx.Session.ChannelMessageSend(ctx.Channel.ID, "I cannot ban this user!")
			DeleteMessageWithTime(ctx, msg.ID, 7500)
			break
		}

		// Target username
		target := member.Username + "#" + member.Discriminator

		// Author username
		author := ctx.Event.Message.Author.Username + "#" + ctx.Event.Message.Author.Discriminator

		// Creates DM channel between bot and target
		channel, err := ctx.Session.UserChannelCreate(member.ID)

		// Sends ban message to channel the command was instantiated in
		ctx.Session.ChannelMessageSend(ctx.Channel.ID, fmt.Sprintf("`%s` has been banned by `%s`", target, author))

		// Logs ban to redis database
		LogBan(ctx, member, reason)

		// Exits the loop if the user has DMs blocked
		if err != nil {
			break
		}

		// Sends a DM to the user with the ban information if the user can accept DMs
		if reason != "N/A" {
			_, err = ctx.Session.ChannelMessageSend(channel.ID, fmt.Sprintf("You have been banned by `%s` with reason `%s`", author, reason))

			if err != nil {
				log.Println(err)
				break
			}
		} else {
			_, err = ctx.Session.ChannelMessageSend(channel.ID, fmt.Sprintf("You have been banned by `%s` without a reason.", author))

			if err != nil {
				log.Println(err)
				break
			}
		}
	}
}

// Lock :
// Overrides default @everyone permission and prevents "SEND MESSAGES" permission.
func Lock(ctx Context) {

	// Default @everyone role that cannot change index
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

// Unlock :
// Overrides default @everyone permission and allows "SEND MESSAGES" permission.
func Unlock(ctx Context) {

	// Default @everyone role that cannot change index
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

// Check :
// Check the user's warnings, mutes, kicks, bans, nicknames, and usernames from the redis database.
func Check(ctx Context) {

	// Fetch users from message content
	members := FetchMessageContentUsers(ctx, strings.Join(ctx.Args, ctx.Command.ArgsDelim))

	// Type of check to look for
	checkType := strings.Join(ctx.Args, ctx.Command.ArgsDelim)

	// Returns if a user cannot be found in the message, deletes delayed response
	if len(members) == 0 {
		msg, err := ctx.Session.ChannelMessageSend(ctx.Channel.ID, "I cannot find that user!")

		if err != nil {
			log.Println(err)
		}

		DeleteMessageWithTime(ctx, msg.ID, 7500)
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

	switch checkType {
	case "WARN":
		fallthrough
	case "WARNS":
		fallthrough
	case "WARNING":
		fallthrough
	case "WARNINGS":

		for _, member := range members {
			for _, usr := range g.GuildUser {
				if usr.Member.User.ID == member.ID {

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
		for _, member := range members {
			for _, usr := range g.GuildUser {
				if usr.User.ID == member.ID {

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
		for _, member := range members {
			for _, usr := range g.GuildUser {
				if usr.User.ID == member.ID {

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
		for _, member := range members {
			for _, usr := range g.GuildUser {
				if usr.User.ID == member.ID {

					embed := &discordgo.MessageEmbed{
						Title:       fmt.Sprintf("%s#%s / %s", member.Username, member.Discriminator, member.ID),
						Color:       RandomInt(0, 16777215),
						Description: fmt.Sprintf("Run `%scheck <@member|ID|Name> [warnings/mutes/kicks/bans]` for a complete list of information.", g.GuildPrefix),
						Thumbnail: &discordgo.MessageEmbedThumbnail{
							URL:    member.AvatarURL("2048"),
							Width:  2048,
							Height: 2048,
						},
						Fields: []*discordgo.MessageEmbedField{
							&discordgo.MessageEmbedField{
								Name:   fmt.Sprintf("❯ Total Warnings"),
								Value:  fmt.Sprintf("%d", len(usr.Warnings)),
								Inline: false,
							},
							&discordgo.MessageEmbedField{
								Name:   fmt.Sprintf("❯ Total Mutes"),
								Value:  fmt.Sprintf("%d", len(usr.Mutes)),
								Inline: false,
							},
							&discordgo.MessageEmbedField{
								Name:   fmt.Sprintf("❯ Total Kicks"),
								Value:  fmt.Sprintf("%d", len(usr.Kicks)),
								Inline: false,
							},
							&discordgo.MessageEmbedField{
								Name:   fmt.Sprintf("❯ Total Bans"),
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
