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
var c []data.Channel

var commands = []*discordgo.ApplicationCommand{
	{
		Name:        "hello-world",
		Description: "For testing purposes",
	},
	{
		Name:        "leaderboard",
		Description: "Current Leaderboard",
	},
}

var commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"leaderboard": leaderboard,
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

	if err = json.Unmarshal(b, &c); err != nil {
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
