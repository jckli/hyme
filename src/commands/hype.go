package commands

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"github.com/disgoorg/disgolink/v2/disgolink"
	"github.com/disgoorg/disgolink/v2/lavalink"
	"github.com/disgoorg/snowflake/v2"
	"github.com/TopiSenpai/dgo-paginator"
	"github.com/jckli/hyme/src/music"
	"github.com/jckli/hyme/src/utils"
	"math/rand"
	"strconv"
	"time"
)

func HypePlaylist(s *discordgo.Session, i *discordgo.InteractionCreate, bot *music.Bot, manager *paginator.Manager) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	userid := i.Member.User.ID
	voiceState, err := bot.Session.State.VoiceState(i.GuildID, userid)
	if err != nil {
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Embeds: utils.ErrorEmbed("You are not in a voice channel. Please join one and try again."),
		})
		return
	}
	var player disgolink.Player
	player = bot.Lavalink.ExistingPlayer(snowflake.MustParse(i.GuildID))
	if player != nil {
		curTrack := player.Track()
		if curTrack != nil && player.Paused() {
			player.Update(context.TODO(), lavalink.WithPaused(false))
		}
	} else {
		player = bot.Lavalink.Player(snowflake.MustParse(i.GuildID))
	}
	queue := bot.Players.Get(i.GuildID)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var toPlay *lavalink.Track
	bot.Lavalink.BestNode().LoadTracks(ctx, "https://open.spotify.com/playlist/5O1zB17DBI71eThy3IzXfP?si=6c4b3a3d6aa4412e", disgolink.NewResultHandler(
		func(track lavalink.Track) {},
		func(playlist lavalink.Playlist) {
			rand.Shuffle(len(playlist.Tracks), func(i, j int) {
				playlist.Tracks[i], playlist.Tracks[j] = playlist.Tracks[j], playlist.Tracks[i]
			})
			if player.Track() == nil {
				s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
					Embeds: utils.SuccessEmbed("Playing the `"+ playlist.Info.Name +"` playlist: `"+ strconv.Itoa(len(playlist.Tracks)) +"` tracks." ),
				})
				toPlay = &playlist.Tracks[0]
				queue.Add(playlist.Tracks[1:]...)
			} else {
				s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
					Embeds: utils.SuccessEmbed("Added the `"+ playlist.Info.Name +"` playlist to queue: `"+ strconv.Itoa(len(playlist.Tracks)) +"` tracks." ),
				})
				queue.Add(playlist.Tracks...)
			}
		},
		func(tracks []lavalink.Track) {},
		func() {},
		func(err error) {
			s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
				Embeds: utils.ErrorEmbed("Error loading: `"+ err.Error() +"`"),
			})
		},
	))
	if toPlay == nil {
		return
	}
	err2 := bot.Session.ChannelVoiceJoinManual(i.GuildID, voiceState.ChannelID, false, true)
	if err2 != nil {
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Embeds: utils.ErrorEmbed("Couldn't join the voice channel."),
		})
		return
	}
	player.Update(context.TODO(), lavalink.WithTrack(*toPlay))
}