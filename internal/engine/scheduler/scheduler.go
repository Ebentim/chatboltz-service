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

					// Create a detached context for the step execution so it isn't killed immediately on scheduler shutdown.
					// We add a hard timeout (e.g. 5 minutes) to prevent zombies.
					stepCtx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
					defer cancel()

					// Start heartbeat ticker
					hbDone := make(chan struct{})
					go func() {
						hbTicker := time.NewTicker(30 * time.Second)
						defer hbTicker.Stop()
						for {
							select {
							case <-hbDone:
								return
							case <-stepCtx.Done():
								return
							case <-hbTicker.C:
								if err := store.HeartbeatStep(stepCtx, s.ID); err != nil {
									log.Printf("scheduler: heartbeat failed for step %s: %v", s.ID, err)
								}
							}
						}
					}()

					// run step
					res, err := exec.RunStep(stepCtx, s)
					
					// Stop heartbeat
					close(hbDone)

					if err != nil {
						log.Printf("executor error for step %s: %v", s.ID, err)
						// update step with failure
						s.Status = "failed"
						s.Error = &[]string{err.Error()}[0] // hack to get pointer to string
						if updateErr := store.UpdateStep(context.Background(), s); updateErr != nil {
							log.Printf("scheduler: failed to update failed step %s: %v", s.ID, updateErr)
						}
						return
					}
					// on success persist result and mark completed
					s.Result = res.Output
					s.Status = "completed"
					if updateErr := store.UpdateStep(context.Background(), s); updateErr != nil {
						log.Printf("scheduler: failed to update completed step %s: %v", s.ID, updateErr)
					}
					// Slightly different context for update to ensure it persists even during shutdown
				}(step)
			}
		}
	}()

	return done, nil
}
