package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"github.com/bwmarrin/discordgo"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	session, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	session.Identify.Intents = discordgo.IntentsGuildMessages
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return
	}
	session.AddHandler(ReadyEvent)
	err = session.Open()
	if err != nil {
		fmt.Println("Error opening connection: ", err)
		return
	}
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