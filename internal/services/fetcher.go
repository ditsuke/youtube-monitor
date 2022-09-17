package services

import (
	"context"
	"github.com/ditsuke/youtube-focus/internal/yt"
	"github.com/rs/zerolog"
	"time"
)

type Fetcher struct {
	Logger   zerolog.Logger
	Client   *yt.Client
	Interval time.Duration
}

// Spawn kicks off the Fetcher service in a new goroutine. The context passed can be used for
// cancellation.
func (f *Fetcher) Spawn(ctx context.Context, query string, tx chan<- []yt.Video) {
	go func() {
		f.Start(ctx, query, tx)
	}()
}

// Start is like Spawn, but blocks the calling goroutine.
func (f *Fetcher) Start(ctx context.Context, query string, tx chan<- []yt.Video) {
	ticker := time.NewTicker(f.Interval)

	for {
		// spawn a goroutine to fetch and send records across the channel
		go f.fetchAndSend(query, tx)

		// block until it's time for the next batch query or the context expires
		select {
		case <-ticker.C:
			continue
		case <-ctx.Done():
			f.Logger.Debug().Str("reason", "context cancellation").Msg("stopping fetch service")
			break
		}
		<-ticker.C
	}
}

func (f *Fetcher) fetchAndSend(query string, tx chan<- []yt.Video) {
	r, err := f.Client.QueryLatestVideos(query)
	if err != nil {
		f.Logger.Warn().AnErr("query videos", err).Msg("")
		return
	}
	tx <- r
}
