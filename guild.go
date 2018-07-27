package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"

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
		Cooldown:        0,
		RunIn:           []string{"Text"},
		Aliases:         []string{},
		UserPermissions: []string{"Bot Owner"},
		ArgsDelim:       " ",
		Usage:           []string{"[Guild Setting]"},
		Description:     "Lists guild configurations.",
	})

	RegisterNewCommand(Command{
		Name:            "set",
		Func:            Set,
		Enabled:         true,
		NSFWOnly:        false,
		IgnoreSelf:      true,
		IgnoreBots:      true,
		Cooldown:        0,
		RunIn:           []string{"Text"},
		Aliases:         []string{},
		UserPermissions: []string{"Bot Owner"},
		ArgsDelim:       " | ",
		Usage:           []string{"<Guild Setting>", "<value>"},
		Description:     "Sets guild configurations.",
	})

	RegisterNewCommand(Command{
		Name:            "resetguildsettings",
		Func:            ResetGuildSettings,
		Enabled:         true,
		NSFWOnly:        false,
		IgnoreSelf:      true,
		IgnoreBots:      true,
		Cooldown:        0,
		RunIn:           []string{"Text"},
		Aliases:         []string{},
		UserPermissions: []string{"Bot Owner", "Administrator"},
		ArgsDelim:       "",
		Usage:           []string{},
		Description:     "Resets guild configurations.",
	})
}

// Settings lists database guild configurations
func Settings(ctx Context) {

	// Fetch guild settings
	data, err := redis.Bytes(p.Do("GET", ctx.Guild.ID))
	if err != nil {
		panic(err.Error())
	}

	var g Guild
	err = json.Unmarshal(data, &g)

	var blc, blu, ar []string
	for _, v := range g.BlacklistedChannels {
		blc = append(blc, v.Name)
	}

	for _, v := range g.BlacklistedUsers {
		blu = append(blu, v.Username+"#"+v.Discriminator)
	}

	for _, v := range g.AutoRole {
		ar = append(ar, v.Name)
	}

	wc := " "
	if g.WelcomeChannel != nil {
		wc = g.WelcomeChannel.Name
	}

	gc := " "
	if g.GoodbyeChannel != nil {
		gc = g.WelcomeChannel.Name
	}

	dc := " "
	if g.MessageDeleteChannel != nil {
		dc = g.MessageDeleteChannel.Name
	}

	ec := " "
	if g.MessageEditChannel != nil {
		dc = g.MessageEditChannel.Name
	}

	str := fmt.Sprintf(
		"== %s Configuration ==\n\n"+
			"Guild Prefix             ::   %s\n"+
			"Blacklisted Channel      ::   %s\n"+
			"Blacklisted Members      ::   %s\n"+
			"Welcome Message          ::   %s\n"+
			"Welcome Channel          ::   %s\n"+
			"Goodbye Message          ::   %s\n"+
			"Goodbye Channel          ::   %s\n"+
			"Message Deleted Channel  ::   %s\n"+
			"Message Edited Channel   ::   %s\n"+
			"Muted Role               ::   %s\n"+
			"Auto Roles               ::   %s",
		g.Guild.Name, g.GuildPrefix, strings.Join(blc, ", "), strings.Join(blu, ", "),
		g.WelcomeMessage, wc, g.GoodbyeMessage, gc, dc, ec, " ", strings.Join(ar, ", "))

	ctx.Session.ChannelMessageSend(ctx.Channel.ID, FormatString(str, "asciidoc"))
}

