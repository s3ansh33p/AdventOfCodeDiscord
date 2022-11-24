package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"dustin-ward/AdventOfCodeBot/bot"
	"dustin-ward/AdventOfCodeBot/data"

	"github.com/bwmarrin/discordgo"
)

const (
	REQUEST_RATE = time.Minute * 15
)

func main() {
	// Initialize Discord Session
	Session, err := bot.InitSession()
	if err != nil {
		log.Fatal(fmt.Errorf("Fatal: %w", err))
	}

	log.Println("Session initialized")

	// Add init handler
	Session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	// Open connection
	if err = Session.Open(); err != nil {
		log.Fatal(fmt.Errorf("Fatal: %w", err))
	}
	defer Session.Close()

	// Register commands
	r, err := bot.RegisterCommands()
	if err != nil {
		log.Fatal(fmt.Errorf("Fatal: %w", err))
	}
	for _, c := range r {
		log.Printf("Command registered: \"%s\" with id: %v", c.Name, c.ID)
	}

	// Continually fetch advent of code data every 15 minutes
	for _, ch := range bot.C {
		go func() {
			for {
				log.Println("Attempting to fetch data for leaderboard " + ch.Leaderboard + "...")
				if err := data.FetchData(ch.Leaderboard, ch.SessionToken, ch.Leaderboard); err != nil {
					log.Println(fmt.Errorf("Error: %w", err))
				} else {
					log.Println("Success!")
				}

				time.Sleep(REQUEST_RATE)
			}
		}()
	}

	// Wait for SIGINT to end program
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop
}
