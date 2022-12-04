package commands

import (
	"github.com/bwmarrin/discordgo"
	"strconv"
)

func Ping(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "🏓 Pong! My ping is " + strconv.Itoa(int(s.HeartbeatLatency().Milliseconds())) + "ms.",
		},
	})
}