package config

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	log "github.com/sirupsen/logrus"
	"tonflow/pkg"
)

const (
	timeLayout = "02.01.06 15:04:05"

	mainnetLiteservers = "https://ton-blockchain.github.io/global.config.json"
	testnetLiteservers = "https://ton-blockchain.github.io/testnet-global.config.json"
)

var Config = struct {
	Production   bool   `env:"PRODUCTION" envDefault:"false"`
	Debug        bool   `env:"DEBUG" envDefault:"true"`
	AppName      string `env:"APP_NAME" envDefault:"Tonflow"`
	BotToken     string
	ProdBotToken string `env:"PROD_BOT_TOKEN,required"`
	DevBotToken  string `env:"DEV_BOT_TOKEN,required"`
	BotAdminID   int64  `env:"BOT_ADMIN_ID,required"`
	LiteServers  string
	Key          string `env:"KEY,required"`
	RedisHost    string `env:"REDIS_HOST,required"`
	RedisPort    string `env:"REDIS_PORT,required"`
	RedisURI     string
	PgHost       string `env:"PG_HOST,required"`
	PgPort       string `env:"PG_PORT,required"`
	PgUser       string `env:"PG_USER,required"`
	PgPassword   string `env:"PG_PASSWORD,required"`
	PgName       string `env:"PG_NAME,required"`
	PgSSL        string `env:"PG_SSL,required"`
	PgURI        string
	PgMigration  bool `env:"PG_MIGRATION,required"`
}{
	LiteServers: testnetLiteservers,
}

func GetConfig() {
	err := env.Parse(&Config)
	if err != nil {
		log.Fatal(err)
	}

	if Config.Production {
		log.SetFormatter(&log.JSONFormatter{})
		Config.BotToken = Config.ProdBotToken
		Config.LiteServers = mainnetLiteservers
	} else {
		Config.BotToken = Config.DevBotToken
		log.SetFormatter(&log.TextFormatter{
			ForceColors:            true,
			DisableLevelTruncation: true,
			PadLevelText:           true,
			FullTimestamp:          true,
			TimestampFormat:        timeLayout,
			QuoteEmptyFields:       true,
		})
		log.SetReportCaller(false)
	}

	Config.RedisURI = fmt.Sprintf("%s:%s", Config.RedisHost, Config.RedisPort)

	Config.PgURI = fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		Config.PgUser, Config.PgPassword, Config.PgHost, Config.PgPort, Config.PgName, Config.PgSSL,
	)

	if Config.Debug {
		log.SetLevel(log.DebugLevel)
	}

	log.Debugf("Got config:\n%v", pkg.AnyPrint(Config))
}
