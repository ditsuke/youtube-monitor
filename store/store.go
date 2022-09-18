package store

import (
	"fmt"
	"github.com/ditsuke/youtube-focus/config"
	"github.com/ditsuke/youtube-focus/internal/interfaces"
	"github.com/ditsuke/youtube-focus/internal/yt"
	"github.com/ditsuke/youtube-focus/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

func GetDSN(user, password, host, db string) string {
	return fmt.Sprintf("user=%s password=%s host=%s dbname=%s", user, password, host, db)
}

func GetDSNFromConfig(cfg config.Config) string {
	return fmt.Sprintf(
		"user=%s password=%s port=%s host=%s dbname=%s",
		cfg.PostgresUser, cfg.PostgresPass,
		cfg.PostgresPort, cfg.PostgresHost,
		cfg.PostgresDB,
	)
}

func GetDB(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		// @todo log + propogate
		return nil
	}

	return db
}

type Store[T any] interface {
	Save([]T)
	Retrieve(limit int)
}

// VideoMetaStore is an abstraction layer for the video meta storage.
// Includes methods to save, retrieve and search records (with pagination capabilities).
type VideoMetaStore struct {
	DB *gorm.DB
}

// Save records to the video store.
func (v *VideoMetaStore) Save(records []model.Video) {
	v.DB.Create(records)
}

// Retrieve a maximum of limit videos published after some time.Time in reverse-chronological
// order (ie: sorted by latest)
// The publishedAfter param can be used for pagination -- by using the published_at
// attribute of the last record in a result, get the next page.
func (v *VideoMetaStore) Retrieve(publishedAfter time.Time, limit int) []model.Video {
	videos := new([]model.Video)
	result := v.DB.
		Order("published_at DESC").
		Limit(limit).
		Find(videos, "published_at >= ?", publishedAfter)

	if result.Error != nil {
		fmt.Println("error: ", result.Error)
	}

	return *videos
}

// Search videos in the store by title and description. Retrieves a maximum of limit videos
// published after some time.Time
// The publishedAfter param can be used for pagination -- by using the published_at
// attribute of the last record in a result, get the next page.
func (v *VideoMetaStore) Search(query string, publishedAfter time.Time, limit int) []model.Video {
	videos := new([]model.Video)
	result := v.DB.
		Order("published_at DESC").
		Limit(limit).
		Where("LOWER(title) LIKE LOWER(?)", "%"+query+"%").
		Or("LOWER(description) LIKE LOWER(?)", "%"+query+"%").
		Find(
			videos,
			"published_at >= ?", publishedAfter,
		)

	// @todo replace with log + error return
	if result.Error != nil {
		fmt.Println("error: ", result.Error)
	}

	return *videos
}
