package models

import (
	"pocket/pkg/db"

	"gorm.io/gorm"
)

func init() {
	db.DB.AutoMigrate(&Account{})
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

type AccountUpstream struct {
	*gorm.Model
	ID        int64
	Name      string
	BaseURL   string
	SubURLMap string

	Account   Account `gorm:"foreignKey:AccountID"`
	AccountID int64
}

func (a *AccountUpstream) TableName() string {
	return "account_upstream"
}

type Upstream struct {
	ServiceName string            `json:"servicename"`
	BaseURL     string            `json:"baseurl"`
	SubURLMap   map[string]string `json:"suburl_map"`
}
