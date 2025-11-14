package entity

import (
	"database/sql/driver"
	"time"

	"github.com/lib/pq"
)

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

func (AgentAppearance) TableName() string {
	return "agent_appearances"
}

func (AgentBehavior) TableName() string {
	return "agent_behaviors"
}

func (AgentChannel) TableName() string {
	return "agent_channels"
}

func (AgentIntegration) TableName() string {
	return "agent_integrations"
}

func (AgentStats) TableName() string {
	return "agent_stats"
}

func (TrainingData) TableName() string {
	return "training_data"
}

type Agent struct {
	ID               string            `json:"id" gorm:"primaryKey;type:varchar(36)"`
	UserId           string            `json:"userId" gorm:"type:varchar(36);not null;index"`
	Name             string            `json:"name" gorm:"type:varchar(255);not null;index"`
	Description      string            `json:"description" gorm:"type:text;not null"`
	AgentType        AgentType         `json:"agent_type" gorm:"type:int;not null;index"`
	AiModelId        string            `json:"ai_model_id" gorm:"type:varchar(36);not null;index"`
	Status           AgentStatus       `json:"status" gorm:"type:varchar(20);not null;index"`
	CreatedAt        string            `json:"created_at" gorm:"not null"`
	UpdatedAt        string            `json:"updated_at" gorm:"not null"`
	User             *Users            `json:"user,omitempty" gorm:"foreignKey:UserId;references:ID;constraint:OnDelete:CASCADE,-:save,-:update"`
	AiModel          *AiModel          `json:"ai_model,omitempty" gorm:"foreignKey:AiModelId;references:ID;constraint:OnDelete:RESTRICT,-:save,-:update"`
	AgentAppearance  *AgentAppearance  `json:"agent_appearance,omitempty" gorm:"foreignKey:AgentId;references:ID;constraint:OnDelete:CASCADE,-:save,-:update"`
	AgentBehavior    *AgentBehavior    `json:"agent_behavior,omitempty" gorm:"foreignKey:AgentId;references:ID;constraint:OnDelete:CASCADE,-:save,-:update"`
	AgentChannel     *AgentChannel     `json:"agent_channel,omitempty" gorm:"foreignKey:AgentId;references:ID;constraint:OnDelete:CASCADE,-:save,-:update"`
	AgentIntegration *AgentIntegration `json:"agent_integration,omitempty" gorm:"foreignKey:AgentId;references:ID;constraint:OnDelete:CASCADE,-:save,-:update"`
	AgentStats       *AgentStats       `json:"agent_stats,omitempty" gorm:"foreignKey:AgentId;references:ID;constraint:OnDelete:CASCADE,-:save,-:update"`
	TrainingData     []TrainingData    `json:"training_data,omitempty" gorm:"foreignKey:AgentId;references:ID;constraint:OnDelete:CASCADE,-:save,-:update"`
}

type AgentUpdate struct {
	Name        *string      `json:"name,omitempty"`
	Description *string      `json:"description,omitempty"`
	AgentType   *AgentType   `json:"agent_type" gorm:"enum('multimodal', 'text', 'audio')"`
	AiModelId   *string      `json:"ai_model_id,omitempty"`
	Status      *AgentStatus `json:"status,omitempty"`
}

