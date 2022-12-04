package music

import (
	"context"
	"os"
	"sync"
	"fmt"
	"strconv"
	"github.com/bwmarrin/discordgo"
	"github.com/disgoorg/disgolink/dgolink"
	"github.com/disgoorg/disgolink/lavalink"
)

type Bot struct {
	Link           *dgolink.Link
	PlayerManagers map[string]*PlayerManager
}

type PlayerManager struct {
	lavalink.PlayerEventAdapter
	Player        lavalink.Player
	Queue         []lavalink.AudioTrack
	QueueMu       sync.Mutex
	RepeatingMode int
}

func InitLink(s *discordgo.Session) *dgolink.Link {
	link := dgolink.New(s)
	return link
}

func (b *Bot) RegisterNodes() {
	secure, _ := strconv.ParseBool(os.Getenv("LAVALINK_SECURE"))
	b.Link.AddNode(context.TODO(), lavalink.NodeConfig{
		Name:        "Chisato",
		Host:        os.Getenv("LAVALINK_HOST"),
		Port:        os.Getenv("LAVALINK_PORT"),
		Password:    os.Getenv("LAVALINK_PASSWORD"),
		Secure:      secure,
	})
}

func (m *PlayerManager) AddQueue(tracks ...lavalink.AudioTrack) {
	m.QueueMu.Lock()
	m.Queue = append(m.Queue, tracks...)
	m.QueueMu.Unlock()
}

func (m *PlayerManager) PopQueue() lavalink.AudioTrack {
	m.QueueMu.Lock()
	defer m.QueueMu.Unlock()
	if len(m.Queue) == 0 {
		return nil
	}
	track := m.Queue[0]
	m.Queue = m.Queue[1:]
	return track
}

func (m *PlayerManager) ClearQueue() {
	m.QueueMu.Lock()
	m.Queue = []lavalink.AudioTrack{}
	m.QueueMu.Unlock()
}

func (m *PlayerManager) GetQueue() []lavalink.AudioTrack {
	m.QueueMu.Lock()
	defer m.QueueMu.Unlock()
	return m.Queue
}

const (
	repeatingModeOff = iota
	repeatingModeSong
	repeatingModeQueue
)

func (m *PlayerManager) OnTrackEnd(player lavalink.Player, track lavalink.AudioTrack, reason lavalink.AudioTrackEndReason) {
	if !reason.MayStartNext() {
		return
	}
	switch m.RepeatingMode {
	case repeatingModeOff:
		if nextTrack := m.PopQueue(); nextTrack != nil {
			if err := player.Play(nextTrack); err != nil {
				fmt.Println("error playing next track:", err)
			}
		}
	case repeatingModeSong:
		if err := player.Play(track.Clone()); err != nil {
			fmt.Println("error playing next track:", err)
		}

	case repeatingModeQueue:
		m.AddQueue(track)
		if nextTrack := m.PopQueue(); nextTrack != nil {
			if err := player.Play(nextTrack); err != nil {
				fmt.Println("error playing next track:", err)
			}
		}
	}
}

