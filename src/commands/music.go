package commands

import (
	"context"
	"regexp"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/disgolink/v3/lavalink"
	"github.com/jckli/hyme/src/dbot"
	"github.com/jckli/hyme/src/music"
	"github.com/jckli/hyme/src/utils"
)

var (
	urlPattern = regexp.MustCompile(
		"^https?://[-a-zA-Z0-9+&@#/%?=~_|!:,.;]*[-a-zA-Z0-9+&@#/%=~_|]?",
	)
	searchPattern = regexp.MustCompile(`^(.{2})(search|isrc):(.+)`)
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

func playHandler(e *handler.CommandEvent, b *dbot.Bot) error {
	data := e.SlashCommandInteractionData()
	query := data.String("query")

	if !urlPattern.MatchString(query) && !searchPattern.MatchString(query) {
		if source, ok := data.OptString("source"); ok {
			query = lavalink.SearchType(source).Apply(query)
		} else {
			query = lavalink.SearchTypeYouTube.Apply(query)
		}
	}

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
		player.Node().LoadTracksHandler(ctx, query, disgolink.NewResultHandler(
			func(track lavalink.Track) {
				loadErr = music.TrackHandler(
					ctx,
					b.Music,
					e,
					*voiceState.ChannelID,
					track,
				)
			},
			func(playlist lavalink.Playlist) {
				loadErr = music.TrackHandler(
					ctx,
					b.Music,
					e,
					*voiceState.ChannelID,
					playlist.Tracks...,
				)
			},
			func(tracks []lavalink.Track) {
				loadErr = music.TrackHandler(
					ctx,
					b.Music,
					e,
					*voiceState.ChannelID,
					tracks[0],
				)
			},
			func() {
				_, loadErr = e.Client().
					Rest().
					UpdateInteractionResponse(
						e.ApplicationID(),
						e.Token(),
						discord.
							NewMessageUpdateBuilder().
							SetEmbeds(utils.ErrorEmbed("No results found for: `"+query+"`.")).
							Build(),
					)
			},
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

var skipCommand = discord.SlashCommandCreate{
	Name:        "skip",
	Description: "Skips the current song",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionInt{
			Name:        "amount",
			Description: "The amount of songs to skip",
			Required:    false,
		},
	},
}

func skipHandler(e *handler.CommandEvent, b *dbot.Bot) error {
	player := b.Music.Lavalink.Player(*e.GuildID())
	queue := b.Music.Players.Get(*e.GuildID())
	if player == nil || player.Track() == nil {
		return e.Respond(
			discord.InteractionResponseTypeCreateMessage,
			discord.NewMessageUpdateBuilder().
				SetEmbeds(utils.ErrorEmbed("I am currently not playing anything.")).
				Build(),
		)
	}

	amount := 1
	if data, ok := e.SlashCommandInteractionData().OptInt("amount"); ok {
		amount = data
	}

	b.Music.MusicLogger.Infof("Skipping %d songs", amount)

	track, ok := queue.Skip(amount)
	if !ok {
		return e.Respond(
			discord.InteractionResponseTypeCreateMessage,
			discord.NewMessageUpdateBuilder().
				SetEmbeds(utils.ErrorEmbed("There are no more songs to skip.")).
				Build(),
		)
	}

	err := player.Update(context.Background(), lavalink.WithTrack(track))
	if err != nil {
		b.Music.MusicLogger.Error(err)
		return e.Respond(
			discord.InteractionResponseTypeCreateMessage,
			discord.NewMessageUpdateBuilder().
				SetEmbeds(utils.ErrorEmbed("An error has occured.")).
				Build(),
		)
	}

	return e.Respond(
		discord.InteractionResponseTypeCreateMessage,
		discord.NewMessageUpdateBuilder().
			SetEmbeds(utils.SkipEmbedHandler(&track, amount)).
			Build(),
	)
}

var queueCommand = discord.SlashCommandCreate{
	Name:        "queue",
	Description: "Shows the current queue",
}

func queueHandler(e *handler.CommandEvent, b *dbot.Bot) error {
	queue := b.Music.Players.Get(*e.GuildID())
	if queue == nil {
		return e.Respond(
			discord.InteractionResponseTypeCreateMessage,
			discord.NewMessageUpdateBuilder().
				SetEmbeds(utils.ErrorEmbed("There is no queue.")).
				Build(),
		)
	}

	return e.Respond(
		discord.InteractionResponseTypeCreateMessage,
		discord.NewMessageUpdateBuilder().
			SetEmbeds().
			Build(),
	)
}
