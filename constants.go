package main

import (
	"os"

	"github.com/joho/godotenv"
)

var basedir string
var hostname string

var BotToken string

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	basedir = os.Getenv("BASEDIR")
	hostname = os.Getenv("HOSTNAME")
	BotToken = os.Getenv("POCKET_BOTTOKEN")

}
