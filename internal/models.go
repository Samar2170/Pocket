package internal

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"pocket/pkg/db"

	"gorm.io/gorm"
)

func init() {
	db.DB.AutoMigrate(&TextContent{}, &ImageContent{})
}

type TextContent struct {
	*gorm.Model
	ID      int64
	Text    string `gorm:"type:text;not null"`
	Account string
}

func getAccountID(account string) int64 {
	jsonFile, err := os.Open("accountmap.json")
	if err != nil {
		log.Fatalf("Error when opening file: %s", err)
	}
	defer jsonFile.Close()
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatalf("Error reading file: %s", err)
	}
	var accountMap map[string]int64
	json.Unmarshal(byteValue, &accountMap)
	return accountMap[account]
}

type ImageContent struct {
	*gorm.Model
	ID       int64
	ImageURL string `gorm:"type:text;not null"`
	Account  string
	Caption  string
}
