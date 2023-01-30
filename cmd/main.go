package main

import (
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"tonflow/bot"
	"tonflow/config"
	"tonflow/storage/postgres"
	"tonflow/storage/redis"
	"tonflow/tonclient"
)

func main() {
	config.GetConfig()

	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, os.Interrupt, syscall.SIGTERM)

	// ton
	tonClient, err := tonclient.NewTonClient(config.Config.LiteServers)
	if err != nil {
		log.Fatalf("failed to init ton service: %v", err)
	}

	// redis
	redisClient, err := redis.NewRedisClient(config.Config.RedisURI)
	if err != nil {
		log.Fatalf("failed to init redis client: %v", err)
	}

	// storage
	storageClient, err := postgres.NewConnection(config.Config.PgURI)
	if err != nil {
		log.Fatalf("failed to init storage: %v", err)
	}

	// telegram bot
	botService, err := bot.NewBot(
		config.Config.BotToken,
		config.Config.BotAdminID,
		tonClient,
		redisClient,
		storageClient,
		config.Config.Debug,
		config.Config.BlockchainTxFee,
		config.Config.Key,
	)
	if err != nil {
		log.Fatalf("failed to init bot service: %v", err)
	}
	log.Debugf("Authorized on Telegram bot @%s", botService.BotName)

	botService.Start()
	log.Infof("%s bot started", config.Config.AppName)

	<-shutdownCh
	botService.Stop()

	log.Debugf("%s stopped", config.Config.AppName)
}
