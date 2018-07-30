package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/gomodule/redigo/redis"
)

/**
 * events.go
 * Chase Weaver
 *
 * This package bundles event commands when they are triggered.
 */

// MessageCreate :
// Triggers on a message that is visible to the bot.
// Handles message and command responses.
func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Loads Message into temp cache
	c.Set(m.ID, m.Message, 0)

	// Default bot prefix
	prefix := conf.Prefix

	// Fetches channel object
	channel, err := s.State.Channel(m.ChannelID)
	if err != nil {
		return
	}

	// Gets the guild prefix from database
	if channel.Type == discordgo.ChannelTypeGuildText {
		guild, err := s.State.Guild(channel.GuildID)

		if err != nil {
			log.Println(err)
		}

		// Registers a new guild if not done already
		RegisterNewGuild(guild)

		data, err := redis.Bytes(p.Do("GET", guild.ID))

		if err != nil {
			log.Println(err)
		}

		var g Guild
		err = json.Unmarshal(data, &g)
		prefix = g.GuildPrefix
	}

	// Checks if message content begins with prefix
	if !strings.HasPrefix(m.Content, prefix) {
		return
	}

	// Give context for command pass-in
	ctx := Context{
		Session: s,
		Event:   m,
		Channel: channel,
		Name:    strings.Split(strings.TrimPrefix(m.Content, prefix), " ")[0],
	}

	// Fetches guild object if text channel is NOT a DM
	if ctx.Channel.Type == discordgo.ChannelTypeGuildText {
		guild, err := s.State.Guild(ctx.Channel.GuildID)
		if err != nil {
			return
		}

		ctx.Guild = guild

		// Registers a new guild if not done already
		RegisterNewGuild(ctx.Guild)

		data, err := redis.Bytes(p.Do("GET", guild.ID))

		if err != nil {
			log.Println(err)
		}

		var g Guild
		err = json.Unmarshal(data, &g)
		prefix = g.GuildPrefix
	}

	// Returns a valid command using a name/alias
	ctx.Command = FetchCommand(ctx.Name)

	// Splits command arguments
	tmp := strings.TrimPrefix(m.Content, prefix)

	// Splits the arguments by the deliminator
	ctx.Args = strings.Split(tmp, ctx.Command.ArgsDelim)[1:]

	// Checks if the config for the command passes all checks and is part of a text channel in a guild
	if ctx.Channel.Type == discordgo.ChannelTypeGuildText && !CommandIsValid(ctx) {
		return
	}

	// Checks if the command can be ran in a DM or not
	if ctx.Channel.Type == discordgo.ChannelTypeDM && !Contains(ctx.Command.RunIn, "DM") {
		return
	}

	// Fetch command funcs from command properties init()
	funcs := map[string]interface{}{
		ctx.Command.Name: ctx.Command.Func,
	}

	// Log commands to console
	LogCommands(ctx)

	// Type in channel
	err = ctx.Session.ChannelTyping(ctx.Channel.ID)

	if err != nil {
		log.Println(err)
	}

	// Call command with args pass-in
	Call(funcs, FetchCommandName(ctx.Name), ctx)
}

// GuildCreate :
// Initializes a new guild when the bot is first added.
func GuildCreate(s *discordgo.Session, m *discordgo.GuildCreate) {

	if GuildExists(m.Guild) {
		return
	}

	// Register new guild
	_, err := RegisterNewGuild(m.Guild)

	if err != nil {
		log.Println(err)
	}

	log.Println(
		fmt.Sprintf(`
			== New Guild Added ==\n
			Guild Name: %s\n
			Guild ID:   %s\n`,
			m.Guild.Name, m.Guild.ID))
}

// GuildDelete :
// Removes a guild when the bot is removed from a guild.
func GuildDelete(s *discordgo.Session, m *discordgo.GuildDelete) {

	if !GuildExists(m.Guild) {
		return
	}

	// Delete guild key
	_, err := DeleteGuild(m.Guild)

	if err != nil {
		log.Println(err)
	}

	log.Println(
		fmt.Sprintf(`
			== Guild Removed ==\n
			Guild Name: %s\n
			Guild ID:   %s\n`,
			m.Guild.Name, m.Guild.ID))
}

