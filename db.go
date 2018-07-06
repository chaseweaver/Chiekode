package main

import (
	"log"
	"os"

	"github.com/gomodule/redigo/redis"
)

// GuildConfig handles configuration information per guild
type GuildConfig struct {
	GuildPrefix string
	GuildID     int64

	ValidChannels []int64

	BlacklistedMembers []int64

	Warnings []int64
	Mutes    []int64
	Kicks    []int64
	Bans     []int64

	WelcomeMessage string
	WelcomeChannel int64

	GoodbyeMessage string
	GoodbyeChannel int64

	Events []string

	DisabledFunctions []string

	BirthdayRole int64
	MutedRole    int64
	AutoRole     int64
}

// DialRedisDatabaseLocal connectes to a local Redis database by port pass-in
func DialRedisDatabaseLocal(net string, add string) error {
	c, err := redis.Dial(net, add)
	if err != nil {
		log.Println(err)
		return err
	}
	defer c.Close()
	return nil
}

// DialRedisDatabaseURL connectes to a Redis database by URL pass-in
func DialRedisDatabaseURL(url string) error {
	c, err := redis.DialURL(os.Getenv(url))
	if err != nil {
		log.Println(err)
		return err
	}
	defer c.Close()
	return nil
}
