package app

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strconv"
	"ton-flow-bot/internal/storage/postgres"
)

type config struct {
	Env        string
	EnvFile    string
	AppName    string
	BotToken   string
	BotAdminID int64
	Debug      bool

	WorkMode    string
	LiteServers string
	Wallet      string

	PG *postgres.Config
}

func getEnvString(key string) (string, error) {
	value, ok := os.LookupEnv(key)
	if !ok {
		return "", fmt.Errorf("failed to get %s value", key)
	}
	return value, nil
}

func getEnvInt64(key string) (int64, error) {
	value, ok := os.LookupEnv(key)
	if !ok {
		return 0, fmt.Errorf("failed to get %s value", key)
	}

	result, err := strconv.ParseInt(value, 10, 64)
	if err == nil {
		return 0, fmt.Errorf("failed to parse int: %w", err)
	}
	return result, nil
}

func getEnvBool(key string) (bool, error) {
	value, ok := os.LookupEnv(key)
	if !ok {
		return false, fmt.Errorf("failed to get %s value", key)
	}
	switch value {
	case "true":
		return true, nil
	case "false":
		return false, nil
	default:
		return false, nil
	}
}

func loadConfig() (*config, error) {
	env := os.Getenv("APP_ENV")
	var envFile string
	if "" == env {
		env = "local"
		envFile = ".env." + env
	}
	if "production" == env {
		envFile = ".env"
	}
	err := godotenv.Load(envFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load env file: %s", err)
	}

	// App vars
	appName, err := getEnvString("APP_NAME")
	if err != nil {
		return nil, err
	}
	botToken, err := getEnvString("BOT_TOKEN")
	if err != nil {
		return nil, err
	}
	botAdminID, err := getEnvInt64("BOT_ADMIN_ID")
	if err != nil {
		return nil, err
	}
	debug, err := getEnvBool("APP_DEBUG")
	if err != nil {
		return nil, err
	}

	// TON vars
	workMode, err := getEnvString("WORK_MODE")
	if err != nil {
		return nil, err
	}

	var liteServers, wallet string
	if workMode == "mainnet" {
		liteServers, err = getEnvString("MAINNET_LITESERVES")
		if err != nil {
			return nil, err
		}
		wallet, err = getEnvString("MAINNET_WALLET")
		if err != nil {
			return nil, err
		}
	} else {
		liteServers, err = getEnvString("TESTNET_LITESERVERS")
		if err != nil {
			return nil, err
		}
		wallet, err = getEnvString("TESTNET_WALLET")
		if err != nil {
			return nil, err
		}
	}

	// PG vars
	pgHost, err := getEnvString("PG_HOST")
	if err != nil {
		return nil, err
	}
	pgPort, err := getEnvString("PG_PORT")
	if err != nil {
		return nil, err
	}
	pgUser, err := getEnvString("PG_USER")
	if err != nil {
		return nil, err
	}
	pgPassword, err := getEnvString("PG_PASSWORD")
	if err != nil {
		return nil, err
	}
	pgName, err := getEnvString("PG_NAME")
	if err != nil {
		return nil, err
	}
	pgSSL, err := getEnvString("PG_SSL")
	if err != nil {
		return nil, err
	}
	pgMigration, err := getEnvBool("PG_MIGRATION")
	if err != nil {
		return nil, err
	}

	return &config{
		Env:         env,
		EnvFile:     envFile,
		AppName:     appName,
		BotToken:    botToken,
		BotAdminID:  botAdminID,
		Debug:       debug,
		WorkMode:    workMode,
		LiteServers: liteServers,
		Wallet:      wallet,
		PG: &postgres.Config{
			Host:      pgHost,
			Port:      pgPort,
			User:      pgUser,
			Password:  pgPassword,
			Name:      pgName,
			SSL:       pgSSL,
			Migration: pgMigration,
		},
	}, nil
}
