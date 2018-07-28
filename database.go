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
		GuildUser            map[string]GuildUser
		BlacklistedUsers     []*discordgo.User
		BlacklistedChannels  []*discordgo.Channel
		AutoRole             []*discordgo.Role
		MutedRole            *discordgo.Role
		DisabledCommands     []Command
	}

	// GuildUser information
	GuildUser struct {
		User      *discordgo.User
		Member    *discordgo.Member
		Age       string
		JoinedAt  string
		Usernames []Usernames
		Nicknames []Nicknames
		Roles     []*discordgo.Role
		Warnings  map[string]Warnings
		Kicks     map[string]Kicks
		Bans      map[string]Bans
		Mutes     map[string]Mutes
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

	// Usernames of user
	Usernames struct {
		Username      string
		Discriminator string
		Time          time.Time
	}

	// Nicknames of user
	Nicknames struct {
		Nickname string
		Time     time.Time
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

// UnpackGuildStruct :
// Fetches guild struct from database.
func UnpackGuildStruct(guildID string) (Guild, error) {

	// Fetch Guild information from redis database
	data, err := redis.Bytes(p.Do("GET", guildID))
	if err != nil {
		log.Println(err)
		return Guild{}, err
	}

	var g Guild
	err = json.Unmarshal(data, &g)

	if err != nil {
		log.Println(err)
	}

	return g, nil
}

// PackGuildStruct :
// Pushes guild struct to database.
func PackGuildStruct(guildID string, g Guild) error {

	serialized, err := json.Marshal(g)

	if err != nil {
		log.Println(err)
		return err
	}

	_, err = p.Do("SET", guildID, serialized)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
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
		Guild:               guild,
		GuildPrefix:         conf.Prefix,
		WelcomeMessage:      "Welcome $MEMBER_MENTION$ to $GUILD_NAME$! Enjoy your stay.",
		GoodbyeMessage:      "Goodbye, `$MEMBER_NAME$`!",
		MemberAddMessage:    "✅ | `$MEMBER_NAME&` (ID: $MEMBER_ID$ | Age: $MEMBER_AGE$) has joinied the guild.",
		MemberRemoveMessage: "❌ | `$MEMBER_NAME&` (ID: $MEMBER_ID$ | Age: $MEMBER_AGE$ | Joined At: $MEMBER_JOINED$) has left the guild.",
		GuildUser:           make(map[string]GuildUser),
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
func RegisterNewUser(user *discordgo.User) GuildUser {

	// Fetch Account creation time
	age, err := CreationTime(user.ID)

	if err != nil {
		log.Println(err)
	}

	// Create a new GuildUser
	nu := GuildUser{
		User:      user,
		Age:       age.Format("01/02/06 03:04:05 PM MST"),
		Usernames: []Usernames{},
		Nicknames: []Nicknames{},
		Roles:     []*discordgo.Role{},
		Warnings:  make(map[string]Warnings),
		Kicks:     make(map[string]Kicks),
		Bans:      make(map[string]Bans),
		Mutes:     make(map[string]Mutes),
	}

	// Add the current username to the GuildUser
	nu.Usernames = append(nu.Usernames, Usernames{
		Username:      user.Username,
		Discriminator: user.Discriminator,
		Time:          time.Now(),
	})

	return nu
}

// LogWarning :
// Logs a warning to a user's record in the redis database.
func LogWarning(ctx Context, mem *discordgo.User, reason string) {

	// Fetch guild information
	g, err := UnpackGuildStruct(ctx.Guild.ID)
	if err != nil {
		log.Println(err)
		return
	}

	// Check for existing GuildUser, map warnings
	if _, ok := g.GuildUser[mem.ID]; ok {

		user := g.GuildUser[mem.ID]
		warning := Warnings{
			AuthorUser: ctx.Event.Author,
			TargetUser: mem,
			Channel:    ctx.Channel,
			Reason:     reason,
			Time:       time.Now(),
		}

		user.Warnings[TimeKey()] = warning
		g.GuildUser[mem.ID] = user

	} else {

		user := RegisterNewUser(mem)
		warning := Warnings{
			AuthorUser: ctx.Event.Author,
			TargetUser: mem,
			Channel:    ctx.Channel,
			Reason:     reason,
			Time:       time.Now(),
		}

		user.Warnings[TimeKey()] = warning
		g.GuildUser[mem.ID] = user

	}

	// Pack guild information
	err = PackGuildStruct(ctx.Guild.ID, g)
	if err != nil {
		log.Println(err)
		return
	}
}

// LogKick :
// Logs a kick to a user's record in the redis database.
func LogKick(ctx Context, mem *discordgo.User, reason string) {

	// Fetch guild information
	g, err := UnpackGuildStruct(ctx.Guild.ID)
	if err != nil {
		log.Println(err)
		return
	}

	// Check for existing GuildUser, map kicks
	if _, ok := g.GuildUser[mem.ID]; ok {

		user := g.GuildUser[mem.ID]
		kick := Kicks{
			AuthorUser: ctx.Event.Author,
			TargetUser: mem,
			Channel:    ctx.Channel,
			Reason:     reason,
			Time:       time.Now(),
		}

		user.Kicks[TimeKey()] = kick
		g.GuildUser[mem.ID] = user

	} else {

		user := RegisterNewUser(mem)
		kick := Kicks{
			AuthorUser: ctx.Event.Author,
			TargetUser: mem,
			Channel:    ctx.Channel,
			Reason:     reason,
			Time:       time.Now(),
		}

		user.Kicks[TimeKey()] = kick
		g.GuildUser[mem.ID] = user

	}

	// Pack guild information
	err = PackGuildStruct(ctx.Guild.ID, g)
	if err != nil {
		log.Println(err)
		return
	}
}

// LogBan :
// Logs a ban to a user's record in the redis database.
func LogBan(ctx Context, mem *discordgo.User, reason string) {

	// Fetch guild information
	g, err := UnpackGuildStruct(ctx.Guild.ID)
	if err != nil {
		log.Println(err)
		return
	}

	// Check for existing GuildUser, map bans
	if _, ok := g.GuildUser[mem.ID]; ok {

		user := g.GuildUser[mem.ID]
		ban := Bans{
			AuthorUser: ctx.Event.Author,
			TargetUser: mem,
			Channel:    ctx.Channel,
			Reason:     reason,
			Time:       time.Now(),
		}

		user.Bans[TimeKey()] = ban
		g.GuildUser[mem.ID] = user

	} else {

		user := RegisterNewUser(mem)
		ban := Bans{
			AuthorUser: ctx.Event.Author,
			TargetUser: mem,
			Channel:    ctx.Channel,
			Reason:     reason,
			Time:       time.Now(),
		}

		user.Bans[TimeKey()] = ban
		g.GuildUser[mem.ID] = user

	}

	// Pack guild information
	err = PackGuildStruct(ctx.Guild.ID, g)
	if err != nil {
		log.Println(err)
		return
	}
}

// FormatKicks :
// Returns a string of kicks.
func FormatKicks(kicks []Kicks) string {

	str := "\n"
	for _, v := range kicks {
		avatar := fmt.Sprintf("%s#%s / %s", v.AuthorUser.Username, v.AuthorUser.Discriminator, v.AuthorUser.ID)
		channel := fmt.Sprintf("<#%s> / %s", v.Channel.ID, v.Channel.ID)

		str = str + fmt.Sprintf(
			"**Author**:\t%s\n"+
				"**Channel**:  %s\n"+
				"**Time**:\t\t%s\n"+
				"**Reason**:   %s\n\n",
			avatar, channel, v.Time.Format("01/02/06 03:04:05 PM MST"), v.Reason)
	}

	return str
}

// FormatBans :
// Returns a string of bans.
func FormatBans(bans []Bans) string {

	str := "\n"
	for _, v := range bans {
		avatar := fmt.Sprintf("%s#%s / %s", v.AuthorUser.Username, v.AuthorUser.Discriminator, v.AuthorUser.ID)
		channel := fmt.Sprintf("<#%s> / %s", v.Channel.ID, v.Channel.ID)

		str = str + fmt.Sprintf(
			"**Author**:\t%s\n"+
				"**Channel**:  %s\n"+
				"**Time**:\t\t%s\n"+
				"**Reason**:   %s\n\n",
			avatar, channel, v.Time.Format("01/02/06 03:04:05 PM MST"), v.Reason)
	}

	return str
}

// FormatUsernames :
// Returns a string of usernames.
func FormatUsernames(usernames []Usernames) string {

	str := "\n"
	for _, v := range usernames {
		str = str + fmt.Sprintf(
			"**Username**:\t%s\n"+
				"**Time**:\t\t%s\n\n",
			v.Username+"#"+v.Discriminator, v.Time.Format("01/02/06 03:04:05 PM MST"))
	}

	return str
}

// FormatNicknames :
// Returns a string of nicknames.
func FormatNicknames(nicknames []Nicknames) string {

	str := "\n"
	for _, v := range nicknames {
		str = str + fmt.Sprintf(
			"**Nickname**:\t%s\n"+
				"**Time**:\t\t%s\n\n",
			v.Nickname, v.Time.Format("01/02/06 03:04:05 PM MST"))
	}

	return str
}
