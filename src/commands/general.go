package commands

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

func HandlePing(e *handler.CommandEvent) error {
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
