package commands

import (
	"fmt"
	"os"
	"github.com/bwmarrin/discordgo"
)

type CommandHandler func(s *discordgo.Session, i *discordgo.InteractionCreate)

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

func InteractionRecieved(s *discordgo.Session, i *discordgo.InteractionCreate) {
	for _, v := range Commands {
		if i.ApplicationCommandData().Name == v.Command.Name {
			v.Handler(s, i)
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
}