package rag

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/alpinesboltltd/boltz-ai/internal/entity"
	"github.com/pinecone-io/go-pinecone/v4/pinecone"
	"google.golang.org/protobuf/types/known/structpb"
)

// PineconeDB implements VectorDB interface using Pinecone's official Go SDK v4.
type PineconeDB struct {
	client *pinecone.Client
	index  *pinecone.IndexConnection
}

// NewPineconeDB creates a new Pinecone vector database client.
func NewPineconeDB(apiKey, indexName string) (*PineconeDB, error) {
	client, err := pinecone.NewClient(pinecone.NewClientParams{ApiKey: apiKey})
	if err != nil {
		return nil, fmt.Errorf("failed to create Pinecone client: %w", err)
	}

	index, err := client.Index(pinecone.NewIndexConnParams{Host: indexName})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to index: %w", err)
	}

	return &PineconeDB{client: client, index: index}, nil
}

// Store upserts a document chunk into Pinecone.
func (p *PineconeDB) Store(chunk *entity.DocumentChunk) error {
	metadataJSON, _ := json.Marshal(chunk.Metadata)

	metadata, err := structpb.NewStruct(map[string]interface{}{
		"agent_id":    chunk.AgentID,
		"document_id": chunk.DocumentID,
		"content":     chunk.Content,
		"metadata":    string(metadataJSON),
	})
	if err != nil {
		return fmt.Errorf("failed to create metadata: %w", err)
	}

	vector := &pinecone.Vector{
		Id:       chunk.ID,
		Values:   &chunk.Embedding,
		Metadata: metadata,
	}

	_, err = p.index.UpsertVectors(context.Background(), []*pinecone.Vector{vector})
	return err
}

// Search performs similarity search in Pinecone.
func (p *PineconeDB) Search(agentID string, embedding []float32, topK int, threshold float32) ([]entity.RetrievedChunk, error) {
	filter, err := structpb.NewStruct(map[string]interface{}{"agent_id": agentID})
	if err != nil {
		return nil, fmt.Errorf("failed to create filter: %w", err)
	}

	resp, err := p.index.QueryByVectorValues(context.Background(), &pinecone.QueryByVectorValuesRequest{
		Vector:          embedding,
		TopK:            uint32(topK),
		IncludeMetadata: true,
		MetadataFilter:  filter,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query Pinecone: %w", err)
	}

	var results []entity.RetrievedChunk
	for _, match := range resp.Matches {
		if match.Score >= threshold {
			var metadata map[string]string
			if metadataStr := match.Vector.Metadata.Fields["metadata"].GetStringValue(); metadataStr != "" {
				json.Unmarshal([]byte(metadataStr), &metadata)
			}

			results = append(results, entity.RetrievedChunk{
				Content:    match.Vector.Metadata.Fields["content"].GetStringValue(),
				Metadata:   metadata,
				Score:      match.Score,
				DocumentID: match.Vector.Metadata.Fields["document_id"].GetStringValue(),
			})
		}
	}

	return results, nil
}

// Delete removes vectors by agent ID filter.
func (p *PineconeDB) Delete(agentID string) error {
	filter, err := structpb.NewStruct(map[string]interface{}{"agent_id": agentID})
	if err != nil {
		return fmt.Errorf("failed to create filter: %w", err)
	}
	err = p.index.DeleteVectorsByFilter(context.Background(), filter)
	return err
}
