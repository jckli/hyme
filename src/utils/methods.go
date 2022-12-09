package utils

import (
	"fmt"
	"time"
	"github.com/bwmarrin/discordgo"
)

func GetCurrentVoiceChannel(userId string, guild *discordgo.Guild, s *discordgo.Session) (channel *discordgo.Channel, err error) {
	for _, vs := range guild.VoiceStates {
		if vs.UserID == userId {
			channel, _ = s.State.Channel(vs.ChannelID)
			return channel, nil
		}
	}
	return nil, fmt.Errorf("user is not in a voice channel")
}

func ConvertMilliToTime(milli int64) string {
	duration := time.Duration(milli) * time.Millisecond
	duration = duration.Round(time.Second)
	duration = duration.Truncate(time.Minute)
	hours := duration / time.Hour
	minutes := duration % time.Hour / time.Minute
	seconds := duration % time.Minute / time.Second
	if hours > 0 {
		return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
	} else {
		if minutes > 0 {
			return fmt.Sprintf("%02d:%02d", minutes, seconds)
		} else {
			return fmt.Sprintf("%02d", seconds)
		}
	}
}