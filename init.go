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
		Locked:          false,
		RunIn:           []string{"Text", "DM"},
		Aliases:         []string{},
		BotPermissions:  []string{},
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
		Locked:          false,
		RunIn:           []string{"Text", "DM"},
		Aliases:         []string{"pfp", "icon"},
		BotPermissions:  []string{},
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
		Locked:          false,
		RunIn:           []string{"Text", "DM"},
		Aliases:         []string{},
		BotPermissions:  []string{},
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
		Locked:          true,
		RunIn:           []string{"Text", "DM"},
		Aliases:         []string{"e"},
		BotPermissions:  []string{},
		UserPermissions: []string{},
		ArgsDelim:       " ",
		ArgsUsage:       "<golang expression>",
		Description:     "Evaluation command for owner only.",
	})
}
