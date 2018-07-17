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
	warningColor = 16776960
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

// DeleteMessageWithTime deletes a message by ID after a given time in milliseconds
func DeleteMessageWithTime(ctx Context, ID string, t float32) {
	time.Sleep(time.Duration(t) * time.Millisecond)

	err := ctx.Session.ChannelMessageDelete(ctx.Channel.ID, ID)

	if err != nil {
		log.Println(err)
	}
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

// ParseMessageContentIDs returns an array of Discord IDs found within a string
func ParseMessageContentIDs(content string) []string {
	re := regexp.MustCompile("([0-9]{18,18})")
	return re.FindAllString(content, -1)
}

// FetchMessageContentUsers returns an array of Discord Users found within a string by ID and Mention with guild restriction
func FetchMessageContentUsers(ctx Context) []*discordgo.User {
	var arr []*discordgo.User
	re := regexp.MustCompile("([0-9]{18,18})")
	for _, v := range re.FindAllString(ctx.Event.Message.Content, -1) {
		mem, err := ctx.Session.User(v)

		if err != nil {
			log.Println(err)
		}

		// Ignore IDs from members NOT part of the guild
		_, err = ctx.Session.GuildMember(ctx.Guild.ID, mem.ID)

		if err == nil {
			arr = append(arr, mem)
		}
	}

	return arr
}

// FetchMessageContentUsersAllGuilds returns an array of Discord Users found within a string by ID and Mention without guild restriction
func FetchMessageContentUsersAllGuilds(ctx Context) []*discordgo.User {
	var arr []*discordgo.User
	re := regexp.MustCompile("([0-9]{18,18})")
	for _, v := range re.FindAllString(ctx.Event.Message.Content, -1) {
		mem, err := ctx.Session.User(v)

		if err != nil {
			log.Println(err)
		}

		arr = append(arr, mem)
	}
	return arr
}

// RemoveMessageIDs removes Mentions and IDs (members, channels, roles, etc) from a string
func RemoveMessageIDs(str string) string {
	tmp := str
	re := regexp.MustCompile("((<!@)[0-9]{18,18}[>])|((<@)[0-9]{18,18}[>])|([0-9]{18,18})|(<@>)|(<!@>)")
	for _, v := range re.FindAllString(str, -1) {
		tmp = strings.Replace(tmp, v, "", -1)
	}
	return strings.TrimSpace(tmp)
}
