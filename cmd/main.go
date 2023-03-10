package main

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/xssnick/tonutils-go/tlb"
	"os"
	"os/signal"
	"syscall"
	"tonflow/blockchain"
	"tonflow/bot"
	. "tonflow/config"
	"tonflow/pkg"
	"tonflow/storage/postgres"
	"tonflow/storage/redis"
)

func main() {
	GetConfig()

	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, os.Interrupt, syscall.SIGTERM)

	// ton blockchain client
	blockchainClient, err := blockchain.NewClient(Config.LiteServers, Config.Production)
	if err != nil {
		log.Fatalf("failed to init ton service: %v", err)
	}

	// redis client
	redisClient, err := redis.NewRedisClient(Config.RedisURI)
	if err != nil {
		log.Fatalf("failed to init redis client: %v", err)
	}

	// storage client
	storageClient, err := postgres.NewConnection(Config.PgURI)
	if err != nil {
		log.Fatalf("failed to init storage: %v", err)
	}
	log.Debug("in memory addresses:\n", pkg.PrintAny(storageClient.GetInMemoryWallets()))

	// telegram bot
	botService, err := bot.NewBot(
		Config.BotToken,
		Config.BotAdminID,
		blockchainClient,
		redisClient,
		storageClient,
		Config.Debug,
		Config.BlockchainTxFee,
	)
	if err != nil {
		log.Fatalf("failed to init bot service: %v", err)
	}
	log.Debugf("Authorized on Telegram bot @%s", botService.BotName)

	botService.Start()
	log.Infof("%s bot started", Config.AppName)

	txCh := make(chan *tlb.Transaction)
	errCh := make(chan error)
	go blockchain.Scan(blockchainClient, storageClient, txCh, errCh)
	go func() {
		for {
			select {
			case tx := <-txCh:
				log.Debug("Transaction:", tx.String())
				botService.Notify(context.Background(), tx)
			case err = <-errCh:
				log.Error(err)
				go blockchain.Scan(blockchainClient, storageClient, txCh, errCh)
			}
		}
	}()

	<-shutdownCh
	botService.Stop()
	// scan.stop()
	log.Debugf("%s stopped", Config.AppName)
	os.Exit(0)
}