// GuildMemberAdd :
// Adds a new member to the guild database.
// Logs member to specified guild channel.
// Welcomes guild member in specified guild channel.
func GuildMemberAdd(s *discordgo.Session, m *discordgo.GuildMemberAdd) {

	// Fetch Guild information from redis database
	g, err := UnpackGuildStruct(m.GuildID)
	if err != nil {
		log.Println(err)
		return
	}

	// Get guild from user ID
	guild, err := s.Guild(m.GuildID)

	if err != nil {
		log.Println(err)
		return
	}

	// Check for User ID in Guild map, register user if missing
	if _, ok := g.GuildUser[m.User.ID]; !ok {
		user := RegisterNewUser(m.User)
		g.GuildUser[m.User.ID] = user
	}

	// Send a formatted message to the welcome channel
	if g.WelcomeChannel != nil && len(g.WelcomeMessage) != 0 {

		// Format welcome message
		msg := FormatWelcomeGoodbyeMessage(guild, m.Member, g.WelcomeMessage)
		_, err := s.ChannelMessageSend(g.WelcomeChannel.ID, msg)

		if err != nil {
			log.Println(err)
			return
		}
	}

	// Send a formatted message to the welcome logger channel
	if g.MemberAddChannel != nil && len(g.MemberAddMessage) != 0 {

		// Format welcome message
		msg := FormatWelcomeGoodbyeMessage(guild, m.Member, g.MemberAddMessage)
		_, err := s.ChannelMessageSend(g.MemberAddChannel.ID, msg)

		if err != nil {
			log.Println(err)
			return
		}
	}
}

// GuildMemberRemove :
// Logs member to specified guild channel.
// Says goodbye to guild member in specified guild channel.
func GuildMemberRemove(s *discordgo.Session, m *discordgo.GuildMemberRemove) {

	// Fetch Guild information from redis database
	g, err := UnpackGuildStruct(m.GuildID)
	if err != nil {
		log.Println(err)
		return
	}

	// Get guild from user ID
	guild, err := s.Guild(m.GuildID)

	if err != nil {
		log.Println(err)
		return
	}

	// Send a formatted message to the goodbye channel
	if g.GoodbyeChannel != nil && len(g.GoodbyeMessage) != 0 {

		// Format goodbye message
		msg := FormatWelcomeGoodbyeMessage(guild, m.Member, g.GoodbyeMessage)
		_, err := s.ChannelMessageSend(g.GoodbyeChannel.ID, msg)

		if err != nil {
			log.Println(err)
			return
		}

	}

	// Send a formatted message to the goodbye logger channel
	if g.MemberRemoveChannel != nil && len(g.MemberRemoveMessage) != 0 {

		// Format goodbye message
		msg := FormatWelcomeGoodbyeMessage(guild, m.Member, g.MemberRemoveMessage)
		_, err := s.ChannelMessageSend(g.MemberRemoveChannel.ID, msg)

		if err != nil {
			log.Println(err)
			return
		}

	}
}

// MessageDelete :
// This should be reworked to include > 1024 character limit
// Logs deleted message to specified guild channel.
func MessageDelete(s *discordgo.Session, m *discordgo.MessageDelete) {

	// Fetch message from cache
	msg, found := c.Get(m.ID)
	if !found {
		return
	}

	// Get cached message object
	mo := msg.(*discordgo.Message)

	// Ignore messages deleted by bots
	if mo.Author.Bot {
		return
	}

	// Fetch Guild information from redis database
	g, err := UnpackGuildStruct(m.GuildID)
	if err != nil {
		log.Println(err)
		return
	}

	// Send deleted message to the guild deleted-channel
	if g.MessageDeleteChannel != nil {

		_, err := s.ChannelMessageSendEmbed(g.MessageDeleteChannel.ID,
			NewEmbed().
				SetTitle("Deleted Message").
				SetColor(deleteColor).
				SetAuthor(fmt.Sprintf("%s#%s / %s", mo.Author.Username, mo.Author.Discriminator, mo.Author.ID), mo.Author.AvatarURL("256"), mo.Author.AvatarURL("2048")).
				AddField("Channel", fmt.Sprintf("<#%s>", m.ChannelID)).
				AddField("Content", mo.Content).
				SetTimestamp(time.Now().Format(time.RFC3339)).MessageEmbed)

		if err != nil {
			log.Println(err)
			return
		}

	}
}

