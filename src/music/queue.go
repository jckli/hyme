package music

import (
	"context"
	"github.com/disgoorg/disgolink/v2/lavalink"
	"math/rand"
)

type QueueType string

type Queue struct {
	Tracks []lavalink.Track
	Type   QueueType
	Cancel context.CancelFunc
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

func (q *Queue) Shuffle() {
	rand.Shuffle(len(q.Tracks), func(i, j int) {
		q.Tracks[i], q.Tracks[j] = q.Tracks[j], q.Tracks[i]
	})
}

func (q *Queue) Remove(index int) {
	q.Tracks = append(q.Tracks[:index], q.Tracks[index+1:]...)
}

func (q *Queue) Move(from, to int) {
	track := q.Tracks[from]
	q.Tracks = append(q.Tracks[:from], q.Tracks[from+1:]...)
	q.Tracks = append(q.Tracks[:to], append([]lavalink.Track{track}, q.Tracks[to:]...)...)
}
