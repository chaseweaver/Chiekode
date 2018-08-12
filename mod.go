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
		Cooldown:        0,
		RunIn:           []string{"Text"},
		Aliases:         []string{},
		UserPermissions: []string{"Bot Owner", "Kick Members"},
		ArgsDelim:       " ",
		Usage:           []string{"<@Member(s)|ID(s)|Name(s)>"},
		Description:     "Warns a member.",
	})

	RegisterNewCommand(Command{
		Name:            "kick",
		Func:            Kick,
		Enabled:         true,
		NSFWOnly:        false,
		IgnoreSelf:      true,
		IgnoreBots:      true,
		Cooldown:        0,
		RunIn:           []string{"Text"},
		Aliases:         []string{},
		UserPermissions: []string{"Bot Owner", "Kick Members"},
		ArgsDelim:       " ",
		Usage:           []string{"<@Member(s)|ID(s)|Name#xxxx(s)>"},
		Description:     "Kicks a member from the guild.",
	})

	RegisterNewCommand(Command{
		Name:            "ban",
		Func:            Ban,
		Enabled:         true,
		NSFWOnly:        false,
		IgnoreSelf:      true,
		IgnoreBots:      true,
		Cooldown:        0,
		RunIn:           []string{"Text"},
		Aliases:         []string{},
		UserPermissions: []string{"Bot Owner", "Ban Members"},
		ArgsDelim:       " ",
		Usage:           []string{"<Member(s)|ID(s)|Name#xxxx(s)>"},
		Description:     "Bans a member from the guild.",
	})

	RegisterNewCommand(Command{
		Name:            "lock",
		Func:            Lock,
		Enabled:         true,
		NSFWOnly:        false,
		IgnoreSelf:      true,
		IgnoreBots:      true,
		Cooldown:        0,
		RunIn:           []string{"Text"},
		Aliases:         []string{"lockdown"},
		UserPermissions: []string{"Bot Owner", "Manage Channels"},
		ArgsDelim:       "",
		Usage:           []string{},
		Description:     "Locks a channel (prevents SEND_MESSAGES) for the default @everyone permission.",
	})

	RegisterNewCommand(Command{
		Name:            "unlock",
		Func:            Unlock,
		Enabled:         true,
		NSFWOnly:        false,
		IgnoreSelf:      true,
		IgnoreBots:      true,
		Cooldown:        0,
		RunIn:           []string{"Text"},
		Aliases:         []string{},
		UserPermissions: []string{"Bot Owner", "Manage Channels"},
		ArgsDelim:       "",
		Usage:           []string{},
		Description:     "Unlocks a channel (grants SEND_MESSAGES) for the default @everyone permission.",
	})

	RegisterNewCommand(Command{
		Name:            "mute",
		Func:            Mute,
		Enabled:         true,
		NSFWOnly:        false,
		IgnoreSelf:      true,
		IgnoreBots:      true,
		Cooldown:        0,
		RunIn:           []string{"Text"},
		Aliases:         []string{},
		UserPermissions: []string{"Bot Owner", "Kick Members"},
		ArgsDelim:       " ",
		Usage:           []string{"<@Member(s)|ID(s)|Name#xxxx(s)>", "[1h30m|2h|30m|etc]", "[reason]"},
		Description:     "Mutes a user with a set role with a reason (optional) and a time (optional).",
	})

	RegisterNewCommand(Command{
		Name:            "check",
		Func:            Check,
		Enabled:         true,
		NSFWOnly:        false,
		IgnoreSelf:      true,
		IgnoreBots:      true,
		Cooldown:        0,
		RunIn:           []string{"Text"},
		Aliases:         []string{},
		UserPermissions: []string{"Bot Owner", "Kick Members"},
		ArgsDelim:       " ",
		Usage:           []string{"<@Member(s)|ID(s)|Name#xxxx(s)>", "[warnings|mutes|kicks|bans|nicknames|usernames]"},
		Description:     "Checks the warnings, mutes, kicks, bans, nicknames, and usernames of a mentioned user.",
	})

	RegisterNewCommand(Command{
		Name:            "clear",
		Func:            Clear,
		Enabled:         true,
		NSFWOnly:        false,
		IgnoreSelf:      true,
		IgnoreBots:      true,
		Cooldown:        0,
		RunIn:           []string{"Text"},
		Aliases:         []string{"reset"},
		UserPermissions: []string{"Bot Owner", "Administrator", "Ban Members", "Kick Members"},
		ArgsDelim:       " ",
		Usage:           []string{"<@Member|ID|Name#xxxx> <warnings|mutes|kicks|bans|usernames|nicknames|all>"},
		Description:     "Clears a guild member's recorded data.",
	})
}

