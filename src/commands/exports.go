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
	skipCommand,
}

func CommandHandlers(b *dbot.Bot) *handler.Mux {
	h := handler.New()

	// General Commands
	h.Command("/ping", PingHandler)
	h.Command("/hyme", InfoHandler)

	// Music Commands
	h.Command("/play", func(e *handler.CommandEvent) error {
		return playHandler(e, b)
	})
	h.Command("/skip", func(e *handler.CommandEvent) error {
		return skipHandler(e, b)
	})

	return h
}
