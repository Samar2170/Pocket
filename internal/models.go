package internal

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"pocket/pkg/db"
	"sync"

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
	byteValue, err := io.ReadAll(file)
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
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		LoadAccountUpstreams()
		wg.Done()
	}()
	wg.Wait()
	return nil
}

type Upstream struct {
	ServiceName string            `json:"servicename"`
	BaseURL     string            `json:"baseurl"`
	SubURLMap   map[string]string `json:"suburl_map"`
}

func LoadAccountUpstreams() error {
	var err error
	var results []struct {
		Username              string
		ID                    int64
		AccountUpstreamExists bool
	}
	err = db.DB.Table("accounts").
		Select("accounts.username as username, accounts.id as id, exists(select 1 from account_upstream where account_upstream.account_id = accounts.id) as account_upstream_exists").
		Scan(&results).Error
	if err != nil {
		return err
	}
	accountUpstreamFile, err := os.Open("accountus.json")
	if err != nil {
		return err
	}
	var upstreams map[string]interface{}
	defer accountUpstreamFile.Close()
	decoder := json.NewDecoder(accountUpstreamFile)
	err = decoder.Decode(&upstreams)
	if err != nil {
		return err
	}
	for _, result := range results {
		if !result.AccountUpstreamExists {
			us := upstreams[result.Username]
			if us == nil {
				continue
			}
			usJson, err := json.Marshal(us)
			if err != nil {
				log.Println(err)
				continue
			}
			var usMap Upstream
			err = json.Unmarshal(usJson, &usMap)
			if err != nil {
				log.Println(err)
				continue
			}
			subUrlMap, err := json.Marshal(usMap.SubURLMap)
			if err != nil {
				log.Println(err)
				continue
			}
			err = db.DB.Create(&AccountUpstream{
				AccountID: result.ID,
				Name:      usMap.ServiceName,
				BaseURL:   usMap.BaseURL,
				SubURLMap: string(subUrlMap),
			}).Error
			if err != nil {
				log.Println(err)
				continue
			}
		}
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

func (t *TextContent) TableName() string {
	return "text_content"
}

func (i *ImageContent) TableName() string {
	return "image_content"
}
