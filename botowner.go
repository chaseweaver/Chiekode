package main

import (
	"log"
	"strings"

	"github.com/novalagung/golpal"
)

/**
 * owner.go
 * Chase Weaver
 *
 * This package bundles commands for the owner of the bot.
 */

func init() {
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
		UserPermissions: []string{"Bot Owner"},
		ArgsDelim:       " ",
		ArgsUsage:       "<golang expression>",
		Description:     "Evaluation command for bot-owner only.",
	})

	RegisterNewCommand(Command{
		Name:            "registerguild",
		Func:            RegisterGuild,
		Enabled:         true,
		NSFWOnly:        false,
		IgnoreSelf:      true,
		IgnoreBots:      true,
		RunIn:           []string{"Text"},
		Aliases:         []string{},
		UserPermissions: []string{"Bot Owner", "Administrator"},
		ArgsDelim:       "",
		ArgsUsage:       "",
		Description:     "Registers a guild in the database if non-existent",
	})

	RegisterNewCommand(Command{
		Name:            "removeguild",
		Func:            RemoveGuild,
		Enabled:         true,
		NSFWOnly:        false,
		IgnoreSelf:      true,
		IgnoreBots:      true,
		RunIn:           []string{"Text"},
		Aliases:         []string{},
		UserPermissions: []string{"Bot Owner"},
		ArgsDelim:       "",
		ArgsUsage:       "",
		Description:     "Removes a guild from the database if existent",
	})
}

// Ping command will return Pong!
func Ping(ctx Context) {
	ctx.Session.ChannelMessageSend(ctx.Channel.ID, "🏓 Pong!")
}

// Eval is the bot's evaluate command for complex functions
func Eval(ctx Context) {
	out, err := golpal.New().AddLibs("strings", "runtime", "fmt", "github.com/bwmarrin/discordgo").ExecuteRaw(strings.Join(ctx.Args, " "))
	if err != nil {
		ctx.Session.ChannelMessageSend(ctx.Channel.ID, "**ERROR**\n"+FormatString(err.Error(), "go"))
	} else {
		ctx.Session.ChannelMessageSend(ctx.Channel.ID, "**RESULT**\n"+FormatString(out, "go"))
	}
}

// RegisterGuild creates a new guild in the database if non-existent
func RegisterGuild(ctx Context) {

	if GuildExists(ctx.Guild) {
		ctx.Session.ChannelMessageSend(ctx.Channel.ID, "```The guild already exists in the databse.```")
		return
	}

	_, err := RegisterNewGuild(ctx.Guild)

	if err != nil {
		log.Println(err)
		return
	}

	ctx.Session.ChannelMessageSend(ctx.Channel.ID, "```Guild has successfully been registered.```")
}

// RemoveGuild deletes a guild in the database if existent
func RemoveGuild(ctx Context) {

	if !GuildExists(ctx.Guild) {
		ctx.Session.ChannelMessageSend(ctx.Channel.ID, "```The guild does not exist in the database.```")
		return
	}

	_, err := DeleteGuild(ctx.Guild)

	if err != nil {
		log.Println(err)
		return
	}

	ctx.Session.ChannelMessageSend(ctx.Channel.ID, "```Guild has successfully been removed.```")
}
