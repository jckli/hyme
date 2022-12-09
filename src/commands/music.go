package commands

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/disgoorg/disgolink/v2/disgolink"
	"github.com/disgoorg/disgolink/v2/lavalink"
	"github.com/disgoorg/snowflake/v2"
	"github.com/jckli/hyme/src/music"
	"github.com/jckli/hyme/src/utils"
)

var (
	urlPattern = regexp.MustCompile("^https?://[-a-zA-Z0-9+&@#/%?=~_|!:,.;]*[-a-zA-Z0-9+&@#/%=~_|]?")
	searchPattern = regexp.MustCompile(`^(.{2})search:(.+)`)
)

func PlayTrack(s *discordgo.Session, i *discordgo.InteractionCreate, bot *music.Bot) {
	// Defer the response, gives more time to process the command
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
	query := i.ApplicationCommandData().Options[0].StringValue()
	if !urlPattern.MatchString(query) && !searchPattern.MatchString(query) {
		query = lavalink.SearchTypeYoutube.Apply(query)
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
	bot.Lavalink.BestNode().LoadTracks(ctx, query, disgolink.NewResultHandler(
		func(track lavalink.Track) {
			s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
				Embeds: utils.SuccessEmbed("Playing track: [`"+ track.Info.Title +"`]("+ *track.Info.URI +")"),
			})
			if player.Track() == nil {
				toPlay = &track
			} else {
				queue.Add(track)
			}
		},
		func(playlist lavalink.Playlist) {
			s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
				Embeds: utils.SuccessEmbed("Playing playlist: `"+ playlist.Info.Name +"` with `"+ strconv.Itoa(len(playlist.Tracks)) +"` tracks." ),
			})
			if player.Track() == nil {
				toPlay = &playlist.Tracks[0]
				queue.Add(playlist.Tracks[1:]...)
			} else {
				queue.Add(playlist.Tracks...)
			}
		},
		func(tracks []lavalink.Track) {
			if player.Track() == nil {
				toPlay = &tracks[0]
				s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
					Embeds: utils.SuccessEmbed("Playing track: [`"+ tracks[0].Info.Title +"`]("+ *tracks[0].Info.URI +")"),
				})
			} else {
				queue.Add(tracks[0])
				s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
					Embeds: utils.SuccessEmbed("Adding track to queue: [`"+ tracks[0].Info.Title +"`]("+ *tracks[0].Info.URI +")"),
				})
			}
		},
		func() {
			s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
				Embeds: utils.ErrorEmbed("No results found for: `" + query + "`."),
			})
		},
		func(err error) {
			s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
				Embeds: utils.ErrorEmbed("Error loading track: `"+ err.Error() +"`"),
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

func Pause(s *discordgo.Session, i *discordgo.InteractionCreate, bot *music.Bot) {
	player := bot.Lavalink.ExistingPlayer(snowflake.MustParse(i.GuildID))
	track := player.Track()
	if player == nil || track == nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: utils.ErrorEmbed("I am not currently playing anything."),
			},
		})
		return
	}
	err := player.Update(context.Background(), lavalink.WithPaused(!player.Paused()))
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: utils.ErrorEmbed("Error pausing/resuming the track."),
			},
		})
		return
	}
	if player.Paused() {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: utils.SuccessEmbed("Paused the player."),
			},
		})
	} else {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: utils.SuccessEmbed("Resumed the player with track: [`"+ player.Track().Info.Title +"`]("+ *player.Track().Info.URI +")"),
			},
		})
	}
}

func Stop(s *discordgo.Session, i *discordgo.InteractionCreate, bot *music.Bot) {
	player := bot.Lavalink.ExistingPlayer(snowflake.MustParse(i.GuildID))
	curTrack := player.Track()
	if player == nil || curTrack == nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: utils.ErrorEmbed("I am not currently playing anything."),
			},
		})
		return
	}
	queue := bot.Players.Get(i.GuildID)
	nextTrack, has := queue.Next()
	if has {
		err1 := player.Update(context.TODO(), lavalink.WithTrack(nextTrack))
		err2 := player.Update(context.TODO(), lavalink.WithPaused(true))
		if err1 != nil || err2 != nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: utils.ErrorEmbed("Error stopping the track."),
				},
			})
			return
		}
	} else {
		// skip the current track and stop player
		err := player.Update(context.TODO(), lavalink.WithNullTrack())
		if err != nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: utils.ErrorEmbed("Error stopping the track."),
				},
			})
			return
		}
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: utils.SuccessEmbed("Skipped the current track and paused the player."),
		},
	})
}

func Disconnect(s *discordgo.Session, i *discordgo.InteractionCreate, bot *music.Bot) {
	player := bot.Lavalink.ExistingPlayer(snowflake.MustParse(i.GuildID))
	if player == nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: utils.ErrorEmbed("I am not connected to a voice channel."),
			},
		})
		return
	}
	err := bot.Session.ChannelVoiceJoinManual(i.GuildID, "", false, false)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: utils.ErrorEmbed("Couldn't disconnect from the voice channel."),
			},
		})
		return
	}
	bot.Lavalink.RemovePlayer(snowflake.MustParse(i.GuildID))
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: utils.SuccessEmbed("Disconnected from the voice channel."),
		},
	})
}

func Queue(s *discordgo.Session, i *discordgo.InteractionCreate, bot *music.Bot) {
	queue := bot.Players.Get(i.GuildID)
	if queue == nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: utils.ErrorEmbed("I am not connected to a voice channel."),
			},
		})
		return
	}
	if len(queue.Tracks) == 0 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: utils.ErrorEmbed("The queue is empty."),
			},
		})
		return
	}
	var description string
	for i, track := range queue.Tracks {
		description += fmt.Sprintf("%d. [`%s`](%s) [%s] \n", i+1, track.Info.Title, *track.Info.URI, utils.ConvertMilliToTime(int64(track.Info.Length)))
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: utils.SuccessEmbed("🔊 Queue `("+ strconv.Itoa(len(queue.Tracks)) +")`:\n"+ description),
		},
	})
}
