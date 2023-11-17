package dbot

import (
	"context"
	"os"
	"strconv"

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
	Music     *music.Music
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

func (b *Bot) Setup(listeners ...bot.EventListener) bot.Client {
	client, err := disgo.New(
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

	b.Music.Client = client

	return client

}

func (b *Bot) ReadyEvent(_ *events.Ready) {
	err := b.Music.Client.SetPresence(
		context.TODO(),
		gateway.WithListeningActivity("HYPE!!!1!!"),
		gateway.WithOnlineStatus(discord.OnlineStatusOnline),
	)
	if err != nil {
		b.Logger.Error("Error while setting presence: ", err)
	}

	b.Logger.Info("Bot presence set successfully.")
}
