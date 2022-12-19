package app

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"park-wallet/internal/service/bot"
	"park-wallet/internal/service/ton"
	"park-wallet/internal/storage/postgres"
	"park-wallet/internal/storage/redis"
	"park-wallet/pkg"
	"syscall"
)

type App struct {
	config     *config
	bot        *bot.Bot
	shutdownCh chan os.Signal
}

func NewApp() (*App, error) {
	// app config
	cfg, err := loadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	// logger
	if cfg.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Logger = log.With().Caller().Logger()
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: pkg.TimeLayoutLOG})
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
	log.Debug().Msg(pkg.AnyPrint(cfg.AppName+" config", cfg))

	// ton service
	tonService, err := ton.NewTon(cfg.LiteServers)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to init ton service")
	}

	// redis
	redisClient, err := redis.NewRedisClient(cfg.Redis)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to init redis client")
	}

	// storage
	st, err := postgres.NewStorage(cfg.PG)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to init storage")
	}

	// telegram bot service
	botService, err := bot.NewBot(cfg.BotToken, cfg.BotAdminID, tonService, redisClient, st, cfg.Debug)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to init bot service")
	}
	log.Debug().Msgf("authorized on Telegram bot %s", botService.BotName)

	// exit channel
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	return &App{
		config:     cfg,
		bot:        botService,
		shutdownCh: quit,
	}, nil
}

func (app *App) Run() {
	app.bot.Start()
	log.Info().Msgf("%s started", app.config.AppName)

	<-app.shutdownCh
	log.Info().Msgf("waiting for all services to stop...")
	app.bot.Stop()
	log.Info().Msg("application stopped")
}
