package usecase

import (
	"github.com/alpinesboltltd/boltz-ai/internal/entity"
	appErrors "github.com/alpinesboltltd/boltz-ai/internal/errors"
	"github.com/alpinesboltltd/boltz-ai/internal/repository"
)

type AgentUsecase struct {
	Agent repository.AgentRepositoryInterface
}

func NewAgentUseCase(agentRepo repository.AgentRepositoryInterface) *AgentUsecase {
	return &AgentUsecase{
		Agent: agentRepo,
	}
}

func (u *AgentUsecase) CreateNewAgent(userId, name, description, aiModelId string, agentType entity.AgentType, status entity.AgentStatus) (*entity.Agent, error) {
	if userId == "" || name == "" {
		return nil, appErrors.NewValidationError("User ID and name are required")
	}

	agent, err := u.Agent.CreateAgent(userId, name, description, aiModelId, agentType, status)
	if err != nil {
		return nil, err
	}
	return agent, nil
}

func (u *AgentUsecase) UpdateAgent(id string, update entity.AgentUpdate) (*entity.Agent, error) {
	if id == "" {
		return nil, appErrors.NewValidationError("Agent ID is required")
	}

	if err := u.Agent.UpdateAgentByID(id, update); err != nil {
		return nil, err
	}
	agent, err := u.Agent.GetAgent(id)
	if err != nil {
		return nil, err
	}
	return agent, nil
}

func (u *AgentUsecase) GetAgent(id string) (*entity.Agent, error) {
	if id == "" {
		return nil, appErrors.NewValidationError("Agent ID is required")
	}

	agent, err := u.Agent.GetAgent(id)
	if err != nil {
		return nil, err
	}
	return agent, nil
}

func (u *AgentUsecase) GetUserAgents(userId string) (*[]entity.Agent, error) {
	if userId == "" {
		return nil, appErrors.NewValidationError("User ID is required")
	}

	agent, err := u.Agent.GetAgentsByUserId(userId)
	if err != nil {
		return nil, err
	}

	return agent, nil
}

func (u *AgentUsecase) DeleteAgent(id string, user_id string) error {
	if id == "" || user_id == "" {
		return appErrors.NewValidationError("Agent ID and user ID is required")
	}
	if err := u.Agent.DeleteAgent(id, user_id); err != nil {
		return err
	}
	return nil
}

func (u *AgentUsecase) DeleteAgentAppearance(agentId string) error {
	if agentId == "" {
		return appErrors.NewValidationError("Agent ID is required")
	}
	return u.Agent.DeleteAgentAppearance(agentId)
}

func (u *AgentUsecase) DeleteAgentBehavior(agentId string) error {
	if agentId == "" {
		return appErrors.NewValidationError("Agent ID is required")
	}
	return u.Agent.DeleteAgentBehavior(agentId)
}

func (u *AgentUsecase) DeleteAgentChannel(agentId string) error {
	if agentId == "" {
		return appErrors.NewValidationError("Agent ID is required")
	}
	return u.Agent.DeleteAgentChannel(agentId)
}

func (u *AgentUsecase) DeleteAgentStats(agentId string) error {
	if agentId == "" {
		return appErrors.NewValidationError("Agent ID is required")
	}
	return u.Agent.DeleteAgentStats(agentId)
}

func (u *AgentUsecase) DeleteAgentIntegration(agentId string) error {
	if agentId == "" {
		return appErrors.NewValidationError("Agent ID is required")
	}
	return u.Agent.DeleteAgentIntegration(agentId)
}

func (u *AgentUsecase) CreateAgentAppearance(appearance entity.AgentAppearance) (*entity.AgentAppearance, error) {
	if appearance.AgentId == "" {
		return nil, appErrors.NewValidationError("Agent ID is required")
	}

	if _, err := u.Agent.GetAgent(appearance.AgentId); err != nil {
		return nil, appErrors.NewNotFoundError("agent does not exist")
	}
	return u.Agent.CreateAgentAppearance(appearance.AgentId, appearance.PrimaryColor, appearance.FontFamily, appearance.ChatIcon, appearance.WelcomeMessage, appearance.Position, appearance.IconSize, appearance.BubbleStyle)
}

func (u *AgentUsecase) CreateAgentBehavior(behavior entity.AgentBehavior) (*entity.AgentBehavior, error) {
	if behavior.AgentId == "" {
		return nil, appErrors.NewValidationError("Agent ID is required")
	}

	if _, err := u.Agent.GetAgent(behavior.AgentId); err != nil {
		return nil, appErrors.NewNotFoundError("agent does not exist")
	}

	sysInstrId := ""
	promptTmplId := ""
	if behavior.SystemInstructionId != nil {
		sysInstrId = *behavior.SystemInstructionId
	}
	if behavior.PromptTemplateId != nil {
		promptTmplId = *behavior.PromptTemplateId
	}

	return u.Agent.CreateAgentBehavior(behavior.AgentId, behavior.FallbackMessage, behavior.OfflineMessage, sysInstrId, promptTmplId, behavior.EnableHumanHandoff, behavior.Temperature, behavior.MaxTokens)
}

func (u *AgentUsecase) CreateAgentChannel(channel entity.AgentChannel) (*entity.AgentChannel, error) {
	if channel.AgentId == "" {
		return nil, appErrors.NewValidationError("Agent ID is required")
	}

	if _, err := u.Agent.GetAgent(channel.AgentId); err != nil {
		return nil, appErrors.NewNotFoundError("agent does not exist")
	}
	return u.Agent.CreateAgentChannel(channel.AgentId, []string(channel.ChannelId))
}

