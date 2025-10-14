package aiprovider

import (
	"context"
	"fmt"

	"github.com/alpinesboltltd/boltz-ai/internal/entity"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type OpenAIProvider struct {
	client *openai.Client
}

func NewOpenAIClient(apiKey string) *OpenAIProvider {
	client := openai.NewClient(option.WithAPIKey(apiKey))
	return &OpenAIProvider{client: &client}
}

// FIXME: Get accurate model capabilities
func (p *OpenAIProvider) GetCapabilities() entity.ModelCapabilities {
	// TODO: fetch model capabilities from the database
	return entity.ModelCapabilities{Text: true, Voice: true, Vision: true}
}

func (p *OpenAIProvider) CompleteConversation(conversation Conversation, config map[string]interface{}) (string, error) {
	messages := EnsureSystemMessage(conversation.Messages)

	model := "gpt-3.5-turbo"
	if m, ok := config["model"].(string); ok && m != "" {
		model = m
	}

	chatCompletion, err := p.client.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{
		Messages: ToOpenAIMessages(messages),
		Model:    model,
	})

	if err != nil {
		return "", err
	}

	if len(chatCompletion.Choices) == 0 {
		return "", fmt.Errorf("no response choices returned")
	}

	return chatCompletion.Choices[0].Message.Content, nil
}

func (p *OpenAIProvider) CompleteConversationStream(conversation Conversation, config map[string]interface{}, callback StreamCallback) error {
	messages := EnsureSystemMessage(conversation.Messages)
	model := "gpt-3.5-turbo"
	if m, ok := config["model"].(string); ok && m != "" {
		model = m
	}

	stream := p.client.Chat.Completions.NewStreaming(context.Background(), openai.ChatCompletionNewParams{
		Messages: ToOpenAIMessages(messages),
		Model:    model,
		// Stream:   openai.Bool(true),
	})

	for stream.Next() {
		chunk := stream.Current()
		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
			if err := callback(chunk.Choices[0].Delta.Content, false); err != nil {
				return err
			}
		}
	}

	if err := stream.Err(); err != nil {
		return err
	}

	return callback("", true)
}

func (p *OpenAIProvider) CompleteMultimodalConversation(messages []MultimodalMessage, config map[string]interface{}) (string, error) {
	// Convert multimodal messages to OpenAI format
	var oaiMessages []openai.ChatCompletionMessageParamUnion
	for _, msg := range messages {
		switch msg.Role {
		case RoleSystem:
			oaiMessages = append(oaiMessages, openai.SystemMessage(msg.Content))
		case RoleUser:
			if msg.HasImageContent() {
				if msg.IsBase64Image() {
					// OpenAI supports base64 images with data URI format
					dataURI := "data:image/jpeg;base64," + msg.MediaBase64
					oaiMessages = append(oaiMessages, openai.UserMessage(msg.Content+" [Image: "+dataURI[:50]+"...]"))
				} else {
					// OpenAI supports image URLs
					oaiMessages = append(oaiMessages, openai.UserMessage(msg.Content+" [Image URL: "+msg.MediaURL+"]"))
				}
			} else {
				oaiMessages = append(oaiMessages, openai.UserMessage(msg.Content))
			}
		case RoleAssistant:
			oaiMessages = append(oaiMessages, openai.AssistantMessage(msg.Content))
		}
	}

	model := "gpt-4o"
	if m, ok := config["model"].(string); ok && m != "" {
		model = m
	}

	chatCompletion, err := p.client.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{
		Messages: oaiMessages,
		Model:    model,
	})

	if err != nil {
		return "", err
	}

	if len(chatCompletion.Choices) == 0 {
		return "", fmt.Errorf("no response choices returned")
	}

	return chatCompletion.Choices[0].Message.Content, nil
}

// Chat with ai tool calls. use with agents with tool configurations
func (p *OpenAIProvider) CompleteConversationWithTools(ctx context.Context, conversation Conversation, model string, tools []openai.ChatCompletionToolParam) (*openai.ChatCompletion, error) {
	messages := EnsureSystemMessage(conversation.Messages)

	if model == "" {
		model = openai.ChatModelGPT3_5Turbo
	}

	params := openai.ChatCompletionNewParams{
		Messages: ToOpenAIMessages(messages),
		Model:    model,
	}

	if len(tools) > 0 {
		params.Tools = tools
	}

	return p.client.Chat.Completions.New(ctx, params)
}

func ToOpenAIMessages(messages []Message) []openai.ChatCompletionMessageParamUnion {
	var result []openai.ChatCompletionMessageParamUnion
	for _, m := range messages {
		switch m.Role {
		case RoleSystem:
			result = append(result, openai.SystemMessage(m.Content))
		case RoleUser:
			result = append(result, openai.UserMessage(m.Content))
		case RoleAssistant:
			result = append(result, openai.AssistantMessage(m.Content))
			// TODO: FIX THIS UP
		// case RoleTool:
		// 	result = append(result, openai.ToolMessage(m.Content))
		// case RoleDeveloper:
		// 	result = append(result, openai.DeveloperMessage(m.Content))
		default:
			panic("unknown role: " + m.Role)
		}
	}
	return result
}
