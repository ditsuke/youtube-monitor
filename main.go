package main

import (
	"context"
	"github.com/ditsuke/youtube-focus/config"
	"github.com/ditsuke/youtube-focus/internal/services"
	"github.com/ditsuke/youtube-focus/internal/yt"
	"github.com/ditsuke/youtube-focus/store"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog"
	"github.com/sethvargo/go-envconfig"
	"time"
)

// Args: query, fetch_interval
// Env: YOUTUBE_API_KEY
func main() {
	log := zerolog.New(zerolog.NewConsoleWriter())

	cfg := config.Config{}
	if err := envconfig.Process(context.Background(), &cfg); err != nil {
		log.Fatal().Err(err).Msg("config from environment")
	}

	client, err := yt.New(
		context.Background(), cfg.YouTubeAPIKey,
		yt.WithLogger(log.With().Str("comp", "yt").Logger()),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("youtube client initialization")
	}

	db := store.GetDB(store.GetDSNFromConfig(cfg))

	c := make(chan []yt.Video)

	fetcher := services.Fetcher[yt.Video]{
		Logger:   log.With().Str("comp", "fetcher").Logger(),
		Interval: time.Second * 15,
		FetchFunc: func() ([]yt.Video, error) {
			return client.QueryLatestVideos("minecraft")
		},
	}

	persister := services.Persister[yt.Video]{
		Logger: log.With().Str("comp", "persister").Logger(),
		Store:  &store.VideoMetaStore{DB: db},
	}

	ctx := context.Background()

	fetcher.Spawn(ctx, c)
	persister.Spawn(ctx, c)

	// let our services run for a couple minutes to test things out
	time.Sleep(time.Minute * 30)
}
