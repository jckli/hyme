package commands

import (
	"github.com/bwmarrin/discordgo"
	"strconv"
	"github.com/jckli/hyme/src/music"
	"github.com/TopiSenpai/dgo-paginator"
)

func Ping(s *discordgo.Session, i *discordgo.InteractionCreate, bot *music.Bot, manager *paginator.Manager) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "🏓 Pong! My ping is " + strconv.Itoa(int(s.HeartbeatLatency().Milliseconds())) + "ms.",
		},
	})
}