package repository

import (
	"errors"
	"time"

	"github.com/alpinesboltltd/boltz-ai/internal/entity"
	appErrors "github.com/alpinesboltltd/boltz-ai/internal/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AgentRepository struct {
	db *gorm.DB
}

func NewAgentRepository(db *gorm.DB) AgentRepositoryInterface {
	return &AgentRepository{db: db}
}

func (r *AgentRepository) CreateAgent(userId, name, description, aiModelId string, agentType entity.AgentType, status entity.AgentStatus) (*entity.Agent, error) {
	agent := &entity.Agent{
		ID:          uuid.New().String(),
		UserId:      userId,
		Name:        name,
		Description: description,
		AgentType:   agentType,
		AiModelId:   aiModelId,
		Status:      status,
		CreatedAt:   time.Now().UTC().Format(time.RFC3339),
		UpdatedAt:   time.Now().UTC().Format(time.RFC3339),
	}

	if err := r.db.Create(agent).Error; err != nil {
		return nil, appErrors.WrapDatabaseError(err, "create agent")
	}

	return agent, nil
}

func (r *AgentRepository) UpdateAgentByID(id string, update entity.AgentUpdate) error {
	if err := r.db.Model(&entity.Agent{}).Where("id = ?", id).Updates(update).Error; err != nil {
		return appErrors.WrapDatabaseError(err, "update agent by ID")
	}
	return nil
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
		CreatedAt:      time.Now().UTC().Format(time.RFC3339),
		UpdatedAt:      time.Now().UTC().Format(time.RFC3339),
	}

	if err := r.db.Create(appearance).Error; err != nil {
		return nil, appErrors.WrapDatabaseError(err, "create agent appearance")
	}
	return appearance, nil
}

func (r *AgentRepository) CreateAgentBehavior(agent_id, fallback_message, Offline_message, system_instruction_id, prompt_template_id string, enable_human_handoff bool, temperature float64, max_tokens int) (*entity.AgentBehavior, error) {
	behavior := &entity.AgentBehavior{
		ID:                 uuid.New().String(),
		AgentId:            agent_id,
		FallbackMessage:    fallback_message,
		EnableHumanHandoff: enable_human_handoff,
		OfflineMessage:     Offline_message,
		Temperature:        temperature,
		MaxTokens:          max_tokens,
		CreatedAt:          time.Now().UTC().Format(time.RFC3339),
		UpdatedAt:          time.Now().UTC().Format(time.RFC3339),
	}

	if system_instruction_id != "" {
		behavior.SystemInstructionId = &system_instruction_id
	}
	if prompt_template_id != "" {
		behavior.PromptTemplateId = &prompt_template_id
	}

	if err := r.db.Create(behavior).Error; err != nil {
		return nil, appErrors.WrapDatabaseError(err, "create agent behavior")
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
		CreatedAt:        time.Now().UTC().Format(time.RFC3339),
		UpdatedAt:        time.Now().UTC().Format(time.RFC3339),
	}

	if err := r.db.Create(stats).Error; err != nil {
		return nil, appErrors.WrapDatabaseError(err, "create agent stats")
	}
	return stats, nil
}

func (r *AgentRepository) CreateAgentChannel(agent_id string, channel_id []string) (*entity.AgentChannel, error) {
	channel := &entity.AgentChannel{
		ID:        uuid.New().String(),
		AgentId:   agent_id,
		ChannelId: entity.StringArray(channel_id),
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
	}

	if err := r.db.Create(channel).Error; err != nil {
		return nil, appErrors.WrapDatabaseError(err, "create agent channel")
	}
	return channel, nil
}

func (r *AgentRepository) CreateAgentIntegrations(agent_id, api_key, api_secret string, integration_id []string, is_active bool) (*entity.AgentIntegration, error) {
	integration := &entity.AgentIntegration{
		ID:            uuid.New().String(),
		AgentId:       agent_id,
		IntegrationId: entity.StringArray(integration_id),
		ApiKey:        &api_key,
		ApiSecret:     &api_secret,
		IsActive:      is_active,
		CreatedAt:     time.Now().UTC().Format(time.RFC3339),
		UpdatedAt:     time.Now().UTC().Format(time.RFC3339),
	}

	if err := r.db.Create(integration).Error; err != nil {
		return nil, appErrors.WrapDatabaseError(err, "create agent integration")
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

func (r *AgentRepository) UpdateAgentChannel(channel *entity.AgentChannel) error {
	channel.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	return r.db.Save(channel).Error
}

func (r *AgentRepository) UpdateAgentIntegration(integration *entity.AgentIntegration) error {
	integration.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	return r.db.Save(integration).Error
}

func (r *AgentRepository) GetAgent(id string) (*entity.Agent, error) {
	var agent entity.Agent
	if err := r.db.Preload("User").Preload("AiModel").Preload("AgentAppearance").Preload("AgentBehavior").Preload("AgentBehavior.SystemInstruction").Preload("AgentBehavior.PromptTemplate").Preload("AgentChannel").Preload("AgentIntegration").Preload("AgentStats").Preload("TrainingData").Where("id = ?", id).First(&agent).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.NewNotFoundError("Agent not found")
		}
		return nil, appErrors.WrapDatabaseError(err, "get agent")
	}
	return &agent, nil
}