// MessageUpdate :
// This should be reworked to include > 1024 character limit
// Logs edited message to specified guild channel.
func MessageUpdate(s *discordgo.Session, m *discordgo.MessageUpdate) {

	// Fetch message from cache
	msg, found := c.Get(m.ID)
	if !found {
		return
	}

	// Get cached message object
	mo := msg.(*discordgo.Message)

	// Ignore messages edited by bots
	if mo.Author.Bot {
		return
	}

	// Fetch Guild information from redis database
	g, err := UnpackGuildStruct(m.GuildID)
	if err != nil {
		log.Println(err)
		return
	}

	// Send edited message to the guild edited-channel
	if g.MessageEditChannel != nil {

		_, err := s.ChannelMessageSendEmbed(g.MessageEditChannel.ID,
			NewEmbed().
				SetTitle("Edited Message").
				SetColor(editColor).
				SetAuthor(fmt.Sprintf("%s#%s / %s", mo.Author.Username, mo.Author.Discriminator, mo.Author.ID), mo.Author.AvatarURL("256"), mo.Author.AvatarURL("2048")).
				AddField("Channel", fmt.Sprintf("<#%s>", m.ChannelID)).
				AddField("Old Content", mo.Content).
				AddField("New Content", m.Content).
				SetTimestamp(time.Now().Format(time.RFC3339)).MessageEmbed)

		if err != nil {
			log.Println(err)
			return
		}

	}
}

// GuildMemberUpdate :
// Logs changes to guild members and saves them to the database
func GuildMemberUpdate(s *discordgo.Session, m *discordgo.GuildMemberUpdate) {

	// Fetch Guild information from redis database
	g, err := UnpackGuildStruct(m.GuildID)
	if err != nil {
		log.Println(err)
		return
	}

	// Check for User ID in Guild map, register user if missing
	user := RegisterNewUser(m.User)
	if _, ok := g.GuildUser[m.User.ID]; ok {
		user = g.GuildUser[m.User.ID]
	}

	// Append Username changes
	if user.Usernames == nil || len(user.Usernames) == 0 {
		username := Usernames{
			Username:      m.User.Username,
			Discriminator: m.User.Discriminator,
			Time:          time.Now(),
		}
		user.Usernames[MakeTimestamp()] = username
	} else {

		// Map keys to array
		var keys []int64
		for k := range user.Usernames {
			keys = append(keys, k)
		}

		// Sort the keys by time
		sort.Slice(keys, func(i, j int) bool {
			return keys[i] < keys[j]
		})

		// Compare last key data with current user
		if user.Usernames[keys[len(keys)-1]].Username != m.User.Username || user.Usernames[keys[len(keys)-1]].Discriminator != m.User.Discriminator {
			username := Usernames{
				Username:      m.User.Username,
				Discriminator: m.User.Discriminator,
				Time:          time.Now(),
			}
			user.Usernames[MakeTimestamp()] = username
		}
	}

	// Append Nickname changes
	if user.Nicknames == nil || len(user.Nicknames) == 0 {
		nick := m.Nick

		if nick == "" {
			nick = "RESET NICKNAME"
		}

		nickname := Nicknames{
			Nickname: nick,
			Time:     time.Now(),
		}

		user.Nicknames[MakeTimestamp()] = nickname
	} else {

		// Map keys to array
		var keys []int64
		for k := range user.Nicknames {
			keys = append(keys, k)
		}

		// Sort the keys by time
		sort.Slice(keys, func(i, j int) bool {
			return keys[i] < keys[j]
		})

		nick := m.Nick

		// Compare last key data with current user
		if user.Nicknames[keys[len(keys)-1]].Nickname != nick {
			if nick == "" {
				nick = "RESET NICKNAME"
			}

			nickname := Nicknames{
				Nickname: nick,
				Time:     time.Now(),
			}

			user.Nicknames[MakeTimestamp()] = nickname
		}
	}

	g.GuildUser[m.User.ID] = user

	err = PackGuildStruct(m.GuildID, g)
	if err != nil {
		log.Println(err)
		return
	}
}
