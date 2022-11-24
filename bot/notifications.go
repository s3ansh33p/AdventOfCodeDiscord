package bot

import (
	"fmt"
	"log"
	"time"
)

const (
	ProblemUrl = "https://adventofcode.com/2022/day/"
)

func problemNotification() {
	day := time.Now().AddDate(0, 0, 1).Day() + 1
	for _, ch := range C {
		log.Println("Info: sending day", day, "notification in channel", ch.ChannelId)
		messageString := fmt.Sprintf(
			"🎄 <@&%s> 🎄\nThe problem for Day %d will be released soon! (<t:%d:R>)\nYou can see the problem statement here when its up: https://adventofcode.com/2022/day/%d",
			ch.RoleId,
			day,
			time.Now().Unix()+(int64(1)*60),
			day,
		)
		_, err := s.ChannelMessageSend(ch.ChannelId, messageString)
		if err != nil {
			log.Println(fmt.Errorf("Error: unable to send notification: %w", err))
		}
	}
	return
}
