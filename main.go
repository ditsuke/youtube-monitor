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
	"github.com/sethvargo/go-envconfig"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	// key attached to service sub-loggers
	service = "service"

	// Graceful shutdown timeout
	shutdownTimeoutSeconds = 15
)

type superCtx struct {
	ctx    context.Context
	store  *store.VideoMetaStore
	logger *zerolog.Logger
	cfg    config.Config
}

func main() {
	logger := zerolog.New(zerolog.NewConsoleWriter())

	cfg := config.Config{}
	if err := envconfig.Process(context.Background(), &cfg); err != nil {
		logger.Fatal().Err(err).Msg("config from environment")
	}
	logger.Info().Msg(fmt.Sprintf("config=%+v", cfg))

	db, err := cfg.GetDB()
	if err != nil {
		logger.Fatal().Err(err).Str("operation", "db-connect").Msg("failed")
	}
	videoStore := &store.VideoMetaStore{DB: db}

	ctx, ctxCancel := context.WithCancel(context.Background())
	defer ctxCancel()

	spawnBackgroundServices(superCtx{
		ctx:    ctx,
		logger: &logger,
		cfg:    cfg,
		store:  videoStore,
	})

	server := api.Server{
		Cfg:    cfg,
		Logger: logger.With().Str(service, "api-server").Logger(),
	}

	logger.Info().Msg("starting server...")
	go server.StartServer(ctx)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	// Block until receive os signal
	<-sig

	// done := make(chan struct{})
	doneCtx, doneDoneCtx := context.WithTimeout(context.Background(), shutdownTimeoutSeconds)
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
func spawnBackgroundServices(s superCtx) {
	c := make(chan []yt.Video)
	ytClient, err := yt.New(s.cfg.YouTubeAPIKeys,
		yt.WithLogger(s.logger.With().Str("comp", "yt").Logger()),
	)
	if err != nil {
		s.logger.Fatal().Err(err).Msg("youtube client initialization")
	}

	fetcher := services.Fetcher[yt.Video]{
		Logger:   s.logger.With().Str(service, "video-fetcher").Logger(),
		Interval: time.Duration(s.cfg.YouTubePollInterval) * time.Second,
		FetchFunc: func() ([]yt.Video, error) {
			return ytClient.QueryLatestVideos(s.cfg.YouTubeVideoQuery)
		},
	}

	persister := services.Persister[yt.Video]{
		Logger: s.logger.With().Str("comp", "persister").Logger(),
		Store:  s.store,
	}

	fetcher.Spawn(s.ctx, c)
	persister.Spawn(s.ctx, c)
}
