package handler

import (
	"async-email-worker/internal/helper"

	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
)

type JobHandler struct {
	Redis  *redis.Client
	Logger *slog.Logger
}

func NewJobHandler(r *redis.Client, logger *slog.Logger) *JobHandler {
	return &JobHandler{
		Redis:  r,
		Logger: logger,
	}
}

func (j *JobHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ctx := r.Context()

	if id == "" {
		helper.SendResponse(w, http.StatusBadRequest, "Fail", "Invalid ID atau Format JSON tidak valid", nil)
		return
	}
	val, err := j.Redis.Get(ctx, "job_status:"+id).Result()
	if err != nil {
		helper.SendResponse(w, http.StatusNotFound, "Fail", "Job tidak ditemukan", nil)
		return
	}

	helper.SendResponse(w, http.StatusOK, "OK", "Job Found", map[string]any{
		"status": val,
	})
}
