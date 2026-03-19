package main

import (
	"async-email-worker/internal/config"
	"async-email-worker/internal/worker"
	"async-email-worker/pkg/redis"
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Error("configuration file not found")
		os.Exit(1)
	} else if cfg.ApiPort == "" {
		logger.Error("APP_PORT")
		os.Exit(1)

	}

	rdb := redis.NewRedisClient(cfg.RedisAddress)

	w := worker.NewWorker(rdb, logger, cfg.QueueName, cfg.WorkerNumber, cfg.MaxRetry)

	ctx, cancel := context.WithCancel(context.Background())

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	var wgWorkers sync.WaitGroup

	go func() {
		sig := <-sigChan
		logger.Info("Shutdown signal received", "signal", sig.String())
		cancel()
	}()

	wgWorkers.Add(1)
	go func() {
		defer wgWorkers.Done()
		w.Start(ctx)
	}()

	logger.Info("Creating Workers")

	<-ctx.Done()

	logger.Info("Waiting worker to finish...")
	wgWorkers.Wait()

	logger.Info("Worker stopped gracefully")
}
