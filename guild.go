package main

func init() {
	RegisterNewCommand(Command{
		Name:            "settings",
		Func:            Settings,
		Enabled:         true,
		NSFWOnly:        false,
		IgnoreSelf:      true,
		IgnoreBots:      true,
		RunIn:           []string{"Text"},
		Aliases:         []string{},
		UserPermissions: []string{"Bot Owner"},
		ArgsDelim:       " ",
		ArgsUsage:       "[settings name]",
		Description:     "Lists guild configurations.",
	})

	RegisterNewCommand(Command{
		Name:            "set",
		Func:            Set,
		Enabled:         true,
		NSFWOnly:        false,
		IgnoreSelf:      true,
		IgnoreBots:      true,
		RunIn:           []string{"Text"},
		Aliases:         []string{},
		UserPermissions: []string{"Bot Owner"},
		ArgsDelim:       " ",
		ArgsUsage:       "<settings name> <value>",
		Description:     "Sets guild configurations.",
	})
}

// Settings lists database guild configurations
func Settings(ctx Context) {
	/*
		data, err := redis.Bytes(p.Do("GET", ctx.Guild.ID))
		if err != nil {
			panic(err.Error())
		}

		var g Guild
		err = json.Unmarshal(data, &g)

		key := strings.ToUpper(strings.Join(ctx.Args, " "))
		if len(key) == 0 {
			str := fmt.Sprintf(
				"== %s Configuration ==\n\n"+
					"Guild Prefix          ::   %s\n"+
					"Blacklisted Channel   ::   %d\n"+
					"Blacklisted Members   ::   %d\n"+
					"Welcome Message       ::   %s\n"+
					"Welcome Channel       ::   %s\n"+
					"Goodbye Message       ::   %s\n"+
					"Goodbye Channel       ::   %s\n"+
					"Events                ::   %d\n"+
					"Disabled Commands     ::   %d\n"+
					"Birthday Role         ::   %s\n"+
					"Muted Role            ::   %s\n"+
					"Auto Roles            ::   %d",
				g.Guild.Name, g.GuildPrefix, len(g.BlacklistedChannels), len(g.BlacklistedMembers),
				g.WelcomeMessage, g.WelcomeChannel, g.GoodbyeMessage, g.GoodbyeChannel, len(g.Events),
				len(g.DisabledCommands), g.BirthdayRole, g.MutedRole, len(g.AutoRole))

			ctx.Session.ChannelMessageSend(ctx.Channel.ID, FormatString(str, "asciidoc"))
			return
		}

		var set, val string
		switch key {
		case "PREFIX":
			fallthrough
		case "GUILD PREFIX":
			val = "GUILD PREFIX"
			set = g.GuildPrefix
		case "BLACKLISTED CHANNEL":
			fallthrough
		case "BLACKLISTED CHANNELS":
			var str string
			for _, v := range g.BlacklistedChannels {
				str += v.Name + ", "
			}

			str = TrimSuffix(str, ", ")

			val = "BLACKLISTED CHANNELS"
			set = str
		case "BLACKLISTED MEMBER":
			fallthrough
		case "BLACKLISTED MEMBERS":
			var str string
			for _, v := range g.BlacklistedMembers {
				str += v.Username + ", "
			}

			str = TrimSuffix(str, ", ")

			val = "BLACKLISTED MEMBERS"
			set = str
		case "WELCOME MESSAGE":
			val = "WELCOME MESSAGE"
			set = g.WelcomeMessage
		case "WELCOME CHANNEL":
			val = "WELCOME CHANNEL"
			set = g.WelcomeChannel
		case "GOODBYE MESSAGE":
			val = "GOODBYE MESSAGE"
			set = g.GoodbyeMessage
		case "GOODBYE CHANNEL":
			val = "GOODBYE CHANNEL"
			set = g.GoodbyeChannel
		case "EVENTS":
			val = "EVENTS"
			set = strings.Join(g.Events, ", ")
		case "DISABLED":
			fallthrough
		case "DISABLED COMMANDS":
			val = "DISABLED COMMANDS"
			set = strings.Join(g.DisabledCommands, ", ")
		case "BIRTHDAY":
			fallthrough
		case "BIRTHDAY ROLE":
			val = "BIRTHDAY ROLE"
			set = g.BirthdayRole
		case "MUTED":
			fallthrough
		case "MUTED ROLE":
			val = "MUTED ROLE"
			set = g.MutedRole
		case "AUTO":
			fallthrough
		case "AUTO ROLE":
			fallthrough
		case "AUTO ROLES":
			val = "AUTO ROLES"
			set = strings.Join(g.AutoRole, ", ")
		default:
			ctx.Session.ChannelMessageSend(ctx.Channel.ID, fmt.Sprintf("I could not find the guild setting `%s`", key))
			return
		}

		if len(set) == 0 {
			ctx.Session.ChannelMessageSend(ctx.Channel.ID, fmt.Sprintf("`%s` does not have a set value.", val))
		} else {
			ctx.Session.ChannelMessageSend(ctx.Channel.ID, fmt.Sprintf("`%s` has the value of: %s", val, set))
		}
	*/
}

