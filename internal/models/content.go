package models

import "gorm.io/gorm"

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
