package rag

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/alpinesboltltd/boltz-ai/internal/entity"
	"github.com/pinecone-io/go-pinecone/v4/pinecone"
)

// PineconeDB implements VectorDB interface using Pinecone's official Go SDK v4.
type PineconeDB struct {
	client *pinecone.Client
	index  *pinecone.Index
}

// NewPineconeDB creates a new Pinecone vector database client.
func NewPineconeDB(apiKey, indexName string) (*PineconeDB, error) {
	client, err := pinecone.NewClient(pinecone.NewClientParams{ApiKey: apiKey})
	if err != nil {
		return nil, fmt.Errorf("failed to create Pinecone client: %w", err)
	}

	index := client.Index(indexName)
	return &PineconeDB{client: client, index: index}, nil
}

// Store upserts a document chunk into Pinecone.
func (p *PineconeDB) Store(chunk *entity.DocumentChunk) error {
	metadataJSON, _ := json.Marshal(chunk.Metadata)

	vector := &pinecone.Vector{
		Id:     chunk.ID,
		Values: chunk.Embedding,
		Metadata: map[string]interface{}{
			"agent_id":    chunk.AgentID,
			"document_id": chunk.DocumentID,
			"content":     chunk.Content,
			"metadata":    string(metadataJSON),
		},
	}

	_, err := p.index.UpsertVectors(context.Background(), []*pinecone.Vector{vector})
	return err
}

// Search performs similarity search in Pinecone.
func (p *PineconeDB) Search(agentID string, embedding []float32, topK int, threshold float32) ([]entity.RetrievedChunk, error) {
	filter := map[string]interface{}{"agent_id": agentID}

	resp, err := p.index.QueryByVectorValues(context.Background(), &pinecone.QueryByVectorValuesRequest{
		Vector:          embedding,
		TopK:            uint32(topK),
		IncludeMetadata: true,
		Filter:          filter,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query Pinecone: %w", err)
	}

	var results []entity.RetrievedChunk
	for _, match := range resp.Matches {
		if match.Score >= threshold {
			var metadata map[string]string
			if metadataStr, ok := match.Metadata["metadata"].(string); ok {
				json.Unmarshal([]byte(metadataStr), &metadata)
			}

			results = append(results, entity.RetrievedChunk{
				Content:    match.Metadata["content"].(string),
				Metadata:   metadata,
				Score:      match.Score,
				DocumentID: match.Metadata["document_id"].(string),
			})
		}
	}

	return results, nil
}

// Delete removes vectors by agent ID filter.
func (p *PineconeDB) Delete(agentID string) error {
	filter := map[string]interface{}{"agent_id": agentID}
	_, err := p.index.DeleteVectorsByFilter(context.Background(), filter)
	return err
}
