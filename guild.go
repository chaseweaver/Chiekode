package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/gomodule/redigo/redis"
)

func init() {
	RegisterNewCommand(Command{
		Name:            "settings",
		Func:            Settings,
		Enabled:         true,
		NSFWOnly:        false,
		IgnoreSelf:      true,
		IgnoreBots:      true,
		RunIn:           []string{"Text"},
		Aliases:         []string{},
		UserPermissions: []string{"Bot Owner"},
		ArgsDelim:       " ",
		ArgsUsage:       "[settings name]",
		Description:     "Lists guild configurations.",
	})

	RegisterNewCommand(Command{
		Name:            "set",
		Func:            Set,
		Enabled:         true,
		NSFWOnly:        false,
		IgnoreSelf:      true,
		IgnoreBots:      true,
		RunIn:           []string{"Text"},
		Aliases:         []string{},
		UserPermissions: []string{"Bot Owner"},
		ArgsDelim:       " | ",
		ArgsUsage:       "<settings name> <value>",
		Description:     "Sets guild configurations.",
	})
}

// Settings lists database guild configurations
func Settings(ctx Context) {

	p := pool.Get()
	data, err := redis.Bytes(p.Do("GET", ctx.Guild.ID))
	if err != nil {
		panic(err.Error())
	}

	var g Guild
	err = json.Unmarshal(data, &g)

	key := strings.Join(ctx.Args[0:], " ")

	if len(ctx.Args) == 0 {
		str := fmt.Sprintf(
			"== %s Configuration ==\n\n"+
				"Guild Prefix ::              %s\n"+
				"Blacklisted Channels IDs ::  %d\n"+
				"Blacklisted Members ::       %d\n"+
				"Welcome Message ::           %s\n"+
				"Welcome Channel ID ::        %s\n"+
				"Goodbye Message ::		        %s\n"+
				"Goodbye Channel ID ::        %s\n"+
				"Events ::                    %d\n"+
				"Disabled Commands ::         %d\n"+
				"Birthday Role ID ::          %s\n"+
				"Muted Role ID ::             %s\n"+
				"Auto Role IDs ::             %d",
			g.GuildName, g.GuildPrefix, len(g.BlacklistedChannels), len(g.BlacklistedMembers),
			g.WelcomeMessage, g.WelcomeChannel, g.GoodbyeMessage, g.GoodbyeChannel, len(g.Events),
			len(g.DisabledCommands), g.BirthdayRole, g.MutedRole, len(g.AutoRole))

		ctx.Session.ChannelMessageSend(ctx.Channel.ID, FormatString(str, "asciidoc"))
		return
	}

	switch key {
	case "Guild Prefix":
		ctx.Session.ChannelMessageSend(ctx.Channel.ID, fmt.Sprintf("The Guild Prefix is set to `%s`", g.GuildPrefix))
	case "Runnable Channels IDs":
	case "Blacklisted Members":
	case "Welcome Message":
		ctx.Session.ChannelMessageSend(ctx.Channel.ID, fmt.Sprintf("The Welcome Message is set to `%s`", g.WelcomeMessage))
	case "Welcome Channel ID":
		channel, err := ctx.Session.Channel(g.WelcomeChannel)
		if err != nil {
			log.Println(err)
		}
		ctx.Session.ChannelMessageSend(ctx.Channel.ID, fmt.Sprintf("The Welcome Channel is set to `%s`", channel.Name))
	case "Goodbye Message":
		ctx.Session.ChannelMessageSend(ctx.Channel.ID, fmt.Sprintf("The Goodbye Message is set to `%s`", g.GoodbyeMessage))
	case "Goodbye Channel ID":
		channel, err := ctx.Session.Channel(g.GoodbyeChannel)
		if err != nil {
			log.Println(err)
		}
		ctx.Session.ChannelMessageSend(ctx.Channel.ID, fmt.Sprintf("The Goodbye Channel is set to `%s`", channel.Name))
	default:
		ctx.Session.ChannelMessageSend(ctx.Channel.ID, fmt.Sprintf("I could not find the guild setting `%s`", key))
	}
}

// Set allows configration of database guild settings
func Set(ctx Context) {

	p := pool.Get()
	data, err := redis.Bytes(p.Do("GET", ctx.Guild.ID))
	if err != nil {
		log.Println(err)
	}

	var g Guild
	err = json.Unmarshal(data, &g)

	if err != nil {
		log.Println(err)
	}

	if len(ctx.Args) == 0 {
		ctx.Session.ChannelMessageSend(ctx.Channel.ID, fmt.Sprintf("Invalid settings key. Run `%s`help for more information.", g.GuildPrefix))
		return
	}

	switch ctx.Args[0] {
	case "Guild Prefix":
		g.GuildPrefix = ctx.Args[1]
	}

	log.Println(g.GuildPrefix)

	serialized, err := json.Marshal(g)

	if err != nil {
		log.Println(err)
	}

	_, err = p.Do("SET", ctx.Guild.ID, serialized)
	if err != nil {
		log.Println(err)
		ctx.Session.ChannelMessageSend(ctx.Channel.ID, "**ERROR**\n"+FormatString(err.Error(), "doc"))
		return
	}

	ctx.Session.ChannelMessageSend(ctx.Channel.ID, fmt.Sprintf("Guild Prefix has successfully been changed to `%s`", ctx.Args[1]))
}