// Warn :
// Warn a user by ID / Name#xxxx / Mention, logs it to the redis database.
func Warn(ctx Context) {

	// Fetch users from message content, returns list of members and the remaining string with the member removed
	members, reason := FetchMessageContentUsersString(ctx, strings.Join(ctx.Args, ctx.Command.ArgsDelim))

	// Returns if a user cannot be found in the message, deletes command message, then deletes delayed response
	if len(members) == 0 {
		msg, err := ctx.Session.ChannelMessageSend(ctx.Channel.ID, "❌ | I cannot find that user!")

		if err != nil {
			return
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

	// Fetch Guild information from redis database
	g, guildErr := UnpackGuildStruct(ctx.Guild.ID)
	if guildErr != nil {
		log.Println(guildErr)
	}

	// Warns all members found within the message, logs warning to redis database
	for _, member := range members {

		// Prevent someone from warning the bot
		if member.ID == ctx.Session.State.User.ID {
			msg, err := ctx.Session.ChannelMessageSend(ctx.Channel.ID, "❌ | I will not warn myself!")

			if err != nil {
				return
			}

			DeleteMessageWithTime(ctx, msg.ID, 7500)
			return
		}

		// Target username
		target := member.Username + "#" + member.Discriminator

		// Author username
		author := ctx.Event.Message.Author.Username + "#" + ctx.Event.Message.Author.Discriminator

		// Sends warn message to channel the command was instantiated in
		ctx.Session.ChannelMessageSend(ctx.Channel.ID, fmt.Sprintf("`%s` has been warned by `%s`", target, author))

		// Logs warning to redis database
		LogWarning(ctx, member, reason)

		// Send logs to Guild Moderation Channel
		if guildErr != nil && g.ModerationLogsChannel != nil {
			ctx.Session.ChannelMessageSendEmbed(g.ModerationLogsChannel.ID,
				NewEmbed().
					SetTitle("Member Warned").
					SetColor(warningColor).
					SetAuthor(fmt.Sprintf("%s#%s / %s", member.Username, member.Discriminator, member.ID), member.AvatarURL("256"), member.AvatarURL("2048")).
					AddField("Author", fmt.Sprintf("%s#%s / %s", ctx.Event.Author.Username, ctx.Event.Author.Discriminator, ctx.Event.Author.ID)).
					AddField("Channel", fmt.Sprintf("<#%s>", ctx.Channel.ID)).
					AddField("Reason", reason).
					SetTimestamp(time.Now().Format(time.RFC3339)).MessageEmbed)
		}

		tr := fmt.Sprintf("with reason `%s`", reason)
		if reason == "N/A" {
			tr = "without a reason"
		}

		// Creates DM channel between bot and target
		channel, err := ctx.Session.UserChannelCreate(member.ID)

		// Sends a DM to the user with the warning information if the user can accept DMs
		if err == nil {
			ctx.Session.ChannelMessageSend(channel.ID, fmt.Sprintf("You have been warned by `%s` %s.", author, tr))
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
		msg, err := ctx.Session.ChannelMessageSend(ctx.Channel.ID, "❌ | I cannot find that user!")

		if err != nil {
			log.Println(err)
			return
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

	// Fetch Guild information from redis database
	g, guildErr := UnpackGuildStruct(ctx.Guild.ID)
	if guildErr != nil {
		log.Println(guildErr)
	}

	// Kicks all members found within the message, logs warning to redis database
	for _, member := range members {

		// Prevent someone from kicking the bot
		if member.ID == ctx.Session.State.User.ID {
			msg, err := ctx.Session.ChannelMessageSend(ctx.Channel.ID, "❌ | I will not kick myself!")

			if err != nil {
				return
			}

			DeleteMessageWithTime(ctx, msg.ID, 7500)
			return
		}

		// Target username
		target := member.Username + "#" + member.Discriminator

		// Author username
		author := ctx.Event.Message.Author.Username + "#" + ctx.Event.Message.Author.Discriminator

		// Sends kick message to channel the command was instantiated in
		ctx.Session.ChannelMessageSend(ctx.Channel.ID, fmt.Sprintf("`%s` has been kicked by `%s`", target, author))

		// Logs kick to redis database
		LogKick(ctx, member, reason)

		// Send logs to Guild Moderation Channel
		if guildErr != nil && g.ModerationLogsChannel != nil {
			ctx.Session.ChannelMessageSendEmbed(g.ModerationLogsChannel.ID,
				NewEmbed().
					SetTitle("Member Kicked").
					SetColor(kickColor).
					SetAuthor(fmt.Sprintf("%s#%s / %s", member.Username, member.Discriminator, member.ID), member.AvatarURL("256"), member.AvatarURL("2048")).
					AddField("Author", fmt.Sprintf("%s#%s / %s", ctx.Event.Author.Username, ctx.Event.Author.Discriminator, ctx.Event.Author.ID)).
					AddField("Channel", fmt.Sprintf("<#%s>", ctx.Channel.ID)).
					AddField("Reason", reason).
					SetTimestamp(time.Now().Format(time.RFC3339)).MessageEmbed)
		}

		tr := fmt.Sprintf("with reason `%s`", reason)
		if reason == "N/A" {
			tr = "without a reason"
		}

		// Creates DM channel between bot and target
		channel, err := ctx.Session.UserChannelCreate(member.ID)

		// Sends a DM to the user with the kick information if the user can accept DMs
		if err == nil {
			ctx.Session.ChannelMessageSend(channel.ID, fmt.Sprintf("You have been kicked by `%s` %s.", author, tr))
		}

		// Kicks the guild member with given reason
		err = ctx.Session.GuildMemberDeleteWithReason(ctx.Guild.ID, member.ID, reason)

		if err != nil {
			msg, err := ctx.Session.ChannelMessageSend(ctx.Channel.ID, "❌ | I cannot kick this user!")

			if err != nil {
				return
			}

			DeleteMessageWithTime(ctx, msg.ID, 7500)
			break
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
		msg, err := ctx.Session.ChannelMessageSend(ctx.Channel.ID, "❌ | I cannot find that user!")

		if err != nil {
			return
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

	// Fetch Guild information from redis database
	g, guildErr := UnpackGuildStruct(ctx.Guild.ID)
	if guildErr != nil {
		log.Println(err)
	}

	// Bans all members found within the message, logs warning to redis database
	for _, member := range members {

		// Prevent someone from banning the bot
		if member.ID == ctx.Session.State.User.ID {
			msg, err := ctx.Session.ChannelMessageSend(ctx.Channel.ID, "❌ | I will not ban myself!")

			if err != nil {
				return
			}

			DeleteMessageWithTime(ctx, msg.ID, 7500)
			return
		}

		// Target username
		target := member.Username + "#" + member.Discriminator

		// Author username
		author := ctx.Event.Message.Author.Username + "#" + ctx.Event.Message.Author.Discriminator

		// Sends ban message to channel the command was instantiated in
		ctx.Session.ChannelMessageSend(ctx.Channel.ID, fmt.Sprintf("`%s` has been banned by `%s`", target, author))

		// Logs ban to redis database
		LogBan(ctx, member, reason)

		// Send logs to Guild Moderation Channel
		if guildErr != nil && g.ModerationLogsChannel != nil {
			ctx.Session.ChannelMessageSendEmbed(g.ModerationLogsChannel.ID,
				NewEmbed().
					SetTitle("Member Banned").
					SetColor(banColor).
					SetAuthor(fmt.Sprintf("%s#%s / %s", member.Username, member.Discriminator, member.ID), member.AvatarURL("256"), member.AvatarURL("2048")).
					AddField("Author", fmt.Sprintf("%s#%s / %s", ctx.Event.Author.Username, ctx.Event.Author.Discriminator, ctx.Event.Author.ID)).
					AddField("Channel", fmt.Sprintf("<#%s>", ctx.Channel.ID)).
					AddField("Reason", reason).
					SetTimestamp(time.Now().Format(time.RFC3339)).MessageEmbed)
		}

		tr := fmt.Sprintf("with reason `%s`", reason)
		if reason == "N/A" {
			tr = "without a reason"
		}

		// Creates DM channel between bot and target
		channel, err := ctx.Session.UserChannelCreate(member.ID)

		// Sends a DM to the user with the ban information if the user can accept DMs
		if err == nil {
			ctx.Session.ChannelMessageSend(channel.ID, fmt.Sprintf("You have been banned by `%s` %s.", author, tr))
		}

		// Bans the guild member with given reason, deletes 0 messages
		err = ctx.Session.GuildBanCreateWithReason(ctx.Guild.ID, member.ID, reason, 0)

		if err != nil {
			msg, err := ctx.Session.ChannelMessageSend(ctx.Channel.ID, "❌ | I cannot ban this user!")

			if err != nil {
				return
			}

			DeleteMessageWithTime(ctx, msg.ID, 7500)
			break
		}
	}
}

// Lock :
// Overrides default @everyone permission and prevents "SEND_MESSAGES" permission.
func Lock(ctx Context) {

	DeleteMessageWithTime(ctx, ctx.Event.Message.ID, 0)
	guild, err := ctx.Session.Guild(ctx.Guild.ID)

	if err != nil {
		log.Println(err)
		return
	}

	// Find @everyone role
	var everyone *discordgo.Role
	for _, v := range guild.Roles {
		if v.Name == "@everyone" {
			everyone = v
		}
	}

	// Return if the bot is unable to find the @everyone role
	if everyone == nil {
		ctx.Session.ChannelMessageSend(ctx.Channel.ID, fmt.Sprintf("❌ | I cannot find the `@everyone` role!"))
		return
	}

	// Find and append new permission overwrites for the channel for the @everyone role
	var allow, deny int
	for _, v := range ctx.Channel.PermissionOverwrites {
		if v.ID == everyone.ID {
			allow = v.Allow
			deny = v.Deny | discordgo.PermissionSendMessages
		}
	}

	// Apply new permissions
	err = ctx.Session.ChannelPermissionSet(ctx.Channel.ID, everyone.ID, "0", allow, deny)

	if err != nil {
		log.Println(err)
		ctx.Session.ChannelMessageSend(ctx.Channel.ID, fmt.Sprintf("❌ | An error has occured!"))
		return
	}

	ctx.Session.ChannelMessageSend(ctx.Channel.ID, fmt.Sprintf("🔒 | This channel is now under lockdown!"))
}

// Unlock :
// Overrides default @everyone permission and allows "SEND_MESSAGES" permission.
func Unlock(ctx Context) {

	DeleteMessageWithTime(ctx, ctx.Event.Message.ID, 0)
	guild, err := ctx.Session.Guild(ctx.Guild.ID)

	if err != nil {
		log.Println(err)
		return
	}

	// Find @everyone role
	var everyone *discordgo.Role
	for _, v := range guild.Roles {
		if v.Name == "@everyone" {
			everyone = v
		}
	}

	// Return if the bot is unable to find the @everyone role
	if everyone == nil {
		ctx.Session.ChannelMessageSend(ctx.Channel.ID, fmt.Sprintf("❌ | I cannot find the `@everyone` role!"))
		return
	}

	// Find and append new permission overwrites for the channel for the @everyone role
	var allow, deny int
	for _, v := range ctx.Channel.PermissionOverwrites {
		if v.ID == everyone.ID {
			allow = v.Allow | discordgo.PermissionSendMessages
			deny = v.Deny
		}
	}

	// Apply new permissions
	err = ctx.Session.ChannelPermissionSet(ctx.Channel.ID, everyone.ID, "0", allow, deny)

	if err != nil {
		log.Println(err)
		ctx.Session.ChannelMessageSend(ctx.Channel.ID, fmt.Sprintf("❌ | An error has occured!"))
		return
	}

	ctx.Session.ChannelMessageSend(ctx.Channel.ID, fmt.Sprintf("🔓 | This channel is no longer under lockdown!"))
}

// Mute :
// Adds the guild's "mute" role to a member for a set time (if given)
func Mute(ctx Context) {

	// Delete command message
	DeleteMessageWithTime(ctx, ctx.Event.Message.ID, 0)

	// Fetch guild information
	g, err := UnpackGuildStruct(ctx.Guild.ID)
	if err != nil {
		log.Println(err)
		return
	}

	// Check if the guild role is set
	if g.MutedRole == nil {
		ctx.Session.ChannelMessageSend(ctx.Channel.ID, fmt.Sprintf("❌ | You do not have a muted role set up! Please configure one using `%sset muted role%s<@Role|Name|ID>`", g.GuildPrefix, commands["set"].ArgsDelim))
		return
	}

	// Check to see if the set guild role exists / still exists within the context of the guild (in case of deletion, etc.)
	role, err := ctx.Session.State.Role(ctx.Guild.ID, g.MutedRole.ID)

	if err != nil {
		log.Println(err)
		ctx.Session.ChannelMessageSend(ctx.Channel.ID, "❌ | There was an error with the set Muted Role. Perhaps it has been changed since the last time it was configured.")
		return
	}

	// Fetch users from message content, returns list of members and the remaining string with the member removed
	members, reason := FetchMessageContentUsersString(ctx, strings.Join(ctx.Args, ctx.Command.ArgsDelim))
	length, _ := time.ParseDuration(reason)

	// Retuns if a user cannot be found in the message, deletes delayed response
	if len(members) == 0 {
		msg, err := ctx.Session.ChannelMessageSend(ctx.Channel.ID, "❌ | I cannot find that user!")

		if err != nil {
			return
		}

		DeleteMessageWithTime(ctx, msg.ID, 7500)
		return
	}

	// If no reason is specified, set one for the database logger
	if len(reason) == 0 {
		reason = "N/A"
	}

	// Bans all members found within the message, logs warning to redis database
	for _, member := range members {

		// Prevent someone from muting the bot
		if member.ID == ctx.Session.State.User.ID {
			msg, err := ctx.Session.ChannelMessageSend(ctx.Channel.ID, "❌ | I will not mute myself!")

			if err != nil {
				return
			}

			DeleteMessageWithTime(ctx, msg.ID, 7500)
			return
		}

		// Target username
		target := member.Username + "#" + member.Discriminator

		// Author username
		author := ctx.Event.Message.Author.Username + "#" + ctx.Event.Message.Author.Discriminator

		// Creates DM channel between bot and target
		channel, err := ctx.Session.UserChannelCreate(member.ID)

		if err != nil {
			log.Println(err)
		}

		// Sends mute message to channel the command was instantiated in
		if err == nil {
			ctx.Session.ChannelMessageSend(ctx.Channel.ID, fmt.Sprintf("`%s` has been muted by `%s`", target, author))
		}

		err = ctx.Session.GuildMemberRoleAdd(ctx.Guild.ID, member.ID, role.ID)

		if err != nil {
			log.Println(err)
			ctx.Session.ChannelMessageSend(ctx.Channel.ID, "❌ | An error has occured!")
			break
		}

		// Logs mutes to redis database
		LogMute(ctx, member, reason, length)

		if g.ModerationLogsChannel != nil {
			ctx.Session.ChannelMessageSendEmbed(g.ModerationLogsChannel.ID,
				NewEmbed().
					SetTitle("Member Mute").
					SetColor(muteColor).
					SetAuthor(fmt.Sprintf("%s#%s / %s", member.Username, member.Discriminator, member.ID), member.AvatarURL("256"), member.AvatarURL("2048")).
					AddField("Author", fmt.Sprintf("%s#%s / %s", ctx.Event.Author.Username, ctx.Event.Author.Discriminator, ctx.Event.Author.ID)).
					AddField("Channel", fmt.Sprintf("<#%s>", ctx.Channel.ID)).
					AddField("Duration", fmt.Sprintf("%v", length)).
					AddField("Reason", reason).
					SetTimestamp(time.Now().Format(time.RFC3339)).MessageEmbed)
		}

		tr := fmt.Sprintf("with reason `%s`", reason)
		if reason == "N/A" {
			tr = "without a reason"
		}

		tl := fmt.Sprintf("for `%v`", length)
		if length == 0 {
			tl = "indefinitely"
		}

		// Sends a DM to the user with the mute information if the user can accept DMs
		ctx.Session.ChannelMessageSend(channel.ID, fmt.Sprintf("You have been muted by `%s` %s %s.", author, tr, tl))
	}
}

// Unmute :
// Removes the guild's "mute" role from a member
func Unmute(ctx Context) {

	// Delete command message
	DeleteMessageWithTime(ctx, ctx.Event.Message.ID, 0)

	// Fetch guild information
	g, err := UnpackGuildStruct(ctx.Guild.ID)
	if err != nil {
		log.Println(err)
		return
	}

	// Check if the guild role is set
	if g.MutedRole == nil {
		ctx.Session.ChannelMessageSend(ctx.Channel.ID, fmt.Sprintf("❌ | You do not have a muted role set up! Please configure one using `%sset muted role%s<@Role|Name|ID>`", g.GuildPrefix, commands["set"].ArgsDelim))
	}

	// Check to see if the set guild role exists / still exists within the context of the guild (in case of deletion, etc.)
	role, err := ctx.Session.State.Role(ctx.Guild.ID, g.MutedRole.ID)

	if err != nil {
		log.Println(err)
		ctx.Session.ChannelMessageSend(ctx.Channel.ID, "❌ | There was an error with the set Muted Role. Perhaps it has been changed since the last time it was configured.")
		return
	}

	// Fetch users from message content, returns list of members and the remaining string with the member removed
	members, reason := FetchMessageContentUsersString(ctx, strings.Join(ctx.Args, ctx.Command.ArgsDelim))

	// Retuns if a user cannot be found in the message, deletes delayed response
	if len(members) == 0 {
		msg, err := ctx.Session.ChannelMessageSend(ctx.Channel.ID, "❌ | I cannot find that user!")

		if err != nil {
			return
		}

		DeleteMessageWithTime(ctx, msg.ID, 7500)
		return
	}

	// If no reason is specified, set one for the database logger
	if len(reason) == 0 {
		reason = "N/A"
	}

	// Bans all members found within the message, logs warning to redis database
	for _, member := range members {

		// Target username
		target := member.Username + "#" + member.Discriminator

		// Author username
		author := ctx.Event.Message.Author.Username + "#" + ctx.Event.Message.Author.Discriminator

		// Creates DM channel between bot and target
		channel, err := ctx.Session.UserChannelCreate(member.ID)

		if err != nil {
			log.Println(err)
		}

		// Sends mute message to channel the command was instantiated in
		if err == nil {
			ctx.Session.ChannelMessageSend(ctx.Channel.ID, fmt.Sprintf("`%s` has been unmuted by `%s`", target, author))
		}

		err = ctx.Session.GuildMemberRoleAdd(ctx.Guild.ID, member.ID, role.ID)

		if err != nil {
			log.Println(err)
			ctx.Session.ChannelMessageSend(ctx.Channel.ID, "❌ | An error has occured!")
			break
		}

		if g.ModerationLogsChannel != nil {
			ctx.Session.ChannelMessageSendEmbed(g.ModerationLogsChannel.ID,
				NewEmbed().
					SetTitle("Member Unmute").
					SetColor(unmuteColor).
					SetAuthor(fmt.Sprintf("%s#%s / %s", member.Username, member.Discriminator, member.ID), member.AvatarURL("256"), member.AvatarURL("2048")).
					AddField("Author", fmt.Sprintf("%s#%s / %s", ctx.Event.Author.Username, ctx.Event.Author.Discriminator, ctx.Event.Author.ID)).
					AddField("Channel", fmt.Sprintf("<#%s>", ctx.Channel.ID)).
					AddField("Reason", reason).
					SetTimestamp(time.Now().Format(time.RFC3339)).MessageEmbed)
		}

		tr := fmt.Sprintf("with reason `%s`", reason)
		if reason == "N/A" {
			tr = "without a reason"
		}

		// Sends a DM to the user with the unmute information if the user can accept DMs
		ctx.Session.ChannelMessageSend(channel.ID, fmt.Sprintf("You have been unmuted by `%s` %s.", author, tr))
	}

}

// Check :
// Check the user's warnings, mutes, kicks, bans, nicknames, and usernames from the redis database.
func Check(ctx Context) {

	// Fetch users from message content
	members, checkType := FetchMessageContentUsersString(ctx, strings.Join(ctx.Args, ctx.Command.ArgsDelim))

	// Type of check to look for
	if len(checkType) > 1 {
		checkType = strings.ToUpper(checkType)[1:]
	}

	// Returns if a user cannot be found in the message, deletes delayed response
	if len(members) == 0 {
		msg, err := ctx.Session.ChannelMessageSend(ctx.Channel.ID, "❌ | I cannot find that user!")

		if err != nil {
			return
		}

		DeleteMessageWithTime(ctx, msg.ID, 7500)
		return
	}

	// Fetch guild information
	g, err := UnpackGuildStruct(ctx.Guild.ID)
	if err != nil {
		log.Println(err)
		return
	}

	switch checkType {
	case "WARNINGS":
		for _, member := range members {
			if _, ok := g.GuildUser[member.ID]; !ok {
				msg, err := ctx.Session.ChannelMessageSend(ctx.Channel.ID, "❌ | I do not have any logs for that user!")

				if err != nil {
					return
				}

				DeleteMessageWithTime(ctx, msg.ID, 7500)
				return
			}

			if len(g.GuildUser[member.ID].Warnings) == 0 {
				msg, err := ctx.Session.ChannelMessageSend(ctx.Channel.ID, "❌ | No warnings found!")
				if err != nil {
					log.Println(err)
					return
				}
				DeleteMessageWithTime(ctx, msg.ID, 5000)
				return
			}

			str := FormatWarnings(g.GuildUser[member.ID].Warnings)
			ctx.Session.ChannelMessageSendEmbed(ctx.Channel.ID,
				NewEmbed().
					SetTitle(fmt.Sprintf("Warning Stats [%d]", len(g.GuildUser[member.ID].Warnings))).
					SetColor(warningColor).
					SetAuthor(fmt.Sprintf("%s#%s / %s", member.Username, member.Discriminator, member.ID), g.GuildUser[member.ID].User.AvatarURL("256"), g.GuildUser[member.ID].User.AvatarURL("2048")).
					SetDescription(str).
					SetTimestamp(time.Now().Format(time.RFC3339)).MessageEmbed)
		}

	case "KICKS":
		for _, member := range members {
			if _, ok := g.GuildUser[member.ID]; !ok {
				msg, err := ctx.Session.ChannelMessageSend(ctx.Channel.ID, "❌ | I do not have any logs for that user!")

				if err != nil {
					log.Println(err)
					return
				}

				DeleteMessageWithTime(ctx, msg.ID, 7500)
				return
			}

			if len(g.GuildUser[member.ID].Kicks) == 0 {
				msg, err := ctx.Session.ChannelMessageSend(ctx.Channel.ID, "❌ | No kicks found!")
				if err != nil {
					log.Println(err)
					return
				}
				DeleteMessageWithTime(ctx, msg.ID, 5000)
				return
			}

			str := FormatKicks(g.GuildUser[member.ID].Kicks)
			ctx.Session.ChannelMessageSendEmbed(ctx.Channel.ID,
				NewEmbed().
					SetTitle(fmt.Sprintf("Kick Stats [%d]", len(g.GuildUser[member.ID].Kicks))).
					SetColor(warningColor).
					SetAuthor(fmt.Sprintf("%s#%s / %s", member.Username, member.Discriminator, member.ID), g.GuildUser[member.ID].User.AvatarURL("256"), g.GuildUser[member.ID].User.AvatarURL("2048")).
					SetDescription(str).
					SetTimestamp(time.Now().Format(time.RFC3339)).MessageEmbed)
		}

	case "BANS":
		for _, member := range members {
			if _, ok := g.GuildUser[member.ID]; !ok {
				msg, err := ctx.Session.ChannelMessageSend(ctx.Channel.ID, "❌ | I do not have any logs for that user!")

				if err != nil {
					log.Println(err)
					return
				}

				DeleteMessageWithTime(ctx, msg.ID, 7500)
				return
			}

			if len(g.GuildUser[member.ID].Bans) == 0 {
				msg, err := ctx.Session.ChannelMessageSend(ctx.Channel.ID, "❌ | No bans found!")
				if err != nil {
					log.Println(err)
					return
				}
				DeleteMessageWithTime(ctx, msg.ID, 5000)
				return
			}

			str := FormatBans(g.GuildUser[member.ID].Bans)
			ctx.Session.ChannelMessageSendEmbed(ctx.Channel.ID,
				NewEmbed().
					SetTitle(fmt.Sprintf("Ban Stats [%d]", len(g.GuildUser[member.ID].Bans))).
					SetColor(warningColor).
					SetAuthor(fmt.Sprintf("%s#%s / %s", member.Username, member.Discriminator, member.ID), g.GuildUser[member.ID].User.AvatarURL("256"), g.GuildUser[member.ID].User.AvatarURL("2048")).
					SetDescription(str).
					SetTimestamp(time.Now().Format(time.RFC3339)).MessageEmbed)
		}
	case "NICKNAMES":
		for _, member := range members {
			if _, ok := g.GuildUser[member.ID]; !ok {
				msg, err := ctx.Session.ChannelMessageSend(ctx.Channel.ID, "❌ | I do not have any logs for that user!")

				if err != nil {
					log.Println(err)
					return
				}

				DeleteMessageWithTime(ctx, msg.ID, 7500)
				return
			}

			if len(g.GuildUser[member.ID].Nicknames) == 0 {
				msg, err := ctx.Session.ChannelMessageSend(ctx.Channel.ID, "❌ | No nicknames found!")
				if err != nil {
					log.Println(err)
					return
				}
				DeleteMessageWithTime(ctx, msg.ID, 5000)
				return
			}

			str := FormatNicknames(g.GuildUser[member.ID].Nicknames)
			ctx.Session.ChannelMessageSendEmbed(ctx.Channel.ID,
				NewEmbed().
					SetTitle(fmt.Sprintf("Nickname Stats [%d]", len(g.GuildUser[member.ID].Nicknames))).
					SetColor(warningColor).
					SetAuthor(fmt.Sprintf("%s#%s / %s", member.Username, member.Discriminator, member.ID), g.GuildUser[member.ID].User.AvatarURL("256"), g.GuildUser[member.ID].User.AvatarURL("2048")).
					SetDescription(str).
					SetTimestamp(time.Now().Format(time.RFC3339)).MessageEmbed)
		}
	case "USERNAMES":
		for _, member := range members {
			if _, ok := g.GuildUser[member.ID]; !ok {
				msg, err := ctx.Session.ChannelMessageSend(ctx.Channel.ID, "❌ | I do not have any logs for that user!")

				if err != nil {
					log.Println(err)
					return
				}

				DeleteMessageWithTime(ctx, msg.ID, 7500)
				return
			}

			if len(g.GuildUser[member.ID].Usernames) == 0 {
				msg, err := ctx.Session.ChannelMessageSend(ctx.Channel.ID, "❌ | No usernames found!")
				if err != nil {
					log.Println(err)
					return
				}
				DeleteMessageWithTime(ctx, msg.ID, 5000)
				return
			}

			str := FormatUsernames(g.GuildUser[member.ID].Usernames)
			ctx.Session.ChannelMessageSendEmbed(ctx.Channel.ID,
				NewEmbed().
					SetTitle(fmt.Sprintf("Username Stats [%d]", len(g.GuildUser[member.ID].Usernames))).
					SetColor(warningColor).
					SetAuthor(fmt.Sprintf("%s#%s / %s", member.Username, member.Discriminator, member.ID), g.GuildUser[member.ID].User.AvatarURL("256"), g.GuildUser[member.ID].User.AvatarURL("2048")).
					SetDescription(str).
					SetTimestamp(time.Now().Format(time.RFC3339)).MessageEmbed)
		}
	default:

		for _, member := range members {
			if _, ok := g.GuildUser[member.ID]; !ok {
				msg, err := ctx.Session.ChannelMessageSend(ctx.Channel.ID, "❌ | I do not have any logs for that user!")

				if err != nil {
					log.Println(err)
					return
				}

				DeleteMessageWithTime(ctx, msg.ID, 7500)
				return
			}

			_, err := ctx.Session.ChannelMessageSendEmbed(ctx.Channel.ID,
				NewEmbed().
					SetTitle(fmt.Sprintf("Run `%scheck <@member|ID|Name#xxxx> [warnings|mutes|kicks|bans|usernames|nicknames]` for a complete list of information.", g.GuildPrefix)).
					SetColor(RandomInt(0, 16777215)).
					SetAuthor(fmt.Sprintf("%s#%s / %s", member.Username, member.Discriminator, member.ID),
						g.GuildUser[member.ID].User.AvatarURL("256"), g.GuildUser[member.ID].User.AvatarURL("2048")).
					SetThumbnail(member.AvatarURL("2048")).
					AddField("❯ Total Warnings", fmt.Sprintf("%d", len(g.GuildUser[member.ID].Warnings))).
					AddField("❯ Total Mutes", fmt.Sprintf("%d", len(g.GuildUser[member.ID].Mutes))).
					AddField("❯ Total Kicks", fmt.Sprintf("%d", len(g.GuildUser[member.ID].Kicks))).
					AddField("❯ Total Bans", fmt.Sprintf("%d", len(g.GuildUser[member.ID].Bans))).
					AddField("❯ Total Nicknames", fmt.Sprintf("%d", len(g.GuildUser[member.ID].Nicknames))).
					AddField("❯ Total Usernames", fmt.Sprintf("%d", len(g.GuildUser[member.ID].Usernames))).
					SetTimestamp(time.Now().Format(time.RFC3339)).MessageEmbed)

			if err != nil {
				log.Println(err)
				return
			}

		}
	}
}

// Clear :
// Clears a GuildUser's recorded information
func Clear(ctx Context) {

	// Fetch users from message content
	members, checkType := FetchMessageContentUsersString(ctx, strings.Join(ctx.Args, ctx.Command.ArgsDelim))

	// Type of check to look for
	if len(checkType) > 1 {
		checkType = strings.ToUpper(checkType)[1:]
	}

	// Returns if a user cannot be found in the message, deletes delayed response
	if len(members) == 0 {
		msg, err := ctx.Session.ChannelMessageSend(ctx.Channel.ID, "❌ | I cannot find that user!")

		if err != nil {
			log.Println(err)
			return
		}

		DeleteMessageWithTime(ctx, msg.ID, 7500)
		return
	}

	member := members[0]

	// Fetch Guild information from redis database
	data, err := redis.Bytes(p.Do("GET", ctx.Guild.ID))
	if err != nil {
		log.Println(err)
		return
	}

	var g Guild
	err = json.Unmarshal(data, &g)

	if err != nil {
		log.Println(err)
	}

	// Check for User ID in Guild map, register user if missing
	if _, ok := g.GuildUser[member.ID]; !ok {
		ctx.Session.ChannelMessageSend(ctx.Channel.ID, "❌ | I do not have any information on that user.")
		return
	}

	user := g.GuildUser[member.ID]

	switch checkType {
	case "WARNINGS":
		user.Warnings = make(map[int64]Warnings)
	case "MUTES":
		user.Mutes = make(map[int64]Mutes)
	case "KICKS":
		user.Kicks = make(map[int64]Kicks)
	case "BANS":
		user.Bans = make(map[int64]Bans)
	case "NICKNAMES":
		user.Nicknames = make(map[int64]Nicknames)
	case "USERNAME":
		fallthrough
	case "USERNAMES":
		user.Usernames = make(map[int64]Usernames)
		username := Usernames{
			Username:      user.User.Username,
			Discriminator: user.User.Discriminator,
			Time:          time.Now(),
		}
		user.Usernames[MakeTimestamp()] = username
	case "ALL":
		user.Warnings = make(map[int64]Warnings)
		user.Mutes = make(map[int64]Mutes)
		user.Kicks = make(map[int64]Kicks)
		user.Bans = make(map[int64]Bans)
		user.Nicknames = make(map[int64]Nicknames)
		user.Usernames = make(map[int64]Usernames)
		user.Usernames = make(map[int64]Usernames)
		username := Usernames{
			Username:      user.User.Username,
			Discriminator: user.User.Discriminator,
			Time:          time.Now(),
		}
		user.Usernames[MakeTimestamp()] = username
	default:
		ctx.Session.ChannelMessageSend(ctx.Channel.ID, "❌ | Please choose a type to clear `<warnings|mutes|kicks|bans|usernames|nicknames|all>`")
		return
	}

	// Set newly modified user back into GuildUser struct
	g.GuildUser[member.ID] = user

	serialized, err := json.Marshal(g)

	if err != nil {
		log.Println(err)
		return
	}

	_, err = p.Do("SET", ctx.Guild.ID, serialized)
	if err != nil {
		log.Println(err)
	}

	ctx.Session.ChannelMessageSend(ctx.Channel.ID, fmt.Sprintf("✅ | %s cleared successfully!", strings.ToTitle(checkType)))
}
