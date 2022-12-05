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
	session.State.TrackVoice = true
	session.Identify.Intents = discordgo.IntentGuilds | discordgo.IntentGuildVoiceStates
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return
	}
	bot := &music.Bot{
		Session: session,
		Players: &music.PlayerManager{
			Queues: make(map[string]*music.Queue),
		},
	}
	session.AddHandler(func(s *discordgo.Session, e *discordgo.VoiceStateUpdate) {
		music.OnVoiceStateUpdate(s, e, bot)
	})
	session.AddHandler(func(s *discordgo.Session, e *discordgo.VoiceServerUpdate) {
		music.OnVoiceServerUpdate(s, e, bot)
	})
	session.AddHandler(ReadyEvent)
	err = session.Open()
	if err != nil {
		fmt.Println("Error opening connection: ", err)
		return
	}
	bot.Lavalink = music.InitLink(session, bot)
	bot.RegisterNodes()
	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		commands.InteractionRecieved(s, i, bot)
	})
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