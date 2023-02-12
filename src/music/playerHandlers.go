package music

import (
	"context"
	"fmt"
	"time"

	"github.com/disgoorg/disgolink/v2/disgolink"
	"github.com/disgoorg/disgolink/v2/lavalink"
	"github.com/disgoorg/snowflake/v2"
	"github.com/disgoorg/log"
)

func (b *Bot) onPlayerPause(player disgolink.Player, event lavalink.PlayerPauseEvent) {
	fmt.Printf("onPlayerPause: %v\n", event)
	queue := b.Players.Get(event.GuildID().String())
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		fmt.Println("Started timeout")
		queue.Cancel = cancel
		defer cancel()
		<- ctx.Done()
		if ctx.Err() == context.DeadlineExceeded {
			fmt.Println("DeadlineExceeded")
			player := b.Lavalink.ExistingPlayer(snowflake.MustParse(event.GuildID().String()))
			player.Update(context.TODO(), lavalink.WithNullTrack())
			queue.Cancel()
			b.Session.ChannelVoiceJoinManual(event.GuildID().String(), "", false, false)
			b.Lavalink.RemovePlayer(snowflake.MustParse(event.GuildID().String()))
			queue.Clear()
		}
	}()
}

func (b *Bot) onPlayerResume(player disgolink.Player, event lavalink.PlayerResumeEvent) {
	fmt.Printf("onPlayerResume: %v\n", event)
	queue := b.Players.Get(event.GuildID().String())
	if queue.Cancel != nil {
		queue.Cancel()
	}
}

func (b *Bot) onTrackStart(player disgolink.Player, event lavalink.TrackStartEvent) {
	fmt.Printf("onTrackStart: %v\n", event)
	queue := b.Players.Get(event.GuildID().String())
	fmt.Println(queue.Cancel)
	if queue.Cancel != nil {
		fmt.Println("Cancel auto disconnect")
		queue.Cancel()
	}
}

func (b *Bot) onTrackEnd(player disgolink.Player, event lavalink.TrackEndEvent) {
	fmt.Printf("onTrackEnd: %v\n", event)

	queue := b.Players.Get(event.GuildID().String())
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
		fmt.Println("Started timeout")
		queue.Cancel = cancel
		defer cancel()
		<- ctx.Done()
		if ctx.Err() == context.DeadlineExceeded {
			fmt.Println("DeadlineExceeded")
			player := b.Lavalink.ExistingPlayer(snowflake.MustParse(event.GuildID().String()))
			player.Update(context.TODO(), lavalink.WithNullTrack())
			queue.Cancel()
			b.Session.ChannelVoiceJoinManual(event.GuildID().String(), "", false, false)
			b.Lavalink.RemovePlayer(snowflake.MustParse(event.GuildID().String()))
			queue.Clear()
		}
	}()

	if !event.Reason.MayStartNext() {
		return
	}
	
	var (
		nextTrack lavalink.Track
		ok bool
	)
	switch queue.Type {
	case QueueTypeNormal:
		nextTrack, ok = queue.Next()

	case QueueTypeRepeatTrack:
		nextTrack = *player.Track()

	case QueueTypeRepeatQueue:
		lastTrack, _ := b.Lavalink.BestNode().DecodeTrack(context.TODO(), event.EncodedTrack)
		queue.Add(*lastTrack)
		nextTrack, ok = queue.Next()
	}

	if !ok {
		return
	}
	if err := player.Update(context.TODO(), lavalink.WithTrack(nextTrack)); err != nil {
		log.Error("Failed to play next track: ", err)
	}
}

func (b *Bot) onTrackException(player disgolink.Player, event lavalink.TrackExceptionEvent) {
	fmt.Printf("onTrackException: %v\n", event)
}

func (b *Bot) onTrackStuck(player disgolink.Player, event lavalink.TrackStuckEvent) {
	fmt.Printf("onTrackStuck: %v\n", event)
}

func (b *Bot) onWebSocketClosed(player disgolink.Player, event lavalink.WebSocketClosedEvent) {
	fmt.Printf("onWebSocketClosed: %v\n", event)
}