package commands

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/disgolink/v3/lavalink"
	"github.com/disgoorg/paginator"
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
				SetEmbeds(utils.ErrorEmbed("An error has occurred.")).
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
	player := b.Music.Lavalink.Player(*e.GuildID())
	if player == nil || player.Track() == nil {
		return e.Respond(
			discord.InteractionResponseTypeCreateMessage,
			discord.NewMessageUpdateBuilder().
				SetEmbeds(utils.ErrorEmbed("I am currently not playing anything.")).
				Build(),
		)
	}
	queue := b.Music.Players.Get(*e.GuildID())
	if queue == nil {
		return e.Respond(
			discord.InteractionResponseTypeCreateMessage,
			discord.NewMessageUpdateBuilder().
				SetEmbeds(utils.ErrorEmbed("There is no queue.")).
				Build(),
		)
	}

	track := player.Track()

	queuePages := []string{}
	var pageText string
	if len(queue.Tracks) == 0 {
		pageText = "No songs in queue."
		queuePages = append(queuePages, pageText)
	} else {
		split := utils.Chunks(queue.Tracks, 10)
		i := 1
		for _, chunk := range split {
			for _, track := range chunk {
				track := fmt.Sprintf(
					"%d. [`%s`](%s) by `%s` [%s]\n",
					i,
					track.Info.Title,
					*track.Info.URI,
					track.Info.Author,
					utils.FormatDuration(track.Info.Length),
				)
				pageText += track
				i++
			}
			queuePages = append(queuePages, pageText)
			pageText = ""

		}
	}

	err := b.Paginator.Create(e.Respond, paginator.Pages{
		ID: e.ID().String(),
		PageFunc: func(page int, embed *discord.EmbedBuilder) {
			utils.QueueEmbedHandler(embed, *track, queuePages[page])
			embed.SetFooterText(
				fmt.Sprintf("Page %d/%d", page+1, len(queuePages)),
			)
		},
		Pages:      len(queuePages),
		Creator:    e.User().ID,
		ExpireMode: paginator.ExpireModeAfterLastUsage,
	}, false)
	if err != nil {
		b.Music.MusicLogger.Error(err)
		return e.Respond(
			discord.InteractionResponseTypeCreateMessage,
			discord.NewMessageUpdateBuilder().
				SetEmbeds(utils.ErrorEmbed("An error has occurred.")).
				Build(),
		)
	}

	return nil
}

var disconnectCommand = discord.SlashCommandCreate{
	Name:        "disconnect",
	Description: "Disconnects the bot from the voice channel",
}

