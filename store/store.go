package store

import (
	"fmt"
	"github.com/ditsuke/youtube-focus/internal/interfaces"
	"github.com/ditsuke/youtube-focus/internal/yt"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
	"time"
)

// VideoMetaStore is an abstraction layer for the video meta storage.
// Includes methods to save, retrieve and search records (with pagination capabilities).
type VideoMetaStore struct {
	Logger zerolog.Logger
	DB     *gorm.DB
}

// interface compliance constraint for VideoMetaStore
var _ interfaces.Store[yt.Video, time.Time] = &VideoMetaStore{}

const OrderReverseChrono = "published_at DESC"

// Save records to the video store.
func (v *VideoMetaStore) Save(records []yt.Video) {
	v.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(records)
}

// Retrieve a maximum of limit videos published after some time.Time in reverse-chronological
// order (ie: sorted by latest)
// The publishedBefore param can be used for pagination -- by using the published_at
// attribute of the last record in a result, get the next batch.
func (v *VideoMetaStore) Retrieve(publishedBefore time.Time, limit int) []yt.Video {
	videos := new([]yt.Video)
	result := v.DB.
		Order(OrderReverseChrono).
		Limit(limit).
		Find(videos, "published_at < ?", publishedBefore)

	if result.Error != nil {
		fmt.Println("error: ", result.Error)
	}

	if videos == nil {
		return []yt.Video{}
	}

	return *videos
}

// Search videos in the store by title and description. Retrieves a maximum of limit videos
// published before some time.Time, sorted by latest first (reverse chronological)
// The publishedBefore param can be used for pagination -- by using the published_at
// attribute of the last record in a result, get the next batch.
func (v *VideoMetaStore) Search(query string, publishedBefore time.Time, limit int) []yt.Video {
	videos := new([]yt.Video)
	result := v.DB.
		Order(OrderReverseChrono).
		Limit(limit).
		Where(v.DB.Where("LOWER(title) LIKE LOWER(?)",
			"%"+query+"%").Or("LOWER(description) LIKE LOWER(?)", "%"+query+"%")).
		Where("published_at <= ?", publishedBefore).
		Find(videos)

	// @todo could use an error return
	if result.Error != nil {
		v.Logger.Error().Err(result.Error).Msg("video query")
	}

	if videos == nil {
		return []yt.Video{}
	}

	return *videos
}

// NaturalSearch searches videos with a special natural-language aware operation, retrieving
// a maximum of limit videos. This method does not support pagination at the moment.
func (v *VideoMetaStore) NaturalSearch(query string, limit int) []yt.Video {
	nlQuery := strings.Join(strings.Split(query, " "), "|")

	videos := new([]yt.Video)
	result := v.DB.
		Order(OrderReverseChrono).
		Limit(limit).
		Where("tsv @@ to_tsquery('english', ?)", nlQuery).
		Find(videos)

	if result.Error != nil {
		v.Logger.Error().Err(result.Error).Msg("natural language query for videos")
	}

	if videos == nil {
		return []yt.Video{}
	}

	return *videos
}
