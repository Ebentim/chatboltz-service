// Package usecase provides business logic for RAG operations.
// This package contains use cases for training agents and retrieving context.
package usecase

import (
	"fmt"

	"github.com/alpinesboltltd/boltz-ai/internal/config"
	"github.com/alpinesboltltd/boltz-ai/internal/entity"
	"github.com/alpinesboltltd/boltz-ai/internal/rag"
	"github.com/alpinesboltltd/boltz-ai/internal/repository"
	"gorm.io/gorm"
)

// TrainingUseCase handles agent training operations including document processing.
// It orchestrates the training workflow from document ingestion to knowledge base creation.
type TrainingUseCase struct {
	// ragService provides RAG operations for document processing and retrieval
	ragService *rag.RAGService
	// agentRepo provides access to agent data
	agentRepo repository.AgentRepositoryInterface
}

// NewTrainingUseCase creates a new training use case with the required dependencies.
// It initializes the RAG service with Cohere client and repository.
//
// Parameters:
//   - cfg: Application configuration containing API keys
//   - db: Database connection for repository operations
//   - agentRepo: Repository for agent operations
//
// Returns:
//   - *TrainingUseCase: Configured training use case
func NewTrainingUseCase(cfg *config.Config, db *gorm.DB, agentRepo repository.AgentRepositoryInterface) (*TrainingUseCase, error) {
	cohere, err := rag.NewCohereClient(cfg.COHERE_API_KEY)
	if err != nil {
		return nil, fmt.Errorf("failed to create Cohere client: %w", err)
	}

	ragRepo := repository.NewRAGRepository(db)
	ragService := rag.NewRAGService(cohere, ragRepo)

	return &TrainingUseCase{
		ragService: ragService,
		agentRepo:  agentRepo,
	}, nil
}

// ProcessDocument processes a single document for an agent's knowledge base.
// This is the primary method for adding new training content to an agent.
//
// Supported document types:
//   - text: Plain text content
//   - pdf: Extracted PDF text
//   - audio: Transcribed audio content
//   - video: Transcribed video content
//   - faq: FAQ content in Q&A format
//   - image: OCR text or image descriptions
//
// Parameters:
//   - agentID: ID of the agent to train
//   - title: Human-readable title for the document
//   - content: Raw text content to process
//   - docType: Type of document for appropriate processing
//   - sourceURL: Optional URL of the original document
//
// Returns:
//   - error: Any error that occurred during processing
func (t *TrainingUseCase) ProcessDocument(agentID, title, content string, docType entity.DocumentType, sourceURL *string) error {
	return t.ragService.ProcessDocument(agentID, title, docType, content, sourceURL)
}

// TrainAgentFromLegacyData migrates legacy training data to the new RAG system.
// This method processes existing TrainingData records and converts them to the new format.
//
// Parameters:
//   - agentID: ID of the agent whose legacy data to migrate
//
// Returns:
//   - error: Any error that occurred during migration
func (t *TrainingUseCase) TrainAgentFromLegacyData(agentID string) error {
	agent, err := t.agentRepo.GetAgent(agentID)
	if err != nil {
		return fmt.Errorf("failed to get agent: %w", err)
	}

	if len(agent.TrainingData) == 0 {
		return fmt.Errorf("no training data found for agent %s", agentID)
	}

	// Process legacy training data
	for _, td := range agent.TrainingData {
		if td.IsActive {
			for _, text := range td.Content {
				if err := t.ragService.ProcessDocument(agentID, text.Title, entity.DocumentTypeText, text.Content, nil); err != nil {
					return fmt.Errorf("failed to process legacy text: %w", err)
				}
			}
		}
	}

	return nil
}

// GetAgentDocuments retrieves all training documents for an agent.
// This includes document metadata and processing status.
//
// Parameters:
//   - agentID: ID of the agent whose documents to retrieve
//
// Returns:
//   - []entity.TrainingDocument: Array of training documents
//   - error: Any error that occurred during retrieval
func (t *TrainingUseCase) GetAgentDocuments(agentID string) ([]entity.TrainingDocument, error) {
	return t.ragService.GetAgentDocuments(agentID)
}

// DeleteAgentTraining removes all training data for an agent.
// This clears the agent's knowledge base completely.
//
// Parameters:
//   - agentID: ID of the agent whose training data to delete
//
// Returns:
//   - error: Any error that occurred during deletion
func (t *TrainingUseCase) DeleteAgentTraining(agentID string) error {
	return t.ragService.DeleteAgentDocuments(agentID)
}
