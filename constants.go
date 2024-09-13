package main

import (
	"os"
	"pocket/pkg/auditlog"

	"github.com/joho/godotenv"
)

var (
	basedir  string
	hostname string

	BotToken string
	Host     string
	Port     string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		auditlog.Errorlogger.Error().Str("error", err.Error()).Msg("Error loading .env file")
	}
	basedir = os.Getenv("BASEDIR")
	hostname = os.Getenv("HOSTNAME")
	BotToken = os.Getenv("POCKET_BOTTOKEN")
	Host = os.Getenv("HOST")
	Port = os.Getenv("PORT")

}
