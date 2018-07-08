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

// CheckValidPrereq checks if command is valid to be ran
func CheckValidPrereq(ctx Context) bool {

	// Enabled
	if !ctx.Command.Enabled {
		return false
	}

	//IgnoreSelf
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

// MemberHasPermission checks if the guild member has the required permission across all roles
func MemberHasPermission(s *discordgo.Session, guildID string, userID string, permission int) (bool, error) {
	mem, err := s.State.Member(guildID, userID)
	if err != nil {
		if mem, err = s.GuildMember(guildID, userID); err != nil {
			return false, err
		}
	}

	// Iterate through the role IDs stored in mem.Roles
	// to check permissions
	for _, roleID := range mem.Roles {
		role, err := s.State.Role(guildID, roleID)
		if err != nil {
			return false, err
		}
		if role.Permissions&permission != 0 {
			return true, nil
		}
	}

	return false, nil
}

// IsNSFWOnly returns true if the command requests an NSFW-only channel and the channel is NSFW only
func IsNSFWOnly(ctx Context, c Command) bool {
	if ctx.Channel.NSFW && c.NSFWOnly {
		return true
	}
	return false
}