func (u *AgentUsecase) CreateAgentIntegration(integration entity.AgentIntegration) (*entity.AgentIntegration, error) {
	if integration.AgentId == "" {
		return nil, appErrors.NewValidationError("Agent ID is required")
	}

	if _, err := u.Agent.GetAgent(integration.AgentId); err != nil {
		return nil, appErrors.NewNotFoundError("agent does not exist")
	}

	apiKey := ""
	apiSecret := ""
	if integration.ApiKey != nil {
		apiKey = *integration.ApiKey
	}
	if integration.ApiSecret != nil {
		apiSecret = *integration.ApiSecret
	}

	return u.Agent.CreateAgentIntegrations(integration.AgentId, apiKey, apiSecret, []string(integration.IntegrationId), integration.IsActive)
}

func (u *AgentUsecase) GetAgentAppearance(agentId string) (*entity.AgentAppearance, error) {
	if agentId == "" {
		return nil, appErrors.NewValidationError("Agent ID is required")
	}
	return u.Agent.GetAgentAppearance(agentId)
}

func (u *AgentUsecase) GetAgentBehavior(agentId string) (*entity.AgentBehavior, error) {
	if agentId == "" {
		return nil, appErrors.NewValidationError("Agent ID is required")
	}
	return u.Agent.GetAgentBehavior(agentId)
}

func (u *AgentUsecase) GetAgentChannel(agentId string) (*entity.AgentChannel, error) {
	if agentId == "" {
		return nil, appErrors.NewValidationError("Agent ID is required")
	}
	return u.Agent.GetAgentChannel(agentId)
}

func (u *AgentUsecase) GetAgentStats(agentId string) (*entity.AgentStats, error) {
	if agentId == "" {
		return nil, appErrors.NewValidationError("Agent ID is required")
	}
	return u.Agent.GetAgentStats(agentId)
}

func (u *AgentUsecase) GetAgentIntegration(agentId string) (*entity.AgentIntegration, error) {
	if agentId == "" {
		return nil, appErrors.NewValidationError("Agent ID is required")
	}
	return u.Agent.GetAgentIntegrations(agentId)
}

func (u *AgentUsecase) UpdateAgentAppearance(agentId string, appearance entity.AgentAppearance) (*entity.AgentAppearance, error) {
	if agentId == "" {
		return nil, appErrors.NewValidationError("Agent ID is required")
	}

	existing, err := u.Agent.GetAgentAppearance(agentId)
	if err != nil {
		return nil, err
	}

	if appearance.PrimaryColor != "" {
		existing.PrimaryColor = appearance.PrimaryColor
	}
	if appearance.FontFamily != "" {
		existing.FontFamily = appearance.FontFamily
	}
	if appearance.ChatIcon != "" {
		existing.ChatIcon = appearance.ChatIcon
	}
	if appearance.WelcomeMessage != "" {
		existing.WelcomeMessage = appearance.WelcomeMessage
	}
	if appearance.Position != "" {
		existing.Position = appearance.Position
	}
	if appearance.IconSize != "" {
		existing.IconSize = appearance.IconSize
	}
	if appearance.BubbleStyle != "" {
		existing.BubbleStyle = appearance.BubbleStyle
	}

	if err := u.Agent.UpdateAgentAppearance(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (u *AgentUsecase) UpdateAgentBehavior(agentId string, behavior entity.AgentBehavior) (*entity.AgentBehavior, error) {
	if agentId == "" {
		return nil, appErrors.NewValidationError("Agent ID is required")
	}

	existing, err := u.Agent.GetAgentBehavior(agentId)
	if err != nil {
		return nil, err
	}

	if behavior.FallbackMessage != "" {
		existing.FallbackMessage = behavior.FallbackMessage
	}
	if behavior.OfflineMessage != "" {
		existing.OfflineMessage = behavior.OfflineMessage
	}
	if behavior.SystemInstructionId != nil {
		existing.SystemInstructionId = behavior.SystemInstructionId
	}
	if behavior.PromptTemplateId != nil {
		existing.PromptTemplateId = behavior.PromptTemplateId
	}
	existing.EnableHumanHandoff = behavior.EnableHumanHandoff
	if behavior.Temperature != 0 {
		existing.Temperature = behavior.Temperature
	}
	if behavior.MaxTokens != 0 {
		existing.MaxTokens = behavior.MaxTokens
	}

	if err := u.Agent.UpdateAgentBehavior(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (u *AgentUsecase) UpdateAgentChannel(agentId string, channel entity.AgentChannel) (*entity.AgentChannel, error) {
	if agentId == "" {
		return nil, appErrors.NewValidationError("Agent ID is required")
	}

	existing, err := u.Agent.GetAgentChannel(agentId)
	if err != nil {
		return nil, err
	}

	if len(channel.ChannelId) > 0 {
		existing.ChannelId = channel.ChannelId
	}

	if err := u.Agent.UpdateAgentChannel(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (u *AgentUsecase) UpdateAgentIntegration(agentId string, integration entity.AgentIntegration) (*entity.AgentIntegration, error) {
	if agentId == "" {
		return nil, appErrors.NewValidationError("Agent ID is required")
	}

	existing, err := u.Agent.GetAgentIntegrations(agentId)
	if err != nil {
		return nil, err
	}

	if len(integration.IntegrationId) > 0 {
		existing.IntegrationId = integration.IntegrationId
	}
	if integration.ApiKey != nil {
		existing.ApiKey = integration.ApiKey
	}
	if integration.ApiSecret != nil {
		existing.ApiSecret = integration.ApiSecret
	}
	existing.IsActive = integration.IsActive

	if err := u.Agent.UpdateAgentIntegration(existing); err != nil {
		return nil, err
	}
	return existing, nil
}
