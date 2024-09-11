package models

import (
	"pocket/pkg/db"
	"time"
)

func init() {
	db.DB.AutoMigrate(&Note{})
}

type Note struct {
	ID          string `gorm:"primaryKey"`
	NoteContent string `gorm:"type:text"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
