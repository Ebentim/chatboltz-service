package store

import (
	"context"
	"log"
	"time"

	eng "github.com/alpinesboltltd/boltz-ai/internal/engine"
	"github.com/alpinesboltltd/boltz-ai/internal/entity"
	"gorm.io/gorm"
)

// PostgresStore is a GORM-backed StateStore implementation (scaffolded).
type PostgresStore struct {
	db *gorm.DB
}

func NewPostgresStore(db *gorm.DB) *PostgresStore {
	return &PostgresStore{db: db}
}

func toEngineRun(e *entity.WorkflowRun) *eng.WorkflowRun {
	if e == nil {
		return nil
	}
	return &eng.WorkflowRun{
		ID: e.ID, WorkflowType: e.WorkflowType, WorkflowVersion: e.WorkflowVersion,
		Status: e.Status, Payload: e.Payload, CreatedAt: e.CreatedAt, UpdatedAt: e.UpdatedAt,
	}
}

func toEntityRun(r *eng.WorkflowRun) *entity.WorkflowRun {
	if r == nil {
		return nil
	}
	return &entity.WorkflowRun{
		ID: r.ID, WorkflowType: r.WorkflowType, WorkflowVersion: r.WorkflowVersion,
		Status: r.Status, Payload: r.Payload,
	}
}

func (s *PostgresStore) CreateRun(ctx context.Context, run *eng.WorkflowRun) error {
	ent := toEntityRun(run)
	return s.db.WithContext(ctx).Create(ent).Error
}