type AgentResponse struct {
	ID          string      `json:"id"`
	UserId      string      `json:"userId"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	AgentType   AgentType   `json:"agent_type"`
	AiModelId   string      `json:"ai_model_id"`
	AiModel     *AiModel    `json:"ai_model,omitempty"`
	Status      AgentStatus `json:"status"`
	CreatedAt   string      `json:"created_at"`
	UpdatedAt   string      `json:"updated_at"`
}

type AgentDetailResponse struct {
	Agent            AgentResponse      `json:"agent"`
	AgentAppearance  *AgentAppearance   `json:"agent_appearance,omitempty"`
	AgentBehavior    *AgentBehavior     `json:"agent_behavior,omitempty"`
	AgentChannel     *AgentChannel      `json:"agent_channel,omitempty"`
	AgentIntegration *AgentIntegration  `json:"agent_integration,omitempty"`
	AgentStats       *AgentStats        `json:"agent_stats,omitempty"`
	TrainingData     []TrainingData     `json:"training_data,omitempty"`
}

type AgentAppearanceUpdate struct {
	PrimaryColor   *string `json:"primary_color,omitempty"`
	FontFamily     *string `json:"font_family,omitempty"`
	ChatIcon       *string `json:"chat_icon,omitempty"`
	WelcomeMessage *string `json:"welcome_message,omitempty"`
	Position       *string `json:"position,omitempty"`
	IconSize       *string `json:"icon_size,omitempty"`
	BubbleStyle    *string `json:"bubble_style,omitempty"`
}

type AgentBehaviorUpdate struct {
	FallbackMessage     *string  `json:"fallback_message,omitempty"`
	EnableHumanHandoff  *bool    `json:"enable_human_handoff,omitempty"`
	OfflineMessage      *string  `json:"offline_message,omitempty"`
	SystemInstructionId *string  `json:"system_instruction_id,omitempty"`
	PromptTemplateId    *string  `json:"prompt_template_id,omitempty"`
	Temperature         *float64 `json:"temperature,omitempty"`
	MaxTokens           *int     `json:"max_tokens,omitempty"`
}

type AgentChannelUpdate struct {
	ChannelId StringArray `json:"channel_id,omitempty"`
}

type AgentIntegrationUpdate struct {
	IntegrationId StringArray `json:"integration_id,omitempty"`
	ApiKey        *string     `json:"api_key,omitempty"`
	ApiSecret     *string     `json:"api_secret,omitempty"`
	IsActive      *bool       `json:"is_active,omitempty"`
}

type AgentAppearance struct {
	ID             string `json:"id" gorm:"primaryKey;type:varchar(36)"`
	AgentId        string `json:"agent_id" gorm:"type:varchar(36);not null;uniqueIndex"`
	PrimaryColor   string `json:"primary_color" gorm:"type:varchar(50);not null"`
	FontFamily     string `json:"font_family" gorm:"type:varchar(100);not null"`
	ChatIcon       string `json:"chat_icon" gorm:"type:varchar(100);not null"`
	WelcomeMessage string `json:"welcome_message" gorm:"type:text;not null"`
	Position       string `json:"position" gorm:"type:varchar(20);not null"`
	IconSize       string `json:"icon_size" gorm:"type:varchar(20);not null"`
	BubbleStyle    string `json:"bubble_style" gorm:"type:varchar(20);not null"`
	CreatedAt      string `json:"created_at" gorm:"not null"`
	UpdatedAt      string `json:"updated_at" gorm:"not null"`
	Agent          *Agent `json:"agent,omitempty" gorm:"foreignKey:AgentId;references:ID;constraint:OnDelete:CASCADE,-:save,-:update"`
}

type AgentBehavior struct {
	ID                  string             `json:"id" gorm:"primaryKey;type:varchar(36)"`
	AgentId             string             `json:"agent_id" gorm:"type:varchar(36);not null;uniqueIndex"`
	FallbackMessage     string             `json:"fallback_message" gorm:"type:text;not null"`
	EnableHumanHandoff  bool               `json:"enable_human_handoff" gorm:"type:boolean;default:false"`
	OfflineMessage      string             `json:"offline_message" gorm:"type:text"`
	SystemInstructionId *string            `json:"system_instruction_id" gorm:"type:varchar(36);index"`
	PromptTemplateId    *string            `json:"prompt_template_id" gorm:"type:varchar(36);index"`
	Temperature         float64            `json:"temperature" gorm:"type:decimal(3,2);default:0.7"`
	MaxTokens           int                `json:"max_tokens" gorm:"type:int;default:2048"`
	CreatedAt           string             `json:"created_at" gorm:"not null"`
	UpdatedAt           string             `json:"updated_at" gorm:"not null"`
	Agent               *Agent             `json:"agent,omitempty" gorm:"foreignKey:AgentId;references:ID;constraint:OnDelete:CASCADE,-:save,-:update"`
	SystemInstruction   *SystemInstruction `json:"system_instruction,omitempty" gorm:"foreignKey:SystemInstructionId;references:ID;constraint:OnDelete:SET NULL,-:save,-:update"`
	PromptTemplate      *PromptTemplate    `json:"prompt_template,omitempty" gorm:"foreignKey:PromptTemplateId;references:ID;constraint:OnDelete:SET NULL,-:save,-:update"`
}

type StringArray []string

func (a StringArray) Value() (driver.Value, error) {
	return pq.Array(a).Value()
}

func (a *StringArray) Scan(value interface{}) error {
	return pq.Array(a).Scan(value)
}

type AgentChannel struct {
	ID        string      `json:"id" gorm:"primaryKey;type:varchar(36)"`
	AgentId   string      `json:"agent_id" gorm:"type:varchar(36);not null;uniqueIndex"`
	ChannelId StringArray `json:"channel_id" gorm:"type:text[];not null"`
	CreatedAt string      `json:"created_at" gorm:"not null"`
	UpdatedAt string      `json:"updated_at" gorm:"not null"`
	Agent     *Agent      `json:"agent,omitempty" gorm:"foreignKey:AgentId;references:ID;constraint:OnDelete:CASCADE,-:save,-:update"`
}
type AgentIntegration struct {
	ID            string      `json:"id" gorm:"primaryKey;type:varchar(36)"`
	AgentId       string      `json:"agent_id" gorm:"type:varchar(36);not null;uniqueIndex"`
	IntegrationId StringArray `json:"integration_id" gorm:"type:text[];not null"`
	ApiKey        *string     `json:"api_key" gorm:"type:varchar(500)"`
	ApiSecret     *string     `json:"api_secret" gorm:"type:varchar(500)"`
	IsActive      bool        `json:"is_active" gorm:"type:boolean;default:false;index"`
	CreatedAt     string      `json:"created_at" gorm:"not null"`
	UpdatedAt     string      `json:"updated_at" gorm:"not null"`
	Agent         *Agent      `json:"agent,omitempty" gorm:"foreignKey:AgentId;references:ID;constraint:OnDelete:CASCADE,-:save,-:update"`
}

type AgentStats struct {
	ID               string    `json:"id" gorm:"primaryKey;type:varchar(36)"`
	AgentId          string    `json:"agent_id" gorm:"type:varchar(36);not null;uniqueIndex"`
	TotalMessages    int       `json:"total_messages" gorm:"type:int;default:0;index"`
	UniqueUsers      int       `json:"unique_users" gorm:"type:int;default:0;index"`
	AverageRating    float64   `json:"average_rating" gorm:"type:decimal(3,2);default:0.0;index"`
	ResponseRate     float64   `json:"response_rate" gorm:"type:decimal(5,2);default:0.0;index"`
	ConversionsCount int       `json:"conversions_count" gorm:"type:int;default:0;index"`
	LastCalculatedAt time.Time `json:"last_calculated_at" gorm:"type:timestamp;index"`
	CreatedAt        string    `json:"created_at" gorm:"not null"`
	UpdatedAt        string    `json:"updated_at" gorm:"not null"`
	Agent            *Agent    `json:"agent,omitempty" gorm:"foreignKey:AgentId;references:ID;constraint:OnDelete:CASCADE,-:save,-:update"`
}

type TrainingData struct {
	ID          string          `json:"id" gorm:"primaryKey;type:varchar(36)"`
	AgentId     string          `json:"agent_id" gorm:"type:varchar(36);not null;index"`
	ContentType string          `json:"content_type" gorm:"type:varchar(50);not null;index"`
	Content     []TrainingTexts `json:"content" gorm:"type:jsonb"`
	IsActive    bool            `json:"is_active" gorm:"type:boolean;default:false;index"`
	CreatedAt   string          `json:"created_at" gorm:"not null"`
	UpdatedAt   string          `json:"updated_at" gorm:"not null"`
	Agent       *Agent          `json:"agent,omitempty" gorm:"foreignKey:AgentId;references:ID;constraint:OnDelete:CASCADE,-:save,-:update"`
}
