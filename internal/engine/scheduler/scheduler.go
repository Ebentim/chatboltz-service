package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/alpinesboltltd/boltz-ai/internal/engine"
)

// Start begins a simple scheduler loop and worker dispatch. This is a scaffold
// to be extended: it shows how to claim steps from the StateStore and run them
// using an Executor.
func Start(ctx context.Context, store engine.StateStore, exec engine.Executor, reg engine.WorkflowRegistry, disp engine.Dispatcher, workerCount int) error {
	// worker semaphore
	sem := make(chan struct{}, workerCount)

	ticker := time.NewTicker(500 * time.Millisecond)
	go func() {
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
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
				go func(s *engine.WorkflowStepRecord) {
					defer func() { <-sem }()
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

	return nil
}
