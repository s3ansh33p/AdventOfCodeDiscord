package bot

import (
	"dustin-ward/AdventOfCodeBot/data"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"time"

	"github.com/bwmarrin/discordgo"
)

func helloworld(s *discordgo.Session, i *discordgo.InteractionCreate) {
	log.Printf("Info: command \"helloworld\" executed from guildId: %s", i.GuildID)
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Hello from the AoC bot 🙂",
		},
	})
}

func leaderboard(s *discordgo.Session, i *discordgo.InteractionCreate) {
	log.Printf("Info: command \"leaderboard\" executed from guildId: %s", i.GuildID)

	channel, err := getChannel(i.GuildID)
	if err != nil {
		log.Println(fmt.Errorf("Error: leaderboard: %v", err))
		respondWithError(s, i, "Your server has not been correctly configured!")
		return
	}

	D, err := data.GetData(channel.Leaderboard)
	if err != nil {
		log.Println(fmt.Errorf("Error: leaderboard: %v", err))
		respondWithError(s, i, "An internal error occured...")
		return
	}

	M := make([]data.User, 0)
	for _, m := range D.Members {
		M = append(M, m)
	}

	sort.Slice(M, func(i, j int) bool {
		if M[i].Stars == M[j].Stars {
			return M[i].LocalScore > M[j].LocalScore
		}
		return M[i].Stars > M[j].Stars
	})

	fields := make([]*discordgo.MessageEmbedField, 0)
	for _, m := range M {
		f := &discordgo.MessageEmbedField{
			Name:  fmt.Sprintf("**%s**", m.Name),
			Value: fmt.Sprintf("⭐ %d ⏳ %d", m.Stars, m.LocalScore),
		}
		fields = append(fields, f)
	}
	fields[0].Name += " 🥇"
	fields[1].Name += " 🥈"
	fields[2].Name += " 🥉"

	embeds := make([]*discordgo.MessageEmbed, 1)
	embeds[0] = &discordgo.MessageEmbed{
		URL:   "https://adventofcode.com/2022/private/view/" + channel.Leaderboard,
		Type:  discordgo.EmbedTypeRich,
		Title: "🎄 2022 Leaderboard 🎄",
		Color: 0x127C06,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Leaderboard as of " + time.Now().Format("2006/01/02 3:4:5pm"),
		},
		Fields: fields,
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: embeds,
		},
	})
	if err != nil {
		log.Println(fmt.Errorf("Warn: %w", err))
	}
}

func configure(s *discordgo.Session, i *discordgo.InteractionCreate) {
	log.Printf("Info: command \"configure\" executed from guildId: %s", i.GuildID)
	options := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(i.ApplicationCommandData().Options))
	for _, opt := range i.ApplicationCommandData().Options {
		options[opt.Name] = opt
	}

	ch := data.Channel{
		GuildId:      i.GuildID,
		ChannelId:    options["channel"].ChannelValue(nil).ID,
		RoleId:       options["role"].RoleValue(nil, i.GuildID).ID,
		Leaderboard:  options["leaderboard"].StringValue(),
		SessionToken: options["session-token"].StringValue(),
	}

	C[i.GuildID] = ch

	b, err := json.Marshal(C)
	if err != nil {
		log.Println(fmt.Errorf("Error: configure: %v", err))
		respondWithError(s, i, "Error: Invalid arguments were supplied...")
		return
	}

	if err := ioutil.WriteFile("./channels.json", b, 0777); err != nil {
		log.Println(fmt.Errorf("Error: configure: %v", err))
		respondWithError(s, i, "Error: Internal server error...")
		return
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Server successfully configured!",
		},
	})
	if err != nil {
		log.Println(fmt.Errorf("Warn: %w", err))
	}
}
