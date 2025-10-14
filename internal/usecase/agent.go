package usecase

import (
	"errors"

	"github.com/alpinesboltltd/boltz-ai/internal/entity"
	"github.com/alpinesboltltd/boltz-ai/internal/repository"
	"gorm.io/gorm"
)

/* TODO:
1. Create agent
2. update agent
3. retrive agent by owner
4. retrive agent by id
5. retrive agents
*/

type AgentUsecase struct {
	Agent repository.AgentRepositoryInterface
}

func NewAgentUseCase(agentRepo repository.AgentRepositoryInterface) *AgentUsecase {
	return &AgentUsecase{
		Agent: agentRepo,
	}
}

func (u *AgentUsecase) CreateNewAgent(userId, name, description, aiModel string, agentType entity.AgentType, capabilities []string, credit_per_1k int, status entity.AgentStatus) (*entity.Agent, error) {
	agent, err := u.Agent.CreateAgent(userId, name, description, aiModel, agentType, capabilities, credit_per_1k, status)
	if err != nil {
		return nil, err
	}
	return agent, nil
}

func (u *AgentUsecase) UpdateAgent(id string, update entity.AgentUpdate) (*entity.Agent, error) {
	if err := u.Agent.UpdateAgentByID(id, update); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("record not found")
		}
		return nil, err
	}
	agent, err := u.Agent.GetAgent(id)
	if err != nil {
		return nil, err
	}
	return agent, nil
}

func (u *AgentUsecase) GetAgent(id string) (*entity.Agent, error) {
	agent, err := u.Agent.GetAgent(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("record not found")
		}
		return nil, err
	}
	return agent, nil
}

func (u *AgentUsecase) GetUserAgents(userId string) (*[]entity.Agent, error) {
	agent, err := u.Agent.GetAgentsByUserId(userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("record not found")
		}
		return nil, err
	}

	return agent, nil
}
