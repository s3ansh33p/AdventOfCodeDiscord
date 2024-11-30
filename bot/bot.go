package bot

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"dustin-ward/AdventOfCodeBot/data"

	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron"
)

var s *discordgo.Session
var C map[string]*data.Channel
var crn *cron.Cron
var adminPerm int64 = 0
var dataDir string = os.Getenv("DATA_DIR")

// Command definitions
var commands = []*discordgo.ApplicationCommand{
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
	{
		Name:                     "start-notifications",
		Description:              "Start the notification process",
		DefaultMemberPermissions: &adminPerm,
	},
	{
		Name:                     "stop-notifications",
		Description:              "Stop the notification process",
		DefaultMemberPermissions: &adminPerm,
	},
	{
		Name:                     "check-notifications",
		Description:              "Check to see if notifications are currently enabled",
		DefaultMemberPermissions: &adminPerm,
	},
}

// Command handlefuncs
var commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"leaderboard":         leaderboard,
	"configure-server":    configure,
	"start-notifications": startCountdown,
	"stop-notifications":  stopCountdown,
	"check-notifications": checkCountdown,
}

func InitSession() (*discordgo.Session, error) {
	// Get token
	token := os.Getenv("AOC_BOT_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("no discord token found. Please set $AOC_BOT_TOKEN")
	}

	// Init discordgo session
	S, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, fmt.Errorf("invalid bot configuration: %w", err)
	}

	// Attach handlers to functions
	S.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		} else {
			log.Printf("Warning: command handler not found: \"%s\"\n", i.ApplicationCommandData().Name)
		}
	})

	// Check for config file
	if _, err := os.Stat(dataDir + "channels.json"); errors.Is(err, os.ErrNotExist) {
		log.Println("Info: no channel config file found")

		C = make(map[string]*data.Channel, 3)
	} else {
		// Read channel configs from file (Not an ideal storage method...)
		b, err := os.ReadFile(dataDir + "channels.json")
		if err != nil {
			return nil, err
		}

		// Populate channel info in local memory
		if err = json.Unmarshal(b, &C); err != nil {
			return nil, err
		}
	}

	s = S
	return S, nil
}

func TakeDown() error {
	log.Println("Shutting Down...")
	crn.Stop()

	// Save channel configurations
	b, err := json.Marshal(C)
	if err != nil {
		return err
	}
	if err = os.WriteFile(dataDir+"channels.json", b, 0777); err != nil {
		return err
	}

	return nil
}

func RegisterCommands() ([]*discordgo.ApplicationCommand, error) {
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, "", v)
		if err != nil {
			return nil, fmt.Errorf("cannot create '%s' command: %w", v.Name, err)
		}
		registeredCommands[i] = cmd
	}
	return registeredCommands, nil
}

func SetupNotifications() error {
	// Cronjob for 11:30pm EST
	TO, err := time.LoadLocation("America/Toronto")
	if err != nil {
		return err
	}
	crn = cron.NewWithLocation(TO)
	if err := crn.AddFunc("0 30 23 * * *", problemNotification); err != nil {
		return err
	}
	crn.Start()
	return nil
}
