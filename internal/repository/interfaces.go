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
	CreateTrainingData(agentID, contentType string, content []entity.TrainingTexts, isActive bool) (*entity.TrainingData, error)
	GetTrainingDataByAgentID(agentID string) ([]entity.TrainingData, error)
	UpdateTrainingData(trainingData *entity.TrainingData) error
	DeleteTrainingData(id string) error
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

// RAGRepositoryInterface defines the contract for RAG data persistence operations.
// This interface abstracts database operations for training documents and chunks.
type RAGRepositoryInterface interface {
	// CreateTrainingDocument creates a new training document record in the database.
	// Parameters:
	//   - doc: Training document to create
	// Returns:
	//   - error: Any error that occurred during creation
	CreateTrainingDocument(doc *entity.TrainingDocument) error

	// GetTrainingDocumentsByAgentID retrieves all training documents for an agent.
	// This includes preloaded chunks for each document.
	// Parameters:
	//   - agentID: ID of the agent whose documents to retrieve
	// Returns:
	//   - []entity.TrainingDocument: Array of documents with chunks
	//   - error: Any error that occurred during retrieval
	GetTrainingDocumentsByAgentID(agentID string) ([]entity.TrainingDocument, error)

	// UpdateTrainingDocument updates an existing training document record.
	// Parameters:
	//   - doc: Training document with updated fields
	// Returns:
	//   - error: Any error that occurred during update
	UpdateTrainingDocument(doc *entity.TrainingDocument) error

	// DeleteTrainingDocument removes a training document and its associated chunks.
	// Parameters:
	//   - docID: ID of the document to delete
	// Returns:
	//   - error: Any error that occurred during deletion
	DeleteTrainingDocument(docID string) error

	// StoreChunks performs batch insertion of document chunks with embeddings.
	// This is optimized for inserting multiple chunks at once.
	// Parameters:
	//   - chunks: Array of document chunks to store
	// Returns:
	//   - error: Any error that occurred during storage
	StoreChunks(chunks []entity.DocumentChunk) error

	// SearchSimilar performs vector similarity search using cosine distance.
	// Only searches through active documents and returns chunks above the threshold.
	// Parameters:
	//   - agentID: ID of the agent whose knowledge base to search
	//   - embedding: Query embedding vector (1024 dimensions)
	//   - topK: Maximum number of chunks to return
	//   - threshold: Minimum similarity score (0.0 to 1.0)
	// Returns:
	//   - []entity.RetrievedChunk: Array of similar chunks with scores
	//   - error: Any error that occurred during search
	SearchSimilar(agentID string, embedding []float32, topK int, threshold float32) ([]entity.RetrievedChunk, error)

	// DeleteChunksByDocumentID removes all chunks associated with a specific document.
	// Parameters:
	//   - docID: ID of the document whose chunks to delete
	// Returns:
	//   - error: Any error that occurred during deletion
	DeleteChunksByDocumentID(docID string) error

	// DeleteChunksByAgentID removes all chunks associated with a specific agent.
	// This is useful for retraining or cleaning up an agent's knowledge base.
	// Parameters:
	//   - agentID: ID of the agent whose chunks to delete
	// Returns:
	//   - error: Any error that occurred during deletion
	DeleteChunksByAgentID(agentID string) error
}

type TokenRepositoryInterface interface {
	CreateToken(email, purpose, token string, expiresAt time.Time) error
	GetToken(email, purpose string) (*entity.Token, error)
}
