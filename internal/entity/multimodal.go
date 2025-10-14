package entity

import "time"

type MediaType string

const (
	MediaTypeText  MediaType = "text"
	MediaTypeAudio MediaType = "audio"
	MediaTypeImage MediaType = "image"
	MediaTypeVideo MediaType = "video"
)

type MultimodalMessage struct {
	ID          string    `json:"id" gorm:"primaryKey;type:varchar(36)"`
	Role        string    `json:"role" gorm:"type:varchar(50);not null"`
	Content     string    `json:"content,omitempty" gorm:"type:text"`
	MediaType   MediaType `json:"media_type" gorm:"type:varchar(50);not null"`
	MediaURL    string    `json:"media_url,omitempty" gorm:"type:text"`
	MediaData   []byte    `json:"media_data,omitempty" gorm:"type:bytea"`
	Timestamp   time.Time `json:"timestamp" gorm:"not null"`
	ProcessedBy string    `json:"processed_by" gorm:"type:varchar(100);not null"` // which provider processed this
}

type ModelCapabilityMap map[string]ModelCapabilities

var DefaultModelCapabilities = ModelCapabilityMap{
	// OpenAI
	"gpt-4o":        {Text: true, Voice: true, Vision: true},
	"gpt-4o-mini":   {Text: true, Voice: false, Vision: true},
	"gpt-4-turbo":   {Text: true, Voice: false, Vision: true},
	"gpt-3.5-turbo": {Text: true, Voice: false, Vision: false},

	// Anthropic
	"claude-3-5-sonnet": {Text: true, Voice: false, Vision: true},
	"claude-3-haiku":    {Text: true, Voice: false, Vision: true},
	"claude-3-opus":     {Text: true, Voice: false, Vision: true},

	// Google
	"gemini-2.0-flash": {Text: true, Voice: true, Vision: true},
	"gemini-1.5-pro":   {Text: true, Voice: false, Vision: true},
	"gemini-1.5-flash": {Text: true, Voice: false, Vision: true},

	// Meta
	"llama-3.2-90b": {Text: true, Voice: false, Vision: true},
	"llama-3.1-70b": {Text: true, Voice: false, Vision: false},

	// Groq
	"llama-3.3-70b-versatile":    {Text: true, Voice: false, Vision: false},
	"llama-3.1-8b-instant":       {Text: true, Voice: false, Vision: false},
	"llama-3.1-70b-versatile":    {Text: true, Voice: false, Vision: false},
	"gemma2-9b-it":               {Text: true, Voice: false, Vision: false},
	"mixtral-8x7b-32768":         {Text: true, Voice: false, Vision: false},
	"whisper-large-v3":           {Text: false, Voice: true, Vision: false},
	"whisper-large-v3-turbo":     {Text: false, Voice: true, Vision: false},
	"distil-whisper-large-v3-en": {Text: false, Voice: true, Vision: false},
}

func (m ModelCapabilityMap) GetCapabilities(model string) ModelCapabilities {
	if caps, exists := m[model]; exists {
		return caps
	}
	return ModelCapabilities{Text: true, Voice: false, Vision: false}
}
