package main

import (
	"context"
	"fmt"
	"github.com/ditsuke/youtube-focus/api"
	"github.com/ditsuke/youtube-focus/config"
	"github.com/ditsuke/youtube-focus/internal/services"
	"github.com/ditsuke/youtube-focus/internal/yt"
	"github.com/ditsuke/youtube-focus/store"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sethvargo/go-envconfig"
	"os"
	"os/signal"
	"time"
)

const Service = "service"

// Args: query, fetch_interval
// Env: YOUTUBE_API_KEY
func main() {
	logger := zerolog.New(zerolog.NewConsoleWriter())

	cfg := config.Config{}
	if err := envconfig.Process(context.Background(), &cfg); err != nil {
		logger.Fatal().Err(err).Msg("config from environment")
	}
	logger.Info().Msg(fmt.Sprintf("config=%+v", cfg))

	ytClient, err := yt.New(cfg.YouTubeAPIKeys,
		yt.WithLogger(logger.With().Str("comp", "yt").Logger()),
	)
	if err != nil {
		logger.Fatal().Err(err).Msg("youtube client initialization")
	}

	videoStore := &store.VideoMetaStore{DB: store.GetDB(store.GetDSNFromConfig(cfg))}

	ctx, ctxCancel := context.WithCancel(context.Background())
	defer ctxCancel()

	spawnBackgroundServices(ctx, ytClient, videoStore, cfg)

	server := api.Server{
		Cfg:    cfg,
		Logger: logger.With().Str(Service, "api-server").Logger(),
	}

	logger.Info().Msg("starting server...")
	go server.StartServer(ctx)

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)

	// Block until receive os signal
	<-sig

	// done := make(chan struct{})
	doneCtx, doneDoneCtx := context.WithTimeout(context.Background(), time.Second*15)
	defer doneDoneCtx()
	go func() {
		logger.Info().Msg("attempting graceful shutdown")
		ctxCancel()
	}()

	<-doneCtx.Done()
	if doneCtx.Err() == context.DeadlineExceeded {
		logger.Fatal().Msg("graceful shutdown timed-out. dying...")
	}
}

// spawnBackgroundServices spawns services to fetch and store the latest videos from YouTube.
func spawnBackgroundServices(ctx context.Context,
	ytClient *yt.Client,
	videoStore *store.VideoMetaStore,
	cfg config.Config) {
	c := make(chan []yt.Video)

	fetcher := services.Fetcher[yt.Video]{
		Logger:   log.With().Str(Service, "video-fetcher").Logger(),
		Interval: time.Duration(cfg.YouTubePollInterval) * time.Second,
		FetchFunc: func() ([]yt.Video, error) {
			return ytClient.QueryLatestVideos(cfg.YouTubeVideoQuery)
		},
	}

	persister := services.Persister[yt.Video]{
		Logger: log.With().Str("comp", "persister").Logger(),
		Store:  videoStore,
	}

	fetcher.Spawn(ctx, c)
	persister.Spawn(ctx, c)
}
