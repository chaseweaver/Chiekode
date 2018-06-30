package main

import (
	"errors"
	"log"
	"reflect"

	"github.com/bwmarrin/discordgo"
)

/**
 * handler.go
 * Chase Weaver
 *
 * This package handles commands/events be initializing them and their
 * properties.
 */

// Command struct per command
type Command struct {
	Name            string
	Func            func(*discordgo.Session, *discordgo.MessageCreate, []string)
	Enabled         bool
	NSFWOnly        bool
	IgnoreSelf      bool
	IgnoreBots      bool
	RunIn           []string
	Aliases         []string
	BotPermissions  []string
	UserPermissions []string
	ArgsDelim       string
	Usage           string
	Description     string
}

var commands = make(map[string]Command)

// CheckValidPrereq checks if command is valid to be ran
func CheckValidPrereq(s *discordgo.Session, m *discordgo.MessageCreate, c Command) bool {
	if !c.Enabled {
		return false
	}

	nc, err := s.Channel(m.ChannelID)
	if err != nil {
		log.Print(err)
	}

	if !nc.NSFW && c.NSFWOnly {
		return false
	}

	if m.Author.ID == s.State.User.ID && c.IgnoreSelf {
		return false
	}

	if m.Author.Bot && m.Author.ID != s.State.User.ID && c.IgnoreBots {
		return false
	}

	return true
}

// RegisterNewCommand creates a new command
func RegisterNewCommand(k string, c Command) {
	commands[k] = c
}

// RemoveCommand deletes a command
func RemoveCommand(k string) {
	delete(commands, k)
}

// HasCommand checks if a command is already mapped
func HasCommand(k string) bool {
	_, ok := commands[k]
	if ok == true {
		return true
	}

	for key := range commands {
		for val := range commands[key].Aliases {
			if commands[key].Aliases[val] == k {
				return true
			}
		}
	}

	return false
}

// FetchCommand returns a valid command
func FetchCommand(k string) Command {
	if HasCommand(commands[k].Name) {
		return commands[k]
	}

	for key := range commands {
		for val := range commands[key].Aliases {
			if commands[key].Aliases[val] == k {
				return commands[key]
			}
		}
	}

	log.Println(commands[k])

	return Command{}
}

// FetchCommandName returns a valid command
func FetchCommandName(k string) string {
	if HasCommand(commands[k].Name) {
		return commands[k].Name
	}

	for key := range commands {
		for val := range commands[key].Aliases {
			if commands[key].Aliases[val] == k {
				return commands[key].Name
			}
		}
	}

	return k
}

// Call func based on name and passes Session, MessageCreate, ...args
func Call(m map[string]interface{}, name string, params ...interface{}) (result []reflect.Value, err error) {
	f := reflect.ValueOf(m[name])
	if len(params) != f.Type().NumIn() {
		err = errors.New("the number of params is not adapted")
		return
	}
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	result = f.Call(in)
	return
}
