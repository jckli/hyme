package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/jckli/hyme/src/music"
	"github.com/jckli/hyme/src/commands"
	_ "github.com/joho/godotenv/autoload"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	session, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	session.Identify.Intents = discordgo.IntentGuilds | discordgo.IntentGuildVoiceStates
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return
	}
	bot := &music.Bot{
		Link:           music.InitLink(session),
		PlayerManagers: map[string]*music.PlayerManager{},
	}
	err = session.Open()
	if err != nil {
		fmt.Println("Error opening connection: ", err)
		return
	}
	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		commands.InteractionRecieved(s, i, bot)
	})
	session.AddHandler(ReadyEvent)
	bot.RegisterNodes()
	commands.CreateCommands(session)
	fmt.Println("Bot is running!")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	fmt.Println("Bot is shutting down!")
}

func ReadyEvent(session *discordgo.Session, event *discordgo.Ready) {
	fmt.Println("Bot is ready!")
	session.UpdateListeningStatus("HYPE!!!1!!")
}