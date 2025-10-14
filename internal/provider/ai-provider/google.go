package aiprovider

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/alpinesboltltd/boltz-ai/internal/entity"
	"google.golang.org/genai"
)

type GoogleAIProvider struct {
	client *genai.Client
}

func NewGoogleAIClient(apiKey string) (*GoogleAIProvider, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{APIKey: apiKey})
	if err != nil {
		return nil, fmt.Errorf("error creating Google AI client: %w", err)
	}
	return &GoogleAIProvider{client: client}, nil
}

func (p *GoogleAIProvider) GetCapabilities() entity.ModelCapabilities {
	return entity.DefaultModelCapabilities.GetCapabilities("gemini-2.0-flash")
}

func (p *GoogleAIProvider) CompleteConversation(conversation Conversation, config map[string]any) (string, error) {
	ctx := context.Background()

	model := "gemini-1.5-flash"
	if m, ok := config["model"].(string); ok && m != "" {
		model = m
	}

	// Build generation config
	genConfig := &genai.GenerateContentConfig{}
	if temp, ok := config["temperature"].(float64); ok {
		tempFloat32 := float32(temp)
		genConfig.Temperature = &tempFloat32
	}
	if maxTokens, ok := config["max_tokens"].(int); ok {
		maxTokensInt32 := int32(maxTokens)
		genConfig.MaxOutputTokens = maxTokensInt32
	}

	// Convert messages to Google AI format
	messages := EnsureSystemMessage(conversation.Messages)
	var contents []*genai.Content
	for _, msg := range messages {
		switch msg.Role {
		case RoleSystem:
			genConfig.SystemInstruction = &genai.Content{
				Role:  string(RoleSystem),
				Parts: []*genai.Part{{Text: msg.Content}},
			}
		case RoleAssistant:
			contents = append(contents, &genai.Content{
				Role:  "model",
				Parts: []*genai.Part{{Text: msg.Content}},
			})
		case RoleUser:
			contents = append(contents, &genai.Content{
				Role:  string(RoleUser),
				Parts: []*genai.Part{{Text: msg.Content}},
			})
		}
	}

	result, err := p.client.Models.GenerateContent(ctx, model, contents, genConfig)
	if err != nil {
		return "", fmt.Errorf("error generating content: %w", err)
	}

	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no response generated")
	}

	return result.Text(), nil
}

func (p *GoogleAIProvider) CompleteConversationStream(conversation Conversation, config map[string]interface{}, callback StreamCallback) error {
	// Google AI doesn't support streaming in this SDK version
	// Fallback to regular completion with chunked delivery
	result, err := p.CompleteConversation(conversation, config)
	if err != nil {
		return err
	}

	// Simulate streaming by chunking response
	words := strings.Fields(result)
	for i, word := range words {
		if err := callback(word+" ", i == len(words)-1); err != nil {
			return err
		}
	}
	return nil
}

func (p *GoogleAIProvider) CompleteMultimodalConversation(messages []MultimodalMessage, config map[string]interface{}) (string, error) {
	ctx := context.Background()

	model := "gemini-1.5-flash"
	if m, ok := config["model"].(string); ok && m != "" {
		model = m
	}

	genConfig := &genai.GenerateContentConfig{}
	if temp, ok := config["temperature"].(float64); ok {
		tempFloat32 := float32(temp)
		genConfig.Temperature = &tempFloat32
	}
	if maxTokens, ok := config["max_tokens"].(int); ok {
		maxTokensInt32 := int32(maxTokens)
		genConfig.MaxOutputTokens = maxTokensInt32
	}

	var contents []*genai.Content
	for _, msg := range messages {
		switch msg.Role {
		case RoleSystem:
			genConfig.SystemInstruction = &genai.Content{
				Role:  string(RoleSystem),
				Parts: []*genai.Part{{Text: msg.Content}},
			}
		case RoleAssistant:
			contents = append(contents, &genai.Content{
				Role:  "model",
				Parts: []*genai.Part{{Text: msg.Content}},
			})
		case RoleUser:
			if msg.HasImageContent() {
				if msg.IsBase64Image() {
					// Google supports base64 images via InlineData
					imageData, _ := base64.StdEncoding.DecodeString(msg.MediaBase64)
					mimeType := DetectImageMimeType(msg.MediaBase64)
					contents = append(contents, &genai.Content{
						Role: string(RoleUser),
						Parts: []*genai.Part{
							{Text: msg.Content},
							{InlineData: &genai.Blob{MIMEType: mimeType, Data: imageData}},
						},
					})
				} else {
					// Google supports image URLs via FileData
					contents = append(contents, &genai.Content{
						Role: string(RoleUser),
						Parts: []*genai.Part{
							{Text: msg.Content},
							{FileData: &genai.FileData{MIMEType: "image/jpeg", FileURI: msg.MediaURL}},
						},
					})
				}
			} else {
				contents = append(contents, &genai.Content{
					Role:  string(RoleUser),
					Parts: []*genai.Part{{Text: msg.Content}},
				})
			}
		}
	}

	result, err := p.client.Models.GenerateContent(ctx, model, contents, genConfig)
	if err != nil {
		return "", fmt.Errorf("error generating content: %w", err)
	}

	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no response generated")
	}

	return result.Text(), nil
}
