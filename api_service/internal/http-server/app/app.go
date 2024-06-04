package app

import (
	"ai-assistant-api/internal/config"
	h "ai-assistant-api/internal/http-server/handlers"
	mwr "ai-assistant-api/internal/http-server/middleware"
	bg "ai-assistant-api/internal/service/background"
	kafka "ai-assistant-api/internal/service/kafka/producer"
	"ai-assistant-api/internal/service/task"
	sql "ai-assistant-api/internal/storage/sqlite"
	"ai-assistant-api/internal/utils/logger/sl"

	"context"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	backgroundService *bg.Service
	backgroundDone    chan struct{}

	server *http.Server
}

func New(cfg *config.Config, log *slog.Logger) (*Server, error) {
	// init deps
	storage, err := sql.New(cfg)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		return nil, err
	}

	kafkaService, err := kafka.New(log, cfg)
	if err != nil {
		log.Error("failed to create kafka service", sl.Err(err))
		return nil, err
	}

	taskService := task.New(log, storage, kafkaService)

	bgService := bg.New(cfg, log, taskService)

	// init router
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(mwr.HTTPMetrics)

	router.Post("/task", h.HandleCreateTask(log, taskService))
	router.Get("/task/{id}", h.HandleGetTask(log, taskService))

	router.Handle("/metrics", promhttp.Handler())

	return &Server{
		server: &http.Server{
			Addr:         cfg.Address,
			Handler:      router,
			ReadTimeout:  cfg.HTTPServer.Timeout,
			WriteTimeout: cfg.HTTPServer.Timeout,
			IdleTimeout:  cfg.HTTPServer.IdleTimeout,
		},
		backgroundService: bgService,
	}, nil
}

func (s *Server) Run() error {
	s.backgroundService.Start()
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.backgroundService.Shutdown()
	return s.server.Shutdown(ctx)
}
