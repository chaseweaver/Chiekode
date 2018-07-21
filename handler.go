package main

import (
	"errors"
	"reflect"

	"github.com/bwmarrin/discordgo"
)

/**
 * handler.go
 * Chase Weaver
 *
 * This package handles commands/events be initializing them and their
 * properties.
 */

type (

	// Context of pass-in per command
	Context struct {
		Session *discordgo.Session
		Event   *discordgo.MessageCreate
		Guild   *discordgo.Guild
		Channel *discordgo.Channel
		Command Command
		Name    string
		Args    []string
	}

	// Command struct per command
	Command struct {
		Name            string
		Func            func(Context)
		Enabled         bool
		NSFWOnly        bool
		IgnoreSelf      bool
		IgnoreBots      bool
		Locked          bool
		RunIn           []string
		Aliases         []string
		BotPermissions  []string
		UserPermissions []string
		ArgsDelim       string
		ArgsUsage       string
		Description     string
	}
)

var commands = make(map[string]Command)

// CommandIsValid :
// Checks if command is valid to be ran.
func CommandIsValid(ctx Context) bool {

	// Enabled
	if !ctx.Command.Enabled {
		return false
	}

	// IgnoreSelf
	if ctx.Event.Author.ID == ctx.Session.State.User.ID && ctx.Command.IgnoreSelf {
		return false
	}

	// IgnoreBots
	if ctx.Event.Author.Bot && ctx.Command.IgnoreBots && ctx.Event.Author.ID != ctx.Session.State.User.ID {
		return false
	}

	// Locked
	if ctx.Command.Locked && ctx.Event.Author.ID != conf.OwnerID {
		return false
	}

	// RunIn Text
	if ctx.Channel.Type == discordgo.ChannelTypeGuildText && !Contains(ctx.Command.RunIn, "Text") {
		return false
	}

	// RunIn DM
	if ctx.Channel.Type == discordgo.ChannelTypeDM && !Contains(ctx.Command.RunIn, "DM") {
		return false
	}

	// UserPermissions
	for key := range ctx.Command.UserPermissions {
		if len(ctx.Command.UserPermissions) != 0 && !MemberHasPermission(ctx, ctx.Command.UserPermissions[key]) {
			return false
		}
	}

	return true
}

// RegisterNewCommand :
// Creates a new command
func RegisterNewCommand(c Command) {
	if !HasCommand(c.Name) {
		commands[c.Name] = c
	}
}

// HasCommand :
// Checks if a command is already mapped.
func HasCommand(k string) bool {
	_, ok := commands[k]
	if ok == true {
		return true
	}

	return false
}

// FetchCommand :
// Returns a valid command.
func FetchCommand(k string) Command {
	if HasCommand(commands[k].Name) {
		return commands[k]
	}

	for key := range commands {
		for val := range commands[key].Aliases {
			if commands[key].Aliases[val] == k {
				return commands[key]
			}
		}
	}

	return Command{}
}

// FetchCommandName :
// Returns a valid command if it exists.
func FetchCommandName(k string) string {
	if HasCommand(commands[k].Name) {
		return commands[k].Name
	}

	for key := range commands {
		for val := range commands[key].Aliases {
			if commands[key].Aliases[val] == k {
				return commands[key].Name
			}
		}
	}

	return k
}

// MemberHasPermission :
// Checks if the guild member has the required permission across all roles.
func MemberHasPermission(ctx Context, perm string) bool {

	// Check for Guild Owner
	if ctx.Guild.OwnerID == ctx.Event.Author.ID {
		return true
	}

	var permission int

	switch perm {
	case "Bot Owner":
		return true
	case "Read Messages":
		permission = discordgo.PermissionReadMessages
	case "Send Messages":
		permission = discordgo.PermissionSendMessages
	case "Send TTS Messages":
		permission = discordgo.PermissionSendTTSMessages
	case "Manage Messages":
		permission = discordgo.PermissionManageMessages
	case "Embed Links":
		permission = discordgo.PermissionEmbedLinks
	case "Attach Files":
		permission = discordgo.PermissionAttachFiles
	case "Read Message History":
		permission = discordgo.PermissionReadMessageHistory
	case "Mention Everyone":
		permission = discordgo.PermissionMentionEveryone
	case "Use External Emojis":
		permission = discordgo.PermissionUseExternalEmojis
	case "Voice Connect":
		permission = discordgo.PermissionVoiceConnect
	case "Voice Speak":
		permission = discordgo.PermissionVoiceSpeak
	case "Voice Mute Members":
		permission = discordgo.PermissionVoiceMuteMembers
	case "Voice Deafen Members":
		permission = discordgo.PermissionVoiceDeafenMembers
	case "Voice Move Members":
		permission = discordgo.PermissionVoiceMoveMembers
	case "Voice Use VAD":
		permission = discordgo.PermissionVoiceUseVAD
	case "Change Nickname":
		permission = discordgo.PermissionChangeNickname
	case "Manage Nicknames":
		permission = discordgo.PermissionManageNicknames
	case "Manage Roles":
		permission = discordgo.PermissionManageRoles
	case "Manage Webhooks":
		permission = discordgo.PermissionManageWebhooks
	case "Manage Emojis":
		permission = discordgo.PermissionManageEmojis
	case "Create Instant Invite":
		permission = discordgo.PermissionCreateInstantInvite
	case "Kick Members":
		permission = discordgo.PermissionKickMembers
	case "Ban Members":
		permission = discordgo.PermissionBanMembers
	case "Administrator":
		permission = discordgo.PermissionAdministrator
	case "Manage Channels":
		permission = discordgo.PermissionManageChannels
	case "Manage Server":
		permission = discordgo.PermissionManageServer
	case "Add Reactions":
		permission = discordgo.PermissionAddReactions
	case "View Audit Logs":
		permission = discordgo.PermissionViewAuditLogs
	case "All Text":
		permission = discordgo.PermissionAllText
	case "All Voice":
		permission = discordgo.PermissionAllVoice
	case "All Channel":
		permission = discordgo.PermissionAllChannel
	case "All":
		permission = discordgo.PermissionAll
	}

	mem, err := ctx.Session.State.Member(ctx.Guild.ID, ctx.Event.Author.ID)
	if err != nil {
		if mem, err = ctx.Session.GuildMember(ctx.Guild.ID, ctx.Event.Author.ID); err != nil {
			return false
		}
	}

	// Iterate through the role IDs stored in mem.Roles
	// to check permissions
	for _, roleID := range mem.Roles {
		role, err := ctx.Session.State.Role(ctx.Guild.ID, roleID)
		if err != nil {
			return false
		}
		if role.Permissions&permission != 0 {
			return true
		}
	}

	return false
}

// Call :
// Func based on name and passes Session, MessageCreate, ...args.
func Call(m map[string]interface{}, name string, params ...interface{}) (result []reflect.Value, err error) {
	f := reflect.ValueOf(m[name])
	if len(params) != f.Type().NumIn() {
		err = errors.New("the number of params is not adapted")
		return
	}
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	result = f.Call(in)
	return
}
