package aiprovider

import (
	"context"
	"fmt"
	"strings"

	"github.com/alpinesboltltd/boltz-ai/internal/entity"
	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/anthropics/anthropic-sdk-go/packages/param"
	"github.com/anthropics/anthropic-sdk-go/shared/constant"
)

type AnthropicProvider struct {
	client *anthropic.Client
}

func NewAnthropicClient(apiKey string) *AnthropicProvider {
	client := anthropic.NewClient(option.WithAPIKey(apiKey))
	return &AnthropicProvider{client: &client}
}

func (p *AnthropicProvider) GetCapabilities() entity.ModelCapabilities {
	return entity.ModelCapabilities{Text: true, Voice: false, Vision: true}
}

func (p *AnthropicProvider) CompleteConversation(conversation Conversation, config map[string]interface{}) (string, error) {
	messages := EnsureSystemMessage(conversation.Messages)

	var systemMsg string
	var chatMessages []anthropic.MessageParam

	for _, msg := range messages {
		if msg.Role == RoleSystem {
			systemMsg = msg.Content
		} else {
			switch msg.Role {
			case RoleUser:
				chatMessages = append(chatMessages, anthropic.NewUserMessage(anthropic.NewTextBlock(msg.Content)))
			case RoleAssistant:
				chatMessages = append(chatMessages, anthropic.NewAssistantMessage(anthropic.NewTextBlock(msg.Content)))
			default:
				continue
			}
		}
	}

	genConfig := anthropic.MessageNewParams{
		Model:     "claude-3-5-haiku-20241022",
		Messages:  chatMessages,
		MaxTokens: 1000,
	}

	if systemMsg != "" {
		genConfig.System = []anthropic.TextBlockParam{
			{Text: systemMsg,
				Type: constant.Text("text")},
		}
	}

	if model, ok := config["model"].(string); ok && model != "" {
		genConfig.Model = anthropic.Model(model)
	}

	if temp, ok := config["temperature"].(float64); ok {
		genConfig.Temperature = param.NewOpt(temp)
	}

	if maxTokens, ok := config["max_tokens"].(int); ok {
		genConfig.MaxTokens = int64(maxTokens)
	}

	message, err := p.client.Messages.New(context.Background(), genConfig)
	if err != nil {
		return "", err
	}

	if len(message.Content) == 0 {
		return "", fmt.Errorf("no response content")
	}

	return message.Content[0].Text, nil
}

func (p *AnthropicProvider) CompleteConversationStream(conversation Conversation, config map[string]interface{}, callback StreamCallback) error {
	// Anthropic doesn't support streaming in this SDK version
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

func (p *AnthropicProvider) CompleteMultimodalConversation(messages []MultimodalMessage, config map[string]interface{}) (string, error) {
	var systemMsg string
	var chatMessages []anthropic.MessageParam

	for _, msg := range messages {
		if msg.Role == RoleSystem {
			systemMsg = msg.Content
		} else {
			switch msg.Role {
			case RoleUser:
				if msg.HasImageContent() {
					if msg.IsBase64Image() {
						// Anthropic supports base64 images
						mimeType := DetectImageMimeType(msg.MediaBase64)
						chatMessages = append(chatMessages, anthropic.NewUserMessage(
							anthropic.NewTextBlock(msg.Content),
							anthropic.NewImageBlock(anthropic.Base64ImageSourceParam{
								Data:      msg.MediaBase64,
								MediaType: anthropic.Base64ImageSourceMediaType(mimeType),
							}),
						))
					} else {
						// Anthropic supports image URLs
						chatMessages = append(chatMessages, anthropic.NewUserMessage(
							anthropic.NewTextBlock(msg.Content),
							anthropic.NewImageBlock(anthropic.URLImageSourceParam{URL: msg.MediaURL}),
						))
					}
				} else {
					chatMessages = append(chatMessages, anthropic.NewUserMessage(anthropic.NewTextBlock(msg.Content)))
				}
			case RoleAssistant:
				chatMessages = append(chatMessages, anthropic.NewAssistantMessage(anthropic.NewTextBlock(msg.Content)))
			default:
				continue
			}
		}
	}

	genConfig := anthropic.MessageNewParams{
		Model:     "claude-3-5-haiku-20241022",
		Messages:  chatMessages,
		MaxTokens: 1000,
	}

	if systemMsg != "" {
		genConfig.System = []anthropic.TextBlockParam{
			{Text: systemMsg, Type: constant.Text("text")},
		}
	}

	if model, ok := config["model"].(string); ok && model != "" {
		genConfig.Model = anthropic.Model(model)
	}

	if temp, ok := config["temperature"].(float64); ok {
		genConfig.Temperature = param.NewOpt(temp)
	}

	if maxTokens, ok := config["max_tokens"].(int); ok {
		genConfig.MaxTokens = int64(maxTokens)
	}

	message, err := p.client.Messages.New(context.Background(), genConfig)
	if err != nil {
		return "", err
	}

	if len(message.Content) == 0 {
		return "", fmt.Errorf("no response content")
	}

	return message.Content[0].Text, nil
}
