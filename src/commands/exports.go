package commands

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var CommandList = []discord.ApplicationCommandCreate{
	pingCommand,
	infoCommand,
}

func CommandHandlers() *handler.Mux {
	h := handler.New()

	h.Command("/ping", PingHandler)
	h.Command("/hyme", InfoHandler)

	return h
}
