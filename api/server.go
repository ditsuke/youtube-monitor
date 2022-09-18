package api

import (
	"context"
	"github.com/ditsuke/youtube-focus/api/handlers"
	"github.com/ditsuke/youtube-focus/config"
	"github.com/ditsuke/youtube-focus/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/ironstar-io/chizerolog"
	"github.com/rs/zerolog"
	"net"
	"net/http"
	"time"
)

const GracefulShutdownDeadline = 15 * time.Second

type Server struct {
	Cfg    config.Config
	Logger zerolog.Logger
}

func (s *Server) StartServer(ctx context.Context) {
	r := chi.NewRouter()
	routeLogger := s.Logger.With().Str("part", "router").Logger()

	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(chizerolog.LoggerMiddleware(&routeLogger))
	RegisterRoutes(s.Cfg, r)

	server := http.Server{
		Addr:    net.JoinHostPort(s.Cfg.ServerHost, s.Cfg.ServerPort),
		Handler: r,
	}

	// Spawn the server in a new goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.Logger.Fatal().Err(err).Timestamp().Msg("http server crashed")
		}
	}()

	// Block until the context expires
	<-ctx.Done()
	shutdownCtx, shutdownCtxCancel := context.WithTimeout(context.Background(),
		GracefulShutdownDeadline)
	defer shutdownCtxCancel()

	// Attempt to gracefully shut the server down
	if err := server.Shutdown(shutdownCtx); err != nil {
		s.Logger.Fatal().Err(err).Msg("attempted server shutdown")
	}
}

func RegisterRoutes(cfg config.Config, m *chi.Mux) {
	videoSvc := handlers.New(store.VideoMetaStore{DB: store.GetDB(store.GetDSNFromConfig(cfg))})

	m.Get("/videos", videoSvc.Get)
	m.Get("/videos/search", videoSvc.Search)
}
