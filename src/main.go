package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/disgoorg/disgo/bot"
	"github.com/jckli/hyme/src/commands"
	"github.com/jckli/hyme/src/dbot"
	"github.com/jckli/hyme/src/music"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	hyme := dbot.New("v0.0.1")

	h := commands.CommandHandlers()

	client := hyme.Setup(
		h,
		bot.NewListenerFunc(hyme.ReadyEvent),
	)

	hyme.Music = music.CreateMusic(client)

	var err error
	if hyme.Config.DevMode {
		hyme.Logger.Info(
			"Running in dev mode. Syncing commands to server ID: ",
			hyme.Config.DevServerID,
		)
		_, err = client.Rest().
			SetGuildCommands(client.ApplicationID(), hyme.Config.DevServerID, commands.CommandList)
	} else {
		hyme.Logger.Info(
			"Running in global mode. Syncing commands globally.",
		)
		_, err = client.Rest().SetGlobalCommands(client.ApplicationID(), commands.CommandList)
	}
	if err != nil {
		hyme.Logger.Errorf("Failed to sync commands: %s", err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.OpenGateway(ctx)
	if err != nil {
		hyme.Logger.Fatal("Error while connecting: ", err)
	}
	defer client.Close(context.TODO())

	music.RegisterNodes(ctx, hyme.Music)

	hyme.Logger.Info("Bot is ready!")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
	hyme.Logger.Info("Shutting down...")
}
