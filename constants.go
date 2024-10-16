package main

import (
	"os"
	"pocket/pkg/auditlog"

	"github.com/joho/godotenv"
)

var (
	hostname string
	FXBDir   string
	BotToken string
	Host     string
	Port     string
)

func init() {
	envPath := ".env"
	err := godotenv.Load(envPath)
	if err != nil {
		auditlog.Errorlogger.Error().Str("error", err.Error()).Msg("Error loading .env file")
	}
	FXBDir = os.Getenv("FXBDIR")
	hostname = os.Getenv("HOSTNAME")
	BotToken = os.Getenv("POCKET_BOTTOKEN")
	Host = os.Getenv("HOSTNAME")
	Port = os.Getenv("PORT")
}
