package music

import (
	"context"
	"strconv"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgolink/v3/lavalink"
	"github.com/disgoorg/snowflake/v2"
	"github.com/jckli/hyme/src/utils"
)

func TrackHandler(
	ctx context.Context,
	b *Music,
	e *handler.CommandEvent,
	channelId snowflake.ID,
	tracks ...lavalink.Track,
) error {
	_, ok := e.Client().Caches().VoiceState(*e.GuildID(), e.ApplicationID())
	if !ok {
		err := e.Client().
			UpdateVoiceState(context.Background(), *e.GuildID(), &channelId, false, true)
		if err != nil {
			e.Client().
				Rest().
				UpdateInteractionResponse(
					e.ApplicationID(),
					e.Token(),
					discord.
						NewMessageUpdateBuilder().
						SetEmbeds(utils.ErrorEmbed("An error has occured")).
						Build(),
				)
			return err
		}
	}

	player := b.Lavalink.Player(*e.GuildID())
	queue := b.Players.Get(*e.GuildID())

	var description string

	if player.Track() == nil {
		track := tracks[0]
		tracks = tracks[1:]

		err := player.Update(ctx, lavalink.WithTrack(track))
		if err != nil {
			e.Client().
				Rest().
				UpdateInteractionResponse(
					e.ApplicationID(),
					e.Token(),
					discord.
						NewMessageUpdateBuilder().
						SetEmbeds(utils.ErrorEmbed("An error has occured")).
						Build(),
				)
			return err
		}
		description = "Playing track: [`" + track.Info.Title + "`](" + *track.Info.URI + ")"
	}
	if len(tracks) > 0 || len(queue.Tracks) > 0 {
		queue.Add(tracks...)
		description += "Added `" + strconv.Itoa(
			len(queue.Tracks),
		) + "` tracks to the queue."
	}

	e.Client().
		Rest().
		UpdateInteractionResponse(
			e.ApplicationID(),
			e.Token(),
			discord.
				NewMessageUpdateBuilder().
				SetEmbeds(utils.SuccessEmbed(description)).
				Build(),
		)
	return nil
}
