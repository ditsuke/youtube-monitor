package main

import (
	"context"
	"flag"
	"github.com/ditsuke/youtube-focus/internal/services"
	"github.com/ditsuke/youtube-focus/internal/yt"
	"github.com/ditsuke/youtube-focus/store"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog"
	"os"
	"time"
)

// Args: query, fetch_interval
// Env: YOUTUBE_API_KEY
func main() {
	flag.Parse()

	log := zerolog.New(zerolog.NewConsoleWriter())

	apiKey, ok := os.LookupEnv("YOUTUBE_API_KEY")
	if !ok {
		log.Debug().Msg("no api key in env")
	}

	client, err := yt.New(
		context.Background(), apiKey, yt.WithLogger(log.With().Str("comp", "yt").Logger()),
	)
	if err != nil {
		panic(err)
	}

	db := store.GetDB(store.GetDSN("dt", "dtp", "localhost", "some_db"))

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
	time.Sleep(time.Minute * 2)
}
