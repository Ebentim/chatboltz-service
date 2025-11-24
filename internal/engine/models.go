package engine

import (
	"encoding/json"
	"time"
)

// Minimal models used by the engine scaffold. These are intentionally
// scoped to the engine package to avoid touching existing entity models
// during the scaffold stage.

type WorkflowRun struct {
    ID              string          `json:"id"`
    WorkflowType    string          `json:"workflow_type"`
    WorkflowVersion string          `json:"workflow_version"`
    Status          string          `json:"status"`
    Payload         json.RawMessage `json:"payload"`
    CreatedAt       time.Time       `json:"created_at"`
    UpdatedAt       time.Time       `json:"updated_at"`
}

type WorkflowStepRecord struct {
    ID             string          `json:"id"`
    RunID          string          `json:"run_id"`
    StepName       string          `json:"step_name"`
    Seq            int             `json:"seq"`
    Status         string          `json:"status"`
    Input          json.RawMessage `json:"input"`
    Result         json.RawMessage `json:"result"`
    Attempts       int             `json:"attempts"`
    MaxAttempts    int             `json:"max_attempts"`
    NextAttemptAt  *time.Time      `json:"next_attempt_at"`
    ClaimedAt      *time.Time      `json:"claimed_at"`
    LastHeartbeat  *time.Time      `json:"last_heartbeat"`
    LockOwner      *string         `json:"lock_owner"`
    IdempotencyKey *string         `json:"idempotency_key"`
    Error          *string         `json:"error"`
    CreatedAt      time.Time       `json:"created_at"`
    UpdatedAt      time.Time       `json:"updated_at"`
}

type OutboxEvent struct {
    ID          string          `json:"id"`
    EventType   string          `json:"event_type"`
    Payload     json.RawMessage `json:"payload"`
    State       string          `json:"state"`
    IdempotencyKey *string      `json:"idempotency_key"`
    Published   bool            `json:"published"`
    CreatedAt   time.Time       `json:"created_at"`
}

type StepLog struct {
    ID        string          `json:"id"`
    StepID    string          `json:"step_id"`
    Level     string          `json:"level"`
    Message   string          `json:"message"`
    Meta      json.RawMessage `json:"meta"`
    CreatedAt time.Time       `json:"created_at"`
}
