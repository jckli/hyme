package utils

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func GetCurrentVoiceChannel(userId string, guild *discordgo.Guild, s *discordgo.Session) (channel *discordgo.Channel, err error) {
	fmt.Print(guild.VoiceStates)
	for _, vs := range guild.VoiceStates {
		if vs.UserID == userId {
			channel, _ = s.State.Channel(vs.ChannelID)
			return channel, nil
		}
	}
	return nil, fmt.Errorf("user is not in a voice channel")
}