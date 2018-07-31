package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/gomodule/redigo/redis"
)

type (

	// Guild configuration information per guild
	Guild struct {
		Guild                 *discordgo.Guild
		GuildPrefix           string
		WelcomeMessage        string
		GoodbyeMessage        string
		MemberAddMessage      string
		MemberRemoveMessage   string
		WelcomeChannel        *discordgo.Channel
		GoodbyeChannel        *discordgo.Channel
		MemberAddChannel      *discordgo.Channel
		MemberRemoveChannel   *discordgo.Channel
		MessageEditChannel    *discordgo.Channel
		MessageDeleteChannel  *discordgo.Channel
		ModerationLogsChannel *discordgo.Channel
		GuildUser             map[string]GuildUser
		BlacklistedUsers      []*discordgo.User
		BlacklistedChannels   []*discordgo.Channel
		AutoRole              []*discordgo.Role
		MutedRole             *discordgo.Role
		DisabledCommands      []Command
	}

	// GuildUser information
	GuildUser struct {
		User      *discordgo.User
		Member    *discordgo.Member
		Age       string
		JoinedAt  string
		Muted     Muted
		Usernames map[int64]Usernames
		Nicknames map[int64]Nicknames
		Warnings  map[int64]Warnings
		Kicks     map[int64]Kicks
		Bans      map[int64]Bans
		Mutes     map[int64]Mutes
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

	// Muted information of user
	Muted struct {
		IsMuted       bool
		Time          time.Time
		RemainingTime time.Duration
	}
)

// DialNewPool connectes to a local Redis database by port pass-in.
func DialNewPool(net string, port string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000,
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
		Usernames: make(map[int64]Usernames),
		Nicknames: make(map[int64]Nicknames),
		Warnings:  make(map[int64]Warnings),
		Kicks:     make(map[int64]Kicks),
		Bans:      make(map[int64]Bans),
		Mutes:     make(map[int64]Mutes),
		Muted:     Muted{},
	}

	username := Usernames{
		Username:      user.Username,
		Discriminator: user.Discriminator,
		Time:          time.Now(),
	}

	// Add the current username to the GuildUser
	nu.Usernames[MakeTimestamp()] = username

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

		user.Warnings[MakeTimestamp()] = warning
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

		user.Warnings[MakeTimestamp()] = warning
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

		user.Kicks[MakeTimestamp()] = kick
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

		user.Kicks[MakeTimestamp()] = kick
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

		user.Bans[MakeTimestamp()] = ban
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

		user.Bans[MakeTimestamp()] = ban
		g.GuildUser[mem.ID] = user

	}

	// Pack guild information
	err = PackGuildStruct(ctx.Guild.ID, g)
	if err != nil {
		log.Println(err)
		return
	}
}

// LogMute :
// Logs a mute to a user's record in the redis database.
func LogMute(ctx Context, mem *discordgo.User, reason string, t time.Duration) {

	// Fetch guild information
	g, err := UnpackGuildStruct(ctx.Guild.ID)
	if err != nil {
		log.Println(err)
		return
	}

	// Check for existing GuildUser, map mutes
	if _, ok := g.GuildUser[mem.ID]; ok {

		user := g.GuildUser[mem.ID]
		mute := Mutes{
			AuthorUser: ctx.Event.Author,
			TargetUser: mem,
			Channel:    ctx.Channel,
			Reason:     reason,
			Time:       time.Now(),
			Length:     t,
		}

		user.Mutes[MakeTimestamp()] = mute
		g.GuildUser[mem.ID] = user

	} else {

		user := RegisterNewUser(mem)
		mute := Mutes{
			AuthorUser: ctx.Event.Author,
			TargetUser: mem,
			Channel:    ctx.Channel,
			Reason:     reason,
			Time:       time.Now(),
			Length:     t,
		}

		user.Mutes[MakeTimestamp()] = mute
		g.GuildUser[mem.ID] = user

	}

	// Pack guild information
	err = PackGuildStruct(ctx.Guild.ID, g)
	if err != nil {
		log.Println(err)
		return
	}
}

