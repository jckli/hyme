package music

import (
	"github.com/disgoorg/disgolink/v2/lavalink"
)

type QueueType string

type Queue struct {
	Tracks []lavalink.Track
	Type   QueueType
}

const (
	QueueTypeNormal      QueueType = "normal"
	QueueTypeRepeatTrack QueueType = "repeat_track"
	QueueTypeRepeatQueue QueueType = "repeat_queue"
)

func (q QueueType) String() string {
	switch q {
	case QueueTypeNormal:
		return "Normal"
	case QueueTypeRepeatTrack:
		return "Repeat Track"
	case QueueTypeRepeatQueue:
		return "Repeat Queue"
	default:
		return "Unknown"
	}
}

func (q *Queue) Add(track ...lavalink.Track) {
	q.Tracks = append(q.Tracks, track...)
}

func (q *Queue) Next() (lavalink.Track, bool) {
	if len(q.Tracks) == 0 {
		return lavalink.Track{}, false
	}
	track := q.Tracks[0]
	q.Tracks = q.Tracks[1:]
	return track, true
}

func (q *Queue) Clear() {
	q.Tracks = make([]lavalink.Track, 0)
}
