# Orchestration Engine Design & Implementation Plan

This document describes a production-grade orchestration engine that runs inside the Go monolith (`chatboltz-service`). It follows the project's clean architecture and avoids Kafka, Temporal, Redis, or microservices. The MVP focuses strictly on the Customer Support Representative (CSR) workflow.

Contents (must appear in repository): `plan/orchestration_plan.md`

Requirements satisfied across iterations: deterministic state machines, DB-backed job queue, single-process scheduler loop, goroutines/channels, Postgres durability, idempotent actions, retries/backoff, human-in-the-loop, event routing via DB and in-process signals, structured logs/tracing/metrics.

---

1. High-Level Architecture

---

Summary

- The orchestration engine runs as an internal package inside the monolith (e.g. `internal/engine` and `engine` top-level package). It persists durable state in Postgres and uses a DB-backed queue table plus `SELECT ... FOR UPDATE SKIP LOCKED` for safe scheduling.
- A single scheduler loop polls the DB for pending `workflow_steps`, dispatches them to worker goroutines via an internal queue abstraction, and persists results. Execution is deterministic: workflow state machines decide next steps based on persisted state.

Key components

- StateStore (DB): responsible for workflow runs, steps, logs, retries, outbox events.
- Scheduler: single-process loop that claims pending steps (SKIP LOCKED) and enqueues to worker pool.
- Worker pool: goroutines executing `Executor` actions (use channels and contexts).
- Dispatcher: event router that wakes workflows based on internal events or DB notifications.
- Workflow Registry: register workflow definitions (versioned) and factory for instances.
- Connectors: thin adapter packages for CRM, email, calendar, vector search under `/connectors` mapped to existing integrations (e.g., `internal/integrations/google`, `integrations/twilio`).

Concurrency & safety

- Use database optimistic locking and `SELECT FOR UPDATE SKIP LOCKED` to ensure single-claim semantics for steps.
- Each step execution is idempotent; the engine enforces idempotence via `idempotency_key` and persisted step status.
- Use short-lived transactions for step claim + save; avoid long transactions across external calls.

Durability & crash recovery

- All step state, retry metadata, and logs are persisted in Postgres. Scheduler restarts will resume pending or in-progress steps (in-progress steps have TTL and can be re-queued after heartbeat expiry).

Integration with existing monolith

- Place engine under `internal/engine` (implementation) and public API shim under `engine` if needed. Register connectors with existing `integrations/*` packages.

2. Go Folder Structure

---

Top-level suggested layout (fits clean architecture in repo):

```
/engine                # public package facade (optional)
/internal/engine       # implementation
  /scheduler           # scheduler loop & worker pool
  /executor            # action executors and wrappers
  /store               # Postgres store implementation
  /workflow            # workflow runtime types & registry
  /dispatcher          # event routing & listener
  /models              # Go models mirroring DB schemas (if not in central models)
/workflows             # modular workflows
  /csr                 # Customer Support Representative (MVP) workflow
  /bdr                 # Business Development Representative
  /sdr                 # Sales Development Representative
  /va                  # Virtual Assistant
/connectors            # light glue for external systems
  /crm
  /email
  /calendar
  /vector
/pkg/queue             # DB-backed queue abstraction and helpers
/pkg/store             # common DB helper utilities & migrations
/pkg/executor          # reusable executors, idempotency helpers
/models                # central domain models used across app
/cmd/api               # API endpoints to trigger workflows / human approvals
plan/                  # contains this plan and iterative reviews
```

Note: this layout co-exists with existing `internal/` and `usecase/` layers. You can reuse `repository` or `provider` packages for connectors.

3. Engine Interfaces

---

Below are the production-ready Go interfaces (simplified signatures, add context and error types in code):

