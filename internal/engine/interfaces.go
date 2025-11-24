package engine

import "context"

// Core engine interfaces (lightweight, expand as implementation progresses)

type Workflow interface {
	ID() string
	Version() string
	// Plan returns next step definitions (deterministic, pure)
	Plan(ctx context.Context, run *WorkflowRun) ([]WorkflowStepDef, error)
}

type WorkflowStep interface {
	Name() string
}

type WorkflowStepDef struct {
	StepName string
	Seq      int
	Input    []byte
}

type ExecutionContext struct {
	Run  *WorkflowRun
	Step *WorkflowStepRecord
}

type StepResult struct {
	Success bool
	Output  []byte
}

type StateStore interface {
	CreateRun(ctx context.Context, run *WorkflowRun) error
	LoadRun(ctx context.Context, runID string) (*WorkflowRun, error)
	InsertSteps(ctx context.Context, steps []*WorkflowStepRecord) error
	ClaimNextStep(ctx context.Context, workerID string) (*WorkflowStepRecord, error)
	UpdateStep(ctx context.Context, step *WorkflowStepRecord) error
	AppendLog(ctx context.Context, log *StepLog) error
	EnqueueEvent(ctx context.Context, ev *OutboxEvent) error
	// RequeueStaleSteps inspects in-progress steps whose last heartbeat is older
	// than heartbeatTTL (seconds) and resets them to pending with incremented
	// attempts and appropriate next_attempt_at/backoff. Returns number requeued.
	RequeueStaleSteps(ctx context.Context, heartbeatTTLSeconds int, limit int) (int, error)
}

type Executor interface {
	RunStep(ctx context.Context, step *WorkflowStepRecord) (StepResult, error)
}

type Dispatcher interface {
	Dispatch(ctx context.Context, ev OutboxEvent) error
	Subscribe(eventType string) (<-chan OutboxEvent, error)
}

type WorkflowRegistry interface {
	Register(w Workflow)
	Get(id string) (Workflow, bool)
}

type Logger interface {
	Info(msg string, fields ...Field)
	Error(msg string, fields ...Field)
}

type Field struct {
	K string
	V interface{}
}

type Queue interface {
	Enqueue(ctx context.Context, qname string, payload []byte) error
	Dequeue(ctx context.Context, qname string) ([]byte, error)
}
