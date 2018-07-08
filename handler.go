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

// CommandIsValid checks if command is valid to be ran
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

	// RunIn Text, DM
	if ctx.Channel.Type == discordgo.ChannelTypeGuildText && !Contains(ctx.Command.RunIn, "Text") {
		return false
	} else if ctx.Channel.Type == discordgo.ChannelTypeDM && !Contains(ctx.Command.RunIn, "DM") {
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

// RegisterNewCommand creates a new command
func RegisterNewCommand(c Command) {
	if !HasCommand(c.Name) {
		commands[c.Name] = c
	}
}

// HasCommand checks if a command is already mapped
func HasCommand(k string) bool {
	_, ok := commands[k]
	if ok == true {
		return true
	}

	return false
}

// FetchCommand returns a valid command
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

// FetchCommandName returns a valid command if it exists
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

// MemberHasPermission checks if the guild member has the required permission across all roles
func MemberHasPermission(ctx Context, perm string) bool {
	var permission int

	switch perm {
	case "BotOwner":
		return true
	case "ReadMessages":
		permission = discordgo.PermissionReadMessages
	case "SendMessages":
		permission = discordgo.PermissionSendMessages
	case "SendTTSMessages":
		permission = discordgo.PermissionSendTTSMessages
	case "ManageMessages":
		permission = discordgo.PermissionManageMessages
	case "EmbedLinks":
		permission = discordgo.PermissionEmbedLinks
	case "AttachFiles":
		permission = discordgo.PermissionAttachFiles
	case "ReadMessageHistory":
		permission = discordgo.PermissionReadMessageHistory
	case "MentionEveryone":
		permission = discordgo.PermissionMentionEveryone
	case "UseExternalEmojis":
		permission = discordgo.PermissionUseExternalEmojis
	case "VoiceConnect":
		permission = discordgo.PermissionVoiceConnect
	case "VoiceSpeak":
		permission = discordgo.PermissionVoiceSpeak
	case "VoiceMuteMembers":
		permission = discordgo.PermissionVoiceMuteMembers
	case "VoiceDeafenMembers":
		permission = discordgo.PermissionVoiceDeafenMembers
	case "VoiceMoveMembers":
		permission = discordgo.PermissionVoiceMoveMembers
	case "VoiceUseVAD":
		permission = discordgo.PermissionVoiceUseVAD
	case "ChangeNickname":
		permission = discordgo.PermissionChangeNickname
	case "ManageNicknames":
		permission = discordgo.PermissionManageNicknames
	case "ManageRoles":
		permission = discordgo.PermissionManageRoles
	case "ManageWebhooks":
		permission = discordgo.PermissionManageWebhooks
	case "ManageEmojis":
		permission = discordgo.PermissionManageEmojis
	case "CreateInstantInvite":
		permission = discordgo.PermissionCreateInstantInvite
	case "KickMembers":
		permission = discordgo.PermissionKickMembers
	case "BanMembers":
		permission = discordgo.PermissionBanMembers
	case "Administrator":
		permission = discordgo.PermissionAdministrator
	case "ManageChannels":
		permission = discordgo.PermissionManageChannels
	case "ManageServer":
		permission = discordgo.PermissionManageServer
	case "AddReactions":
		permission = discordgo.PermissionAddReactions
	case "ViewAuditLogs":
		permission = discordgo.PermissionViewAuditLogs
	case "AllText":
		permission = discordgo.PermissionAllText
	case "AllVoice":
		permission = discordgo.PermissionAllVoice
	case "AllChannel":
		permission = discordgo.PermissionAllChannel
	case "All":
		permission = discordgo.PermissionAll
	default:
		return false
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

// Call func based on name and passes Session, MessageCreate, ...args
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