```go
package engine

import "context"

type Workflow interface {
    // ID returns the workflow type identifier, e.g. "csr.v1"
    ID() string
    // Version returns semantic version for workflow definition
    Version() string
    // Plan decides initial steps or next steps given run state
    Plan(ctx context.Context, run *WorkflowRun) ([]WorkflowStepDef, error)
}

type WorkflowStep interface {
    // Execute performs the step action; must be idempotent
    Execute(ctx context.Context, execCtx ExecutionContext) (StepResult, error)
    // Name returns the step name/type
    Name() string
}

type StateStore interface {
    // CreateRun persists a new workflow run
    CreateRun(ctx context.Context, run *WorkflowRun) error
    // LoadRun loads run with steps
    LoadRun(ctx context.Context, runID string) (*WorkflowRun, error)
    // ClaimNextStep claims a pending step (SKIP LOCKED) and returns it
    ClaimNextStep(ctx context.Context, workerID string) (*WorkflowStepRecord, error)
    // UpdateStep updates step status/result
    UpdateStep(ctx context.Context, step *WorkflowStepRecord) error
    // AppendLog persists a step log line
    AppendLog(ctx context.Context, log *StepLog) error
    // EnqueueEvent persists an outbox event
    EnqueueEvent(ctx context.Context, ev *OutboxEvent) error
}

type Executor interface {
    // RunStep runs a WorkflowStep and returns result
    RunStep(ctx context.Context, step *WorkflowStepRecord) (StepResult, error)
}

type Dispatcher interface {
    // Dispatch in-process event to the router
    Dispatch(ctx context.Context, ev Event) error
    // Subscribe returns channel for events of a given type
    Subscribe(eventType string) (<-chan Event, error)
}

type WorkflowRegistry interface {
    Register(w Workflow)
    Get(id string) (Workflow, bool)
}

type Logger interface {
    Info(msg string, fields ...Field)
    Error(msg string, fields ...Field)
}

// Queue abstraction: minimal declarations; implementation uses DB
type Queue interface {
    Enqueue(ctx context.Context, qname string, payload []byte) error
    Dequeue(ctx context.Context, qname string) ([]byte, error)
}
```

Notes:

- `WorkflowStepDef` and `ExecutionContext`/`StepResult` are typed structs in `internal/engine/workflow` used to persist arguments and result metadata.
- Implementations should use contexts for cancellation and deadlines.

4. Data Models

---

SQL schema (Postgres, example migration):

```sql
-- workflow_runs: one row per workflow instance
CREATE TABLE workflow_runs (
  id UUID PRIMARY KEY,
  workflow_type TEXT NOT NULL,
  workflow_version TEXT NOT NULL,
  status TEXT NOT NULL, -- running|completed|failed|paused
  payload JSONB, -- workflow-level context
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

-- workflow_steps: each actionable step
CREATE TABLE workflow_steps (
  id UUID PRIMARY KEY,
  run_id UUID NOT NULL REFERENCES workflow_runs(id) ON DELETE CASCADE,
  step_name TEXT NOT NULL,
  seq INT NOT NULL,
  status TEXT NOT NULL, -- pending|in_progress|completed|failed|waiting
  input JSONB,
  result JSONB,
  attempts INT DEFAULT 0,
  max_attempts INT DEFAULT 5,
  next_attempt_at TIMESTAMP WITH TIME ZONE,
  lock_owner TEXT, -- worker id
  idempotency_key TEXT,
  error TEXT,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

-- outbox events for reliable dispatch to other parts of system
CREATE TABLE outbox_events (
  id UUID PRIMARY KEY,
  event_type TEXT NOT NULL,
  payload JSONB,
  published BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

-- step logs
CREATE TABLE step_logs (
  id UUID PRIMARY KEY,
  step_id UUID REFERENCES workflow_steps(id) ON DELETE CASCADE,
  level TEXT,
  message TEXT,
  meta JSONB,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

-- retry metadata (optional separate table if you want history)
CREATE TABLE retry_meta (
  step_id UUID PRIMARY KEY REFERENCES workflow_steps(id) ON DELETE CASCADE,
  attempts INT DEFAULT 0,
  backoff_seconds INT DEFAULT 0,
  dead_letter BOOLEAN DEFAULT FALSE
);
```

Go structs (in `internal/engine/models`):

