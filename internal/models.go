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

type Account struct {
	*gorm.Model
	Username string
	ID       int64
}

func GetAccountByUsername(username string) (Account, error) {
	var err error
	var account Account
	err = db.DB.Where("username = ?", username).First(&account).Error
	return account, err
}

func LoadAccounts() error {
	var err error
	file, err := os.Open("accountmap.json")
	if err != nil {
		log.Fatalf("failed opening file: for loading accounts: %s", err)
	}
	defer file.Close()
	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("failed reading file: for loading accounts: %s", err)
	}
	var accountMap map[string]int64
	json.Unmarshal(byteValue, &accountMap)
	for k, v := range accountMap {
		err := db.DB.FirstOrCreate(&Account{Username: k, ID: v}).Error
		if err != nil {
			log.Println(err.Error() + "occurred" + k)
		}
		log.Println(k, v)
	}
	return err
}

type TextContent struct {
	*gorm.Model
	ID       int64
	Text     string `gorm:"type:text;not null"`
	Username string

	AccountID int64
	Account   Account `gorm:"foreignKey:AccountID"`
}

type ImageContent struct {
	*gorm.Model
	ID       int64
	ImageURL string `gorm:"type:text;not null"`
	Username string
	Caption  string

	AccountID int64
	Account   Account `gorm:"foreignKey:AccountID"`
}
