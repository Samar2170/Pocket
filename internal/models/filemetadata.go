package models

import (
	"pocket/pkg/db"
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

func GetFileMetaDataById(id string) (FileMetaData, error) {
	var fmd FileMetaData
	err := db.DB.Where("id = ?", id).First(&fmd).Error
	return fmd, err
}
