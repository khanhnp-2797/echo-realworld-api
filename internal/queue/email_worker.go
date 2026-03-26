package queue

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/khanhnp-2797/echo-realworld-api/internal/mailer"
)

const (
	consumerGroup    = "email-workers"
	consumerName     = "worker-1"
	claimMinIdleTime = 30 * time.Second
)

// EmailWorker polls the welcome email stream and calls the mailer.
// Uses Redis Streams consumer groups for at-least-once delivery with automatic
// reclaim of messages that were in-flight during a crash/restart.
type EmailWorker struct {
	rdb    *redis.Client
	mailer mailer.Mailer
}

func NewEmailWorker(rdb *redis.Client, m mailer.Mailer) *EmailWorker {
	return &EmailWorker{rdb: rdb, mailer: m}
}

// Start initialises the consumer group and begins the processing loop.
// Blocks until ctx is cancelled — run in a goroutine: go worker.Start(ctx)
func (w *EmailWorker) Start(ctx context.Context) {
	if err := w.rdb.XGroupCreateMkStream(ctx, welcomeStream, consumerGroup, "$").Err(); err != nil {
		// BUSYGROUP = group already exists — normal after a restart.
		if !strings.Contains(err.Error(), "BUSYGROUP") {
			log.Printf("[email-worker] failed to create consumer group: %v", err)
		}
	}
	log.Printf("[email-worker] started group=%s consumer=%s stream=%s", consumerGroup, consumerName, welcomeStream)

	// On startup, reclaim any messages that were delivered but never acked
	// (e.g. the server crashed between processing and acking).
	w.reclaimPending(ctx)

	for {
		if ctx.Err() != nil {
			log.Printf("[email-worker] context cancelled, exiting")
			return
		}
		w.poll(ctx)
	}
}

func (w *EmailWorker) poll(ctx context.Context) {
	streams, err := w.rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    consumerGroup,
		Consumer: consumerName,
		Streams:  []string{welcomeStream, ">"},
		Count:    10,
		Block:    5 * time.Second,
	}).Result()
	if err != nil {
		// redis.Nil = BLOCK timeout — normal, just loop again.
		if err != redis.Nil {
			log.Printf("[email-worker] xreadgroup error: %v", err)
		}
		return
	}

	for _, stream := range streams {
		for _, msg := range stream.Messages {
			w.process(ctx, msg)
		}
	}
}

func (w *EmailWorker) process(ctx context.Context, msg redis.XMessage) {
	email, _ := msg.Values["email"].(string)
	username, _ := msg.Values["username"].(string)

	if err := w.mailer.SendWelcome(email, username); err != nil {
		// Do NOT ack — message stays in PEL and will be reclaimed on next restart.
		log.Printf("[email-worker] send failed email=%s id=%s: %v", email, msg.ID, err)
		return
	}

	if err := w.rdb.XAck(ctx, welcomeStream, consumerGroup, msg.ID).Err(); err != nil {
		log.Printf("[email-worker] xack failed id=%s: %v", msg.ID, err)
		return
	}
	log.Printf("[email-worker] sent welcome email=%s id=%s", email, msg.ID)
}

// reclaimPending uses XAUTOCLAIM to take ownership of messages idle for more
// than claimMinIdleTime (i.e. a previous worker processed but never acked them).
func (w *EmailWorker) reclaimPending(ctx context.Context) {
	res, _, err := w.rdb.XAutoClaim(ctx, &redis.XAutoClaimArgs{
		Stream:   welcomeStream,
		Group:    consumerGroup,
		Consumer: consumerName,
		MinIdle:  claimMinIdleTime,
		Start:    "0-0",
	}).Result()
	if err != nil {
		// XAutoClaim requires Redis >= 6.2. On older versions, skip gracefully.
		log.Printf("[email-worker] xautoclaim unavailable, skipping pending reclaim: %v", err)
		return
	}
	if len(res) == 0 {
		return
	}
	log.Printf("[email-worker] reclaiming %d pending messages", len(res))
	for _, msg := range res {
		w.process(ctx, msg)
	}
}
