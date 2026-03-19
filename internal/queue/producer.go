package queue

import (
	"context"
	"encoding/json"

	"async-email-worker/internal/entity"

	"github.com/redis/go-redis/v9"
)

type Queue struct {
	Redis     *redis.Client
	QueueName string
}

func NewQueue(client *redis.Client, qname string) *Queue {
	return &Queue{
		Redis:     client,
		QueueName: qname,
	}
}

func (q *Queue) EnqueueEmail(ctx context.Context, job entity.EmailJob) error {
	data, _ := json.Marshal(job)
	return q.Redis.LPush(ctx, q.QueueName, data).Err()
}
