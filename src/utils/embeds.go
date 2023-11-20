package utils

import (
	"github.com/disgoorg/disgo/discord"
)

func ErrorEmbed(description string) discord.Embed {
	embed := discord.NewEmbedBuilder().
		SetDescription(description).
		SetColor(0xff4f4f).
		Build()
	return embed
}

func SuccessEmbed(description string) discord.Embed {
	embed := discord.NewEmbedBuilder().
		SetDescription(description).
		SetColor(0xa4849a).
		Build()
	return embed
}
