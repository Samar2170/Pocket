package models

import (
	"time"
)

type FileMetaData struct {
	ID          string `gorm:"primaryKey"`
	OgFileName  string
	NewFileName string `gorm:"unique"`
	FilePath    string
	Extension   string
	Category    string
	Size        int
	SizeInMB    float64

	CreatedAt time.Time
	UpdatedAt time.Time
}
