package services

import (
	"context"
	"github.com/ditsuke/youtube-focus/internal/yt"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Persister struct {
	Logger zerolog.Logger
	DB     *gorm.DB
}

// Spawn kicks off the Persister service in a new goroutine.
// Context expiration can be used to stop the service.
func (p *Persister) Spawn(ctx context.Context, rx <-chan []yt.Video) {
	go p.Start(ctx, rx)
}

// Start is like Spawn but blocks the calling goroutine.
func (p *Persister) Start(ctx context.Context, rx <-chan []yt.Video) {
	for {
		select {
		case records := <-rx:
			p.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(records)
		case <-ctx.Done():
		}
	}
}
