package commands

import (
	"fmt"
	"os"
	"github.com/bwmarrin/discordgo"
	"github.com/jckli/hyme/src/music"
)

type CommandHandler func(s *discordgo.Session, i *discordgo.InteractionCreate, bot *music.Bot)

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

func InteractionRecieved(s *discordgo.Session, i *discordgo.InteractionCreate, bot *music.Bot) {
	for _, v := range Commands {
		if i.ApplicationCommandData().Name == v.Command.Name {
			v.Handler(s, i, bot)
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
	"pause": New(&discordgo.ApplicationCommand{
		Name: "pause",
		Type: discordgo.ChatApplicationCommand,
		Description: "Pauses the player",
	}, Pause),
}