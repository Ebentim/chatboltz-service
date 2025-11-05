// Package rag provides the core RAG (Retrieval-Augmented Generation) functionality.
// This package includes the Cohere client for generating embeddings using the official Go SDK v2.
package rag

import (
	"context"
	"fmt"

	cohere "github.com/cohere-ai/cohere-go/v2"
	client "github.com/cohere-ai/cohere-go/v2/client"
)

// CohereClient provides access to Cohere's embedding API using the official Go SDK v2.
// It acts as a semantic gateway, converting text into vector embeddings for similarity search.
// Supports multilingual content including text, audio transcripts, video transcripts, and image descriptions.
type CohereClient struct {
	// client is the official Cohere Go SDK client
	client *client.Client
}

// NewCohereClient creates a new Cohere client with the provided API key.
// The client uses embed-multilingual-v3.0 model which produces 1024-dimensional embeddings
// and supports 100+ languages for global applications.
//
// Parameters:
//   - apiKey: Your Cohere API key from https://dashboard.cohere.ai/
//
// Returns:
//   - *CohereClient: Configured client ready for embedding requests
//   - error: Any error that occurred during client initialization
func NewCohereClient(apiKey string) (*CohereClient, error) {
	c := client.NewClient(client.WithToken(apiKey))
	return &CohereClient{client: c}, nil
}

// Embed generates vector embeddings for the provided texts using Cohere's official Go SDK v2.
// This method uses the embed-multilingual-v3.0 model for language-agnostic 1024-dimensional vectors.
// Supports content from images (OCR/descriptions), audio (transcripts), video (transcripts), and text.
//
// The multilingual model supports:
//   - 100+ languages including English, Spanish, French, German, Chinese, Japanese, Arabic, etc.
//   - Mixed-language content within the same embedding space
//   - Consistent vector dimensions across all languages
//
// Parameters:
//   - texts: Array of text strings to embed (max 96 per request)
//   - inputType: Type of input - "search_document" for training content, "search_query" for user queries
//
// Returns:
//   - [][]float32: Array of embeddings, each with 1024 dimensions
//   - error: Any error that occurred during the API request
//
// Example:
//
//	embeddings, err := client.Embed(["Hello world", "Hola mundo", "こんにちは世界"], "search_document")
func (c *CohereClient) Embed(texts []string, inputType string) ([][]float32, error) {
	// Convert input type to SDK enum
	var embedInputType cohere.EmbedInputType
	switch inputType {
	case "search_document":
		embedInputType = cohere.EmbedInputTypeSearchDocument
	case "search_query":
		embedInputType = cohere.EmbedInputTypeSearchQuery
	default:
		embedInputType = cohere.EmbedInputTypeSearchDocument
	}

	req := &cohere.EmbedRequest{
		Texts:     texts,
		Model:     cohere.String("embed-multilingual-v3.0"), // Language-agnostic model
		InputType: &embedInputType,
	}

	resp, err := c.client.Embed(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("failed to generate embeddings: %w", err)
	}

	// Convert response embeddings to float32
	var embeddings [][]float32
	for _, embedding := range resp.EmbeddingsFloats.Embeddings {
		var floatEmbedding []float32
		for _, val := range embedding {
			floatEmbedding = append(floatEmbedding, float32(val))
		}
		embeddings = append(embeddings, floatEmbedding)
	}

	return embeddings, nil
}
