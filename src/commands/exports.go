package commands

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/jckli/hyme/src/dbot"
)

var CommandList = []discord.ApplicationCommandCreate{
	pingCommand,
	infoCommand,
	playCommand,
}

func CommandHandlers(b *dbot.Bot) *handler.Mux {
	h := handler.New()

	h.Command("/ping", PingHandler)
	h.Command("/hyme", InfoHandler)

	h.Command("/play", func(e *handler.CommandEvent) error {
		return playHandler(e, b)
	})

	return h
}
