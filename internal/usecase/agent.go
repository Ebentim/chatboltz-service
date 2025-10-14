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

func (u *AgentUsecase) CreateNewAgent(userId, name, description, aiModel string, agentType entity.AgentType, capabilities []string, credit_per_1k int, status entity.AgentStatus) (*entity.Agent, error) {
	if userId == "" || name == "" {
		return nil, appErrors.NewValidationError("User ID and name are required")
	}

	agent, err := u.Agent.CreateAgent(userId, name, description, aiModel, agentType, capabilities, credit_per_1k, status)
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
