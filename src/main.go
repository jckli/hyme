package main

import (
	"context"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/jckli/hyme/src/commands"
	"github.com/jckli/hyme/src/music"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/log"
	"github.com/disgoorg/paginator"
	"github.com/disgoorg/snowflake/v2"

	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	Token       string
	DevMode     bool
	DevServerID snowflake.ID
}

type Bot struct {
	Logger    log.Logger
	Version   string
	Paginator *paginator.Manager
	Config    Config
	Client    bot.Client
	Music     *music.MusicBot
}

func New(version string) *Bot {
	devServerID, _ := strconv.Atoi(os.Getenv("DEV_SERVER_ID"))

	logger := log.New(log.Ldate | log.Ltime | log.Lshortfile)
	logger.SetLevel(2)
	logger.Infof("Starting bot version: %s", version)

	return &Bot{
		Logger:    logger,
		Version:   version,
		Paginator: paginator.New(),
		Config: Config{
			Token:       os.Getenv("TOKEN"),
			DevMode:     os.Getenv("DEV_MODE") == "true",
			DevServerID: snowflake.ID(devServerID),
		},
	}
}

func (b *Bot) Setup(listeners ...bot.EventListener) {
	var err error
	b.Client, err = disgo.New(
		b.Config.Token,
		bot.WithLogger(b.Logger),
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(
				gateway.IntentGuilds,
				gateway.IntentGuildVoiceStates,
			),
		),
		bot.WithEventListeners(listeners...),
		bot.WithEventListeners(b.Paginator),
		bot.WithCacheConfigOpts(cache.WithCaches(cache.FlagVoiceStates)),
	)
	if err != nil {
		b.Logger.Fatal("Error while building DisGo client: ", err)
	}

}

func (b *Bot) ReadyEvent(_ *events.Ready) {
	err := b.Client.SetPresence(
		context.TODO(),
		gateway.WithListeningActivity("HYPE!!!1!!"),
		gateway.WithOnlineStatus(discord.OnlineStatusOnline),
	)
	if err != nil {
		b.Logger.Error("Error while setting presence: ", err)
	}

	b.Logger.Info("Bot presence set successfully.")
}

func main() {
	hyme := New("v0.0.1")

	h := commands.CommandHandlers()

	hyme.Setup(
		h,
		bot.NewListenerFunc(hyme.ReadyEvent),
	)

	var err error
	if hyme.Config.DevMode {
		hyme.Logger.Info(
			"Running in dev mode. Syncing commands to server ID: ",
			hyme.Config.DevServerID,
		)
		_, err = hyme.Client.Rest().
			SetGuildCommands(hyme.Client.ApplicationID(), hyme.Config.DevServerID, commands.CommandList)
	} else {
		hyme.Logger.Info(
			"Running in global mode. Syncing commands globally.",
		)
		_, err = hyme.Client.Rest().SetGlobalCommands(hyme.Client.ApplicationID(), commands.CommandList)
	}
	if err != nil {
		hyme.Logger.Errorf("Failed to sync commands: %s", err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = hyme.Client.OpenGateway(ctx)
	if err != nil {
		hyme.Logger.Fatal("Error while connecting: ", err)
	}
	defer hyme.Client.Close(context.TODO())

	hyme.Logger.Info("Bot is ready!")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
	hyme.Logger.Info("Shutting down...")
}
