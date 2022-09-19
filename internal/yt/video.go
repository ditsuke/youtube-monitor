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

func (Video) TableName() string {
	return "videos"
}

type VideoFull struct {
	gorm.Model
	Video
	// `tsvector` for postgres native full-text search
	// @todo include the description column, perhaps with a reduced weight.
	TSV string `gorm:"->;type:tsvector GENERATED ALWAYS AS (to_tsvector('english', title)) STORED;default:(-)"`
}

func (VideoFull) TableName() string {
	return "videos"
}
