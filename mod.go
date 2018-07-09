package main

import (
	"strings"
)

/**
 * mod.go
 * Chase Weaver
 *
 * This package bundles commands for the moderation commands of the bot.
 */

func init() {
	RegisterNewCommand(Command{
		Name:            "kick",
		Func:            Kick,
		Enabled:         true,
		NSFWOnly:        false,
		IgnoreSelf:      true,
		IgnoreBots:      true,
		RunIn:           []string{"Text"},
		Aliases:         []string{},
		UserPermissions: []string{"Administrator", "KickMembers"},
		ArgsDelim:       " ",
		ArgsUsage:       "<golang expression>",
		Description:     "Kicks a member via mention or ID.",
	})
}

// Kick command will kick a mentioned user and log it to the redis database
func Kick(ctx Context) {
	for key := range ctx.Event.Message.Mentions {
		mem := ctx.Event.Message.Mentions[key].ID
		ctx.Session.GuildMemberDeleteWithReason(ctx.Guild.ID, mem, strings.Join(ctx.Args, " "))
	}
}
