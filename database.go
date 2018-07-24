package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/gomodule/redigo/redis"
)

type (

	// Guild configuration information per guild
	Guild struct {
		Guild                *discordgo.Guild
		GuildPrefix          string
		WelcomeMessage       string
		GoodbyeMessage       string
		MemberAddMessage     string
		MemberRemoveMessage  string
		WelcomeChannel       *discordgo.Channel
		GoodbyeChannel       *discordgo.Channel
		MemberAddChannel     *discordgo.Channel
		MemberRemoveChannel  *discordgo.Channel
		MessageEditChannel   *discordgo.Channel
		MessageDeleteChannel *discordgo.Channel
		GuildUser            []GuildUser
		BlacklistedUsers     []*discordgo.User
		BlacklistedChannels  []*discordgo.Channel
		AutoRole             []*discordgo.Role
		MutedRole            *discordgo.Role
		DisabledCommands     []Command
	}

	// GuildUser information
	GuildUser struct {
		User              *discordgo.User
		Member            *discordgo.Member
		Age               string
		JoinedAt          string
		PreviousUsernames []string
		PreviousNicknames []string
		Roles             []*discordgo.Role
		Warnings          []Warnings
		Kicks             []Kicks
		Bans              []Bans
		Mutes             []Mutes
	}

	// Warnings information for a user
	Warnings struct {
		AuthorUser *discordgo.User
		TargetUser *discordgo.User
		Channel    *discordgo.Channel
		Reason     string
		Time       time.Time
	}

	// Kicks information for a user
	Kicks struct {
		AuthorUser *discordgo.User
		TargetUser *discordgo.User
		Channel    *discordgo.Channel
		Reason     string
		Time       time.Time
	}

	// Bans information for a user
	Bans struct {
		AuthorUser *discordgo.User
		TargetUser *discordgo.User
		Channel    *discordgo.Channel
		Reason     string
		Time       time.Time
	}

	// Mutes information for a user
	Mutes struct {
		AuthorUser *discordgo.User
		TargetUser *discordgo.User
		Channel    *discordgo.Channel
		Reason     string
		Time       time.Time
		Length     time.Duration
	}
)

// DialNewPool connectes to a local Redis database by port pass-in.
func DialNewPool(net string, port string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000, // max number of connections
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial(net, port)
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
}

// DialNewPoolURL connectes to a Redis database by URL pass-in.
func DialNewPoolURL(url string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000, // max number of connections
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialURL(os.Getenv(url))
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
}

// DeleteGuild :
// Removes a guild from the database.
func DeleteGuild(guild *discordgo.Guild) (interface{}, error) {
	n, err := p.Do("DEL", guild.ID)
	if err != nil {
		log.Println(err)
		return n, err
	}
	return n, nil
}

// GuildExists :
// Checks if guild key exists and returns true, otherwise false.
func GuildExists(guild *discordgo.Guild) bool {
	n, err := p.Do("EXISTS", guild.ID)
	if err != nil {
		log.Println(err)
	}

	if n != 1 {
		return true
	}

	return false
}

// RegisterNewGuild :
// Creates a key with a guild ID and given values.
func RegisterNewGuild(guild *discordgo.Guild) (interface{}, error) {

	// Initialize guild prefix with configuration default
	g := &Guild{
		Guild:          guild,
		GuildPrefix:    conf.Prefix,
		WelcomeMessage: "Welcome $MEMBER_MENTION$ to $GUILD_NAME$! Enjoy your stay.",
		GoodbyeMessage: "Goodbye, `$MEMBER_NAME$`!",
	}

	serialized, err := json.Marshal(g)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	n, err := p.Do("SETNX", guild.ID, serialized)
	if err != nil {
		log.Println(err)
	}

	return n, nil
}

// RegisterNewUser :
// Creates a new guild user with user defaults.
func RegisterNewUser(ctx Context, user *discordgo.User) GuildUser {

	// Fetch Account creation time
	age, err := CreationTime(user.ID)

	if err != nil {
		log.Println(err)
	}

	return GuildUser{
		User:              user,
		Age:               age.Format("01/02/06 03:04:05 PM MST"),
		PreviousUsernames: []string{},
		PreviousNicknames: []string{},
		Roles:             []*discordgo.Role{},
		Warnings:          []Warnings{},
		Kicks:             []Kicks{},
		Bans:              []Bans{},
		Mutes:             []Mutes{},
	}

}

// LogWarning :
// Logs a warning to a user's record in the redis database.
func LogWarning(ctx Context, mem *discordgo.User, reason string) {

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

	found := false

	for u := range g.GuildUser {
		if g.GuildUser[u].User.ID == mem.ID {
			found = true
		}
	}

	if found {
		for u := range g.GuildUser {
			if g.GuildUser[u].User.ID == mem.ID {
				g.GuildUser[u].Warnings = append(g.GuildUser[u].Warnings, Warnings{
					AuthorUser: ctx.Event.Author,
					TargetUser: mem,
					Channel:    ctx.Channel,
					Reason:     reason,
					Time:       time.Now(),
				})
			}
		}
	} else {
		newUser := RegisterNewUser(ctx, mem)
		newUser.Warnings = append(newUser.Warnings, Warnings{
			AuthorUser: ctx.Event.Author,
			TargetUser: mem,
			Channel:    ctx.Channel,
			Reason:     reason,
			Time:       time.Now(),
		})
		g.GuildUser = append(g.GuildUser, newUser)
	}

	serialized, err := json.Marshal(g)

	if err != nil {
		log.Println(err)
		return
	}

	_, err = p.Do("SET", ctx.Guild.ID, serialized)
	if err != nil {
		log.Println(err)
	}
}

