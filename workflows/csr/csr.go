package csr

import (
	"context"
	"fmt"

	"github.com/alpinesboltltd/boltz-ai/internal/engine"
)

// CSR workflow scaffold. Implements engine.Workflow and provides a deterministic Plan.
type CSRWorkflow struct{}

func New() *CSRWorkflow { return &CSRWorkflow{} }

func (w *CSRWorkflow) ID() string { return "csr" }

func (w *CSRWorkflow) Version() string { return "v1" }

func (w *CSRWorkflow) Plan(ctx context.Context, run *engine.WorkflowRun) ([]engine.WorkflowStepDef, error) {
	// Minimal deterministic plan: create a single fetch_ticket step for a new run.
	if run == nil {
		return nil, fmt.Errorf("run is nil")
	}
	step := engine.WorkflowStepDef{StepName: "fetch_ticket", Seq: 1, Input: run.Payload}
	return []engine.WorkflowStepDef{step}, nil
}
