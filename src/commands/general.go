package commands

import (
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var pingCommand = discord.SlashCommandCreate{
	Name:        "ping",
	Description: "Pong!",
}

func PingHandler(e *handler.CommandEvent) error {
	var ping string
	if e.Client().HasGateway() {
		ping = e.Client().Gateway().Latency().String()
	}

	embed := discord.NewEmbedBuilder().
		SetTitle("Pong! 🏓").
		SetDescription("My ping is " + ping).
		SetColor(0xa4849a).
		SetTimestamp(e.CreatedAt()).
		Build()

	return e.Respond(
		discord.InteractionResponseTypeCreateMessage,
		discord.NewMessageCreateBuilder().SetEmbeds(embed).Build(),
	)

}

var infoCommand = discord.SlashCommandCreate{
	Name:        "hyme",
	Description: "Get basic info about Hyme",
}

func InfoHandler(e *handler.CommandEvent) error {
	var (
		guildCount  int
		memberCount int
	)
	e.Client().Caches().GuildsForEach(func(guild discord.Guild) {

		guildCount++
		memberCount += guild.MemberCount
	})

	description := fmt.Sprintf(
		"ohashi's music bot. written in go.\nsupports youtube, spotify, soundcloud, deezer, and bandcamp.\n\nserver count: %d\nuser count: %d",
		guildCount,
		memberCount,
	)

	botUser, _ := e.Client().Caches().SelfUser()

	embed := discord.NewEmbedBuilder().
		SetTitle("Hyme").
		SetDescription(description).
		SetColor(0xa4849a).
		SetAuthor("Hyme", "", *botUser.AvatarURL()).
		Build()

	return e.Respond(
		discord.InteractionResponseTypeCreateMessage,
		discord.NewMessageCreateBuilder().SetEmbeds(embed).Build(),
	)
}
