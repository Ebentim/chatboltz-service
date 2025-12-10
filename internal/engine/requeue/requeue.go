package requeue

import (
	"context"
	"log"
	"time"

	"github.com/alpinesboltltd/boltz-ai/internal/engine"
)

// StartRequeueMonitor starts a background ticker that calls into the StateStore
// to requeue stale in-progress steps. heartbeatTTLSeconds controls how old a
// heartbeat must be to consider a step stale. interval controls how often the
// monitor runs.
func StartRequeueMonitor(ctx context.Context, store engine.StateStore, interval time.Duration, heartbeatTTLSeconds int, batchSize int) {
	if store == nil {
		log.Printf("requeue: no store provided, monitor disabled")
		return
	}
	if interval <= 0 {
		interval = 30 * time.Second // sensible default
	}
	ticker := time.NewTicker(interval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				n, err := store.RequeueStaleSteps(ctx, heartbeatTTLSeconds, batchSize)
				if err != nil {
					log.Printf("requeue: error requeueing stale steps: %v", err)
					continue
				}
				if n > 0 {
					log.Printf("requeue: requeued %d stale steps", n)
				}
			}
		}
	}()
}
