package usecase

import (
	"fmt"
	"strings"

	"github.com/alpinesboltltd/boltz-ai/internal/config"
	"github.com/alpinesboltltd/boltz-ai/internal/entity"
	"github.com/alpinesboltltd/boltz-ai/internal/rag"
	"github.com/alpinesboltltd/boltz-ai/internal/repository"
	"gorm.io/gorm"
)

// RAGRetrieverUseCase handles context retrieval for LLM queries.
// It provides semantic search capabilities across an agent's knowledge base.
type RAGRetrieverUseCase struct {
	// ragService provides RAG operations for context retrieval
	ragService *rag.RAGService
}

// NewRAGRetrieverUseCase creates a new RAG retriever use case.
// It initializes the RAG service with Cohere client and repository.
//
// Parameters:
//   - cfg: Application configuration containing API keys
//   - db: Database connection for repository operations
//
// Returns:
//   - *RAGRetrieverUseCase: Configured RAG retriever use case
func NewRAGRetrieverUseCase(cfg *config.Config, db *gorm.DB) (*RAGRetrieverUseCase, error) {
	// Trim whitespace from API keys (handles trailing newlines from .env)
	cohereKey := strings.TrimSpace(cfg.COHERE_API_KEY)
	openaiKey := strings.TrimSpace(cfg.OPENAI_API_KEY)
	pineconeKey := strings.TrimSpace(cfg.PINECONE_API_KEY)

	cohere, err := rag.NewCohereClient(cohereKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create Cohere client: %w", err)
	}

	mediaProcessor := rag.NewOpenAIMediaProcessor(openaiKey)
	ragRepo := repository.NewRAGRepository(db)

	var vectorDB rag.VectorDB
	if cfg.VECTOR_DB_TYPE == "pinecone" && pineconeKey != "" {
		vectorDB, err = rag.NewPineconeDB(pineconeKey, cfg.PINECONE_INDEX_NAME)
		if err != nil {
			return nil, fmt.Errorf("failed to create Pinecone client: %w", err)
		}
	}

	ragService := rag.NewRAGService(cohere, ragRepo, mediaProcessor, vectorDB, cfg.VECTOR_DB_TYPE)

	return &RAGRetrieverUseCase{
		ragService: ragService,
	}, nil
}

// RetrieveContext performs semantic search to find relevant context for a user query.
// It uses default parameters (TopK=5, Threshold=0.7) for optimal results.
//
// The retrieval process:
// 1. Generate embedding for user query using Cohere
// 2. Search for similar chunks in agent's knowledge base
// 3. Return combined context and individual chunks
//
// Parameters:
//   - agentID: ID of the agent whose knowledge base to search
//   - userQuery: The user's question or search text
//
// Returns:
//   - *entity.RAGResponse: Response with context and retrieved chunks
//   - error: Any error that occurred during retrieval
func (r *RAGRetrieverUseCase) RetrieveContext(agentID, userQuery string) (*entity.RAGResponse, error) {
	query := entity.RAGQuery{
		Query:     userQuery,
		AgentID:   agentID,
		TopK:      5,
		Threshold: 0.7,
	}

	response, err := r.ragService.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve context: %w", err)
	}

	return response, nil
}

// BuildLLMPrompt constructs a prompt for the LLM that includes system instructions,
// relevant context from RAG, and the user query. This creates a grounded prompt
// that helps prevent hallucination.
//
// Prompt structure:
// - System instructions (agent behavior, role, etc.)
// - Relevant context (if available from RAG)
// - User query
//
// Parameters:
//   - userQuery: The user's original question
//   - systemInstruction: System-level instructions for the agent
//   - ragResponse: Retrieved context from RAG (can be nil)
//
// Returns:
//   - string: Complete prompt ready for LLM processing
func (r *RAGRetrieverUseCase) BuildLLMPrompt(userQuery, systemInstruction string, ragResponse *entity.RAGResponse) string {
	if ragResponse == nil || ragResponse.Context == "" {
		return fmt.Sprintf("%s\n\nUser Query: %s", systemInstruction, userQuery)
	}

	return fmt.Sprintf("%s\n\nRelevant Context:\n%s\n\nUser Query: %s",
		systemInstruction, ragResponse.Context, userQuery)
}
