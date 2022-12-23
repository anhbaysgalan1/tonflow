package bot

import (
	"context"
	tgBotAPI "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"sync"
	"tonflow/internal/service/ton"
	"tonflow/internal/storage"
)

type Bot struct {
	BotName   string
	adminID   int64
	api       *tgBotAPI.BotAPI
	ton       *ton.Ton
	redis     storage.TemporaryStorage
	storage   storage.Storage
	cryptoKey string
	stopCh    chan struct{}
}

func NewBot(token string, admin int64, ton *ton.Ton, redisClient storage.TemporaryStorage, storage storage.Storage, debug bool, cryptoKey string) (*Bot, error) {
	botAPI, err := tgBotAPI.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	botAPI.Debug = debug

	l := logger{}
	err = tgBotAPI.SetLogger(l)
	if err != nil {
		return nil, err
	}

	return &Bot{
		BotName:   botAPI.Self.UserName,
		adminID:   admin,
		api:       botAPI,
		ton:       ton,
		redis:     redisClient,
		storage:   storage,
		cryptoKey: cryptoKey,
		stopCh:    make(chan struct{}),
	}, nil

}

func (bot *Bot) Start() {
	u := tgBotAPI.NewUpdate(0)
	u.Timeout = 60
	updates := bot.api.GetUpdatesChan(u)

	var wg sync.WaitGroup

	go func() {
		for update := range updates {
			wg.Add(1)
			up := update
			go func() {
				defer wg.Done()
				bot.handleUpdate(context.Background(), up)
			}()
		}
	}()

	go func() {
		select {
		case <-bot.stopCh:
			wg.Wait()
			bot.stopCh <- struct{}{}
		}
	}()

}

func (bot *Bot) Stop() {
	bot.stopCh <- struct{}{}
	<-bot.stopCh
}
