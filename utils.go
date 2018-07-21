package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	warningColor = 16383942
	muteColor    = 40447
	kickColor    = 54527
	banColor     = 16711684
)

/**
 * utils.go
 * Chase Weaver
 *
 * This package handles various utilities for shorthands and logging.
 */

// RandomInt generates a random int between [x,y]
func RandomInt(min, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(max-min) + min
}

// FormatString adds string formatting (i.e. asciidoc)
func FormatString(s string, t string) string {
	return fmt.Sprintf("```%s\n"+s+"```", t)
}

// TrimSuffix removes a string from the end of another string
func TrimSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		s = s[:len(s)-len(suffix)]
	}
	return s
}

// DeleteMessageWithTime deletes a message by ID after a given time in milliseconds
func DeleteMessageWithTime(ctx Context, ID string, t float32) {
	time.Sleep(time.Duration(t) * time.Millisecond)

	err := ctx.Session.ChannelMessageDelete(ctx.Channel.ID, ID)

	if err != nil {
		log.Println(err)
	}
}

// Wait will delay execution base on time in milliseconds
func Wait(t float32) {
	time.Sleep(time.Duration(t) * time.Millisecond)
}

// Round will take an input and round it to the nearest unit number
func Round(x, unit float64) float64 {
	return math.Round(x/unit) * unit
}

// LogCommands logs commands being run
func LogCommands(ctx Context) {
	if ctx.Channel.Type == discordgo.ChannelTypeGuildText {
		log.Printf(
			"\n"+
				"User:      %s / %s\n"+
				"Guild:     %s / %s\n"+
				"Channel:   %s / %s\n"+
				"Command:   %s\n"+
				"Args:      (%d)%s"+
				"\n\n",
			ctx.Event.Author.Username+"#"+ctx.Event.Author.Discriminator, ctx.Event.Author.ID,
			ctx.Guild.Name, ctx.Guild.ID, ctx.Channel.Name, ctx.Channel.ID,
			ctx.Name, len(ctx.Args), ctx.Args)
	} else {
		log.Printf(
			"\n"+
				"User:      %s / %s\n"+
				"DM:        %s\n"+
				"Command:   %s\n"+
				"Args:      %s"+
				"\n\n",
			ctx.Event.Author.Username+"#"+ctx.Event.Author.Discriminator,
			ctx.Event.Author.ID, ctx.Channel.ID, ctx.Name, ctx.Args)
	}
}

// Contains checks if element is in array
func Contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

// CreationTime returns the time a snowflake was created
func CreationTime(ID string) (t time.Time, err error) {
	i, err := strconv.ParseInt(ID, 10, 64)
	if err != nil {
		return
	}
	timestamp := (i >> 22) + 1420070400000
	t = time.Unix(timestamp/1000, 0)
	return
}

// FetchMessageContentUsers :
// Returns an array of Discord Users found within a string by ID / Name / Mention (guild restriction)
func FetchMessageContentUsers(ctx Context, msg string) []*discordgo.User {
	var arr []*discordgo.User
	re := regexp.MustCompile("([0-9]{18,18})")

	for _, v := range re.FindAllString(msg, -1) {
		usr, err := ctx.Session.User(v)

		if err != nil {
			break
		}
		_, err = ctx.Session.GuildMember(ctx.Guild.ID, usr.ID)

		if err != nil {
			break
		}

		arr = append(arr, usr)
	}

	g, err := ctx.Session.State.Guild(ctx.Guild.ID)

	if err != nil {
		return arr
	}

	// Find users by username
	for _, m := range g.Members {
		if strings.Contains(ctx.Event.Message.Content, m.User.Username+"#"+m.User.Discriminator) {
			arr = append(arr, m.User)
		}
	}

	return arr
}

