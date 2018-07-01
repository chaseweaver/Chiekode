package main

/**
 * init.go
 * Chase Weaver
 *
 * This package is meant to handle initialization of commands
 * and command properties for quick configuration.
 */

func init() {
	RegisterNewCommand("help", Command{
		Name:            "help",
		Func:            Help,
		Enabled:         true,
		NSFWOnly:        false,
		IgnoreSelf:      true,
		IgnoreBots:      true,
		RunIn:           []string{"GuildText", "DM"},
		Aliases:         []string{},
		BotPermissions:  []string{},
		UserPermissions: []string{},
		ArgsDelim:       " ",
		ArgsUsage:       "[command]",
		Description:     "Displays a helpful help menu.",
	})

	RegisterNewCommand("avatar", Command{
		Name:            "avatar",
		Func:            Avatar,
		Enabled:         true,
		NSFWOnly:        false,
		IgnoreSelf:      true,
		IgnoreBots:      true,
		RunIn:           []string{"GuildText", "DM"},
		Aliases:         []string{"pfp", "icon"},
		BotPermissions:  []string{},
		UserPermissions: []string{},
		ArgsDelim:       " ",
		ArgsUsage:       "[@member]",
		Description:     "Fetches the avatar/pfp for the requested member.",
	})

	RegisterNewCommand("ping", Command{
		Name:            "ping",
		Func:            Ping,
		Enabled:         true,
		NSFWOnly:        true,
		IgnoreSelf:      true,
		IgnoreBots:      true,
		RunIn:           []string{"GuildText", "DM"},
		Aliases:         []string{"t1", "t2"},
		BotPermissions:  []string{},
		UserPermissions: []string{},
		ArgsDelim:       " ",
		ArgsUsage:       "",
		Description:     "Test command",
	})
}