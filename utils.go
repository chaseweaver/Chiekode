package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	warningColor = 16383942
	muteColor    = 40447
	unmuteColor  = 4387935
	kickColor    = 54527
	banColor     = 16711684
	deleteColor  = 4378356
	editColor    = 4387980
)

// Embed hold the embed struct
type Embed struct {
	*discordgo.MessageEmbed
}

// Constants for message embed character limits
const (
	EmbedLimitTitle       = 256
	EmbedLimitDescription = 2048
	EmbedLimitFieldValue  = 1024
	EmbedLimitFieldName   = 256
	EmbedLimitField       = 25
	EmbedLimitFooter      = 2048
	EmbedLimit            = 4000
)

/**
 * utils.go
 * Chase Weaver
 *
 * This package handles various utilities for shorthands and logging.
 */

// NewEmbed :
// Rreturns a new embed object
func NewEmbed() *Embed {
	return &Embed{&discordgo.MessageEmbed{}}
}

// SetTitle :
// Sets embed title
// [title]
func (e *Embed) SetTitle(name string) *Embed {
	e.Title = name
	return e
}

// SetDescription :
// Sets embed description
// [description]
func (e *Embed) SetDescription(description string) *Embed {
	if len(description) > 2048 {
		description = description[:2048]
	}
	e.Description = description
	return e
}

// AddField :
// Adds a field embed to array
// [name] [value]
func (e *Embed) AddField(name, value string) *Embed {
	if len(value) > 1024 {
		value = value[:1024]
	}

	if len(name) > 1024 {
		name = name[:1024]
	}

	e.Fields = append(e.Fields, &discordgo.MessageEmbedField{
		Name:  name,
		Value: value,
	})

	return e

}

// SetFooter :
// Sets embed footer
// [iconURL] [text] [proxyURL]
func (e *Embed) SetFooter(args ...string) *Embed {
	iconURL := ""
	text := ""
	proxyURL := ""

	switch {
	case len(args) > 2:
		proxyURL = args[2]
		fallthrough
	case len(args) > 1:
		iconURL = args[1]
		fallthrough
	case len(args) > 0:
		text = args[0]
	case len(args) == 0:
		return e
	}

	e.Footer = &discordgo.MessageEmbedFooter{
		IconURL:      iconURL,
		Text:         text,
		ProxyIconURL: proxyURL,
	}

	return e
}

// SetImage :
// Sets embed image
// [URL] [proxyURL]
func (e *Embed) SetImage(args ...string) *Embed {
	var URL string
	var proxyURL string

	if len(args) == 0 {
		return e
	}
	if len(args) > 0 {
		URL = args[0]
	}
	if len(args) > 1 {
		proxyURL = args[1]
	}
	e.Image = &discordgo.MessageEmbedImage{
		URL:      URL,
		ProxyURL: proxyURL,
	}
	return e
}

// SetThumbnail :
// Sets embed thumbnail
// [URL] [proxyURL]
func (e *Embed) SetThumbnail(args ...string) *Embed {
	var URL string
	var proxyURL string

	if len(args) == 0 {
		return e
	}
	if len(args) > 0 {
		URL = args[0]
	}
	if len(args) > 1 {
		proxyURL = args[1]
	}
	e.Thumbnail = &discordgo.MessageEmbedThumbnail{
		URL:      URL,
		ProxyURL: proxyURL,
	}
	return e
}

// SetAuthor :
// Sets embed author
// [name] [iconURL] [URL] [proxyURL]
func (e *Embed) SetAuthor(args ...string) *Embed {
	var (
		name     string
		iconURL  string
		URL      string
		proxyURL string
	)

	if len(args) == 0 {
		return e
	}
	if len(args) > 0 {
		name = args[0]
	}
	if len(args) > 1 {
		iconURL = args[1]
	}
	if len(args) > 2 {
		URL = args[2]
	}
	if len(args) > 3 {
		proxyURL = args[3]
	}

	e.Author = &discordgo.MessageEmbedAuthor{
		Name:         name,
		IconURL:      iconURL,
		URL:          URL,
		ProxyIconURL: proxyURL,
	}

	return e
}

