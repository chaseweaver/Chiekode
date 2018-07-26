package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	cache "github.com/patrickmn/go-cache"
	"github.com/tkanos/gonfig"
)

// Configuration file contents
type Configuration struct {
	Prefix      string
	OwnerID     string
	BotToken    string
	DatabaseURL string
}

var conf = Configuration{}
var err = gonfig.GetConf("config.json", &conf)
var pool = DialNewPool("tcp", ":6379")
var p = pool.Get()

// Create a cache with a default expiration time of 15 minutes, and which
// purges expired items every 20 minutes
var c = cache.New(15*time.Minute, 20*time.Minute)

func main() {

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + conf.BotToken)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the MessageCreate func as a callback for MessageCreate events.
	dg.AddHandler(MessageCreate)

	// Register the GuildCreate for initializing guild databases
	dg.AddHandler(GuildCreate)

	// Register the GuildDelete for removing guild databases
	dg.AddHandler(GuildDelete)

	// Register the GuildMemberAdd for initializing guild members, welcoming members
	dg.AddHandler(GuildMemberAdd)

	// Register the GuildMemberRemove for saying goodbye to members
	dg.AddHandler(GuildMemberRemove)

	// Register the MessageDelete for logging deleted messages
	dg.AddHandler(MessageDelete)

	// Register the MessageUpdate for logging edited messages
	dg.AddHandler(MessageUpdate)

	// Register the GuildMemberUpdate for tracking username and nickname changes
	dg.AddHandler(GuildMemberUpdate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Nagato is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the redis session.
	p.Close()

	// Cleanly close down the Discord session.
	dg.Close()
}
