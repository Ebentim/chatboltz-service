package entity

import "time"

type Role int

const (
	Assistant Role = iota
	User
	Human
)

type MessageType int

const (
	Text MessageType = iota
	Audio
	Image
	Video
)

var MessageTypeString = map[MessageType]string{
	Text:  "text",
	Audio: "audio",
	Image: "image",
	Video: "video",
}

var MessageRole = map[Role]string{
	Assistant: "assistant",
	User:      "user",
	Human:     "human",
}

type Conversation struct {
	Id               string    `json:"id" gorm:"primaryKey;type:varchar(36)"`
	AgentId          string    `json:"agent_id" gorm:"type:varchar(36);not null"`
	Platform         string    `json:"platform" gorm:"type:varchar(50);not null"`
	ClientId         string    `json:"client_id" gorm:"type:varchar(36);not null"`
	Title            string    `json:"title" gorm:"type:varchar(255);not null"`
	CreatedAt        time.Time `json:"created_at" gorm:"not null"`
	Status           string    `json:"status" gorm:"type:varchar(50);not null"`
	EscalatedToHuman bool      `json:"escalated_to_human" gorm:"type:boolean;default:false"`
	EscalationReason string    `json:"escalation_reason" gorm:"type:text"`
}

type Message struct {
	Id              string    `json:"id" gorm:"primaryKey;type:varchar(36)"`
	ConversationId  string    `json:"conversation_id" gorm:"type:varchar(36);not null"`
	UserId          *string   `json:"user_id" gorm:"type:varchar(36)"`
	Role            string    `json:"role" gorm:"type:varchar(50);not null"`
	Text            string    `json:"text" gorm:"type:text;not null"`
	Timestamp       time.Time `json:"timestamp" gorm:"not null"`
	ConfidenceScore float64   `json:"confidence_score" gorm:"type:decimal(5,4);default:0.0"`
}

type MessageMetadata struct {
	Id          string    `json:"id" gorm:"primaryKey;type:varchar(36)"`
	MessageId   string    `json:"message_id" gorm:"type:varchar(36);not null"`
	Sentiment   string    `json:"sentiment" gorm:"type:varchar(50)"`
	Intent      string    `json:"intent" gorm:"type:varchar(100)"`
	MessageType string    `json:"message_type" gorm:"type:varchar(50);not null"`
	CreatedAt   time.Time `json:"created_at" gorm:"not null"`
	TokenCount  int       `json:"token_count" gorm:"type:int;default:0"`
}
