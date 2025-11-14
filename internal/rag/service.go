package rag

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/alpinesboltltd/boltz-ai/internal/entity"
	"github.com/alpinesboltltd/boltz-ai/internal/repository"
	"github.com/google/uuid"
)

// RAGService provides high-level RAG operations for training agents and retrieving context.
// It orchestrates document processing, embedding generation, and similarity search.
type RAGService struct {
	// processor handles document chunking and embedding generation
	processor *ContentProcessor
	// repo provides data access for training documents and chunks
	repo repository.RAGRepositoryInterface
	// vectorDB handles vector storage (Pinecone or pgvector)
	vectorDB VectorDB
	// vectorDBType determines which vector DB is used
	vectorDBType string
}

// NewRAGService creates a new RAG service with the provided dependencies.
//
// Parameters:
//   - cohere: Configured Cohere client for embeddings
//   - repo: Repository interface for data persistence
//   - vectorDB: Vector database implementation (optional, uses repo if nil)
//   - vectorDBType: Type of vector DB ("pgvector" or "pinecone")
//
// Returns:
//   - *RAGService: Configured RAG service ready for use
func NewRAGService(cohere *CohereClient, repo repository.RAGRepositoryInterface, mediaProcessor MediaProcessor, vectorDB VectorDB, vectorDBType string) *RAGService {
	return &RAGService{
		processor:    NewContentProcessor(cohere, mediaProcessor),
		repo:         repo,
		vectorDB:     vectorDB,
		vectorDBType: vectorDBType,
	}
}

