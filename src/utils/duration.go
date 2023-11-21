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
