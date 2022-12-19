package bot

import (
	"context"
	goRedis "github.com/go-redis/redis/v9"
	telegramBotAPI "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"park-wallet/internal/service/ton"
	"park-wallet/internal/storage"
	"sync"
)

type Bot struct {
	BotName string
	adminID int64
	api     *telegramBotAPI.BotAPI
	ton     *ton.Ton
	redis   *goRedis.Client
	storage storage.Storage
	stopCh  chan struct{}
}

func NewBot(token string, admin int64, ton *ton.Ton, redisClient *goRedis.Client, storage storage.Storage, debug bool) (*Bot, error) {
	botAPI, err := telegramBotAPI.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	botAPI.Debug = debug

	l := logger{}
	err = telegramBotAPI.SetLogger(l)
	if err != nil {
		return nil, err
	}

	return &Bot{
		BotName: botAPI.Self.UserName,
		adminID: admin,
		api:     botAPI,
		ton:     ton,
		redis:   redisClient,
		storage: storage,
		stopCh:  make(chan struct{}),
	}, nil

}

func (bot *Bot) Start() {
	u := telegramBotAPI.NewUpdate(0)
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
