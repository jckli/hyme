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

func MainEmbed(title string, description string, footer string, s *discordgo.Session, i *discordgo.InteractionCreate) *discordgo.MessageEmbed {
	guild, _ := s.Guild(i.GuildID)
	if footer == "" {
		embeds := discordgo.MessageEmbed{
			Type: "rich",
			Color: 0xa4849a,
			Title: title,
			Description: description,
			Author: &discordgo.MessageEmbedAuthor{
				Name: guild.Name,
				IconURL: guild.IconURL(),
			},
		}
		return &embeds
	}
	embeds := discordgo.MessageEmbed{
		Type: "rich",
		Color: 0xa4849a,
		Title: title,
		Description: description,
		Footer: &discordgo.MessageEmbedFooter{
			Text: footer,
		},
		Author: &discordgo.MessageEmbedAuthor{
			Name: guild.Name,
			IconURL: guild.IconURL(),
		},
	}
	return &embeds
}