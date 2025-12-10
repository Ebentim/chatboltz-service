package scheduler

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/alpinesboltltd/boltz-ai/internal/engine"
)

// Start begins a simple scheduler loop and worker dispatch. It returns a
// done channel that will be closed once the scheduler stops and all in-flight
// workers have finished. This allows callers to wait for graceful shutdown.
func Start(ctx context.Context, store engine.StateStore, exec engine.Executor, reg engine.WorkflowRegistry, disp engine.Dispatcher, workerCount int) (<-chan struct{}, error) {
	// worker semaphore
	sem := make(chan struct{}, workerCount)
	var wg sync.WaitGroup

	done := make(chan struct{})

	ticker := time.NewTicker(500 * time.Millisecond)
	go func() {
		defer func() {
			// wait for workers to finish
			wg.Wait()
			ticker.Stop()
			close(done)
		}()

		for {
			select {
			case <-ctx.Done():
				// stop accepting new work and wait for in-flight workers
				return
			case <-ticker.C:
				// try claim a step
				step, err := store.ClaimNextStep(ctx, "scheduler")
				if err != nil {
					log.Printf("scheduler: claim error: %v", err)
					continue
				}
				if step == nil {
					continue
				}

				sem <- struct{}{}
				wg.Add(1)
				go func(s *engine.WorkflowStepRecord) {
					defer func() { <-sem; wg.Done() }()
					// run step
					res, err := exec.RunStep(ctx, s)
					if err != nil {
						log.Printf("executor error: %v", err)
						// update step with failure - scaffold
						s.Status = "failed"
						_ = store.UpdateStep(ctx, s)
						return
					}
					// on success persist result and mark completed
					s.Result = res.Output
					s.Status = "completed"
					_ = store.UpdateStep(ctx, s)
					// Optionally call workflow.Plan via registry to enqueue next steps (left as future work)
				}(step)
			}
		}
	}()

	return done, nil
}
