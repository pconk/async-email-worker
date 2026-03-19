package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"async-email-worker/internal/entity"
	"async-email-worker/internal/helper"
	"async-email-worker/internal/queue"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type EmailHandler struct {
	Queue    *queue.Queue
	Logger   *slog.Logger
	Validate *validator.Validate
}

func NewEmailHandler(q *queue.Queue, l *slog.Logger) *EmailHandler {
	return &EmailHandler{
		Queue:    q,
		Logger:   l,
		Validate: validator.New(),
	}
}

func (h *EmailHandler) SendEmail(w http.ResponseWriter, r *http.Request) {
	var job entity.EmailJob

	if err := json.NewDecoder(r.Body).Decode(&job); err != nil {
		helper.SendResponse(w, http.StatusBadRequest, "Fail", "Invalid ID atau Format JSON tidak valid", nil)
		return
	}

	if err := h.Validate.Struct(job); err != nil {
		// Panggil helper untuk merapikan error
		prettyErrors := helper.FormatValidationError(err)

		// Kirim response dengan status 400 tapi data berisi detail errornya
		helper.SendResponse(w, http.StatusBadRequest, "Validation Error", "Beberapa field tidak valid", prettyErrors)
		return
	}

	job.ID = uuid.New().String()
	job.Retry = 0

	err := h.Queue.EnqueueEmail(r.Context(), job)
	if err != nil {
		helper.SendResponse(w, http.StatusInternalServerError, "Fail", "Gagal queue redis", nil)
		return
	}

	helper.SendResponse(w, http.StatusAccepted, "OK", "Email queued", job)

}
