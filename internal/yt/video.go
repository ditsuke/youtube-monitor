package yt

import (
	"gorm.io/gorm"
	"time"
)

type Video struct {
	gorm.Model
	VideoId      string `gorm:"unique;not null"`
	Title        string
	Description  string
	PublishedAt  time.Time
	ThumbnailUrl string
}

func (Video) TableName() string {
	return "videos"
}
