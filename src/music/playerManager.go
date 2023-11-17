package music

import (
	"github.com/disgoorg/disgolink/v3/lavalink"
	"github.com/disgoorg/snowflake/v2"
)

type PlayerManager struct {
	Queues map[snowflake.ID]*Queue
}

func (p *PlayerManager) Get(guildID snowflake.ID) *Queue {
	queue, ok := p.Queues[guildID]
	if !ok {
		queue = &Queue{
			Tracks: make([]lavalink.Track, 0),
			Type:   QueueTypeNormal,
		}
		p.Queues[guildID] = queue
	}
	return queue
}

func (p *PlayerManager) Delete(guildID snowflake.ID) {
	delete(p.Queues, guildID)
}