// SetURL :
// Sets embed URL
// [URL]
func (e *Embed) SetURL(URL string) *Embed {
	e.URL = URL
	return e
}

// SetTimestamp :
// Sets embed timestamp
func (e *Embed) SetTimestamp(time string) *Embed {
	e.Timestamp = time
	return e
}

// SetColor :
// Sets embed color
// [color]
func (e *Embed) SetColor(clr int) *Embed {
	e.Color = clr
	return e
}

// InlineAllFields :
// Sets all fields in the embed to be inline
func (e *Embed) InlineAllFields() *Embed {
	for _, v := range e.Fields {
		v.Inline = true
	}
	return e
}

// Truncate :
// Truncates any embed value over the character limit
func (e *Embed) Truncate() *Embed {
	e.TruncateDescription()
	e.TruncateFields()
	e.TruncateFooter()
	e.TruncateTitle()
	return e
}

// TruncateFields :
// Truncates fields that are too long
func (e *Embed) TruncateFields() *Embed {
	if len(e.Fields) > 25 {
		e.Fields = e.Fields[:EmbedLimitField]
	}

	for _, v := range e.Fields {

		if len(v.Name) > EmbedLimitFieldName {
			v.Name = v.Name[:EmbedLimitFieldName]
		}

		if len(v.Value) > EmbedLimitFieldValue {
			v.Value = v.Value[:EmbedLimitFieldValue]
		}

	}
	return e
}

// TruncateDescription :
// Truncates description
func (e *Embed) TruncateDescription() *Embed {
	if len(e.Description) > EmbedLimitDescription {
		e.Description = e.Description[:EmbedLimitDescription]
	}
	return e
}

// TruncateTitle :
// Truncates title
func (e *Embed) TruncateTitle() *Embed {
	if len(e.Title) > EmbedLimitTitle {
		e.Title = e.Title[:EmbedLimitTitle]
	}
	return e
}

// TruncateFooter :
// Truncates footer
func (e *Embed) TruncateFooter() *Embed {
	if e.Footer != nil && len(e.Footer.Text) > EmbedLimitFooter {
		e.Footer.Text = e.Footer.Text[:EmbedLimitFooter]
	}
	return e
}

// RandomInt :
// Generates a random int between [x,y].
func RandomInt(min, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(max-min) + min
}

// FormatString :
// Adds string formatting (i.e. asciidoc).
func FormatString(s string, t string) string {
	return fmt.Sprintf("```%s\n"+s+"```", t)
}

// FormatWelcomeGoodbyeMessage :
// Replaces string contents with User / Guild Names and IDs.
func FormatWelcomeGoodbyeMessage(g *discordgo.Guild, m *discordgo.Member, s string) string {
	msg := s

	// Username#xxxx
	if strings.Contains(msg, "$MEMBER_NAME$") {
		msg = strings.Replace(msg, "$MEMBER_NAME$", m.User.Username+"#"+m.User.Discriminator, -1)
	}

	// @Member
	if strings.Contains(s, "$MEMBER_MENTION$") {
		msg = strings.Replace(msg, "$MEMBER_MENTION$", "<@"+m.User.ID+">", -1)
	}

	// Member ID
	if strings.Contains(s, "$MEMBER_ID$") {
		msg = strings.Replace(msg, "$MEMBER_ID$", m.User.ID, -1)
	}

	// Member Age
	if strings.Contains(s, "$MEMBER_AGE$") {

		// Fetch creation time of user
		t, err := CreationTime(m.User.ID)

		if err != nil {
			log.Println(err)
		}

		msg = strings.Replace(msg, "$MEMBER_AGE$", t.Format("01/02/06 03:04:05 PM MST"), -1)
	}

	// Member Joined
	if strings.Contains(s, "$MEMBER_JOINED$") {

		// Fetch joined time of user
		t, err := time.Parse(time.RFC3339Nano, m.JoinedAt)

		if err != nil {
			log.Println(err)
		}

		msg = strings.Replace(msg, "$MEMBER_JOINED$", t.Format("01/02/06 03:04:05 PM MST"), -1)
	}

	// Guild Name
	if strings.Contains(s, "$GUILD_NAME$") {
		msg = strings.Replace(msg, "$GUILD_NAME$", g.Name, -1)
	}

	// Guild ID
	if strings.Contains(s, "$GUILD_ID$") {
		msg = strings.Replace(msg, "$GUILD_ID$", g.ID, -1)
	}

	return msg
}

