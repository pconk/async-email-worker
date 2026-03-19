package worker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"sync"
	"time"

	"async-email-worker/internal/entity"
	"async-email-worker/internal/service"

	"github.com/redis/go-redis/v9"
)

type Worker struct {
	Redis        *redis.Client
	Logger       *slog.Logger
	QueueName    string
	WorkerNumber int
	MaxRetry     int
	emailService *service.EmailService
}

func NewWorker(r *redis.Client, l *slog.Logger, qname string, wnumber int, retry int) *Worker {
	return &Worker{
		Redis:        r,
		Logger:       l,
		QueueName:    qname,
		WorkerNumber: wnumber,
		MaxRetry:     retry,
		emailService: service.NewEmailService(l),
	}
}

func (w *Worker) Start(ctx context.Context) {
	var wg sync.WaitGroup

	for i := 1; i <= w.WorkerNumber; i++ {
		wg.Add(1)

		go func(workerID int) {
			defer wg.Done()
			w.process(ctx, workerID)
		}(i)
	}

	// Tunggu semua worker selesai
	wg.Wait()

	w.Logger.Info("All workers stopped")
}

func (w *Worker) dequeue(ctx context.Context) (entity.EmailJob, error) {
	result, err := w.Redis.BRPop(ctx, 5*time.Second, w.QueueName).Result()

	var job entity.EmailJob

	if err != nil {
		// Cek apakah ini cuma timeout biasa (kosong)
		if errors.Is(err, redis.Nil) {
			return job, fmt.Errorf("timeout") // Kita bungkus jadi custom error
		}
		return job, err
	}

	if len(result) < 2 {
		return job, fmt.Errorf("invalid redis response")
	}

	err = json.Unmarshal([]byte(result[1]), &job)
	if err != nil {
		w.Logger.Error("failed to unmarshal job", "error", err)
		return job, err
	}

	return job, nil
}

func (w *Worker) process(ctx context.Context, workerID int) {
	w.Logger.Info("Worker started", "id", workerID)

	for {
		// cek context sebelum dequeue
		if ctx.Err() != nil {
			w.Logger.Info("Worker stopping", "id", workerID)
			return
		}

		job, err := w.dequeue(ctx)
		if err != nil {
			// 1. Kalau karena context cancel → stop
			if ctx.Err() != nil {
				w.Logger.Info("Worker stopping (ctx cancelled)", "id", workerID)
				return
			}
			// 2. Kalau cuma timeout (antrean kosong), abaikan & loop lagi (silent)
			if err.Error() == "timeout" {
				continue
			}

			w.Logger.Error("failed to dequeue", "error", err, "worker_id", workerID)

			// Kasih jeda dikit biar gak spamming kalau Redis beneran down
			time.Sleep(1 * time.Second)
			continue
		}

		w.handleJob(ctx, job, workerID)
	}
}

func (w *Worker) handleJob(ctx context.Context, job entity.EmailJob, workerID int) {

	status, _ := w.Redis.Get(ctx, "job_status:"+job.ID).Result()

	if status == "success" {
		return // skip
	}

	w.updateStatus(ctx, job.ID, "processing")

	err := w.emailService.SendEmail(job, w.MaxRetry)
	if err != nil {

		w.Logger.Error("send failed",
			"worker_id", workerID,
			"job_id", job.ID,
			"retry", job.Retry,
			"error", err,
		)

		job.Retry++

		if job.Retry >= w.MaxRetry {
			w.moveToDeadQueue(ctx, job)
			w.updateStatus(ctx, job.ID, "failed")
			return
		}

		// exponential backoff
		delay := time.Duration(math.Pow(2, float64(job.Retry))) * time.Second
		time.Sleep(delay)

		w.requeue(ctx, job)
		w.updateStatus(ctx, job.ID, "retrying")
		return
	}

	w.updateStatus(ctx, job.ID, "success")

	w.Logger.Info("job success",
		"worker_id", workerID,
		"job_id", job.ID,
	)
}

func (w *Worker) requeue(ctx context.Context, job entity.EmailJob) {
	data, _ := json.Marshal(job)
	w.Redis.LPush(ctx, w.QueueName, data)
}

func (w *Worker) moveToDeadQueue(ctx context.Context, job entity.EmailJob) {
	data, _ := json.Marshal(job)
	w.Redis.LPush(ctx, "email_dead_queue", data)
}

func (w *Worker) updateStatus(ctx context.Context, id string, status string) {
	key := "job_status:" + id
	w.Redis.Set(ctx, key, status, time.Hour)
}
