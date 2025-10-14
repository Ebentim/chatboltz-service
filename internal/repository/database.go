package repository

import (
	"fmt"

	"github.com/alpinesboltltd/boltz-ai/internal/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(databaseURL string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.AutoMigrate(
		&entity.Users{},
		&entity.Agent{},
		&entity.AgentAppearance{},
		&entity.AgentBehavior{},
		&entity.AgentChannel{},
		&entity.AgentIntegration{},
		&entity.AgentStats{},
		&entity.TrainingData{},
	); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return db, nil
}