```go
type WorkflowRun struct {
  ID string `json:"id"`
  WorkflowType string `json:"workflow_type"`
  WorkflowVersion string `json:"workflow_version"`
  Status string `json:"status"`
  Payload json.RawMessage `json:"payload"`
  CreatedAt time.Time `json:"created_at"`
  UpdatedAt time.Time `json:"updated_at"`
}

type WorkflowStepRecord struct {
  ID string `json:"id"`
  RunID string `json:"run_id"`
  StepName string `json:"step_name"`
  Seq int `json:"seq"`
  Status string `json:"status"`
  Input json.RawMessage `json:"input"`
  Result json.RawMessage `json:"result"`
  Attempts int `json:"attempts"`
  MaxAttempts int `json:"max_attempts"`
  NextAttemptAt *time.Time `json:"next_attempt_at"`
  LockOwner *string `json:"lock_owner"`
  IdempotencyKey *string `json:"idempotency_key"`
  Error *string `json:"error"`
  CreatedAt time.Time `json:"created_at"`
  UpdatedAt time.Time `json:"updated_at"`
}

type OutboxEvent struct {
  ID string `json:"id"`
  EventType string `json:"event_type"`
  Payload json.RawMessage `json:"payload"`
  Published bool `json:"published"`
  CreatedAt time.Time `json:"created_at"`
}

type StepLog struct {
  ID string `json:"id"`
  StepID string `json:"step_id"`
  Level string `json:"level"`
  Message string `json:"message"`
  Meta json.RawMessage `json:"meta"`
  CreatedAt time.Time `json:"created_at"`
}
```

5. Orchestration Execution Flow

---

Trigger -> schedule -> execute -> persist loop

1. Triggering a workflow

   - Client calls `POST /workflows` (in `cmd/api`) or an internal usecase triggers `StateStore.CreateRun` with initial payload.
   - Engine stores `workflow_runs` and enqueues initial `workflow_steps` (status = `pending`) as per workflow `Plan`.

2. Scheduler wakes up

   - Scheduler periodically polls `workflow_steps` where `status = 'pending'` and `next_attempt_at <= now()` using a claim query:
     ```sql
     UPDATE workflow_steps
     SET lock_owner = $1, status = 'in_progress', updated_at = now()
     WHERE id = (
       SELECT id FROM workflow_steps
       WHERE status = 'pending' AND next_attempt_at <= now()
       ORDER BY seq, created_at
       FOR UPDATE SKIP LOCKED
       LIMIT 1
     ) RETURNING *;
     ```
   - The claim is atomic, returning a single row for the worker to process.

3. Step execution

   - Worker unmarshals step input into typed struct and calls registered `WorkflowStep` implementations via `Executor.RunStep`.
   - Executions must be idempotent: steps receive `idempotency_key` and must store outputs by updating `workflow_steps.result` and `status = completed`.

4. Persisting results

   - Worker updates `workflow_steps` (attempts, status, result, error) in a transaction.
   - Worker appends `step_logs` and emits `outbox_events` for external systems.

5. Scheduling next steps

   - After persisting result, engine calls `Workflow.Plan` to determine next steps. New `workflow_steps` inserted with `pending` status.
   - The engine may publish in-process events via `Dispatcher.Dispatch` or create `outbox_events` for external systems.

6. Retries and failures

   - If a step fails, increment `attempts` and set `next_attempt_at` = now() + backoff.
   - If `attempts >= max_attempts`, set status to `failed` and optionally move to dead-letter.

7. Human approvals

   - Steps that require human action are saved with `status = waiting` or `paused` and the engine emits an `outbox_event` + in-process event to notify UI/ops.
   - A REST API endpoint (e.g., `POST /workflows/{run}/approve`) sets the step status to `pending` again with updated payload.

8. Crash recovery

   - Steps in `in_progress` longer than a heartbeat TTL (e.g. 5m) are considered abandoned and reset to `pending` for re-claim. Alternatively, assign a `last_heartbeat` and let scheduler requeue.

9. Workflow Registration Examples

---

Create a registry at application start (e.g., in `cmd/api` or `internal/app`) and register workflows:

```go
reg := engine.NewRegistry()
reg.Register(csr.NewWorkflow()) // csr.v1
reg.Register(bdr.NewWorkflow())
// etc
```

Workflows implement the `Workflow` interface and return initial `WorkflowStepDef`s for `Plan`.

7. Example CSR/BDR/SDR/VA Workflows

---

CSR (MVP)

- Purpose: handle customer support ticket: triage, context retrieval (vector search), draft response via AI, send email, verify resolution.

Step list (logical):

1. fetch_ticket (collect ticket metadata)
2. retrieve_context (vector search and knowledge base)
3. draft_response (call LLM provider via `provider/ai-provider`)
4. human_review (optional) -> WAIT
5. send_response (email/send via connectors)
6. confirm_resolution (poll or await customer reply)

