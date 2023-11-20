package commands

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgolink/v3/lavalink"
	"github.com/jckli/hyme/src/dbot"
)

var playCommand = discord.SlashCommandCreate{
	Name:        "play",
	Description: "Plays a song",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionString{
			Name:        "query",
			Description: "The song or playlist to play",
			Required:    true,
		},
		discord.ApplicationCommandOptionString{
			Name:        "source",
			Description: "The source to search from",
			Choices: []discord.ApplicationCommandOptionChoiceString{
				{
					Name:  "YouTube",
					Value: string(lavalink.SearchTypeYouTube),
				},
				{
					Name:  "Spotify",
					Value: "spsearch",
				},
				{
					Name:  "SoundCloud",
					Value: string(lavalink.SearchTypeSoundCloud),
				},
				{
					Name:  "Deezer",
					Value: "dzsearch",
				},
				{
					Name:  "Deezer ISRC",
					Value: "dzisrc",
				},
			},
		},
	},
}

func playHandler(e *handler.CommandEvent, b *dbot.Bot) {
	err := e.DeferCreateMessage(false)
	if err != nil {
		b.Music.MusicLogger.Error(err)
		return
	}

}
