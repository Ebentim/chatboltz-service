package entity

import "time"

type TrainingTexts struct {
	Title    string   `json:"title"`
	Content  string   `json:"content"`
	Keywords []string `json:"keywords"`
}

type AgentType int

const (
	Multimodal AgentType = iota
	TextOnly
	VoiceOnly
)

var AgentTypeString = map[AgentType]string{
	Multimodal: "multimodal",
	TextOnly:   "text",
	VoiceOnly:  "audio",
}

type AgentStatus string

const (
	Active   AgentStatus = "active"
	Inactive AgentStatus = "inactive"
	Draft    AgentStatus = "draft"
)

func (Agent) TableName() string {
	return "agents"
}

type Agent struct {
	ID           string      `json:"id" gorm:"primaryKey;type:varchar(36)"`
	UserId       string      `json:"userId" gorm:"type:varchar(36); not null"`
	Name         string      `json:"name" gorm:"type:varchar(36);not null"`
	Description  string      `json:"description" gorm:"type:text;not null"`
	AgentType    AgentType   `json:"agent_type" gorm:"enum('multimodal', 'text', 'audio')"` //FIXME: USE THE CAR TYPE ENUMS
	AiModel      string      `json:"ai_model" gorm:"type:varchar(255);not null"`
	AiProvider   string      `json:"ai_provider" gorm:"type:varchar(30);not null"`
	Capabilities []string    `json:"capabilities" gorm:"type:text[];not null"` // ["text", "voice", "vision"]
	CreditsPer1k int         `json:"credits_per_1k" gorm:"type:int;not null"`
	Status       AgentStatus `json:"status" gorm:"enum('active', 'inactive', 'draft'); not null"`
	CreatedAt    string      `json:"created_at" gorm:"not null"`
	UpdatedAt    string      `json:"updated_at" gorm:"not null"`
}

type AgentUpdate struct {
	Name         *string      `json:"name,omitempty"`
	Description  *string      `json:"description,omitempty"`
	AgentType    *AgentType   `json:"agent_type" gorm:"enum('multimodal', 'text', 'audio')"` //FIXME: USE THE CAR TYPE ENUMS
	AiModel      *string      `json:"ai_model,omitempty"`
	AiProvider   *string      `json:"ai_provider,omitempty"`
	Capabilities *[]string    `json:"capabilities,omitempty"`
	CreditsPer1k *int         `json:"credits_per_1k,omitempty"`
	Status       *AgentStatus `json:"status,omitempty"`
}

type AgentAppearance struct {
	ID             string `json:"id" gorm:"primaryKey;type:varchar(36)"`
	AgentId        string `json:"agent_id" gorm:"type:varchar(36);not null"`
	PrimaryColor   string `json:"primary_color" gorm:"type:varchar(36);not null"`
	FontFamily     string `json:"font_family" gorm:"type:varchar(36);not null"`
	ChatIcon       string `json:"chat_icon" gorm:"type:varchar(36);not null"`
	WelcomeMessage string `json:"welcome_message" gorm:"type:text;not null"`
	Position       string `json:"position" gorm:"enum('bottom-left', 'bottom-right');not null"`
	IconSize       string `json:"icon_size" gorm:"enum('small', 'medium', 'large');not null"`
	BubbleStyle    string `json:"bubble_style" gorm:"enum('square', 'round');not null"`
	CreatedAt      string `json:"created_at" gorm:"not null"`
	UpdatedAt      string `json:"updated_at" gorm:"not null"`
}

type AgentBehavior struct {
	ID                  string  `json:"id" gorm:"primaryKey;type:varchar(36)"`
	AgentId             string  `json:"agent_id" gorm:"type:varchar(36);not null"`
	FallbackMessage     string  `json:"fallback_message" gorm:"type:text;not null"`
	EnableHumanHandoff  bool    `json:"enable_human_handoff" gorm:"type:boolean;default:false"`
	OfflineMessage      string  `json:"offline_message"`
	SystemInstructionId string  `json:"system_instruction_id"`
	PromptTemplateId    string  `json:"prompt_template_id"`
	Temperature         float64 `json:"temperature"`
	MaxTokens           int     `json:"max_tokens"`
	CreatedAt           string  `json:"created_at"`
	UpdatedAt           string  `json:"updated_at"`
}

// Agent Channel describes where the agent is deployed
type AgentChannel struct {
	ID        string   `json:"id" gorm:"primaryKey;type:varchar(36)"`
	AgentId   string   `json:"agent_id" gorm:"type:varchar(36);not null"`
	ChannelId []string `json:"channel_id" gorm:"type:text[];not null"`
	CreatedAt string   `json:"created_at" gorm:"not null"`
	UpdatedAt string   `json:"updated_at" gorm:"not null"`
}
type AgentIntegration struct {
	ID            string   `json:"id" gorm:"primaryKey;type:varchar(36)"`
	AgentId       string   `json:"agent_id" gorm:"type:varchar(36);not null"`
	IntegrationId []string `json:"integration_id" gorm:"type:text[];not null"`
	ApiKey        *string  `json:"api_key" gorm:"type:varchar(255)"`
	ApiSecret     *string  `json:"api_secret" gorm:"type:varchar(255)"`
	IsActive      bool     `json:"is_active" gorm:"type:boolean;default:false"`
	CreatedAt     string   `json:"created_at" gorm:"not null"`
	UpdatedAt     string   `json:"updated_at" gorm:"not null"`
}

type AgentStats struct {
	ID               string    `json:"id" gorm:"primaryKey;type:varchar(36)"`
	AgentId          string    `json:"agent_id" gorm:"type:varchar(36);not null"`
	TotalMessages    int       `json:"total_messages" gorm:"type:int;default:0"`
	UniqueUsers      int       `json:"unique_users" gorm:"type:int;default:0"`
	AverageRating    float64   `json:"average_rating" gorm:"type:decimal(3,2);default:0.0"`
	ResponseRate     float64   `json:"response_rate" gorm:"type:decimal(5,2);default:0.0"`
	ConversionsCount int       `json:"conversions_count" gorm:"type:int;default:0"`
	LastCalculatedAt time.Time `json:"last_calculated_at" gorm:"type:timestamp"`
	CreatedAt        string    `json:"created_at" gorm:"not null"`
	UpdatedAt        string    `json:"updated_at" gorm:"not null"`
}

type TrainingData struct {
	ID          string          `json:"id" gorm:"primaryKey;type:varchar(36)"`
	AgentId     string          `json:"agent_id" gorm:"type:varchar(36);not null"`
	ContentType string          `json:"content_type" gorm:"type:varchar(50);not null"`
	Content     []TrainingTexts `json:"content" gorm:"type:jsonb"`
	IsActive    bool            `json:"is_active" gorm:"type:boolean;default:false"`
	CreatedAt   string          `json:"created_at" gorm:"not null"`
	UpdatedAt   string          `json:"updated_at" gorm:"not null"`
}
