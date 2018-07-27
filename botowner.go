package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
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
		Cooldown:        30,
		RunIn:           []string{"Text", "DM"},
		Aliases:         []string{},
		UserPermissions: []string{},
		ArgsDelim:       "",
		Usage:           []string{},
		Description:     "Pong! Responds with the heartbeat.",
	})

	RegisterNewCommand(Command{
		Name:            "test",
		Func:            Test,
		Enabled:         true,
		NSFWOnly:        false,
		IgnoreSelf:      true,
		IgnoreBots:      true,
		Cooldown:        0,
		RunIn:           []string{"DM", "Text"},
		Aliases:         []string{},
		UserPermissions: []string{"Bot Owner"},
		ArgsDelim:       " ",
		Usage:           []string{},
		Description:     "Bot-owner testing function.",
	})
}

// Ping :
// Command will return Pong! with the last heartbeat.
func Ping(ctx Context) {
	t := time.Now().Sub(ctx.Session.LastHeartbeatAck) / 1000
	ctx.Session.ChannelMessageSend(ctx.Channel.ID, fmt.Sprintf("🏓 Pong! Heatbeat: `%s`", t))
}

// Test :
// Bot owner's test command
func Test(ctx Context) {

	member := FetchMessageContentUsers(ctx, strings.Join(ctx.Args, ctx.Command.ArgsDelim))[0]

	// Fetch Guild information from redis database
	data, err := redis.Bytes(p.Do("GET", ctx.Guild.ID))
	if err != nil {
		log.Println(err)
		return
	}

	var g Guild
	err = json.Unmarshal(data, &g)

	if err != nil {
		log.Println(err)
	}

	if _, ok := g.GuildUser[member.ID]; ok {
		log.Println(g.GuildUser)
	} else {
		log.Println("Not found!")
	}
}
