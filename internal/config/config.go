package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	GRPCPort   int
	AppID      int
	AppHash    string
	LogLevel   string
	SessionDir string
}

func Load() *Config {
	godotenv.Load()

	port := 50051
	if p := os.Getenv("GRPC_PORT"); p != "" {
		if i, err := strconv.Atoi(p); err == nil {
			port = i
		}
	}

	appID, _ := strconv.Atoi(os.Getenv("TELEGRAM_APP_ID"))
	appHash := os.Getenv("TELEGRAM_APP_HASH")

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	sessionDir := os.Getenv("SESSION_DIR")
	if sessionDir == "" {
		sessionDir = "./sessions"
	}

	return &Config{
		GRPCPort:   port,
		AppID:      appID,
		AppHash:    appHash,
		LogLevel:   logLevel,
		SessionDir: sessionDir,
	}
}
