package entity

type Channel int

type AuthType int

const (
	Bearer AuthType = iota
	Basic
	ApiKey
	none
)

var AuthTypeString = map[AuthType]string{
	Bearer: "bearer",
	Basic:  "basic",
	ApiKey: "api_key",
	none:   "none",
}

const (
	Slack Channel = iota
	Email
	Telegram
	WhatsApp
	Facebook
	WebsiteWidget
)

var ChannelString = map[Channel]string{
	Slack:         "slack",
	Email:         "email",
	Telegram:      "telegram",
	WhatsApp:      "whatsapp",
	Facebook:      "facebook",
	WebsiteWidget: "website_widget",
}

type Integration int

const (
	Shopify Integration = iota
	Zapier
	GoogleSheets
	GoogleDocs
	Notion
)

var IntegrationString = map[Integration]string{
	Shopify:      "shopify",
	Zapier:       "zapier",
	GoogleSheets: "google_sheets",
	GoogleDocs:   "google_docs",
	Notion:       "notion"}

type SystemInstruction struct {
	Id      string `json:"id" gorm:"primaryKey;type:varchar(36)"`
	Title   string `json:"title" gorm:"type:varchar(255);not null"`
	Content string `json:"content" gorm:"type:text;not null"`
}

type PromptTemplate struct {
	Id      string `json:"id" gorm:"primaryKey;type:varchar(36)"`
	Title   string `json:"title" gorm:"type:varchar(255);not null"`
	Content string `json:"content" gorm:"type:text;not null"`
}

type Channels struct {
	Id   string `json:"id" gorm:"primaryKey;type:varchar(36)"`
	Name string `json:"name" gorm:"type:varchar(100);not null"`
}

type Integrations struct {
	Id   string `json:"id" gorm:"primaryKey;type:varchar(36)"`
	Name string `json:"name" gorm:"type:varchar(100);not null"`
}

type ApiFunctions struct {
	Id            string            `json:"id" gorm:"primaryKey;type:varchar(36)"`
	Name          string            `json:"name" gorm:"type:varchar(255);not null"`
	Description   string            `json:"description" gorm:"type:text;not null"`
	Method        string            `json:"method" gorm:"type:varchar(10);not null"`
	Url           string            `json:"url" gorm:"type:text;not null"`
	AuthType      string            `json:"auth_type" gorm:"type:varchar(50);not null"`
	Headers       map[string]string `json:"headers" gorm:"type:jsonb"`
	Params        map[string]string `json:"params" gorm:"type:jsonb"`
	Body          string            `json:"body" gorm:"type:text"`
	ErrorHandling ApiErrorHandling  `json:"error_handling" gorm:"embedded"`
	CreatedAt     string            `json:"created_at" gorm:"not null"`
	UpdatedAt     string            `json:"updated_at" gorm:"not null"`
}

type ApiErrorHandling struct {
	Id              string  `json:"id" gorm:"primaryKey;type:varchar(36)"`
	Retries         int     `json:"retries" gorm:"type:int;default:3"`
	RetryDelay      int     `json:"retry_delay" gorm:"type:int;default:1000"`
	Timeout         int     `json:"timeout" gorm:"type:int;default:30000"`
	FallbackMessage *string `json:"fallback_message" gorm:"type:text"`
	CreatedAt       string  `json:"created_at" gorm:"not null"`
	UpdatedAt       string  `json:"updated_at" gorm:"not null"`
}