func disconnectHandler(e *handler.CommandEvent, b *dbot.Bot) error {
	player := b.Music.Lavalink.Player(*e.GuildID())
	if player == nil {
		return e.Respond(
			discord.InteractionResponseTypeCreateMessage,
			discord.NewMessageUpdateBuilder().
				SetEmbeds(utils.ErrorEmbed("I am currently not in a voice channel.")).
				Build(),
		)
	}

	voiceState, ok := e.Client().Caches().VoiceState(*e.GuildID(), e.User().ID)
	if !ok {
		return e.Respond(
			discord.InteractionResponseTypeCreateMessage,
			discord.NewMessageUpdateBuilder().
				SetEmbeds(utils.ErrorEmbed("You are not in the same voice channel as me.")).
				Build(),
		)
	}
	if *voiceState.ChannelID != *player.ChannelID() {
		return e.Respond(
			discord.InteractionResponseTypeCreateMessage,
			discord.NewMessageUpdateBuilder().
				SetEmbeds(utils.ErrorEmbed("You are not in the same voice channel as me.")).
				Build(),
		)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := e.Client().UpdateVoiceState(ctx, *e.GuildID(), nil, false, false)
	if err != nil {
		b.Music.MusicLogger.Error(err)
		return e.Respond(
			discord.InteractionResponseTypeCreateMessage,
			discord.NewMessageUpdateBuilder().
				SetEmbeds(utils.ErrorEmbed("An error has occurred.")).
				Build(),
		)
	}

	player.Update(ctx, lavalink.WithNullTrack())
	b.Music.Players.Delete(*e.GuildID())

	return e.Respond(
		discord.InteractionResponseTypeCreateMessage,
		discord.NewMessageUpdateBuilder().
			SetEmbeds(utils.ErrorEmbed("Disconnected from the voice channel.")).
			Build(),
	)
}

var stopCommand = discord.SlashCommandCreate{
	Name:        "stop",
	Description: "Skips to the next song and pauses the player",
}

func stopHandler(e *handler.CommandEvent, b *dbot.Bot) error {
	player := b.Music.Lavalink.Player(*e.GuildID())
	if player == nil || player.Track() == nil {
		return e.Respond(
			discord.InteractionResponseTypeCreateMessage,
			discord.NewMessageUpdateBuilder().
				SetEmbeds(utils.ErrorEmbed("I am currently not playing anything.")).
				Build(),
		)
	}

	voiceState, vsok := e.Client().
		Caches().
		VoiceState(*e.GuildID(), e.User().ID)
	if !vsok {
		return e.Respond(
			discord.InteractionResponseTypeCreateMessage,
			discord.NewMessageUpdateBuilder().
				SetEmbeds(utils.ErrorEmbed("You are not in the same voice channel as me.")).
				Build(),
		)
	}
	if *voiceState.ChannelID != *player.ChannelID() {
		return e.Respond(
			discord.InteractionResponseTypeCreateMessage,
			discord.NewMessageUpdateBuilder().
				SetEmbeds(utils.ErrorEmbed("You are not in the same voice channel as me.")).
				Build(),
		)
	}

	var (
		track lavalink.Track
		ok    bool
	)
	queue := b.Music.Players.Get(*e.GuildID())
	if len(queue.Tracks) == 0 {
		player.Update(context.Background(), lavalink.WithNullTrack())
	} else {
		track, ok = queue.Skip(1)
		if !ok {
			return e.Respond(
				discord.InteractionResponseTypeCreateMessage,
				discord.NewMessageUpdateBuilder().
					SetEmbeds(utils.ErrorEmbed("An error has occurred.")).
					Build(),
			)
		}

		err := player.Update(context.Background(), lavalink.WithTrack(track))
		if err != nil {
			b.Music.MusicLogger.Error(err)
			return e.Respond(
				discord.InteractionResponseTypeCreateMessage,
				discord.NewMessageUpdateBuilder().
					SetEmbeds(utils.ErrorEmbed("An error has occurred.")).
					Build(),
			)
		}
		err = player.Update(context.Background(), lavalink.WithPaused(true))
		if err != nil {
			b.Music.MusicLogger.Error(err)
			return e.Respond(
				discord.InteractionResponseTypeCreateMessage,
				discord.NewMessageUpdateBuilder().
					SetEmbeds(utils.ErrorEmbed("An error has occurred.")).
					Build(),
			)
		}
	}

	return e.Respond(
		discord.InteractionResponseTypeCreateMessage,
		discord.NewMessageUpdateBuilder().
			SetEmbeds(utils.StopEmbedHandler(&track)).
			Build(),
	)
}

var pauseCommand = discord.SlashCommandCreate{
	Name:        "pause",
	Description: "Pauses the player",
}

func pauseHandler(e *handler.CommandEvent, b *dbot.Bot) error {
	player := b.Music.Lavalink.Player(*e.GuildID())
	if player == nil || player.Track() == nil {
		return e.Respond(
			discord.InteractionResponseTypeCreateMessage,
			discord.NewMessageUpdateBuilder().
				SetEmbeds(utils.ErrorEmbed("I am currently not playing anything.")).
				Build(),
		)
	}

	voiceState, vsok := e.Client().
		Caches().
		VoiceState(*e.GuildID(), e.User().ID)
	if !vsok {
		return e.Respond(
			discord.InteractionResponseTypeCreateMessage,
			discord.NewMessageUpdateBuilder().
				SetEmbeds(utils.ErrorEmbed("You are not in the same voice channel as me.")).
				Build(),
		)
	}
	if *voiceState.ChannelID != *player.ChannelID() {
		return e.Respond(
			discord.InteractionResponseTypeCreateMessage,
			discord.NewMessageUpdateBuilder().
				SetEmbeds(utils.ErrorEmbed("You are not in the same voice channel as me.")).
				Build(),
		)
	}

	err := player.Update(context.Background(), lavalink.WithPaused(true))
	if err != nil {
		b.Music.MusicLogger.Error(err)
		return e.Respond(
			discord.InteractionResponseTypeCreateMessage,
			discord.NewMessageUpdateBuilder().
				SetEmbeds(utils.ErrorEmbed("An error has occurred.")).
				Build(),
		)
	}

	return e.Respond(
		discord.InteractionResponseTypeCreateMessage,
		discord.NewMessageUpdateBuilder().
			SetEmbeds(utils.SuccessEmbed("Successfully paused the player.")).
			Build(),
	)
}
