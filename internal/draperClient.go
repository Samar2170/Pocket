package internal

import (
	"log"
	"os"
	"pocket/pkg/apiclient"

	"github.com/joho/godotenv"
)

var apiKey string

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println(err)
	}
	apiKey = os.Getenv("DRAPER_API_KEY")
}

type DraperClient struct {
	client *apiclient.Client
}

func NewDraperClient() *DraperClient {
	return &DraperClient{
		client: apiclient.NewClient("http://127.0.0.1:8000", map[string]string{
			"APIKEY": apiKey,
		}),
	}
}
