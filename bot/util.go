package bot

import (
	"dustin-ward/AdventOfCodeBot/data"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

func getChannel(guildId string) (*data.Channel, error) {
	ch, ok := C[guildId]
	if !ok {
		return nil, fmt.Errorf("Error: channel not found")
	}
	return &ch, nil
}

func respondWithError(s *discordgo.Session, i *discordgo.InteractionCreate, message string) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	})

	if err != nil {
		log.Fatal(fmt.Errorf("Fatal: %v", err))
	}
}