// FormatWarnings :
// Returns a string of warnings.
func FormatWarnings(warnings map[int64]Warnings) string {

	var keys []int64
	for k := range warnings {
		keys = append(keys, k)
	}

	// Sort the keys by time
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	str := "\n"

	for _, v := range keys {
		avatar := fmt.Sprintf("%s#%s / %s", warnings[v].AuthorUser.Username, warnings[v].AuthorUser.Discriminator, warnings[v].AuthorUser.ID)
		channel := fmt.Sprintf("<#%s> / %s", warnings[v].Channel.ID, warnings[v].Channel.ID)

		str = str + fmt.Sprintf(
			"**Author**:\t%s\n"+
				"**Channel**:  %s\n"+
				"**Time**:\t\t%s\n"+
				"**Reason**:   %s\n\n",
			avatar, channel, warnings[v].Time.Format("01/02/06 03:04:05 PM MST"), warnings[v].Reason)
	}

	return str
}

// FormatKicks :
// Returns a string of kicks.
func FormatKicks(kicks map[int64]Kicks) string {

	var keys []int64
	for k := range kicks {
		keys = append(keys, k)
	}

	// Sort the keys by time
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	str := "\n"

	for _, v := range keys {
		avatar := fmt.Sprintf("%s#%s / %s", kicks[v].AuthorUser.Username, kicks[v].AuthorUser.Discriminator, kicks[v].AuthorUser.ID)
		channel := fmt.Sprintf("<#%s> / %s", kicks[v].Channel.ID, kicks[v].Channel.ID)

		str = str + fmt.Sprintf(
			"**Author**:\t%s\n"+
				"**Channel**:  %s\n"+
				"**Time**:\t\t%s\n"+
				"**Reason**:   %s\n\n",
			avatar, channel, kicks[v].Time.Format("01/02/06 03:04:05 PM MST"), kicks[v].Reason)
	}

	return str
}

// FormatBans :
// Returns a string of bans.
func FormatBans(bans map[int64]Bans) string {

	var keys []int64
	for k := range bans {
		keys = append(keys, k)
	}

	// Sort the keys by time
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	str := "\n"

	for _, v := range keys {
		avatar := fmt.Sprintf("%s#%s / %s", bans[v].AuthorUser.Username, bans[v].AuthorUser.Discriminator, bans[v].AuthorUser.ID)
		channel := fmt.Sprintf("<#%s> / %s", bans[v].Channel.ID, bans[v].Channel.ID)

		str = str + fmt.Sprintf(
			"**Author**:\t%s\n"+
				"**Channel**:  %s\n"+
				"**Time**:\t\t%s\n"+
				"**Reason**:   %s\n\n",
			avatar, channel, bans[v].Time.Format("01/02/06 03:04:05 PM MST"), bans[v].Reason)
	}

	return str
}

// FormatUsernames :
// Returns a string of usernames.
func FormatUsernames(usernames map[int64]Usernames) string {

	var keys []int64
	for k := range usernames {
		keys = append(keys, k)
	}

	// Sort the keys by time
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	str := "\n"

	for _, v := range keys {
		str = str + fmt.Sprintf(
			"**Username**:\t%s\n"+
				"**Time**:\t\t%s\n\n",
			usernames[v].Username+"#"+usernames[v].Discriminator, usernames[v].Time.Format("01/02/06 03:04:05 PM MST"))
	}

	return str
}

// FormatNicknames :
// Returns a string of nicknames.
func FormatNicknames(nicknames map[int64]Nicknames) string {

	var keys []int64
	for k := range nicknames {
		keys = append(keys, k)
	}

	// Sort the keys by time
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	str := "\n"
	for _, v := range keys {
		str = str + fmt.Sprintf(
			"**Nickname**:\t%s\n"+
				"**Time**:\t\t%s\n\n",
			nicknames[v].Nickname, nicknames[v].Time.Format("01/02/06 03:04:05 PM MST"))
	}

	return str
}
