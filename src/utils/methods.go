package utils

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/disgoorg/disgolink/v2/lavalink"
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

func ConvertMilliToTime(millis int64) string {
	hours := millis / (1000 * 60 * 60)
	millis -= hours * (1000 * 60 * 60)
	minutes := millis / (1000 * 60)
	millis -= minutes * (1000 * 60)
	seconds := millis / 1000
	result := ""
	if hours > 0 {
		result += fmt.Sprintf("%d:", hours)
	}
	if hours > 0 || minutes > 0 {
		if !(hours > 0) {
			result += fmt.Sprintf("%d:", minutes)
		} else {
			result += fmt.Sprintf("%02d:", minutes)
		}
	}
	if !(hours > 0) && !(minutes > 0) {
		result += fmt.Sprintf("0:%d", seconds)
	} else {
		result += fmt.Sprintf("%02d", seconds)
	}

	return result
}

func Chunks(arr []lavalink.Track, size int) [][]lavalink.Track {
	var chunks [][]lavalink.Track
	for i := 0; i < len(arr); i += size {
		end := i + size
		if end > len(arr) {
			end = len(arr)
		}
		chunks = append(chunks, arr[i:end])
	}
	return chunks
}

func FormatPosition(position lavalink.Duration) string {
	if position == 0 {
		return "0:00"
	}
	return fmt.Sprintf("%d:%02d", position.Minutes(), position.SecondsPart())
}