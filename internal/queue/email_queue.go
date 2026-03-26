package queue

import (
	"context"

	"github.com/redis/go-redis/v9"
)

const welcomeStream = "mq:email:welcome"

// EmailQueue enqueues transactional email jobs durably via Redis Streams.
type EmailQueue interface {
	EnqueueWelcome(ctx context.Context, email, username string) error
}

type redisEmailQueue struct {
	rdb *redis.Client
}

func NewRedisEmailQueue(rdb *redis.Client) EmailQueue {
	return &redisEmailQueue{rdb: rdb}
}

func (q *redisEmailQueue) EnqueueWelcome(ctx context.Context, email, username string) error {
	return q.rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: welcomeStream,
		Values: map[string]any{
			"email":    email,
			"username": username,
		},
	}).Err()
}
