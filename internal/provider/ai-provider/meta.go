package aiprovider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/alpinesboltltd/boltz-ai/internal/entity"
)

const metaAiUrl = "https://api.llama-api.com/chat/completions" //FIXME: Example endpoint

type MetaProvider struct {
	apiKey  string
	baseURL string
}

type MetaRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type MetaResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func NewMetaClient(apiKey string) *MetaProvider {
	return &MetaProvider{
		apiKey:  apiKey,
		baseURL: metaAiUrl,
	}
}

func (p *MetaProvider) GetCapabilities() entity.ModelCapabilities {
	return entity.ModelCapabilities{Text: true, Voice: false, Vision: false}
}

func (p *MetaProvider) CompleteMultimodalConversation(messages []MultimodalMessage, config map[string]interface{}) (string, error) {
	// Convert to regular conversation since Meta doesn't support multimodal
	// Add image descriptions to text content
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

func (p *MetaProvider) CompleteConversation(conversation Conversation, config map[string]interface{}) (string, error) {
	messages := EnsureSystemMessage(conversation.Messages)

	model := "llama-3.1-70b"
	if m, ok := config["model"].(string); ok && m != "" {
		model = m
	}

	reqBody := MetaRequest{
		Model:    model,
		Messages: messages,
		Stream:   false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %w", err)
	}

	req, err := http.NewRequest("POST", p.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API error: %s", string(body))
	}

	var metaResp MetaResponse
	if err := json.Unmarshal(body, &metaResp); err != nil {
		return "", fmt.Errorf("error unmarshaling response: %w", err)
	}

	if len(metaResp.Choices) == 0 {
		return "", fmt.Errorf("no response choices returned")
	}

	return metaResp.Choices[0].Message.Content, nil
}

func (p *MetaProvider) CompleteConversationStream(conversation Conversation, config map[string]interface{}, callback StreamCallback) error {
	// Meta doesn't support streaming
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
