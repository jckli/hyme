package utils

import (
	"strconv"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgolink/v3/lavalink"
)

func ErrorEmbed(description string) discord.Embed {
	embed := discord.NewEmbedBuilder().
		SetDescription(description).
		SetColor(0xff4f4f).
		Build()
	return embed
}

func SuccessEmbed(description string) discord.Embed {
	embed := discord.NewEmbedBuilder().
		SetDescription(description).
		SetColor(0xa4849a).
		Build()
	return embed
}

func PlayEmbedHandler(track *lavalink.Track) discord.Embed {
	embed := discord.NewEmbedBuilder().
		SetColor(0xa4849a).
		SetTitle("Playing Track")

	var description string
	if track.Info.Author != "" {
		description = "[`" + track.Info.Title + "`](" + *track.Info.URI + ") by `" + track.Info.Author + "`\n[" + FormatDuration(
			track.Info.Length,
		) + "]"
	} else {
		description = "[`" + track.Info.Title + "`](" + *track.Info.URI + ")\n[" + FormatDuration(track.Info.Length) + "]"
	}
	embed.SetDescription(description)

	if track.Info.ArtworkURL != nil {
		embed.SetThumbnail(*track.Info.ArtworkURL)
	}

	return embed.Build()
}

func SkipEmbedHandler(track *lavalink.Track, amount int) discord.Embed {
	embed := discord.NewEmbedBuilder().
		SetColor(0xa4849a).
		SetTitle("Skipped " + strconv.Itoa(amount) + " Track(s)")

	description := "Now Playing:\n"
	if track.Info.Author != "" {
		description += "[`" + track.Info.Title + "`](" + *track.Info.URI + ") by `" + track.Info.Author + "`\n[" + FormatDuration(
			track.Info.Length,
		) + "]"
	} else {
		description += "[`" + track.Info.Title + "`](" + *track.Info.URI + ")\n[" + FormatDuration(track.Info.Length) + "]"
	}
	embed.SetDescription(description)

	if track.Info.ArtworkURL != nil {
		embed.SetThumbnail(*track.Info.ArtworkURL)
	}

	return embed.Build()
}

func QueueEmbedHandler(
	embed *discord.EmbedBuilder,
	track lavalink.Track,
	queue string,
) {
	embed.
		SetTitle("Queue")

	description := "**Now Playing:**\n"
	if track.Info.Author != "" {
		description += "[`" + track.Info.Title + "`](" + *track.Info.URI + ") by `" + track.Info.Author + "`\n[" + FormatDuration(
			track.Info.Length,
		) + "]"
	} else {
		description += "[`" + track.Info.Title + "`](" + *track.Info.URI + ")\n[" + FormatDuration(track.Info.Length) + "]"
	}
	embed.SetDescription(description)

	if track.Info.ArtworkURL != nil {
		embed.SetThumbnail(*track.Info.ArtworkURL)
	}

	embed.AddField("Up Next", queue, false)
}

func StopEmbedHandler(
	track *lavalink.Track,
) discord.Embed {
	embed := discord.NewEmbedBuilder().
		SetColor(0xa4849a)

	if track != nil && track.Info.Title != "" {
		if track.Info.Author != "" {
			embed.
				SetTitle("Stopped Playing and Paused Player").
				SetDescription("Skipped to: [`" + track.Info.Title + "`](" + *track.Info.URI + ") by `" + track.Info.Author + "`\n[" + FormatDuration(track.Info.Length) + "]")
		} else {
			embed.
				SetTitle("Stopped Playing and Paused Player").
				SetDescription("Skipped to: [`" + track.Info.Title + "`](" + *track.Info.URI + ")\n[" + FormatDuration(track.Info.Length) + "]")
		}
		if track.Info.ArtworkURL != nil {
			embed.SetThumbnail(*track.Info.ArtworkURL)
		}
	} else {
		embed.
			SetTitle("Stopped Playing").
			SetDescription("Skipped current song.")
	}

	return embed.Build()
}

func ResumeEmbedHandler(
	track *lavalink.Track,
) discord.Embed {
	embed := discord.NewEmbedBuilder().
		SetColor(0xa4849a)

	if track.Info.Author != "" {
		embed.
			SetTitle("Resumed Playing").
			SetDescription("[`" + track.Info.Title + "`](" + *track.Info.URI + ") by `" + track.Info.Author + "`\n[" + FormatDuration(track.Info.Position) + "/" + FormatDuration(track.Info.Length) + "]")
	} else {
		embed.
			SetTitle("Resumed Playing").
			SetDescription("[`" + track.Info.Title + "`](" + *track.Info.URI + ")\n[" + FormatDuration(track.Info.Position) + "/" + FormatDuration(track.Info.Length) + "]")
	}
	if track.Info.ArtworkURL != nil {
		embed.SetThumbnail(*track.Info.ArtworkURL)
	}

	return embed.Build()
}

func MoveEmbedHandler(
	track *lavalink.Track,
	from int,
	to int,
) discord.Embed {
	embed := discord.NewEmbedBuilder().
		SetColor(0xa4849a)

	if track.Info.Author != "" {
		embed.
			SetTitle("Moved Track").
			SetDescription("Successfully moved [`" + track.Info.Title + "`](" + *track.Info.URI + ") by `" + track.Info.Author + "` [" + FormatDuration(track.Info.Length) + "] from position `" + strconv.Itoa(from) + "` to position `" + strconv.Itoa(to) + "` in the queue")
	} else {
		embed.
			SetTitle("Moved Track").
			SetDescription("Successfully moved [`" + track.Info.Title + "`](" + *track.Info.URI + ") [" + FormatDuration(track.Info.Length) + "] from position `" + strconv.Itoa(from) + "` to position `" + strconv.Itoa(to) + "` in the queue")
	}
	if track.Info.ArtworkURL != nil {
		embed.SetThumbnail(*track.Info.ArtworkURL)
	}

	return embed.Build()
}

func RemoveEmbedHandler(
	track *lavalink.Track,
	position int,
) discord.Embed {
	embed := discord.NewEmbedBuilder().
		SetColor(0xa4849a)

	if track.Info.Author != "" {
		embed.
			SetTitle("Removed Track").
			SetDescription("Successfully removed [`" + track.Info.Title + "`](" + *track.Info.URI + ") by `" + track.Info.Author + "` [" + FormatDuration(track.Info.Length) + "] from position `" + strconv.Itoa(position) + "` in the queue")
	} else {
		embed.
			SetTitle("Removed Track").
			SetDescription("Successfully removed [`" + track.Info.Title + "`](" + *track.Info.URI + ") [" + FormatDuration(track.Info.Length) + "] from position `" + strconv.Itoa(position) + "` in the queue")
	}
	if track.Info.ArtworkURL != nil {
		embed.SetThumbnail(*track.Info.ArtworkURL)
	}

	return embed.Build()
}
