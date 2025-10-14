package aiprovider

import "github.com/alpinesboltltd/boltz-ai/internal/entity"

type Role string

const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
	RoleTool      Role = "tool"
	RoleDeveloper Role = "developer"
)

type Message struct {
	Role    Role   `json:"role"`
	Content string `json:"content"`
}

type Conversation struct {
	Messages []Message `json:"messages"`
}

type MultimodalMessage struct {
	Role        Role   `json:"role"`
	Content     string `json:"content"`
	MediaType   string `json:"media_type,omitempty"`
	MediaURL    string `json:"media_url,omitempty"`
	MediaData   []byte `json:"media_data,omitempty"`
	MediaBase64 string `json:"media_base64,omitempty"`
}

type StreamCallback func(chunk string, done bool) error

type LLMProvider interface {
	CompleteConversation(conversation Conversation, config map[string]interface{}) (string, error)
	CompleteMultimodalConversation(messages []MultimodalMessage, config map[string]interface{}) (string, error)
	CompleteConversationStream(conversation Conversation, config map[string]interface{}, callback StreamCallback) error
	GetCapabilities() entity.ModelCapabilities
}

type TTSProvider interface {
	TextToSpeech(text string, config map[string]interface{}) ([]byte, error)
}

type STTProvider interface {
	SpeechToText(audio []byte, config map[string]interface{}) (string, error)
}

func (c *Conversation) addMessage(message string, role Role) {
	c.Messages = append(c.Messages, Message{
		Role:    role,
		Content: message,
	})
}

func (c *Conversation) addModelResponse(message string) {
	c.addMessage(message, RoleAssistant)
}

func (c *Conversation) addUserMessage(message string) {
	c.addMessage(message, RoleUser)
}

func (c *Conversation) addSystemMessage(message string) {
	c.addMessage(message, RoleSystem)
}

func EnsureSystemMessage(messages []Message) []Message {
	if len(messages) == 0 || messages[0].Role != RoleSystem {
		return append([]Message{{Role: RoleSystem, Content: "You are a helpful assistant."}}, messages...)
	}
	return messages
}

// HasImageContent checks if message contains image data (URL or base64)
func (m *MultimodalMessage) HasImageContent() bool {
	return m.MediaType == "image" && (m.MediaURL != "" || m.MediaBase64 != "")
}

// IsBase64Image checks if message contains base64 image data
func (m *MultimodalMessage) IsBase64Image() bool {
	return m.MediaType == "image" && m.MediaBase64 != ""
}