BDR / SDR / VA: similar multi-step sequences with steps for contact lookup (CRM), plan outreach, send email, schedule meeting.

Sample CSR step input/output (json):

```json
{
  "ticket_id": "123",
  "customer_id": "abc",
  "context": {...}
}
```

8. Scheduler & Retry Logic

---

Scheduler design

- Single-process scheduler loop with back-off polling, e.g. tick every 500ms but also use LISTEN/NOTIFY to wake up when new rows inserted.
- Use `ClaimNextStep` implemented with `FOR UPDATE SKIP LOCKED` to allow safe concurrency if we later run multiple workers in the same process.

Retry/backoff

- Exponential backoff: next_attempt_seconds = base \* 2^(attempts-1) with jitter (e.g., 0.1-0.4) and cap (e.g., 1h).
- `max_attempts` default = 5; configurable per-step.
- On exceeding `max_attempts`, mark step `failed` and create dead-letter event in `outbox_events`.

Idempotency

- Steps must be idempotent. Enforce by storing `idempotency_key` and checking if a previously completed result exists before executing.
- For external side effects (email, CRM write), use connector-level idempotency (store provider ids in step result).

Deduplication

- Deduplicate triggers by combining `run_id + step_name + idempotency_key` uniqueness index (DB constraint) when creating steps.

9. Event Routing Design

---

Mechanisms (no Kafka)

- DB-outbox + publisher: persist `outbox_events` and a background publisher reads and routes them (internal in-process or via connector). Mark `published` when done.
- Postgres `LISTEN/NOTIFY` optional: when an important event is written, we `NOTIFY orchestration_event, payload` and scheduler or dispatcher listens to wake instantly.
- In-process channels: dispatcher maintains channels for subscribers and broadcasts events.

Flow

- When a step produces an event (e.g., `email_sent`), worker inserts an `outbox_events` row and calls `pg_notify('orchestration_event', json)` (or uses dispatcher.Dispatch directly). The dispatcher will route to interested workflow runs.

Security and delivery

- Outbox ensures durability. Publisher retries until published. Use idempotency for consumer actions.

10. Human-in-the-Loop

---

Approach

- Steps that require human decision are persisted with `status = waiting` and `metadata.wait_reason` (e.g., "approve response").
- The engine emits an `outbox_event` and an in-process `Dispatcher` event so UI and notification subsystems can surface the task.
- The API exposes endpoints for listing pending human tasks and for completing/approving them. The API handler calls `StateStore.UpdateStep` to set the step input (approved payload) and mark `status = pending` so scheduler can pick it up.

Timeouts & Escalation

- Steps may have `human_timeout` metadata; a background monitor escalates or assigns to fallback (e.g., manager) after TTL.

Auditability

- Store full history in `step_logs` and maintain `retry_meta` records for human decisions and timestamps.

11. Observability

---

Logging

- Use structured logging (e.g. `uber/zap`) with fields: `workflow_run`, `step_id`, `step_name`, `attempt`, `worker_id`, `error`.
- Example fields: `workflow_type`, `workflow_version`, `run_id`, `step_id`, `seq`, `attempts`.

Tracing

- Use OpenTelemetry for traces within the monolith. Start span when scheduler claims step and propagate to executors and connectors.

Metrics

- Counters: `workflow_runs_started_total`, `workflow_steps_executed_total`, `workflow_steps_failed_total`.
- Gauges: `workflow_steps_pending`, `workers_active`.
- Histograms: `step_execution_latency_seconds`, `step_retry_delay_seconds`.

Dashboards

- Grafana dashboards for throughput, failure rate by workflow type, retry distribution, human-waiting counts.

12. Implementation Roadmap

---

Iterative approach (we will perform 5 refinements to reach the best options):

- Round 1: baseline design and validation (this doc + basic APIs).
- Round 2: tighten concurrency primitives and idempotency mechanisms; add outbox patterns.
- Round 3: improve event routing (LISTEN/NOTIFY vs polling tradeoffs) and human-in-loop UX.
- Round 4: add observability/tracing and error handling improvements.
- Round 5: finalize versioning, audit, and rollout strategy.

MVP (scope: CSR only)

