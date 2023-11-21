package music

import (
	"context"
	"github.com/disgoorg/disgolink/v3/lavalink"
	"math/rand"
	"time"
)

type Queue struct {
	Tracks []lavalink.Track
	Type   QueueType
	Cancel context.CancelFunc
}

type QueueType string

const (
	QueueTypeNormal      QueueType = "normal"
	QueueTypeRepeatTrack QueueType = "repeat_track"
	QueueTypeRepeatQueue QueueType = "repeat_queue"
)

func (q *Queue) Next() (lavalink.Track, bool) {
	if len(q.Tracks) == 0 {
		return lavalink.Track{}, false
	}
	track := q.Tracks[0]
	q.Tracks = q.Tracks[1:]
	return track, true
}

func (q *Queue) Add(track ...lavalink.Track) {
	q.Tracks = append(q.Tracks, track...)
}

func (q *Queue) Shuffle() {
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)
	rng.Shuffle(len(q.Tracks), func(i, j int) {
		q.Tracks[i], q.Tracks[j] = q.Tracks[j], q.Tracks[i]
	})
}

func (q *Queue) Skip(amount int) (lavalink.Track, bool) {
	if len(q.Tracks) == 0 {
		return lavalink.Track{}, false
	}
	if amount > len(q.Tracks) {
		amount = len(q.Tracks)
	}
	track := q.Tracks[amount-1]
	q.Tracks = q.Tracks[amount:]
	return track, true
}

func (q *Queue) Clear() {
	q.Tracks = make([]lavalink.Track, 0)
}
