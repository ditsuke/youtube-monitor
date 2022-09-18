package yt

import (
	"gorm.io/gorm"
	"time"
)

type Video struct {
	VideoId      string `gorm:"unique;not null"`
	Title        string
	Description  string
	PublishedAt  time.Time
	ThumbnailUrl string
}

type VideoFull struct {
	gorm.Model
	Video
}

func (VideoFull) TableName() string {
	return "videos"
}