// Set allows configration of database guild settings
func Set(ctx Context) {
	/*
		data, err := redis.Bytes(p.Do("GET", ctx.Guild.ID))
		if err != nil {
			log.Println(err)
		}

		var g Guild
		err = json.Unmarshal(data, &g)

		if err != nil {
			log.Println(err)
		}

		key := strings.ToUpper(strings.Join(ctx.Args, " "))
		if len(key) == 0 {
			ctx.Session.ChannelMessageSend(ctx.Channel.ID, fmt.Sprintf("Invalid settings key. Run `%shelp` for more information.", g.GuildPrefix))
			return
		}

		var set, val string
		switch key {
		case "PREFIX":
			fallthrough
		case "GUILD PREFIX":
			val = "GUILD PREFIX"
			set = g.GuildPrefix
		case "BLACKLISTED CHANNEL":
			fallthrough
		case "BLACKLISTED CHANNELS":
			val = "BLACKLISTED CHANNELS"
			set = strings.Join(g.BlacklistedChannels, ", ")
		case "BLACKLISTED MEMBER":
			fallthrough
		case "BLACKLISTED MEMBERS":
			val = "BLACKLISTED MEMBERS"
			set = strings.Join(g.BlacklistedMembers, ", ")
		case "WELCOME MESSAGE":
			val = "WELCOME MESSAGE"
			set = g.WelcomeMessage
		case "WELCOME CHANNEL":
			val = "WELCOME CHANNEL"
			set = g.WelcomeChannel
		case "GOODBYE MESSAGE":
			val = "GOODBYE MESSAGE"
			set = g.GoodbyeMessage
		case "GOODBYE CHANNEL":
			val = "GOODBYE CHANNEL"
			set = g.GoodbyeChannel
		case "EVENTS":
			val = "EVENTS"
			set = strings.Join(g.Events, ", ")
		case "DISABLED":
			fallthrough
		case "DISABLED COMMANDS":
			val = "DISABLED COMMANDS"
			set = strings.Join(g.DisabledCommands, ", ")
		case "BIRTHDAY":
			fallthrough
		case "BIRTHDAY ROLE":
			val = "BIRTHDAY ROLE"
			set = g.BirthdayRole
		case "MUTED":
			fallthrough
		case "MUTED ROLE":
			val = "MUTED ROLE"
			set = g.MutedRole
		case "AUTO":
			fallthrough
		case "AUTO ROLE":
			fallthrough
		case "AUTO ROLES":
			val = "AUTO ROLES"
			set = strings.Join(g.AutoRole, ", ")
		default:
			ctx.Session.ChannelMessageSend(ctx.Channel.ID, fmt.Sprintf("I could not find the guild setting `%s`", key))
			return
		}

		serialized, err := json.Marshal(g)

		if err != nil {
			log.Println(err)
		}

		_, err = p.Do("SET", ctx.Guild.ID, serialized)
		if err != nil {
			log.Println(err)
			return
		}*/

}
