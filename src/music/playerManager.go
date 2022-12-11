package music

import (
	"github.com/disgoorg/disgolink/v2/lavalink"
)

type PlayerManager struct {
	Queues map[string]*Queue 
}

func (p *PlayerManager) Get(guildId string) *Queue {
	queue, ok := p.Queues[guildId]
	if !ok {
		queue = &Queue{
			Tracks: make([]lavalink.Track, 0),
			Type:   QueueTypeNormal,
		}
		p.Queues[guildId] = queue
	}
	return queue
}

func (p *PlayerManager) Remove(guildId string) {
	delete(p.Queues, guildId)
}    