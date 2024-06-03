package main

import (
	"ai-assistant-api/internal/config"
	"ai-assistant-api/internal/http-server/app"
	"ai-assistant-api/internal/utils/flags"
	"ai-assistant-api/internal/utils/logger/sl"
	"errors"

	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const (
	envLocal = "local"
	envProd  = "prod"
)

func main() {
	flags := flags.MustParseFlags()

	cfg := config.MustLoad(flags.ConfigPath)

	log := setupLogger(cfg.Env)
	log.Info("starting ai-assistant-api", slog.String("env", cfg.Env))

	srv, err := app.New(cfg, log)
	if err != nil {
		log.Error("failed to craete app", sl.Err(err))
		os.Exit(1)
	}
	log.Info("server created successfully")

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := srv.Run()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("failed to start server")
		}
	}()

	log.Info("starting server", slog.String("address", cfg.Address))

	<-done
	log.Info("stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to stop server", sl.Err(err))
		return
	}

	log.Info("server stopped")
}

func setupLogger(env string) (log *slog.Logger) {
	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return
}
