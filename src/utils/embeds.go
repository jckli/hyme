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