// LogKick :
// Logs a kick to a user's record in the redis database.
func LogKick(ctx Context, mem *discordgo.User, reason string) {

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

	found := false

	for u := range g.GuildUser {
		if g.GuildUser[u].User.ID == mem.ID {
			found = true
		}
	}

	if found {
		for u := range g.GuildUser {
			if g.GuildUser[u].User.ID == mem.ID {
				g.GuildUser[u].Kicks = append(g.GuildUser[u].Kicks, Kicks{
					AuthorUser: ctx.Event.Author,
					TargetUser: mem,
					Channel:    ctx.Channel,
					Reason:     reason,
					Time:       time.Now(),
				})
			}
		}
	} else {
		newUser := RegisterNewUser(ctx, mem)
		newUser.Kicks = append(newUser.Kicks, Kicks{
			AuthorUser: ctx.Event.Author,
			TargetUser: mem,
			Channel:    ctx.Channel,
			Reason:     reason,
			Time:       time.Now(),
		})
		g.GuildUser = append(g.GuildUser, newUser)
	}

	serialized, err := json.Marshal(g)

	if err != nil {
		log.Println(err)
		return
	}

	_, err = p.Do("SET", ctx.Guild.ID, serialized)
	if err != nil {
		log.Println(err)
	}
}

// LogBan :
// Logs a ban to a user's record in the redis database.
func LogBan(ctx Context, mem *discordgo.User, reason string) {

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

	found := false

	for u := range g.GuildUser {
		if g.GuildUser[u].User.ID == mem.ID {
			found = true
		}
	}

	if found {
		for u := range g.GuildUser {
			if g.GuildUser[u].User.ID == mem.ID {
				g.GuildUser[u].Bans = append(g.GuildUser[u].Bans, Bans{
					AuthorUser: ctx.Event.Author,
					TargetUser: mem,
					Channel:    ctx.Channel,
					Reason:     reason,
					Time:       time.Now(),
				})
			}
		}
	} else {
		newUser := RegisterNewUser(ctx, mem)
		newUser.Bans = append(newUser.Bans, Bans{
			AuthorUser: ctx.Event.Author,
			TargetUser: mem,
			Channel:    ctx.Channel,
			Reason:     reason,
			Time:       time.Now(),
		})
		g.GuildUser = append(g.GuildUser, newUser)
	}

	serialized, err := json.Marshal(g)

	if err != nil {
		log.Println(err)
		return
	}

	_, err = p.Do("SET", ctx.Guild.ID, serialized)
	if err != nil {
		log.Println(err)
	}
}

// FormatWarning :
// Returns a string of warnings.
func FormatWarning(warnings []Warnings) string {

	str := "\n"
	for _, v := range warnings {
		avatar := fmt.Sprintf("%s#%s / %s", v.AuthorUser.Username, v.AuthorUser.Discriminator, v.AuthorUser.ID)
		channel := fmt.Sprintf("<#%s> / %s", v.Channel.ID, v.Channel.ID)

		str = str + fmt.Sprintf(
			"**Author**:\t%s\n"+
				"**Channel**:  %s\n"+
				"**Time**:\t\t%s\n"+
				"**Reason**:\t%s\n\n",
			avatar, channel, v.Time.Format("01/02/06 03:04:05 PM MST"), v.Reason)
	}

	return str
}

// FormatKick :
// Returns a string of kicks.
func FormatKick(kicks []Kicks) string {

	str := "\n"
	for _, v := range kicks {
		avatar := fmt.Sprintf("%s#%s / %s", v.AuthorUser.Username, v.AuthorUser.Discriminator, v.AuthorUser.ID)
		channel := fmt.Sprintf("<#%s> / %s", v.Channel.ID, v.Channel.ID)

		str = str + fmt.Sprintf(
			"**Author**:\t%s\n"+
				"**Channel**:  %s\n"+
				"**Time**:\t\t%s\n"+
				"**Reason**:\t%s\n\n",
			avatar, channel, v.Time.Format("01/02/06 03:04:05 PM MST"), v.Reason)
	}

	return str
}

// FormatBan :
// Returns a string of bans.
func FormatBan(bans []Bans) string {

	str := "\n"
	for _, v := range bans {
		avatar := fmt.Sprintf("%s#%s / %s", v.AuthorUser.Username, v.AuthorUser.Discriminator, v.AuthorUser.ID)
		channel := fmt.Sprintf("<#%s> / %s", v.Channel.ID, v.Channel.ID)

		str = str + fmt.Sprintf(
			"**Author**:\t%s\n"+
				"**Channel**:  %s\n"+
				"**Time**:\t\t%s\n"+
				"**Reason**:\t%s\n\n",
			avatar, channel, v.Time.Format("01/02/06 03:04:05 PM MST"), v.Reason)
	}

	return str
}
