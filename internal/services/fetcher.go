package services

import (
	"context"
	"github.com/rs/zerolog"
	"time"
)

type Fetcher[T any] struct {
	Logger    zerolog.Logger
	FetchFunc func() ([]T, error)
	Interval  time.Duration
}

// Spawn kicks off the Fetcher service in a new goroutine. The context passed can be used for
// cancellation.
func (f *Fetcher[T]) Spawn(ctx context.Context, tx chan<- []T) {
	go func() {
		f.Start(ctx, tx)
	}()
}

// Start is like Spawn, but blocks the calling goroutine.
func (f *Fetcher[T]) Start(ctx context.Context, tx chan<- []T) {
	ticker := time.NewTicker(f.Interval)

	for {
		f.Logger.Debug().Msg("fetching...")
		// spawn a goroutine to fetch and send records across the channel
		go f.fetchAndSend(tx)

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

func (f *Fetcher[T]) fetchAndSend(tx chan<- []T) {
	r, err := f.FetchFunc()
	if err != nil {
		f.Logger.Warn().AnErr("fetch items", err).Msg("")
		return
	}
	tx <- r
}
