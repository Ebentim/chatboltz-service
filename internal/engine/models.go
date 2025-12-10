package engine

import (
	canonical "github.com/alpinesboltltd/boltz-ai/internal/engine/models"
)

// Minimal models used by the engine scaffold.
// We alias the canonical models to avoid split-brain types.

type WorkflowRun = canonical.WorkflowRun
type WorkflowStepRecord = canonical.WorkflowStepRecord
type OutboxEvent = canonical.OutboxEvent
type StepLog = canonical.StepLog
