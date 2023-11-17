package music

import (
	"context"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/disgolink/v3/lavalink"
	"github.com/disgoorg/log"
	"os"
	"strconv"
)

type Music struct {
	Client      bot.Client
	MusicLogger log.Logger
	Lavalink    disgolink.Client
	Players     *PlayerManager
}

func InitLink(b *Music) disgolink.Client {
	link := disgolink.New(
		b.Client.ApplicationID(),
		disgolink.WithListenerFunc(
			func(player disgolink.Player, event lavalink.PlayerPauseEvent) {
				onPlayerPause(player, event, b)
			},
		),
		disgolink.WithListenerFunc(
			func(player disgolink.Player, event lavalink.PlayerResumeEvent) {
				onPlayerResume(player, event, b)
			},
		),
		disgolink.WithListenerFunc(
			func(player disgolink.Player, event lavalink.TrackStartEvent) {
				onTrackStart(player, event, b)
			},
		),
		disgolink.WithListenerFunc(
			func(player disgolink.Player, event lavalink.TrackEndEvent) {
				onTrackEnd(player, event, b)
			},
		),
		disgolink.WithListenerFunc(
			func(player disgolink.Player, event lavalink.TrackExceptionEvent) {
				onTrackException(player, event, b)
			},
		),
		disgolink.WithListenerFunc(
			func(player disgolink.Player, event lavalink.TrackStuckEvent) {
				onTrackStuck(player, event, b)
			},
		),
		disgolink.WithListenerFunc(
			func(player disgolink.Player, event lavalink.WebSocketClosedEvent) {
				onWebSocketClosed(player, event, b)
			},
		),
		disgolink.WithListenerFunc(
			func(player disgolink.Player, event lavalink.UnknownEvent) {
				onUnknownEvent(player, event, b)
			},
		),
	)

	return link

}

func RegisterNodes(ctx context.Context, b *Music) disgolink.Node {
	logger := log.New(log.Ldate | log.Ltime | log.Lshortfile)
	logger.SetLevel(1)
	b.MusicLogger = logger

	secure, _ := strconv.ParseBool(os.Getenv("LAVALINK_SECURE"))
	name := os.Getenv("LAVALINK_NAME")
	password := os.Getenv("LAVALINK_PASSWORD")
	host := os.Getenv("LAVALINK_HOST")
	port := os.Getenv("LAVALINK_PORT")

	node, err := b.Lavalink.AddNode(ctx, disgolink.NodeConfig{
		Name:     name,
		Address:  host + ":" + port,
		Password: password,
		Secure:   secure,
	})
	if err != nil {
		b.MusicLogger.Fatal(err)
	}
	version, err := node.Version(ctx)
	if err != nil {
		b.MusicLogger.Fatal(err)
	}

	b.MusicLogger.Infof(
		"Connected to lavalink node. Session ID: %s, Version: %s",
		node.SessionID(),
		version,
	)

	return node
}
