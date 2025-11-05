package repository

import (
	"time"

	"github.com/alpinesboltltd/boltz-ai/internal/entity"
)

type UserRepositoryInterface interface {
	CreateUser(firebaseUID, name, email string) (*entity.Users, error)
	GetUserByFirebaseUID(firebaseUID string) (*entity.Users, error)
	GetUserByEmail(email string) (*entity.Users, error)
	GetUserByID(id string) (*entity.Users, error)
	UpdateUser(user *entity.Users) error
	DeleteUser(id string) error
	ListUsers() ([]*entity.Users, error)
}

type AgentRepositoryInterface interface {
	CreateAgent(userId, name, description, aiModel, aiProvider string, agentType entity.AgentType, credit_per_1k int, status entity.AgentStatus) (*entity.Agent, error)
	CreateAgentAppearance(agent_id, primary_color, font_family, chat_icon, welcome_message, position, icon_size, bubble_style string) (*entity.AgentAppearance, error)
	CreateAgentBehavior(agent_id, fallback_message, Offline_message, system_instruction_id, prompt_template_id string, enable_human_handoff bool, temperature float64, max_tokens int) (*entity.AgentBehavior, error)
	CreateAgentChannel(agent_id string, channel_id []string) (*entity.AgentChannel, error)
	CreateAgentStats(agent_id string, total_messages, unique_users, conversions_count int, average_rating, response_rate float64, last_calculated_at time.Time) (*entity.AgentStats, error)
	CreateAgentIntegrations(agent_id, api_key, api_secret string, integration_id []string, is_active bool) (*entity.AgentIntegration, error)
	UpdateAgent(agent *entity.Agent, changes map[string]interface{}) error
	UpdateAgentByID(id string, update entity.AgentUpdate) error
	UpdateAgentAppearance(appearance *entity.AgentAppearance) error
	UpdateAgentBehavior(behavior *entity.AgentBehavior) error
	UpdateAgentChannel(channel *entity.AgentChannel) error
	UpdateAgentStats(stats *entity.AgentStats) error
	UpdateAgentIntegration(integration *entity.AgentIntegration) error
	GetAgent(id string) (*entity.Agent, error)
	GetAgentWithDetails(id string) (*entity.Agent, error)
	GetAgentsByUserId(userId string) (*[]entity.Agent, error)
	GetAgentAppearance(agent_id string) (*entity.AgentAppearance, error)
	GetAgentBehavior(agent_id string) (*entity.AgentBehavior, error)
	GetAgentChannel(agent_id string) (*entity.AgentChannel, error)
	GetAgentStats(agent_id string) (*entity.AgentStats, error)
	GetAgentIntegrations(agent_id string) (*entity.AgentIntegration, error)
	DeleteAgent(agent_id, user_id string) error
	DeleteAgentAppearance(agent_id string) error
	DeleteAgentBehavior(agent_id string) error
	DeleteAgentChannel(agent_id string) error
	DeleteAgentStats(agent_id string) error
	DeleteAgentIntegration(agent_id string) error
	ListAllAgents() (*[]entity.Agent, error)
}

type SystemRepositoryInterface interface {
	CreateSystemInstruction(title, content, createdBy string, templateId *string) (*entity.SystemInstruction, error)
	GetSystemInstruction(id string) (*entity.SystemInstruction, error)
	UpdateSystemInstruction(instruction *entity.SystemInstruction) error
	DeleteSystemInstruction(id string) error
	ListSystemInstructions() (*[]entity.SystemInstruction, error)
	CreatePromptTemplate(title, content string) (*entity.PromptTemplate, error)
	GetPromptTemplate(id string) (*entity.PromptTemplate, error)
	ListPromptTemplates() (*[]entity.PromptTemplate, error)
}