// ProcessDocument processes a document for an agent's knowledge base.
// This method creates a training document record, processes the content into chunks,
// generates embeddings, and stores everything in the database.
//
// The processing flow:
// 1. Create TrainingDocument record
// 2. Process content into chunks based on document type
// 3. Generate embeddings for all chunks
// 4. Store chunks with embeddings in database
// 5. Mark document as processed
//
// Parameters:
//   - agentID: ID of the agent to train
//   - title: Human-readable title for the document
//   - docType: Type of document (text, pdf, audio, video, faq, image)
//   - content: Raw text content to process
//   - sourceURL: Optional URL of the original document
//
// Returns:
//   - error: Any error that occurred during processing
func (r *RAGService) ProcessDocument(agentID, title string, docType entity.DocumentType, content string, sourceURL *string) error {
	// Create training document
	doc := &entity.TrainingDocument{
		ID:           uuid.New().String(),
		AgentID:      agentID,
		Title:        title,
		DocumentType: docType,
		SourceURL:    sourceURL,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := r.repo.CreateTrainingDocument(doc); err != nil {
		return fmt.Errorf("failed to create training document: %w", err)
	}

	// Process content into chunks
	chunks, err := r.processor.ProcessDocument(doc, content)
	if err != nil {
		return fmt.Errorf("failed to process document: %w", err)
	}

	// Add IDs and timestamps to chunks
	for i := range chunks {
		chunks[i].ID = uuid.New().String()
		chunks[i].CreatedAt = time.Now()
		chunks[i].UpdatedAt = time.Now()
	}

	// Store chunks based on vector DB type
	if r.vectorDBType == "pinecone" && r.vectorDB != nil {
		// Store metadata in PostgreSQL
		if err := r.repo.StoreChunksMetadataOnly(chunks); err != nil {
			return fmt.Errorf("failed to store chunks metadata: %w", err)
		}
		// Store vectors in Pinecone
		for _, chunk := range chunks {
			if err := r.vectorDB.Store(&chunk); err != nil {
				return fmt.Errorf("failed to store chunk in Pinecone: %w", err)
			}
		}
	} else {
		// Store everything in PostgreSQL (pgvector)
		if err := r.repo.StoreChunks(chunks); err != nil {
			return fmt.Errorf("failed to store chunks: %w", err)
		}
	}

	// Update document as processed
	now := time.Now()
	doc.ProcessedAt = &now
	return r.repo.UpdateTrainingDocument(doc)
}

// Query performs semantic search across an agent's knowledge base.
// It generates an embedding for the query and finds the most similar chunks.
//
// The query flow:
// 1. Apply default values for TopK (5) and Threshold (0.7)
// 2. Generate embedding for the user query
// 3. Perform similarity search in vector database
// 4. Combine retrieved chunks into context string
//
// Parameters:
//   - query: RAG query with user question and search parameters
//
// Returns:
//   - *entity.RAGResponse: Response with context and individual chunks
//   - error: Any error that occurred during the query
func (r *RAGService) Query(query entity.RAGQuery) (*entity.RAGResponse, error) {
	if query.TopK == 0 {
		query.TopK = 5
	}
	if query.Threshold == 0 {
		query.Threshold = 0.7
	}

	// Generate embedding for query
	embeddings, err := r.processor.cohere.Embed([]string{query.Query}, "search_query")
	if err != nil {
		return nil, fmt.Errorf("failed to generate query embedding: %w", err)
	}

	// Search for similar chunks based on vector DB type
	var chunks []entity.RetrievedChunk
	if r.vectorDBType == "pinecone" && r.vectorDB != nil {
		// Search in Pinecone
		pineconeResults, err := r.vectorDB.Search(query.AgentID, embeddings[0], query.TopK, query.Threshold)
		if err != nil {
			return nil, fmt.Errorf("failed to search Pinecone: %w", err)
		}
		// Get full metadata from PostgreSQL
		var chunkIDs []string
		for _, result := range pineconeResults {
			chunkIDs = append(chunkIDs, result.DocumentID)
		}
		if len(chunkIDs) > 0 {
			dbChunks, err := r.repo.GetChunksByIDs(chunkIDs)
			if err != nil {
				return nil, fmt.Errorf("failed to get chunks metadata: %w", err)
			}
			// Merge Pinecone scores with PostgreSQL metadata
			for _, pineconeResult := range pineconeResults {
				for _, dbChunk := range dbChunks {
					if dbChunk.ID == pineconeResult.DocumentID {
						chunks = append(chunks, entity.RetrievedChunk{
							Content:      dbChunk.Content,
							Metadata:     dbChunk.Metadata,
							Score:        pineconeResult.Score,
							DocumentID:   dbChunk.DocumentID,
							DocumentType: entity.DocumentType(dbChunk.Metadata["type"]),
						})
						break
					}
				}
			}
		}
	} else {
		// Search in PostgreSQL (pgvector)
		chunks, err = r.repo.SearchSimilar(query.AgentID, embeddings[0], query.TopK, query.Threshold)
		if err != nil {
			return nil, fmt.Errorf("failed to search similar chunks: %w", err)
		}
	}

	// Build context from retrieved chunks
	var contextParts []string
	for _, chunk := range chunks {
		contextParts = append(contextParts, chunk.Content)
	}

	return &entity.RAGResponse{
		Context: strings.Join(contextParts, "\n\n"),
		Chunks:  chunks,
		Query:   query.Query,
	}, nil
}

// DeleteAgentDocuments removes all training documents and chunks for an agent.
// This is useful when retraining an agent or cleaning up data.
//
// Parameters:
//   - agentID: ID of the agent whose documents should be deleted
//
// Returns:
//   - error: Any error that occurred during deletion
func (r *RAGService) DeleteAgentDocuments(agentID string) error {
	return r.repo.DeleteChunksByAgentID(agentID)
}

// ProcessMediaFile processes media files (images, audio, video) using MediaProcessor
func (r *RAGService) ProcessMediaFile(agentID, title string, docType entity.DocumentType, fileData []byte, mimeType string, sourceURL *string) error {
	reader := bytes.NewReader(fileData)
	content, err := r.processor.ProcessMediaToText(reader, docType, mimeType)
	if err != nil {
		return fmt.Errorf("failed to process media file: %w", err)
	}

	return r.ProcessDocument(agentID, title, docType, content, sourceURL)
}

// GetAgentDocuments retrieves all training documents for an agent.
// This includes document metadata and associated chunks.
//
// Parameters:
//   - agentID: ID of the agent whose documents to retrieve
//
// Returns:
//   - []entity.TrainingDocument: Array of training documents with chunks
//   - error: Any error that occurred during retrieval
func (r *RAGService) GetAgentDocuments(agentID string) ([]entity.TrainingDocument, error) {
	return r.repo.GetTrainingDocumentsByAgentID(agentID)
}
