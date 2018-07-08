package main

/**
 * init.go
 * Chase Weaver
 *
 * This package is meant to handle initialization of commands
 * and command properties for quick configuration.
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
		ArgsUsage:       "[@member]",
		Description:     "Fetches the avatar/pfp for the requested member.",
	})

	RegisterNewCommand(Command{
		Name:            "ping",
		Func:            Ping,
		Enabled:         true,
		NSFWOnly:        false,
		IgnoreSelf:      true,
		IgnoreBots:      true,
		RunIn:           []string{"Text", "DM"},
		Aliases:         []string{},
		UserPermissions: []string{},
		ArgsDelim:       "",
		ArgsUsage:       "",
		Description:     "Pong!",
	})

	RegisterNewCommand(Command{
		Name:            "eval",
		Func:            Eval,
		Enabled:         true,
		NSFWOnly:        false,
		IgnoreSelf:      true,
		IgnoreBots:      true,
		RunIn:           []string{"Text", "DM"},
		Aliases:         []string{"e"},
		UserPermissions: []string{"BotOwner"},
		ArgsDelim:       " ",
		ArgsUsage:       "<golang expression>",
		Description:     "Evaluation command for bot-owner only.",
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
		UserPermissions: []string{"Administrator", "KickMembers"},
		ArgsDelim:       " ",
		ArgsUsage:       "<golang expression>",
		Description:     "Kicks a member via mention or ID.",
	})
}
