package app

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"ton-flow-bot/internal/storage"
	"ton-flow-bot/internal/storage/postgres"
	"ton-flow-bot/pkg"
)

type App struct {
	config  *config
	bot     *tgbotapi.BotAPI
	storage storage.Storage
}

func NewApp() (*App, error) {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	bot, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to init BotAPI instance config")
	}

	if cfg.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Logger = log.With().Caller().Logger()
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: pkg.TimeLayoutLOG})
		bot.Debug = true
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		bot.Debug = false
	}
	log.Debug().Msg(pkg.AnyPrint(cfg.AppName+" config", cfg))

	log.Debug().Msgf("Authorized on account %s", bot.Self.UserName)

	st, err := postgres.NewPGStorage(cfg.PG)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to init storage")
	}

	return &App{
		config:  cfg,
		bot:     bot,
		storage: st,
	}, nil
}

func (app *App) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := app.bot.GetUpdatesChan(u)

	for update := range updates {
		app.handleUpdate(update)
	}

	app.storage.Close()
}
