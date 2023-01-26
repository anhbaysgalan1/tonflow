package bot

import (
	"context"
	tgBotAPI "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
	"sync"
	"tonflow/storage"
	"tonflow/tonclient"
)

type Bot struct {
	BotName   string
	adminID   int64
	api       *tgBotAPI.BotAPI
	ton       *tonclient.TonClient
	redis     storage.Cache
	storage   storage.Storage
	cryptoKey string
	stopCh    chan struct{}
}

type logger struct {
}

func (l logger) Println(v ...interface{}) {
	log.Debugln(v...)
}

func (l logger) Printf(format string, v ...interface{}) {
	log.Debugf(format, v...)
}

func NewBot(
	token string,
	admin int64,
	ton *tonclient.TonClient,
	redisClient storage.Cache,
	storage storage.Storage,
	debug bool,
	cryptoKey string,
) (*Bot, error) {
	botAPI, err := tgBotAPI.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	botAPI.Debug = debug

	err = tgBotAPI.SetLogger(logger{})
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
	wg := new(sync.WaitGroup)

	go func() {
		for update := range updates {
			up := update
			go func() {
				wg.Add(1)
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
