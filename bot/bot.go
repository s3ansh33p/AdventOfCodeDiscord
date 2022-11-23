package bot

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

var s *discordgo.Session

var commands = []*discordgo.ApplicationCommand{
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
		return nil, fmt.Errorf("invalid bot configuration: %v", err)
	}

	S.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		} else {
			log.Printf("Warning: command handler not found: \"%s\"\n", i.ApplicationCommandData().Name)
		}
	})

	s = S
	return S, nil
}

func RegisterCommands() ([]*discordgo.ApplicationCommand, error) {
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, "", v)
		if err != nil {
			return nil, fmt.Errorf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}
	return registeredCommands, nil
}
