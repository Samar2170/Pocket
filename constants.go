package main

import (
	"os"

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
		panic(err)
	}
	basedir = os.Getenv("BASEDIR")
	hostname = os.Getenv("HOSTNAME")
	BotToken = os.Getenv("POCKET_BOTTOKEN")
	Host = os.Getenv("HOST")
	Port = os.Getenv("PORT")

}
