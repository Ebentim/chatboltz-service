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

	// Run pre-migration fixes
	if err := runPreMigrationFixes(db); err != nil {
		return nil, fmt.Errorf("failed to run pre-migration fixes: %w", err)
	}

	if err := db.AutoMigrate(
		&entity.Users{},
		&entity.PromptTemplate{},
		&entity.SystemInstruction{},
		&entity.Channels{},
		&entity.Integrations{},
		&entity.ApiFunctions{},
		&entity.Agent{},
		&entity.AgentAppearance{},
		&entity.AgentBehavior{},
		&entity.AgentChannel{},
		&entity.AgentIntegration{},
		&entity.AgentStats{},
		&entity.TrainingData{},
		&entity.TrainingDocument{},
		&entity.DocumentChunk{},
		&entity.Workspace{},
		&entity.WorkspaceMember{},
	); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return db, nil
}

func runPreMigrationFixes(db *gorm.DB) error {
	// Clean up invalid foreign key references
	db.Exec("UPDATE agent_behaviors SET system_instruction_id = NULL WHERE system_instruction_id IS NOT NULL AND system_instruction_id NOT IN (SELECT id FROM system_instructions WHERE id IS NOT NULL)")
	db.Exec("UPDATE agent_behaviors SET prompt_template_id = NULL WHERE prompt_template_id IS NOT NULL AND prompt_template_id NOT IN (SELECT id FROM prompt_templates WHERE id IS NOT NULL)")

	// Fix array columns if they exist and are not proper arrays
	var exists bool

	// Check and fix agent_channels.channel_id
	db.Raw("SELECT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'agent_channels' AND column_name = 'channel_id' AND data_type != 'ARRAY')").Scan(&exists)
	if exists {
		db.Exec("ALTER TABLE agent_channels ALTER COLUMN channel_id TYPE text[] USING CASE WHEN channel_id ~ '^\\{.*\\}$' THEN channel_id::text[] ELSE string_to_array(channel_id, ',') END")
	}

	// Check and fix agent_integrations.integration_id
	db.Raw("SELECT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'agent_integrations' AND column_name = 'integration_id' AND data_type != 'ARRAY')").Scan(&exists)
	if exists {
		db.Exec("ALTER TABLE agent_integrations ALTER COLUMN integration_id TYPE text[] USING CASE WHEN integration_id ~ '^\\{.*\\}$' THEN integration_id::text[] ELSE string_to_array(integration_id, ',') END")
	}

	return nil
}