- DB models and migrations for `workflow_runs`, `workflow_steps`, `outbox_events`, `step_logs`, `retry_meta`.
- Engine runtime skeleton in `internal/engine`.
- Scheduler loop and worker pool; `pkg/queue` for DB-backed queue operations.
- Implement CSR workflow in `workflows/csr` with steps: fetch_ticket, retrieve_context, draft_response, human_review (optional), send_response, confirm_resolution.
- Retry/backoff and idempotency enforcement.
- Structured logging and basic Prometheus metrics.

Beta

- Add BDR/SDR/VA workflows.
- Implement connector library (`connectors/*`) for CRM, Email, Calendar, Vector.
- Human-in-the-loop flows and UI task endpoints.
- Enhanced metrics and basic tracing.

Production

- Workflow versioning and migrations support for changing workflow definitions.
- Event replay safety (processing idempotently using outbox and idempotency keys).
- Scale-out readiness: while single-process now, design allows multiple processes in a single DB to share steps (via SKIP LOCKED) with leader election for singleton tasks.
- Rollout strategy: feature flags per workflow type, canary runs for new workflow versions, and tools to migrate runs across versions.

Integration notes (fit to existing repo)

- Reuse `internal/integrations/*` packages as connectors or add thin adapters under `connectors/*`.
- Expose API endpoints in `cmd/api` to: start run, list runs/steps, human approve/reject, view logs.
- Add unit tests under `workflows/csr` and integration tests that run against a test Postgres instance (use `testcontainers` if desired).

Operational checklist for MVP

- DB migration applied
- Acquire DB credentials and connection pooling (pgx, max connections tuned)
- Configure metrics endpoint and logging sink
- Configure secrets for external connectors
- Run e2e test for CSR workflow including human-in-loop simulated approval

Appendix: Example claim query (Postgres, safe for concurrency)

```sql
WITH c AS (
  SELECT id FROM workflow_steps
  WHERE status = 'pending' AND (next_attempt_at IS NULL OR next_attempt_at <= now())
  ORDER BY seq, created_at
  FOR UPDATE SKIP LOCKED
  LIMIT 1
)
UPDATE workflow_steps ws
SET status = 'in_progress', lock_owner = $1, updated_at = now()
FROM c
WHERE ws.id = c.id
RETURNING ws.*;
```

Design tradeoffs and reasoning

- Avoiding external systems reduces operational complexity and latency; DB-backed queue + SKIP LOCKED is robust for single-process or co-located workers.
- LISTEN/NOTIFY is used for low-latency wakeups but not relied on for durability — outbox ensures durable message handoff.
- Determinism achieved by persisting every state transition and letting workflow `Plan` be pure function of run state + inputs (avoid non-deterministic side-effects inside `Plan`).

Next steps

- Implement migrations and the `StateStore` interface in `internal/engine/store/postgres`.
- Implement the scheduler and worker pool in `internal/engine/scheduler`.
- Scaffold `workflows/csr` with typed steps and unit tests.

---

Revision log

- v0.1: initial comprehensive design and roadmap (MVP CSR). Iteration plan: 5 review passes to refine concurrency, idempotency, and event routing.

Contact / Ownership

- Proposed owner: backend platform team (names/team in your org)
- Suggested reviewers: architect, lead backend engineer, SRE

---

## Iterative Refinements (Pass 1 of 5)

This first pass tightens the design to improve determinism, recovery, and idempotency. The changes below will be applied to the next design iteration and to the initial scaffold:

- Add heartbeat / claim metadata to `workflow_steps` to make requeueing deterministic.

  - `claimed_at TIMESTAMP WITH TIME ZONE NULL`
  - `last_heartbeat TIMESTAMP WITH TIME ZONE NULL`
  - `lock_owner TEXT NULL`
  - Rationale: allows safe re-queueing of abandoned `in_progress` steps without relying on approximate timeouts only.

- Add uniqueness and dedup constraints to prevent duplicate step creation:

  - Unique index on `(run_id, step_name, idempotency_key)` where `idempotency_key IS NOT NULL`.
  - Unique constraint on `(run_id, seq)` to prevent conflicting sequence inserts.

- Make `Plan` pure and deterministic:

  - `Plan(ctx, run)` must compute next logical steps only from `run` and persisted data. Side-effects (sending emails, calls) must be performed exclusively in `Execute` steps.
  - Rationale: deterministic plans simplify reasoning during replay and version upgrades.

