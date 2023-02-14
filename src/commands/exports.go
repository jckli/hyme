package commands

import (
	"fmt"
	"os"
	"github.com/bwmarrin/discordgo"
	"github.com/jckli/hyme/src/music"
	"github.com/TopiSenpai/dgo-paginator"
)

type CommandHandler func(s *discordgo.Session, i *discordgo.InteractionCreate, bot *music.Bot, manager *paginator.Manager)

type Command struct {
	Command 	*discordgo.ApplicationCommand
	Handler 	CommandHandler
}

func CreateCommands(s *discordgo.Session) {
	for _, v := range Commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, os.Getenv("DISCORD_GUILD_ID"), v.Command)
		if err != nil {
			fmt.Println("Error creating command: ", err)
			return
		}
		fmt.Println("Created command: ", cmd.Name)
	}
}

func DeleteCommands(s *discordgo.Session) {
	commands, _ := s.ApplicationCommands(s.State.User.ID, os.Getenv("DISCORD_GUILD_ID"))
	for _, v := range commands {
		s.ApplicationCommandDelete(s.State.User.ID, os.Getenv("DISCORD_GUILD_ID"), v.ID)
	}
}

func InteractionRecieved(s *discordgo.Session, i *discordgo.InteractionCreate, bot *music.Bot, manager *paginator.Manager) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}
	for _, v := range Commands {
		if i.ApplicationCommandData().Name == v.Command.Name {
			v.Handler(s, i, bot, manager)
		}
	}
}

func New(command *discordgo.ApplicationCommand, handler CommandHandler) *Command {
	return &Command{
		Command: command,
		Handler: handler,
	}
}

var Commands = map[string]*Command{
	"ping": New(&discordgo.ApplicationCommand{
		Name: "ping",
		Type: discordgo.ChatApplicationCommand,
		Description: "Pong!",
	}, Ping),
	"hyme": New(&discordgo.ApplicationCommand{
		Name: "hyme",
		Type: discordgo.ChatApplicationCommand,
		Description: "Basic information about Hyme",
	}, Info),
	"play": New(&discordgo.ApplicationCommand{
		Name: "play",
		Type: discordgo.ChatApplicationCommand,
		Description: "Play a song",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type: discordgo.ApplicationCommandOptionString,
				Name: "query",
				Description: "Song or playlist to play",
				Required: true,
			},
		},
	}, PlayTrack),
	"disconnect": New(&discordgo.ApplicationCommand{
		Name: "disconnect",
		Type: discordgo.ChatApplicationCommand,
		Description: "Disconnect the bot from the voice channel",
	}, Disconnect),
	"connect": New(&discordgo.ApplicationCommand{
		Name: "connect",
		Type: discordgo.ChatApplicationCommand,
		Description: "Connect the bot to the voice channel",
	}, Connect),
	"pause": New(&discordgo.ApplicationCommand{
		Name: "pause",
		Type: discordgo.ChatApplicationCommand,
		Description: "Pauses/resumes the player",
	}, Pause),
	"stop": New(&discordgo.ApplicationCommand{
		Name: "stop",
		Type: discordgo.ChatApplicationCommand,
		Description: "Skips the current track and pauses the player",
	}, Stop),
	"queue": New(&discordgo.ApplicationCommand{
		Name: "queue",
		Type: discordgo.ChatApplicationCommand,
		Description: "Shows the current queue",
	}, Queue),
	"shuffle": New(&discordgo.ApplicationCommand{
		Name: "shuffle",
		Type: discordgo.ChatApplicationCommand,
		Description: "Shuffles the queue",
	}, Shuffle),
	"skip": New(&discordgo.ApplicationCommand{
		Name: "skip",
		Type: discordgo.ChatApplicationCommand,
		Description: "Skips the current track",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type: discordgo.ApplicationCommandOptionInteger,
				Name: "to",
				Description: "Track to skip to",
				Required: false,
			},
		},
	}, Skip),
	"nowplaying": New(&discordgo.ApplicationCommand{
		Name: "nowplaying",
		Type: discordgo.ChatApplicationCommand,
		Description: "Shows the currently playing track",
	}, NowPlaying),
	"move": New(&discordgo.ApplicationCommand{
		Name: "move",
		Type: discordgo.ChatApplicationCommand,
		Description: "Moves a track in the queue",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type: discordgo.ApplicationCommandOptionInteger,
				Name: "from",
				Description: "Track to move",
				Required: true,
			},
			{
				Type: discordgo.ApplicationCommandOptionInteger,
				Name: "to",
				Description: "Track to move to",
				Required: true,
			},
		},
	}, Move),
	"remove": New(&discordgo.ApplicationCommand{
		Name: "remove",
		Type: discordgo.ChatApplicationCommand,
		Description: "Removes a track from the queue",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type: discordgo.ApplicationCommandOptionInteger,
				Name: "track",
				Description: "Track to remove",
				Required: true,
			},
		},
	}, Remove),
	"hype": New(&discordgo.ApplicationCommand{
		Name: "hype",
		Type: discordgo.ChatApplicationCommand,
		Description: "Auto-queues the HYPE!!!1!! playlist (a roulette playlist)",
	}, HypePlaylist),
}