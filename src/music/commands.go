package music

import (
	"context"
	"fmt"
	"regexp"
	"github.com/bwmarrin/discordgo"
	"github.com/disgoorg/disgolink/lavalink"
	"github.com/disgoorg/snowflake/v2"
	"github.com/jckli/hyme/src/utils"
)

var (
	urlYoutubePattern = regexp.MustCompile(`^(https?\:\/\/)?(www\.youtube\.com|youtu\.?be)\/.+`)
)

func PlayTrack(s *discordgo.Session, i *discordgo.InteractionCreate, bot *Bot) {
	userid := i.Member.User.ID
	guild, _ := s.State.Guild(i.GuildID)
	// Defer the response, gives more time to process the command
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	voiceChannel, err := utils.GetCurrentVoiceChannel(userid, guild, s)
	if err != nil {
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Embeds: utils.ErrorEmbed("You are not in a voice channel. Please join one and try again."),
		})
		return
	}
	query := i.ApplicationCommandData().Options[0].StringValue()
	if !urlYoutubePattern.MatchString(query) {
		query = "ytsearch:" + query
	}
	fmt.Print(query)
	_ = bot.Link.BestRestClient().LoadItemHandler(context.TODO(), query, lavalink.NewResultHandler(
		func(track lavalink.AudioTrack) {
			bot.Play(s, i, i.GuildID, voiceChannel.ID, track)
		},
		func(playlist lavalink.AudioPlaylist) {
			bot.Play(s, i, i.GuildID, voiceChannel.ID, playlist.Tracks()...)
		},
		func(tracks []lavalink.AudioTrack) {
			bot.Play(s, i, i.GuildID, voiceChannel.ID, tracks[0])
		},
		func() {
			s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
				Embeds: utils.ErrorEmbed("No results found for: `" + query + "`."),
			})
		},
		func(ex lavalink.FriendlyException) {
			s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
				Embeds: utils.ErrorEmbed("An error occurred while loading the track."),
			})
		},
	))


}

func (b *Bot) Play(s *discordgo.Session, i *discordgo.InteractionCreate, guildId string, vcId string, tracks ...lavalink.AudioTrack) {
	err := s.ChannelVoiceJoinManual(guildId, vcId, false, true)
	if err != nil {
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Embeds: utils.ErrorEmbed("Couldn't join the voice channel."),
		})
		return
	}
	manager, status := b.PlayerManagers[guildId]
	if !status {
		manager = &PlayerManager{
			Player: b.Link.Player(snowflake.MustParse(guildId)),
			RepeatingMode: repeatingModeOff,
		}
		b.PlayerManagers[guildId] = manager
		manager.Player.AddListener(manager)
	}
	manager.AddQueue(tracks...)
	track := manager.PopQueue()
	err2 := manager.Player.Play(track)
	if err2 != nil {
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Embeds: utils.ErrorEmbed("Couldn't play the track."),
		})
		return
	}
	s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Embeds: utils.SuccessEmbed("Playing track: `" + track.Info().Title + "`."),
	})
}