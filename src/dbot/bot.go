package dbot

import (
	"context"
	"os"
	"strconv"
	"time"

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
		Logger:  logger,
		Version: version,
		Paginator: paginator.New(
			paginator.WithEmbedColor(0xa4849a),
			paginator.WithButtonsConfig(
				paginator.ButtonsConfig{
					First: paginator.DefaultConfig().ButtonsConfig.First,
					Back:  paginator.DefaultConfig().ButtonsConfig.Back,
					Stop:  nil,
					Next:  paginator.DefaultConfig().ButtonsConfig.Next,
					Last:  paginator.DefaultConfig().ButtonsConfig.Last,
				},
			),
			paginator.WithCleanupInterval(5*time.Minute),
		),
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
				gateway.IntentGuildMembers,
				gateway.IntentGuildVoiceStates,
			),
		),
		bot.WithEventListenerFunc(b.onVoiceStateUpdate),
		bot.WithEventListenerFunc(b.onVoiceServerUpdate),
		bot.WithEventListeners(listeners...),
		bot.WithEventListeners(b.Paginator),
		bot.WithCacheConfigOpts(
			cache.WithCaches(cache.FlagVoiceStates, cache.FlagGuilds),
		),
	)
	if err != nil {
		b.Logger.Fatal("Error while building DisGo client: ", err)
	}

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

func (b *Bot) onVoiceStateUpdate(event *events.GuildVoiceStateUpdate) {
	if event.VoiceState.UserID != b.Music.Client.ApplicationID() {
		return
	}
	b.Music.Lavalink.OnVoiceStateUpdate(
		context.TODO(),
		event.VoiceState.GuildID,
		event.VoiceState.ChannelID,
		event.VoiceState.SessionID,
	)
	if event.VoiceState.ChannelID == nil {
		b.Music.Players.Delete(event.VoiceState.GuildID)
	}
}

func (b *Bot) onVoiceServerUpdate(event *events.VoiceServerUpdate) {
	b.Music.Lavalink.OnVoiceServerUpdate(
		context.TODO(),
		event.GuildID,
		event.Token,
		*event.Endpoint,
	)
}
