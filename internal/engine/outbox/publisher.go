package outbox

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/alpinesboltltd/boltz-ai/internal/entity"
	"github.com/alpinesboltltd/boltz-ai/internal/provider/smtp"
	"gorm.io/gorm"
)

// Publisher publishes outbox events reliably. For now supports `email_send` events.
func StartPublisher(ctx context.Context, db *gorm.DB, smtpClient *smtp.Client, batchSize int, interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				processBatch(ctx, db, smtpClient, batchSize)
			}
		}
	}()
}

func processBatch(ctx context.Context, db *gorm.DB, smtpClient *smtp.Client, batchSize int) {
	// Claim batch
	var rows []entity.OutboxEvent
	tx := db.WithContext(ctx).Begin()
	if tx.Error != nil {
		log.Printf("outbox: begin tx error: %v", tx.Error)
		return
	}
	// Select pending events with FOR UPDATE SKIP LOCKED
	if err := tx.Raw(`SELECT * FROM outbox_events WHERE state = 'pending' ORDER BY created_at LIMIT ? FOR UPDATE SKIP LOCKED`, batchSize).Scan(&rows).Error; err != nil {
		tx.Rollback()
		log.Printf("outbox: select pending error: %v", err)
		return
	}
	if len(rows) == 0 {
		tx.Rollback()
		return
	}
	ids := make([]string, 0, len(rows))
	for _, r := range rows {
		ids = append(ids, r.ID)
	}
	// mark in_flight
	if err := tx.Exec(`UPDATE outbox_events SET state='in_flight' WHERE id IN (?)`, ids).Error; err != nil {
		tx.Rollback()
		log.Printf("outbox: mark in_flight error: %v", err)
		return
	}
	if err := tx.Commit().Error; err != nil {
		log.Printf("outbox: commit claim error: %v", err)
		return
	}

	// process each event outside transaction
	for _, ev := range rows {
		switch ev.EventType {
		case "email_send":
			var payload struct {
				To      string `json:"to"`
				Subject string `json:"subject"`
				Body    string `json:"body"`
				HTML    string `json:"html,omitempty"`
			}
			if err := json.Unmarshal(ev.Payload, &payload); err != nil {
				markFailed(db, ev.ID, err)
				continue
			}
			var err error
			if payload.HTML != "" {
				err = smtpClient.SendHTML(payload.To, payload.Subject, payload.HTML, payload.Body)
			} else {
				err = smtpClient.Send(payload.To, payload.Subject, payload.Body)
			}
			if err != nil {
				markFailed(db, ev.ID, err)
				continue
			}
			if err := db.WithContext(ctx).Exec(`UPDATE outbox_events SET state='published', published = true WHERE id = ?`, ev.ID).Error; err != nil {
				log.Printf("outbox: mark published error: %v", err)
			}
		default:
			// unknown event types: mark published to avoid loops
			if err := db.WithContext(ctx).Exec(`UPDATE outbox_events SET state='published', published = true WHERE id = ?`, ev.ID).Error; err != nil {
				log.Printf("outbox: mark published unknown type error: %v", err)
			}
		}
	}
}

func markFailed(db *gorm.DB, id string, err error) {
	log.Printf("outbox: publish failed for %s: %v", id, err)
	_ = db.Exec(`UPDATE outbox_events SET state='failed' WHERE id = ?`, id).Error
}
