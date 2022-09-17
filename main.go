package main

import (
	"context"
	"flag"
	"github.com/ditsuke/youtube-focus/internal/services"
	"github.com/ditsuke/youtube-focus/internal/store"
	"github.com/ditsuke/youtube-focus/internal/yt"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog"
	"os"
	"time"
)

// Args: query, fetch_interval
// Env: YOUTUBE_API_KEY
func main() {
	flag.Bool("debug", false, "enable debug logging")

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

	fetcher := services.Fetcher{
		Logger:   log.With().Str("comp", "fetcher").Logger(),
		Interval: time.Second * 15,
		Client:   client,
	}

	persister := services.Persister{
		Logger: log.With().Str("comp", "persister").Logger(),
		DB:     db,
	}

	ctx := context.Background()

	fetcher.Spawn(ctx, "minecraft", c)
	persister.Spawn(ctx, c)

	// let our services run for a couple minutes to test things out
	time.Sleep(time.Minute * 2)
}