// FetchMessageContentUsersString :
// Returns an array of Discord Users found within a string by ID / Name / Mention (guild restriction)
// Returns the array of users removed from the original message string
func FetchMessageContentUsersString(ctx Context, str string) ([]*discordgo.User, string) {

	var arr []*discordgo.User
	re := regexp.MustCompile("([0-9]{18,18})")
	msg := str

	// Add members by ID from regexp, removes IDs from message string
	for _, v := range re.FindAllString(msg, -1) {
		usr, err := ctx.Session.User(v)

		if err != nil {
			break
		}

		_, err = ctx.Session.GuildMember(ctx.Guild.ID, usr.ID)

		if err != nil {
			break
		}

		msg = strings.Replace(msg, v, "", -1)
		arr = append(arr, usr)
	}

	g, err := ctx.Session.State.Guild(ctx.Guild.ID)

	if err != nil {
		return nil, msg
	}

	// Find users by username, removes Name#xxxx from message string
	for _, m := range g.Members {
		if strings.Contains(ctx.Event.Message.Content, m.User.Username+"#"+m.User.Discriminator) {
			arr = append(arr, m.User)
			msg = strings.Replace(msg, m.User.Username+"#"+m.User.Discriminator, "", -1)
		}
	}

	// Remove regex from message
	rm := regexp.MustCompile("((<@)[0-9]{18,18}[>])|((<@!)[0-9]{18,18}[>])")
	for _, v := range rm.FindAllString(msg, -1) {
		msg = strings.Replace(msg, v, "", -1)
	}

	return arr, msg
}

// FetchMessageContentUsersAllGuilds :
// Returns an array of Discord Users found within a string by ID without guild restriction
func FetchMessageContentUsersAllGuilds(ctx Context, msg string) []*discordgo.User {
	var arr []*discordgo.User
	re := regexp.MustCompile("([0-9]{18,18})")

	for _, v := range re.FindAllString(msg, -1) {
		usr, err := ctx.Session.User(v)

		if err != nil {
			break
		}

		arr = append(arr, usr)
	}
	return arr
}

// FetchMessageContentChannels :
// Returns an array of Discord Channels found within a string by ID and Mention with guild restriction
func FetchMessageContentChannels(ctx Context, msg string) []*discordgo.Channel {
	var arr []*discordgo.Channel
	re := regexp.MustCompile("([0-9]{18,18})")
	channel, err := ctx.Session.GuildChannels(ctx.Guild.ID)

	if err != nil {
		return arr
	}

	// Add channels by ID/Mention
	for _, v := range re.FindAllString(msg, -1) {
		for _, c := range channel {
			if v == c.ID {
				arr = append(arr, c)
			}
		}
	}

	// Add channel by name, case sensitive
	for _, c := range channel {
		if strings.Contains(msg, c.Name) {
			arr = append(arr, c)
		}
	}

	return arr
}

// FetchMessageContentRoles :
// Returns an array of Discord Roles found within a string by ID, Mention, and Name with guild restriction
func FetchMessageContentRoles(ctx Context, msg string) []*discordgo.Role {
	var arr []*discordgo.Role
	re := regexp.MustCompile("([0-9]{18,18})")
	role, err := ctx.Session.GuildRoles(ctx.Guild.ID)

	if err != nil {
		return arr
	}

	// Add role by ID/Mention
	for _, v := range re.FindAllString(msg, -1) {
		for _, r := range role {
			if v == r.ID {
				arr = append(arr, r)
			}
		}
	}

	// Add role by name, case sensitive
	for _, r := range role {
		if strings.Contains(msg, r.Name) {
			arr = append(arr, r)
		}
	}

	return arr
}

// FetchUsersChannelsRoles :
// Returns  a joint of all Users, Channels, and Roles, and returns and removes any occurances found in the message
func FetchUsersChannelsRoles(ctx Context, msg string) ([]*discordgo.User, []*discordgo.Channel, []*discordgo.Role, string) {
	str := msg
	users := FetchMessageContentUsers(ctx, msg)
	channels := FetchMessageContentChannels(ctx, msg)
	roles := FetchMessageContentRoles(ctx, msg)
	re := regexp.MustCompile("((<#)[0-9]{18,18}[>])|((<@&)[0-9]{18,18}[>])|((<@!)[0-9]{18,18}[>])|((<@)[0-9]{18,18}[>])|([0-9]{18,18})|(<@>)|(<!@>)")

	// Remove Regex from message
	for _, v := range re.FindAllString(msg, -1) {
		str = strings.Replace(str, v, "", -1)
	}

	// Fetch Users by name, remove them from message
	for _, u := range users {
		if strings.Contains(msg, u.Username+"#"+u.Discriminator) {
			str = strings.Replace(str, u.Username+"#"+u.Discriminator, "", -1)
		}
	}

	// Fetch Roles by name, remove them from message
	for _, r := range roles {
		if strings.Contains(msg, r.Name) {
			str = strings.Replace(str, r.Name, "", -1)
		}
	}

	// Fetch Channels by name, remove them from message
	for _, c := range channels {
		if strings.Contains(msg, c.Name) {
			str = strings.Replace(str, c.Name, "", -1)
		}
	}

	return users, channels, roles, strings.TrimSpace(str)
}
