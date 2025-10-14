/*
Package aiprovider implements Groq AI provider integration.

Groq provides extremely fast inference for various AI models including:
- Chat completion models (Llama, Mixtral, Gemma)
- Audio transcription models (Whisper)

Supported models and their capabilities:
- llama-3.3-70b-versatile: Text generation
- llama-3.1-8b-instant: Fast text generation
- llama-3.1-70b-versatile: High-quality text generation
- gemma2-9b-it: Instruction-tuned text generation
- mixtral-8x7b-32768: Mixture of experts text generation
- whisper-large-v3: Audio transcription
- whisper-large-v3-turbo: Fast audio transcription
- distil-whisper-large-v3-en: English-optimized audio transcription

Usage:

	provider, err := NewGroqAIClient(apiKey, "https://api.groq.com/openai/v1")
	if err != nil {
		log.Fatal(err)
	}

	conversation := Conversation{
		Messages: []Message{
			{Role: RoleUser, Content: "Hello!"},
		},
	}

	config := map[string]interface{}{
		"model": "llama-3.3-70b-versatile",
	}

	response, err := provider.CompleteConversation(conversation, config)
*/
package aiprovider

import (
	"context"
	"fmt"

	"github.com/alpinesboltltd/boltz-ai/internal/entity"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type GroqAIProvider struct {
	client *openai.Client
}

func NewGroqAIClient(apiKey, groqUrl string) (*GroqAIProvider, error) {
	client := openai.NewClient(option.WithAPIKey(apiKey), option.WithBaseURL(groqUrl))
	return &GroqAIProvider{client: &client}, nil
}

func (p *GroqAIProvider) GetCapabilities() entity.ModelCapabilities {
	// Return general capabilities - specific model capabilities are handled by DefaultModelCapabilities
	return entity.ModelCapabilities{Text: true, Voice: true, Vision: false}
}

func (p *GroqAIProvider) CompleteConversation(conversation Conversation, config map[string]interface{}) (string, error) {
	messages := EnsureSystemMessage(conversation.Messages)

	model := "llama-3.3-70b-versatile"
	if m, ok := config["model"].(string); ok && m != "" {
		model = m
	}

	// Configure additional parameters specific to Groq
	params := openai.ChatCompletionNewParams{
		Messages: ToOpenAIMessages(messages),
		Model:    model,
	}

	// Set optional parameters
	if temperature, ok := config["temperature"].(float64); ok {
		params.Temperature = openai.Float(temperature)
	}
	if maxTokens, ok := config["max_tokens"].(int); ok {
		params.MaxTokens = openai.Int(int64(maxTokens))
	}
	if topP, ok := config["top_p"].(float64); ok {
		params.TopP = openai.Float(topP)
	}
	// Note: Stop sequences can be configured but we'll keep it simple for now
	// The openai-go library handles the stop parameter differently

	chatCompletion, err := p.client.Chat.Completions.New(context.Background(), params)
	if err != nil {
		return "", fmt.Errorf("Groq API error: %w", err)
	}

	if len(chatCompletion.Choices) == 0 {
		return "", fmt.Errorf("no response choices returned from Groq")
	}

	return chatCompletion.Choices[0].Message.Content, nil
}

func (p *GroqAIProvider) CompleteConversationStream(conversation Conversation, config map[string]interface{}, callback StreamCallback) error {
	messages := EnsureSystemMessage(conversation.Messages)
	model := "llama-3.3-70b-versatile"
	if m, ok := config["model"].(string); ok && m != "" {
		model = m
	}

	// Configure parameters for streaming
	params := openai.ChatCompletionNewParams{
		Messages: ToOpenAIMessages(messages),
		Model:    model,
	}

	// Set optional parameters
	if temperature, ok := config["temperature"].(float64); ok {
		params.Temperature = openai.Float(temperature)
	}
	if maxTokens, ok := config["max_tokens"].(int); ok {
		params.MaxTokens = openai.Int(int64(maxTokens))
	}

	stream := p.client.Chat.Completions.NewStreaming(context.Background(), params)

	for stream.Next() {
		chunk := stream.Current()
		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
			if err := callback(chunk.Choices[0].Delta.Content, false); err != nil {
				return fmt.Errorf("callback error: %w", err)
			}
		}
	}

	if err := stream.Err(); err != nil {
		return fmt.Errorf("Groq streaming error: %w", err)
	}

	return callback("", true)
}

func (p *GroqAIProvider) CompleteMultimodalConversation(messages []MultimodalMessage, config map[string]interface{}) (string, error) {
	// Groq primarily supports text and audio, but not vision
	// For images, we'll convert to text descriptions
	conv := Conversation{}
	for _, msg := range messages {
		content := msg.Content
		if msg.HasImageContent() {
			if msg.IsBase64Image() {
				content += " [User provided a base64 encoded image]"
			} else {
				content += " [User provided an image from URL: " + msg.MediaURL + "]"
			}
		}
		conv.Messages = append(conv.Messages, Message{
			Role:    Role(msg.Role),
			Content: content,
		})
	}
	return p.CompleteConversation(conv, config)
}

// Audio-specific methods for Groq's Whisper support
func (p *GroqAIProvider) TranscribeAudio(audioData []byte, config map[string]interface{}) (string, error) {
	// Note: This is a placeholder implementation
	// The actual Groq audio API implementation would depend on their specific audio endpoints
	// For now, we'll return a placeholder response
	return "[Audio transcription not implemented - use Groq's audio transcription API]", fmt.Errorf("audio transcription not yet implemented for Groq")
}

func (p *GroqAIProvider) GenerateSpeech(text string, config map[string]interface{}) ([]byte, error) {
	// Note: Groq doesn't currently support text-to-speech generation
	// This would need to be handled by a TTS fallback service
	return nil, fmt.Errorf("text-to-speech not supported by Groq")
}
