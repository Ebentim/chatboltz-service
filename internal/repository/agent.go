package repository

import (
	"time"

	"github.com/alpinesboltltd/boltz-ai/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AgentRepository struct {
	db *gorm.DB
}

func NewAgentRepository(db *gorm.DB) AgentRepositoryInterface {
	return &AgentRepository{db: db}
}

func (r *AgentRepository) CreateAgent(userId, name, description, aiModel string, agentType entity.AgentType, capabilities []string, credit_per_1k int, status entity.AgentStatus) (*entity.Agent, error) {
	agent := &entity.Agent{
		ID:           uuid.New().String(),
		UserId:       userId,
		Name:         name,
		Description:  description,
		AgentType:    agentType,
		AiModel:      aiModel,
		Capabilities: capabilities,
		CreditsPer1k: credit_per_1k,
		Status:       status,
		CreatedAt:    time.Now().UTC().Format(time.RFC3339),
		UpdatedAt:    time.Now().UTC().Format(time.RFC3339),
	}

	if err := r.db.Create(agent).Error; err != nil {
		return nil, err
	}

	return agent, nil
}

func (r *AgentRepository) UpdateAgentByID(id string, update entity.AgentUpdate) error {
	return r.db.Model(&entity.Agent{}).Where("id = ?", id).Updates(update).Error
}

func (r *AgentRepository) CreateAgentAppearance(agent_id, primary_color, font_family, chat_icon, welcome_message, position, icon_size, bubble_style string) (*entity.AgentAppearance, error) {
	appearance := &entity.AgentAppearance{
		ID:             uuid.New().String(),
		AgentId:        agent_id,
		PrimaryColor:   primary_color,
		FontFamily:     font_family,
		ChatIcon:       chat_icon,
		WelcomeMessage: welcome_message,
		Position:       position,
		IconSize:       icon_size,
		BubbleStyle:    bubble_style,
	}

	if err := r.db.Create(appearance).Error; err != nil {
		return nil, err
	}
	return appearance, nil
}

func (r *AgentRepository) CreateAgentBehavior(agent_id, fallback_message, Offline_message, system_instruction_id, prompt_template_id string, enable_human_handoff bool, temperature float64, max_tokens int) (*entity.AgentBehavior, error) {
	behavior := &entity.AgentBehavior{
		ID:                  uuid.New().String(),
		AgentId:             agent_id,
		FallbackMessage:     fallback_message,
		EnableHumanHandoff:  enable_human_handoff,
		OfflineMessage:      Offline_message,
		SystemInstructionId: system_instruction_id,
		PromptTemplateId:    prompt_template_id,
		Temperature:         temperature,
		MaxTokens:           max_tokens,
	}

	if err := r.db.Create(behavior).Error; err != nil {
		return nil, err
	}
	return behavior, nil
}

func (r *AgentRepository) CreateAgentStats(agent_id string, total_messages, unique_users, conversions_count int, average_rating, response_rate float64, last_calculated_at time.Time) (*entity.AgentStats, error) {
	stats := &entity.AgentStats{
		ID:               uuid.New().String(),
		AgentId:          agent_id,
		TotalMessages:    total_messages,
		UniqueUsers:      unique_users,
		ConversionsCount: conversions_count,
		AverageRating:    average_rating,
		ResponseRate:     response_rate,
		LastCalculatedAt: last_calculated_at,
	}

	if err := r.db.Create(stats).Error; err != nil {
		return nil, err
	}
	return stats, nil
}

func (r *AgentRepository) CreateAgentIntegrations(agent_id, api_key, api_secret string, integration_id []string, is_active bool) (*entity.AgentIntegration, error) {
	integration := &entity.AgentIntegration{
		ID:            uuid.New().String(),
		AgentId:       agent_id,
		IntegrationId: integration_id,
		ApiKey:        &api_key,
		ApiSecret:     &api_secret,
		IsActive:      is_active,
	}

	if err := r.db.Create(integration).Error; err != nil {
		return nil, err
	}
	return integration, nil
}

func (r *AgentRepository) UpdateAgent(agent *entity.Agent, changes map[string]interface{}) error {
	changes["updated_at"] = time.Now().UTC().Format(time.RFC3339)
	return r.db.Model(&entity.Agent{}).Where("id = ?", agent.ID).Updates(changes).Error
}

func (r *AgentRepository) UpdateAgentAppearance(appearance *entity.AgentAppearance) error {
	appearance.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	return r.db.Save(appearance).Error
}

func (r *AgentRepository) UpdateAgentBehavior(behavior *entity.AgentBehavior) error {
	behavior.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	return r.db.Save(behavior).Error
}

func (r *AgentRepository) UpdateAgentStats(stats *entity.AgentStats) error {
	stats.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	return r.db.Save(stats).Error
}

func (r *AgentRepository) UpdateAgentIntegration(integration *entity.AgentIntegration) error {
	integration.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	return r.db.Save(integration).Error
}

func (r *AgentRepository) GetAgent(id string) (*entity.Agent, error) {
	var agent entity.Agent
	if err := r.db.Where("id = ?", id).First(&agent).Error; err != nil {
		return nil, err
	}
	return &agent, nil
}

// TODO: GET AGENT BY OWNER ID
func (r *AgentRepository) GetAgentsByUserId(userId string) (*[]entity.Agent, error) {
	var agents []entity.Agent
	if err := r.db.Where("user_id = ?", userId).Find(&agents).Error; err != nil {
		return nil, err
	}
	return &agents, nil
}

func (r *AgentRepository) GetAgentAppearance(agent_id string) (*entity.AgentAppearance, error) {
	var appearance entity.AgentAppearance
	if err := r.db.Where("agent_id = ?", agent_id).First(&appearance).Error; err != nil {
		return nil, err
	}
	return &appearance, nil

}

func (r *AgentRepository) GetAgentBehavior(agent_id string) (*entity.AgentBehavior, error) {
	var behavior entity.AgentBehavior
	if err := r.db.Where("agent_id = ?", agent_id).First(&behavior).Error; err != nil {
		return nil, err
	}
	return &behavior, nil
}

func (r *AgentRepository) GetAgentStats(agent_id string) (*entity.AgentStats, error) {
	var stats entity.AgentStats
	if err := r.db.Where("agent_id = ?", agent_id).First(&stats).Error; err != nil {
		return nil, err
	}
	return &stats, nil
}

func (r *AgentRepository) GetAgentIntegrations(agent_id string) (*entity.AgentIntegration, error) {
	var integration entity.AgentIntegration
	if err := r.db.Where("agent_id = ?", agent_id).First(&integration).Error; err != nil {
		return nil, err
	}
	return &integration, nil
}

func (r *AgentRepository) DeleteAgent(agent_id string) error {
	if err := r.db.Where("id = ?", agent_id).Delete(&entity.Agent{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *AgentRepository) ListAllAgents() (*[]entity.Agent, error) {
	var agents []entity.Agent
	if err := r.db.Find(&agents).Error; err != nil {
		return nil, err
	}
	return &agents, nil

}
