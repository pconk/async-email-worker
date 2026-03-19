package service

import (
	"fmt"
	"log/slog"

	"async-email-worker/internal/entity"
)

type EmailService struct {
	Logger *slog.Logger
}

func NewEmailService(l *slog.Logger) *EmailService {
	return &EmailService{Logger: l}

}

func (e *EmailService) SendEmail(job entity.EmailJob, maxRetry int) error {
	e.Logger.Info("sending email", "job", job)

	// simulasi gagal
	if job.Retry < (maxRetry - 1) {
		return fmt.Errorf("temporary error")
	}

	return nil
}
