package main

import (
	"async-email-worker/internal/config"
	"async-email-worker/internal/handler"
	"async-email-worker/internal/middleware"
	"async-email-worker/internal/queue"
	"async-email-worker/pkg/redis"
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func main() {
	// Setup JSON Logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Error("configuration file not found")
		os.Exit(1)
	}

	rdb := redis.NewRedisClient(cfg.RedisAddress)
	q := queue.NewQueue(rdb, cfg.QueueName)

	emailHandler := handler.NewEmailHandler(q, logger)
	healthHandler := handler.NewHealthHandler(rdb, logger)
	jobHandler := handler.NewJobHandler(rdb, logger)

	r := chi.NewRouter()

	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.RealIP)
	r.Use(middleware.Logger(logger))

	r.Get("/health", healthHandler.Check)
	r.Post("/email", emailHandler.SendEmail)
	r.Get("/jobs/{id}", jobHandler.GetStatus)

	server := &http.Server{
		Addr:    ":" + cfg.ApiPort,
		Handler: r,
	}

	// run server di goroutine
	go func() {
		logger.Info("Server started", "port", cfg.ApiPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("server shutdown failed", "error", err)
	}

	if err := rdb.Close(); err != nil {
		logger.Error("failed to close redis", "error", err)
	}

	logger.Info("Server exited properly")
}