func (s *PostgresStore) LoadRun(ctx context.Context, runID string) (*eng.WorkflowRun, error) {
	var ent entity.WorkflowRun
	if err := s.db.WithContext(ctx).First(&ent, "id = ?", runID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return toEngineRun(&ent), nil
}

func (s *PostgresStore) InsertSteps(ctx context.Context, steps []*eng.WorkflowStepRecord) error {
	if len(steps) == 0 {
		return nil
	}
	ents := make([]*entity.WorkflowStep, 0, len(steps))
	now := time.Now()
	for _, st := range steps {
		ents = append(ents, &entity.WorkflowStep{
			ID: st.ID, RunID: st.RunID, StepName: st.StepName, Seq: st.Seq,
			Status: st.Status, Input: st.Input, Result: st.Result, Attempts: st.Attempts,
			MaxAttempts: st.MaxAttempts, NextAttemptAt: st.NextAttemptAt, CreatedAt: now, UpdatedAt: now,
		})
	}
	return s.db.WithContext(ctx).Create(&ents).Error
}

func (s *PostgresStore) ClaimNextStep(ctx context.Context, workerID string) (*eng.WorkflowStepRecord, error) {
	// Use a transaction and raw SQL to perform SELECT ... FOR UPDATE SKIP LOCKED + UPDATE ... RETURNING
	var out entity.WorkflowStep
	tx := s.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

		// raw SQL using WITH c as (...) update ... returning *
		query := `WITH c AS (
  SELECT id FROM workflow_steps
  WHERE status = 'pending' AND (next_attempt_at IS NULL OR next_attempt_at <= now())
  ORDER BY seq, created_at
  FOR UPDATE SKIP LOCKED
  LIMIT 1
)
UPDATE workflow_steps ws
SET status = 'in_progress', lock_owner = ?, claimed_at = now(), last_heartbeat = now(), updated_at = now()
FROM c
WHERE ws.id = c.id
RETURNING ws.*;`

	// Execute the query via GORM and scan into the entity struct.
	// Use RowsAffected to detect no-result case and return (nil, nil).
	res := tx.Raw(query, workerID).Scan(&out)
	if res.Error != nil {
		tx.Rollback()
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		tx.Rollback()
		return nil, nil
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// map to engine model
	rec := &eng.WorkflowStepRecord{
		ID: out.ID, RunID: out.RunID, StepName: out.StepName, Seq: out.Seq, Status: out.Status,
		Input: out.Input, Result: out.Result, Attempts: out.Attempts, MaxAttempts: out.MaxAttempts,
		NextAttemptAt: out.NextAttemptAt, ClaimedAt: out.ClaimedAt, LastHeartbeat: out.LastHeartbeat,
		LockOwner: out.LockOwner, IdempotencyKey: out.IdempotencyKey, Error: out.Error,
		CreatedAt: out.CreatedAt, UpdatedAt: out.UpdatedAt,
	}
	return rec, nil
}

func (s *PostgresStore) UpdateStep(ctx context.Context, step *eng.WorkflowStepRecord) error {
	// Map to entity and update
	ent := &entity.WorkflowStep{ID: step.ID}
	// Use map updates to avoid overwriting fields unintentionally
	updates := map[string]interface{}{
		"status":          step.Status,
		"result":          step.Result,
		"attempts":        step.Attempts,
		"next_attempt_at": step.NextAttemptAt,
		"lock_owner":      step.LockOwner,
		"error":           step.Error,
		"updated_at":      time.Now(),
	}
	return s.db.WithContext(ctx).Model(ent).Updates(updates).Error
}

func (s *PostgresStore) AppendLog(ctx context.Context, logRec *eng.StepLog) error {
	ent := &entity.StepLog{
		ID: logRec.ID, StepID: logRec.StepID, Level: logRec.Level, Message: logRec.Message, Meta: logRec.Meta,
	}
	return s.db.WithContext(ctx).Create(ent).Error
}

func (s *PostgresStore) EnqueueEvent(ctx context.Context, ev *eng.OutboxEvent) error {
	ent := &entity.OutboxEvent{
		ID: ev.ID, EventType: ev.EventType, Payload: ev.Payload, State: ev.State, Published: ev.Published, IdempotencyKey: ev.IdempotencyKey,
	}
	return s.db.WithContext(ctx).Create(ent).Error
}

// RequeueStaleSteps finds workflow_steps stuck in 'in_progress' whose
// last_heartbeat is older than heartbeatTTLSeconds, increments attempts,
// sets next_attempt_at (simple linear backoff), and moves them back to
// 'pending' unless attempts >= max_attempts in which case mark failed.
func (s *PostgresStore) RequeueStaleSteps(ctx context.Context, heartbeatTTLSeconds int, limit int) (int, error) {
	if limit <= 0 {
		limit = 100
	}
	// Run inside a transaction to select-for-update the candidates
	tx := s.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return 0, tx.Error
	}

	var rows []entity.WorkflowStep
	sel := tx.Raw(`
		SELECT * FROM workflow_steps
		WHERE status = 'in_progress' AND (last_heartbeat IS NULL OR last_heartbeat < now() - (? * INTERVAL '1 second'))
		ORDER BY last_heartbeat
		LIMIT ?
		FOR UPDATE SKIP LOCKED
	`, heartbeatTTLSeconds, limit).Scan(&rows)
	if sel.Error != nil {
		tx.Rollback()
		return 0, sel.Error
	}
	if len(rows) == 0 {
		tx.Rollback()
		return 0, nil
	}

	now := time.Now()
	requeued := 0
	const maxBackoffSeconds = 3600 // 1 hour cap

	for _, r := range rows {
		newAttempts := r.Attempts + 1
		// determine max attempts (use entity default if zero)
		maxAttempts := r.MaxAttempts
		if maxAttempts == 0 {
			maxAttempts = 5
		}

		if newAttempts >= maxAttempts {
			// mark failed/dead-letter
			if err := tx.Exec(`UPDATE workflow_steps SET status='failed', attempts = ?, updated_at = ? WHERE id = ?`, newAttempts, now, r.ID).Error; err != nil {
				tx.Rollback()
				return requeued, err
			}
			continue
		}

		// exponential backoff: base * 2^(attempts-1)
		base := heartbeatTTLSeconds
		if base <= 0 {
			base = 5
		}
		backoff := base * (1 << (newAttempts - 1))
		if backoff > maxBackoffSeconds {
			backoff = maxBackoffSeconds
		}

		if err := tx.Exec(`UPDATE workflow_steps SET status='pending', attempts = ?, next_attempt_at = now() + (? * INTERVAL '1 second'), updated_at = ? WHERE id = ?`, newAttempts, backoff, now, r.ID).Error; err != nil {
			tx.Rollback()
			return requeued, err
		}
		requeued++
	}

	if err := tx.Commit().Error; err != nil {
		return requeued, err
	}
	// log summary
	if requeued > 0 {
		log.Printf("requeue: requeued %d stale steps (heartbeatTTL=%ds)", requeued, heartbeatTTLSeconds)
	}
	return requeued, nil
}

func (s *PostgresStore) HeartbeatStep(ctx context.Context, stepID string) error {
	// Update last_heartbeat to now()
	return s.db.WithContext(ctx).Model(&entity.WorkflowStep{}).Where("id = ?", stepID).Update("last_heartbeat", time.Now()).Error
}
