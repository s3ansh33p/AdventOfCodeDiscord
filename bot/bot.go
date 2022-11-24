package bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"dustin-ward/AdventOfCodeBot/data"

	"github.com/bwmarrin/discordgo"
)

var s *discordgo.Session
var C map[string]data.Channel
var adminPerm int64 = 0

var commands = []*discordgo.ApplicationCommand{
	{
		Name:                     "hello-world",
		Description:              "For testing purposes",
		DefaultMemberPermissions: &adminPerm,
	},
	{
		Name:        "leaderboard",
		Description: "Current Leaderboard",
	},
	{
		Name:                     "configure-server",
		Description:              "Configure the AdventOfCode bot to work with your leaderboard and server",
		DefaultMemberPermissions: &adminPerm,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionChannel,
				Name:        "channel",
				Description: "Channel to post in",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionRole,
				Name:        "role",
				Description: "Advent of Code role to be pinged",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "leaderboard",
				Description: "Id for your private Advent of Code leaderboard",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "session-token",
				Description: "Valid session token of one member of your private leaderboard",
				Required:    true,
			},
		},
	},
}

var commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"leaderboard":      leaderboard,
	"hello-world":      helloworld,
	"configure-server": configure,
}

func InitSession() (*discordgo.Session, error) {
	S, err := discordgo.New("Bot " + BotToken)
	if err != nil {
		return nil, fmt.Errorf("invalid bot configuration: %w", err)
	}

	S.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		} else {
			log.Printf("Warning: command handler not found: \"%s\"\n", i.ApplicationCommandData().Name)
		}
	})

	b, err := ioutil.ReadFile("./channels.json")
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(b, &C); err != nil {
		return nil, err
	}

	s = S
	return S, nil
}

func RegisterCommands() ([]*discordgo.ApplicationCommand, error) {
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, "", v)
		if err != nil {
			return nil, fmt.Errorf("Cannot create '%s' command: %w", v.Name, err)
		}
		registeredCommands[i] = cmd
	}
	return registeredCommands, nil
}
