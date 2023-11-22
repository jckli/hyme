package commands

import (
	"context"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/disgolink/v3/lavalink"
	"github.com/jckli/hyme/src/dbot"
	"github.com/jckli/hyme/src/music"
	"github.com/jckli/hyme/src/utils"
	"math/rand"
	"time"
)

var hypeCommand = discord.SlashCommandCreate{
	Name:        "hype",
	Description: "Plays the HYPE!!!11!! playlist",
}

func hypeHandler(e *handler.CommandEvent, b *dbot.Bot) error {
	voiceState, ok := e.Client().Caches().VoiceState(*e.GuildID(), e.User().ID)
	if !ok {
		return e.Respond(
			discord.InteractionResponseTypeCreateMessage,
			discord.NewMessageCreateBuilder().
				SetEmbeds(utils.ErrorEmbed("You are not in a voice channel. Please join one and try again.")).
				Build(),
		)
	}

	err := e.DeferCreateMessage(false)
	if err != nil {
		b.Music.MusicLogger.Error(err)
		return err
	}

	player := b.Music.Lavalink.Player(*e.GuildID())

	go func() {
		var loadErr error
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		player.Node().
			LoadTracksHandler(ctx, "https://open.spotify.com/playlist/5O1zB17DBI71eThy3IzXfP?si=165751781e114864", disgolink.NewResultHandler(
				func(track lavalink.Track) {},
				func(playlist lavalink.Playlist) {
					source := rand.NewSource(time.Now().UnixNano())
					rng := rand.New(source)
					rng.Shuffle(len(playlist.Tracks), func(i, j int) {
						playlist.Tracks[i], playlist.Tracks[j] = playlist.Tracks[j], playlist.Tracks[i]
					})

					loadErr = music.TrackHandler(
						ctx,
						b.Music,
						e,
						*voiceState.ChannelID,
						playlist.Tracks...,
					)
				},
				func(tracks []lavalink.Track) {},
				func() {},
				func(err error) {
					_, loadErr = e.Client().
						Rest().
						UpdateInteractionResponse(
							e.ApplicationID(),
							e.Token(),
							discord.
								NewMessageUpdateBuilder().
								SetEmbeds(utils.ErrorEmbed("Error loading track: `"+err.Error()+"`")).
								Build(),
						)
				},
			))
		if loadErr != nil {
			b.Music.MusicLogger.Error(loadErr)
		}
	}()

	return nil
}