- Improve scheduler claim query (explicit fields):

  - Claim transaction sets `status='in_progress'`, `lock_owner`, `claimed_at=now()`, and `last_heartbeat=now()`.
  - Requeue query moves steps from `in_progress` back to `pending` only when `last_heartbeat < now() - heartbeat_ttl` and `attempts < max_attempts`.

- Heartbeat semantics & worker liveness:

  - Workers should update `last_heartbeat` at a short interval (e.g., 30s) for long-running steps.
  - A separate monitor (scheduler subroutine) will reset stale claims; this is safer than relying solely on process termination detection.

- Connector side-effect idempotency keys:

  - Each connector operation (email send, CRM update) must accept an `idempotency_key` and return a stable provider-side identifier to be stored in `workflow_steps.result`.
  - This ensures safe retries without duplicate external actions.

- Add explicit `workflow_versions` metadata to runs and an upgrade policy:

  - `workflow_runs.workflow_version` must be compared to registry version before executing `Plan` or `Execute`.
  - If run version mismatches, scheduler will either: (a) run with a compatibility shim, (b) mark for migration, or (c) process under previous version if still registered.

- Schema additions (concise SQL snippets to include in next migration):

```sql
ALTER TABLE workflow_steps
  ADD COLUMN claimed_at TIMESTAMP WITH TIME ZONE,
  ADD COLUMN last_heartbeat TIMESTAMP WITH TIME ZONE;

CREATE UNIQUE INDEX ux_workflow_steps_run_step_idempotency ON workflow_steps(run_id, step_name, idempotency_key) WHERE idempotency_key IS NOT NULL;
CREATE UNIQUE INDEX ux_workflow_steps_run_seq ON workflow_steps(run_id, seq);
```

Action items for Pass 1 implementation

- Update GORM entity structs to include `ClaimedAt` and `LastHeartbeat` fields.
- Implement worker heartbeat update in `internal/engine/scheduler` worker loop.
- Add requeue routine that resets stale `in_progress` steps back to `pending` with backoff applied.
- Add DB migration file (or GORM model update) to include new columns and indexes.

Notes on iteration process

- We will perform four more passes focused on:
  1. Event routing optimization (LISTEN/NOTIFY tradeoffs, batching)
  2. Stronger versioning and migration tooling
  3. Performance tuning (bulk claims, batching plan evaluations)
  4. Observability & operational runbooks

## Iterative Refinements (Pass 2 of 5) — Event Routing Optimization

Goals

- Minimize latency between step completion and subsequent step scheduling while keeping durability and avoiding tight DB polling loops.
- Ensure outbox-based reliable delivery, deduplicated consumers, and scalable publisher behavior.

Decisions

- Use DB Outbox as the single source of truth for events. All external-facing events are written inside the same transaction that updates step state.
- Implement a hybrid wake strategy:
  - Primary: small poll interval (configurable) for outbox and pending steps to support environments where LISTEN/NOTIFY may be limited.
  - Fast-path: `LISTEN/NOTIFY` to wake a dispatcher/scheduler when available — keep it optional and protected behind feature flag.
- Publisher: a robust in-process outbox publisher that reads un-published events in batches, publishes them to the connector, and marks them published in the same or a follow-up transaction.

Key details

- Outbox rows have states: `pending`, `in_flight`, `published`, `failed`. The publisher moves a batch into `in_flight` (UPDATE ... WHERE id IN ...) to claim them, then attempts delivery.
- Use idempotency keys for event consumers; include an `outbox_events.id` and `idempotency_key` in event payload so receivers can dedupe.
- Batch size and publisher concurrency must be configurable; default batch size = 50, workers = 2.
- For in-process event routing, maintain subscriber channels keyed by `event_type` and use non-blocking sends to avoid stalling the publisher.

Action items — Pass 2

- Add `state` column to `outbox_events` (`pending|in_flight|published|failed`) and `idempotency_key`.
- Implement `outbox/publisher` component (configurable batching, retries, backoff).
- Add optional `pg_notify` calls after writing outbox row — dispatcher listens to channel to wake the scheduler.
- Create tests for publisher durability: simulate transient connector failures, assert publish retried and marked published once.

