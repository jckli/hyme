package music

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"github.com/bwmarrin/discordgo"
	"github.com/disgoorg/disgolink/v2/disgolink"
	"github.com/disgoorg/snowflake/v2"
)

type Bot struct {
	Session  *discordgo.Session
	Lavalink disgolink.Client
	Players *PlayerManager
}

func InitLink(s *discordgo.Session, b *Bot) disgolink.Client {
	link := disgolink.New(snowflake.MustParse(s.State.User.ID),
		disgolink.WithListenerFunc(b.onPlayerPause),
		disgolink.WithListenerFunc(b.onPlayerResume),
		disgolink.WithListenerFunc(b.onTrackStart),
		disgolink.WithListenerFunc(b.onTrackEnd),
		disgolink.WithListenerFunc(b.onTrackException),
		disgolink.WithListenerFunc(b.onTrackStuck),
		disgolink.WithListenerFunc(b.onWebSocketClosed),
	)
	return link
}

func OnVoiceStateUpdate(ctx context.Context, session *discordgo.Session, event *discordgo.VoiceStateUpdate, b *Bot) {
	var guildID *snowflake.ID
	if event.GuildID != "" {
		id := snowflake.MustParse(event.GuildID)
		guildID = &id
	}
	b.Lavalink.OnVoiceStateUpdate(ctx, snowflake.MustParse(event.VoiceState.GuildID), guildID, event.VoiceState.SessionID)
}

func OnVoiceServerUpdate(ctx context.Context, session *discordgo.Session, event *discordgo.VoiceServerUpdate, b *Bot) {
	b.Lavalink.OnVoiceServerUpdate(ctx, snowflake.MustParse(event.GuildID), event.Token, event.Endpoint)
}

func (b *Bot) RegisterNodes() disgolink.Node {
	secure, _ := strconv.ParseBool(os.Getenv("LAVALINK_SECURE"))
	node, err := b.Lavalink.AddNode(context.TODO(), disgolink.NodeConfig{
		Name:        os.Getenv("LAVALINK_NAME"),
		Address:     os.Getenv("LAVALINK_HOST") + ":" + os.Getenv("LAVALINK_PORT"),
		Password:    os.Getenv("LAVALINK_PASSWORD"),
		Secure:      secure,
	})
	if err != nil {
		fmt.Println(err)
	}
	return node
}