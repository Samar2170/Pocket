package main

import (
	"os"
	"pocket/pkg/auditlog"
	"pocket/pkg/utils"

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
	var err error
	envPath := utils.Basedir + "/.env"
	err = godotenv.Load(envPath)
	if err != nil {
		auditlog.Errorlogger.Error().Str("error", err.Error()).Msg("Error loading .env file")
	}
	FXBDir = os.Getenv("FXBDIR")
	hostname = os.Getenv("HOSTNAME")
	BotToken = os.Getenv("POCKET_BOTTOKEN")
	Host = os.Getenv("HOST")
	Port = os.Getenv("PORT")

}
