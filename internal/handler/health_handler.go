package handler

import (
	"async-email-worker/internal/helper"

	"log/slog"
	"net/http"

	"github.com/redis/go-redis/v9"
)

type HealthHandler struct {
	Redis  *redis.Client
	Logger *slog.Logger
}

func NewHealthHandler(r *redis.Client, logger *slog.Logger) *HealthHandler {
	return &HealthHandler{
		Redis:  r,
		Logger: logger,
	}
}

func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	err := h.Redis.Ping(r.Context()).Err()

	if err != nil {
		h.Logger.Error("Health check failed", "error", err.Error())
		helper.SendResponse(w, http.StatusServiceUnavailable, "Error", "Service is unhealthy", nil)
		return
	}

	helper.SendResponse(w, http.StatusOK, "OK", "Service is healthy", nil)
}
