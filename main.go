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
	// resBody, err := data.GetData(data.Jenna2021)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// if err := json.Unmarshal(resBody, &D); err != nil {
	// 	log.Fatal(err)
	// }

	// Leaderboard ID from args
	if len(os.Args) != 2 {
		log.Fatal("Fatal: invalid number of arguments")
	}
	boardId := os.Args[1]

	// Initialize Discord Session
	Session, err := bot.InitSession()
	if err != nil {
		log.Fatal(fmt.Errorf("Fatal: %v", err))
	}

	log.Println("Session initialized")

	// Add init handler
	Session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	// Open connection
	if err = Session.Open(); err != nil {
		log.Fatal(fmt.Errorf("Fatal: %v", err))
	}
	defer Session.Close()

	// Register commands
	r, err := bot.RegisterCommands()
	if err != nil {
		log.Fatal(fmt.Errorf("Fatal: %v", err))
	}
	for _, c := range r {
		log.Printf("Command registered: \"%s\" with id: %v", c.Name, c.ID)
	}

	// Continually fetch advent of code data every 15 minutes
	go func() {
		for {
			log.Println("Attempting to fetch data for leaderboard " + boardId + "...")
			if err := data.FetchData(boardId, boardId); err != nil {
				log.Fatal(fmt.Errorf("Fatal: %v", err))
			}
			log.Println("Success!")

			time.Sleep(REQUEST_RATE)
		}
	}()

	// Wait for SIGINT to end program
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop
}