func (r *AgentRepository) GetAgentsByUserId(userId string) (*[]entity.Agent, error) {
	var agents []entity.Agent
	if err := r.db.Preload("AiModel").Where("user_id = ?", userId).Find(&agents).Error; err != nil {
		return nil, appErrors.WrapDatabaseError(err, "get agents by user ID")
	}
	return &agents, nil
}

func (r *AgentRepository) GetAgentAppearance(agent_id string) (*entity.AgentAppearance, error) {
	var appearance entity.AgentAppearance
	if err := r.db.Where("agent_id = ?", agent_id).First(&appearance).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.NewNotFoundError("Agent appearance not found")
		}
		return nil, appErrors.WrapDatabaseError(err, "get agent appearance")
	}
	return &appearance, nil
}

func (r *AgentRepository) GetAgentBehavior(agent_id string) (*entity.AgentBehavior, error) {
	var behavior entity.AgentBehavior
	if err := r.db.Where("agent_id = ?", agent_id).First(&behavior).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.NewNotFoundError("Agent behavior not found")
		}
		return nil, appErrors.WrapDatabaseError(err, "get agent behavior")
	}
	return &behavior, nil
}

func (r *AgentRepository) GetAgentStats(agent_id string) (*entity.AgentStats, error) {
	var stats entity.AgentStats
	if err := r.db.Where("agent_id = ?", agent_id).First(&stats).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.NewNotFoundError("Agent stats not found")
		}
		return nil, appErrors.WrapDatabaseError(err, "get agent stats")
	}
	return &stats, nil
}

func (r *AgentRepository) GetAgentChannel(agent_id string) (*entity.AgentChannel, error) {
	var channel entity.AgentChannel
	if err := r.db.Where("agent_id = ?", agent_id).First(&channel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.NewNotFoundError("Agent channel not found")
		}
		return nil, appErrors.WrapDatabaseError(err, "get agent channel")
	}
	return &channel, nil
}

func (r *AgentRepository) GetAgentIntegrations(agent_id string) (*entity.AgentIntegration, error) {
	var integration entity.AgentIntegration
	if err := r.db.Where("agent_id = ?", agent_id).First(&integration).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.NewNotFoundError("Agent integration not found")
		}
		return nil, appErrors.WrapDatabaseError(err, "get agent integration")
	}
	return &integration, nil
}

func (r *AgentRepository) DeleteAgent(agent_id, user_id string) error {
	if err := r.db.Where("id = ? AND user_id = ?", agent_id, user_id).Delete(&entity.Agent{}).Error; err != nil {
		return appErrors.WrapDatabaseError(err, "delete agent")
	}
	return nil
}

func (r *AgentRepository) DeleteAgentAppearance(agent_id string) error {
	if err := r.db.Where("agent_id = ?", agent_id).Delete(&entity.AgentAppearance{}).Error; err != nil {
		return appErrors.WrapDatabaseError(err, "delete agent appearance")
	}
	return nil
}

func (r *AgentRepository) DeleteAgentBehavior(agent_id string) error {
	if err := r.db.Where("agent_id = ?", agent_id).Delete(&entity.AgentBehavior{}).Error; err != nil {
		return appErrors.WrapDatabaseError(err, "delete agent behavior")
	}
	return nil
}

func (r *AgentRepository) DeleteAgentChannel(agent_id string) error {
	if err := r.db.Where("agent_id = ?", agent_id).Delete(&entity.AgentChannel{}).Error; err != nil {
		return appErrors.WrapDatabaseError(err, "delete agent channel")
	}
	return nil
}

func (r *AgentRepository) DeleteAgentStats(agent_id string) error {
	if err := r.db.Where("agent_id = ?", agent_id).Delete(&entity.AgentStats{}).Error; err != nil {
		return appErrors.WrapDatabaseError(err, "delete agent stats")
	}
	return nil
}

func (r *AgentRepository) DeleteAgentIntegration(agent_id string) error {
	if err := r.db.Where("agent_id = ?", agent_id).Delete(&entity.AgentIntegration{}).Error; err != nil {
		return appErrors.WrapDatabaseError(err, "delete agent integration")
	}
	return nil
}

func (r *AgentRepository) ListAllAgents() (*[]entity.Agent, error) {
	var agents []entity.Agent
	if err := r.db.Find(&agents).Error; err != nil {
		return nil, appErrors.WrapDatabaseError(err, "list all agents")
	}
	return &agents, nil
}

// Stub implementations to satisfy interface; actual logic should be added elsewhere
func (r *AgentRepository) CreateTrainingData(agentID, contentType string, content []entity.TrainingTexts, isActive bool) (*entity.TrainingData, error) {
	return nil, appErrors.NewNotFoundError("not implemented")
}
func (r *AgentRepository) GetTrainingDataByAgentID(agentID string) ([]entity.TrainingData, error) {
	return nil, appErrors.NewNotFoundError("not implemented")
}
func (r *AgentRepository) UpdateTrainingData(trainingData *entity.TrainingData) error {
	return appErrors.NewNotFoundError("not implemented")
}
func (r *AgentRepository) DeleteTrainingData(id string) error {
	return appErrors.NewNotFoundError("not implemented")
}
