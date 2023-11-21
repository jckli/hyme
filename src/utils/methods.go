package utils

import (
	"fmt"

	"github.com/disgoorg/disgolink/v3/lavalink"
)

func FormatDuration(duration lavalink.Duration) string {
	if duration == 0 {
		return "00:00"
	}
	return fmt.Sprintf("%02d:%02d", duration.Minutes(), duration.SecondsPart())
}

func Chunks(arr []lavalink.Track, size int) [][]lavalink.Track {
	var chunks [][]lavalink.Track
	for i := 0; i < len(arr); i += size {
		end := i + size
		if end > len(arr) {
			end = len(arr)
		}
		chunks = append(chunks, arr[i:end])
	}
	return chunks
}