## Iterative Refinements (Pass 3 of 5) — Workflow Versioning & Migration

Goals

- Enable safe evolution of workflow definitions without corrupting or losing runs in-flight.
- Provide tools to migrate existing `workflow_runs` to new workflow definitions where required.

Decisions

- Each `Workflow` registration is versioned (semantic version `major.minor.patch`). When a run is created its `workflow_version` is pinned.
- Registry holds all active versions. If a run references an unregistered version, the scheduler will refuse to process it and emit an audit event.
- Migration approach options:
  - Compatibility-first: maintain older definitions in code until all runs complete.
  - Migrate-once: provide a one-time migration script that transforms stored run payloads and steps to the new schema.
  - Dual-processing (advanced): keep both versions active and route new runs to the new version.

Key details

- `Plan` and `Execute` must detect the run's version and either call compatible code paths or reject execution.
- Maintain a small compatibility shim layer: `type VersionedWorkflow struct { Version string; Workflow Workflow }` to pair versioned logic.

Action items — Pass 3

- Add clear docs and tests showing migration paths for bumping major/minor versions.
- Add registry APIs to list supported versions and to register compatibility shims.
- Provide a `cli` tool or admin API to apply migration scripts to runs (dry-run mode + apply).

## Iterative Refinements (Pass 4 of 5) — Performance & Scalability Tuning

Goals

- Ensure the single-process scheduler can handle expected workload for MVP and prepare patterns to allow safe scale-up (more processes) later.

Decisions & optimizations

- Bulk claim: allow claiming N steps per transaction instead of 1 for throughput (careful ordering by `seq` and priority). Use `FOR UPDATE SKIP LOCKED` with `LIMIT N`.
- Batching Plan evaluations: when a run finishes a step that schedules many follow-up steps, insert them with a single bulk insert to reduce DB overhead.
- Prepared statements and connection pooling: use `pgx` or configure GORM to reuse prepared statements; set DB pool size matching worker_count + connectors.
- Indexes: ensure indexes on `status`, `next_attempt_at`, `run_id`, `idempotency_key`, and `last_heartbeat` to make scheduler queries fast.
- Prioritization: add optional `priority INT` column on `workflow_steps` for urgent steps.

Key details

- Cap long-running steps by design: prefer moving long polls to background tasks that update step status (to avoid hogging worker goroutines).
- Monitor DB hotspots (locks, deadlocks) and adjust batch sizes accordingly.

Action items — Pass 4

- Add ability to claim batches (`ClaimPendingSteps(limit int)`), and tune batch size via config.
- Add indexes recommended above to migration or GORM models.
- Add a benchmarking harness: synthetic workloads to measure throughput with different batch sizes and worker counts.

## Iterative Refinements (Pass 5 of 5) — Observability, Alerts & Operational Runbooks

Goals

- Provide the operational visibility and runbooks SREs need to operate and recover the orchestration engine.

Observability decisions

- Logs: structured logs with the following fields by default: `workflow_type`, `workflow_version`, `run_id`, `step_id`, `step_name`, `seq`, `attempts`, `worker_id`, `error`, `latency_ms`.
- Traces: instrument scheduler claim + executor RunStep + connector calls with OpenTelemetry spans. Propagate trace context to external calls where possible.
- Metrics (Prometheus): counters and histograms for runs started, steps executed, step latency, retries, dead-letter counts, human-wait counts.
- Dashboards & alerts: create a Grafana dashboard and alerts for high failure rate (>5% over 5m), growing pending steps queue (>100), or rising retry counts.

Runbooks

- Recovery flow for a stuck step: how to inspect `step_logs`, requeue abandoned `in_progress` steps, and re-run specific step manually.
- How to perform a cross-version migration: steps to dry-run a migration script and apply it.
- How to perform a rollback: mark new workflow version as inactive, and re-pin recently created runs to previous version if necessary.

Action items — Pass 5

- Add instrumentation to key code paths and connectors (trace and metrics wrappers).
- Create Prometheus alerts and Grafana dashboards (templates + JSON export) for platform SREs.
- Draft recovery runbook and include as `plan/operational_runbook.md` (or append to this plan).

Iteration complete

- All five refinement passes are recorded in this plan. Each pass includes concrete decisions, schema changes (where applicable), and actionable implementation items to follow when code scaffolding or migrations are performed.
