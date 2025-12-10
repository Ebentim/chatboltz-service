package entity

import "time"

// GORM entities for orchestration engine
type WorkflowRun struct {
	ID              string `gorm:"type:uuid;primaryKey"`
	WorkflowType    string `gorm:"type:text;not null"`
	WorkflowVersion string `gorm:"type:text;not null"`
	Status          string `gorm:"type:text;not null"`
	Payload         []byte `gorm:"type:jsonb"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type WorkflowStep struct {
	ID             string `gorm:"type:uuid;primaryKey"`
	RunID          string `gorm:"type:uuid;index;not null"`
	StepName       string `gorm:"type:text;not null"`
	Seq            int    `gorm:"not null;default:0"`
	Status         string `gorm:"type:text;not null"`
	Input          []byte `gorm:"type:jsonb"`
	Result         []byte `gorm:"type:jsonb"`
	Attempts       int    `gorm:"default:0"`
	MaxAttempts    int    `gorm:"default:5"`
	NextAttemptAt  *time.Time
	ClaimedAt      *time.Time
	LastHeartbeat  *time.Time
	LockOwner      *string `gorm:"type:text"`
	IdempotencyKey *string `gorm:"type:text;index"`
	Error          *string `gorm:"type:text"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type OutboxEvent struct {
	ID             string  `gorm:"type:uuid;primaryKey"`
	EventType      string  `gorm:"type:text;not null"`
	Payload        []byte  `gorm:"type:jsonb"`
	State          string  `gorm:"type:text;default:'pending'"`
	IdempotencyKey *string `gorm:"type:text;index"`
	Published      bool    `gorm:"default:false"`
	CreatedAt      time.Time
}

type StepLog struct {
	ID        string `gorm:"type:uuid;primaryKey"`
	StepID    string `gorm:"type:uuid;index"`
	Level     string `gorm:"type:text"`
	Message   string `gorm:"type:text"`
	Meta      []byte `gorm:"type:jsonb"`
	CreatedAt time.Time
}

type RetryMeta struct {
	StepID        string `gorm:"type:uuid;primaryKey"`
	Attempts      int    `gorm:"default:0"`
	BackoffSecond int    `gorm:"default:0"`
	DeadLetter    bool   `gorm:"default:false"`
}
