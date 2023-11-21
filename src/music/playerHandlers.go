package music

import (
	"context"
	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/disgolink/v3/lavalink"
	"github.com/disgoorg/snowflake/v2"
	"time"
)

func autoLeaveTimeout(
	queue *Queue,
	player disgolink.Player,
	b *Music,
	guildId *snowflake.ID,
) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Minute)
	b.MusicLogger.Infof("Starting auto leave timer for guild: %s", guildId)
	queue.Cancel = cancel
	<-ctx.Done()
	b.MusicLogger.Infof("Auto leave timer finished for guild: %s", guildId)
	player.Update(context.TODO(), lavalink.WithNullTrack())
	b.Client.UpdateVoiceState(context.TODO(), *guildId, nil, false, false)
	queue.Clear()
}

func onPlayerPause(
	player disgolink.Player,
	event lavalink.PlayerPauseEvent,
	b *Music,
) {
	b.MusicLogger.Infof("onPlayerPause: %s", event)
}

func onPlayerResume(
	player disgolink.Player,
	event lavalink.PlayerResumeEvent,
	b *Music,
) {
	b.MusicLogger.Infof("onPlayerResume: %s", event)
}

func onTrackStart(
	player disgolink.Player,
	event lavalink.TrackStartEvent,
	b *Music,
) {
	b.MusicLogger.Infof("onTrackStart: %s", event)
	queue := b.Players.Get(event.GuildID())
	if queue.Cancel != nil {
		queue.Cancel()
		queue.Cancel = nil
	}
}

func onTrackEnd(
	player disgolink.Player,
	event lavalink.TrackEndEvent,
	b *Music,
) {
	b.MusicLogger.Infof("onTrackEnd: %s", event)

	if !event.Reason.MayStartNext() {
		return
	}

	guildId := event.GuildID()
	go autoLeaveTimeout(b.Players.Get(guildId), player, b, &guildId)
	queue := b.Players.Get(guildId)
	var (
		nextTrack lavalink.Track
		ok        bool
	)
	switch queue.Type {
	case QueueTypeNormal:
		nextTrack, ok = queue.Next()
	case QueueTypeRepeatTrack:
		nextTrack = event.Track
	case QueueTypeRepeatQueue:
		queue.Add(event.Track)
		nextTrack, ok = queue.Next()
	}
	if !ok {
		return
	}

	err := player.Update(context.TODO(), lavalink.WithTrack(nextTrack))

	if err != nil {
		b.MusicLogger.Errorf("Failed to play next track: ", err)
	}
}

func onTrackException(
	player disgolink.Player,
	event lavalink.TrackExceptionEvent,
	b *Music,
) {
	b.MusicLogger.Infof("onTrackException: %s", event)
}

func onTrackStuck(
	player disgolink.Player,
	event lavalink.TrackStuckEvent,
	b *Music,
) {
	b.MusicLogger.Infof("onTrackStuck: %s", event)
}

func onWebSocketClosed(
	player disgolink.Player,
	event lavalink.WebSocketClosedEvent,
	b *Music,
) {
	b.MusicLogger.Infof("onWebSocketClosed: %s", event)
}

func onUnknownEvent(
	player disgolink.Player,
	event lavalink.UnknownEvent,
	b *Music,
) {
	b.MusicLogger.Infof(
		"onUnknownEvent: %s, data: %s\n",
		event.Type_,
		string(event.Data),
	)
}
