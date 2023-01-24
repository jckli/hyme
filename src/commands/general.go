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

func Info(s *discordgo.Session, i *discordgo.InteractionCreate, bot *music.Bot, manager *paginator.Manager) {
	serverCount := len(s.State.Guilds)
	userCount := 0
	for _, guild := range s.State.Guilds {
		userCount += guild.MemberCount
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
					{
						Type: "rich",
						Color: 0xa4849a,
						Title: "Hyme",
						Description: `
						ohashi's music bot, written in go
						supports youtube, spotify, soundcloud, and bandcamp
						
						server count: ` + strconv.Itoa(serverCount) + `
						user count: ` + strconv.Itoa(userCount) + `
						`,
						Author: &discordgo.MessageEmbedAuthor{
							Name: "Hyme",
							IconURL: s.State.User.AvatarURL(""),
						},
					},
				},
			},
		},
	)
}