package rag

import (
	"github.com/alpinesboltltd/boltz-ai/internal/entity"
)

// VectorDB defines the interface for vector database operations.
type VectorDB interface {
	// Store upserts a document chunk with its embedding
	Store(chunk *entity.DocumentChunk) error
	// Search performs similarity search and returns matching chunks
	Search(agentID string, embedding []float32, topK int, threshold float32) ([]entity.RetrievedChunk, error)
	// Delete removes all vectors for an agent
	Delete(agentID string) error
}