// TrimSuffix :
// Removes a string from the end of another string.
func TrimSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		s = s[:len(s)-len(suffix)]
	}
	return s
}

// DeleteMessageWithTime :
// Deletes a message by ID after a given time in milliseconds.
func DeleteMessageWithTime(ctx Context, ID string, t float32) {
	time.Sleep(time.Duration(t) * time.Millisecond)

	err := ctx.Session.ChannelMessageDelete(ctx.Channel.ID, ID)

	if err != nil {
		log.Println(err)
	}
}

// Wait :
// Delays execution base on time in milliseconds.
func Wait(t float32) {
	time.Sleep(time.Duration(t) * time.Millisecond)
}

// Round :
// Takes an input and round it to the nearest unit number.
func Round(x, unit float64) float64 {
	return math.Round(x/unit) * unit
}

// LogCommands :
// Logs commands being run.
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

// SliceExists :
// Checks if an element exists within a slice
func SliceExists(slice interface{}, item interface{}) bool {
	s := reflect.ValueOf(slice)

	if s.Kind() != reflect.Slice {
		panic("SliceExists() given a non-slice type")
	}

	for i := 0; i < s.Len(); i++ {
		if s.Index(i).Interface() == item {
			return true
		}
	}

	return false
}

// Contains :
// Checks if element is in array.
func Contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

// CreationTime :
// Returns the time a snowflake was created.
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
// Returns an array of Discord Users found within a string by ID / Name / Mention (guild restriction).
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
// Returns an array of Discord Users found within a string by ID / Name / Mention (guild restriction).
// Returns the array of users removed from the original message string.
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
	rm := regexp.MustCompile("((<@)[0-9]{18,18}[>])|((<@!)[0-9]{18,18}[>])|(<@!>)|(<@>)")
	for _, v := range rm.FindAllString(msg, -1) {
		msg = strings.Replace(msg, v, "", -1)
	}

	return arr, msg
}

// FetchMessageContentUsersAllGuilds :
// Returns an array of Discord Users found within a string by ID without guild restriction.
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
// Returns an array of Discord Channels found within a string by ID and Mention with guild restriction.
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
// Returns an array of Discord Roles found within a string by ID, Mention, and Name with guild restriction.
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
// Returns  a joint of all Users, Channels, and Roles, and returns and removes any occurances found in the message.
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

// SetTimeout :
// Delays execution of a func (non-blocking) for a given time
func SetTimeout(f func(), milliseconds int) {
	timeout := time.Duration(milliseconds) * time.Millisecond
	time.AfterFunc(timeout, f)
}

// UnmuteMember :
// Removes the "mute" role from a user
func UnmuteMember(ctx Context, memberID string) {

	// Fetch Guild information from redis database
	g, err := UnpackGuildStruct(ctx.Guild.ID)
	if err != nil {
		log.Println(err)
		return
	}

	if g.MutedRole == nil {
		return
	}

	// Find "mute" role in database and remove it from the user
	role := g.MutedRole
	err = ctx.Session.GuildMemberRoleRemove(ctx.Guild.ID, memberID, role.ID)

	if err != nil {
		log.Println(err)
		return
	}
}

// MakeTimestamp :
// Creates a unix timestamp
func MakeTimestamp() int64 {
	return time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}