// Set :
// Allows configration of database guild settings
func Set(ctx Context) {

	// Fetch guild settings
	data, err := redis.Bytes(p.Do("GET", ctx.Guild.ID))
	if err != nil {
		log.Println(err)
	}

	var g Guild
	err = json.Unmarshal(data, &g)

	if err != nil {
		log.Println(err)
	}

	// Return if Guild Setting cannot be found
	if len(ctx.Args) <= 1 {
		msg, err := ctx.Session.ChannelMessageSend(ctx.Channel.ID, fmt.Sprintf("Invalid settings key/value. Run `%shelp` for more information.", g.GuildPrefix))

		if err != nil {
			log.Println(err)
		}

		DeleteMessageWithTime(ctx, msg.ID, 7500)
		return
	}

	// Get the Guild Setting to configure
	key := strings.ToUpper(ctx.Args[0])
	val := strings.Join(ctx.Args[1:], ctx.Command.ArgsDelim)

	switch key {
	case "PREFIX":
		fallthrough
	case "GUILD PREFIX":
		g.GuildPrefix = val
	case "BLACKLISTED CHANNEL":
		fallthrough
	case "BLACKLISTED CHANNELS":
		channels := FetchMessageContentChannels(ctx, val)

		if len(channels) == 0 {
			return
		}

		for _, v := range channels {
			g.BlacklistedChannels = append(g.BlacklistedChannels, v)
		}
	case "BLACKLISTED USER":
		fallthrough
	case "BLACKLISTED USERS":
		g.BlacklistedUsers = []*discordgo.User{}
		users := FetchMessageContentUsers(ctx, val)

		if len(users) == 0 {
			return
		}

		for _, v := range users {
			g.BlacklistedUsers = append(g.BlacklistedUsers, v)
		}
	case "WELCOME MESSAGE":
		g.WelcomeMessage = val
	case "WELCOME CHANNEL":
		channels := FetchMessageContentChannels(ctx, val)

		if len(channels) == 0 {
			return
		}

		g.WelcomeChannel = channels[0]
	case "GOODBYE MESSAGE":
		g.GoodbyeMessage = val
	case "GOODBYE CHANNEL":
		channels := FetchMessageContentChannels(ctx, val)

		if len(channels) == 0 {
			return
		}

		g.GoodbyeChannel = channels[0]
	case "MESSAGE DELETED":
		fallthrough
	case "MESSAGE DELETED CHANNEL":
		channels := FetchMessageContentChannels(ctx, val)

		if len(channels) == 0 {
			return
		}

		g.MessageDeleteChannel = channels[0]
	case "MESSAGE EDITED":
		fallthrough
	case "MESSAGE EDITED CHANNEL":
		channels := FetchMessageContentChannels(ctx, val)

		if len(channels) == 0 {
			return
		}

		g.MessageEditChannel = channels[0]
	case "DISABLED":
		fallthrough
	case "DISABLED COMMANDS":
		return
	case "MUTED":
		fallthrough
	case "MUTED ROLE":
		role := FetchMessageContentRoles(ctx, val)

		if len(role) == 0 {
			return
		}

		g.MutedRole = role[0]
	case "AUTO":
		fallthrough
	case "AUTO ROLE":
		fallthrough
	case "AUTO ROLES":
		g.AutoRole = []*discordgo.Role{}
		roles := FetchMessageContentRoles(ctx, val)

		if len(roles) == 0 {
			return
		}

		for _, v := range roles {
			g.AutoRole = append(g.AutoRole, v)
		}
	default:
		ctx.Session.ChannelMessageSend(ctx.Channel.ID, fmt.Sprintf("I could not find the guild setting `%s`", key))
		return
	}

	serialized, err := json.Marshal(g)

	if err != nil {
		log.Println(err)
	}

	_, err = p.Do("SET", ctx.Guild.ID, serialized)
	if err != nil {
		log.Println(err)
		return
	}

}

// ResetGuildSettings :
// Resets the guild to initial settings
func ResetGuildSettings(ctx Context) {

	// Reinitialize guild prefix with configuration defaults
	serialized, err := json.Marshal(&Guild{
		Guild:               ctx.Guild,
		GuildPrefix:         conf.Prefix,
		WelcomeMessage:      "Welcome $MEMBER_MENTION$ to $GUILD_NAME$! Enjoy your stay.",
		GoodbyeMessage:      "Goodbye, `$MEMBER_NAME$`!",
		MemberAddMessage:    "✅ | `$MEMBER_NAME&` (ID: $MEMBER_ID$ | Age: $MEMBER_AGE$) has joinied the guild.",
		MemberRemoveMessage: "❌ | `$MEMBER_NAME&` (ID: $MEMBER_ID$ | Age: $MEMBER_AGE$ | Joined At: $MEMBER_JOINED$) has left the guild.",
		GuildUser:           make(map[string]GuildUser),
	})

	if err != nil {
		log.Println(err)
		return
	}

	_, err = p.Do("SET", ctx.Guild.ID, serialized)
	if err != nil {
		log.Println(err)
		return
	}

	ctx.Session.ChannelMessageSend(ctx.Channel.ID, "✅ | Guild settings have been reset.")
}
