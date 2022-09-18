package services

import (
	"context"
	"github.com/ditsuke/youtube-focus/internal/interfaces"
	"github.com/rs/zerolog"
	"time"
)

type Persister[T any] struct {
	Logger zerolog.Logger
	Store  interfaces.Store[T, time.Time]
}

// Spawn kicks off the Persister service in a new goroutine.
// Context expiration can be used to stop the service.
func (p *Persister[T]) Spawn(ctx context.Context, rx <-chan []T) {
	go p.Start(ctx, rx)
}

// Start is like Spawn but blocks the calling goroutine.
func (p *Persister[T]) Start(ctx context.Context, rx <-chan []T) {
	for {
		select {
		case records := <-rx:
			p.Store.Save(records)
		case <-ctx.Done():
		}
	}
}
