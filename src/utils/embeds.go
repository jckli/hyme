package utils

import (
	"github.com/bwmarrin/discordgo"
)

func ErrorEmbed(description string) []*discordgo.MessageEmbed {
	embeds := []*discordgo.MessageEmbed{
		{
			Type: "rich",
			Color: 0xff4f4f,
			Description: description,
		},
	}
	return embeds
}

func SuccessEmbed(description string) []*discordgo.MessageEmbed {
	embeds := []*discordgo.MessageEmbed{
		{
			Type: "rich",
			Color: 0xa4849a,
			Description: description,
		},
	}
	return embeds
}