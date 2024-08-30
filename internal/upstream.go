package internal

import (
	"pocket/pkg/db"

	"gorm.io/gorm"
)

func init() {
	db.DB.AutoMigrate(&AccountUpstream{})
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
