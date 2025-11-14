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

func (SystemInstruction) TableName() string {
	return "system_instructions"
}

func (PromptTemplate) TableName() string {
	return "prompt_templates"
}

func (Channels) TableName() string {
	return "channels"
}

func (Integrations) TableName() string {
	return "integrations"
}

func (ApiFunctions) TableName() string {
	return "api_functions"
}

type SystemInstruction struct {
	ID         string          `json:"id" gorm:"primaryKey;type:varchar(36)"`
	Title      string          `json:"title" gorm:"type:varchar(255);not null;index"`
	Content    string          `json:"content" gorm:"type:text;not null"`
	TemplateId *string         `json:"template_id" gorm:"type:varchar(36);index"`
	Template   *PromptTemplate `json:"template,omitempty" gorm:"foreignKey:TemplateId;references:ID;constraint:OnDelete:SET NULL,-:save,-:update"`
	CreatedBy  string          `json:"created_by" gorm:"type:varchar(36);not null;index"`
	User       *Users          `json:"user,omitempty" gorm:"foreignKey:CreatedBy;references:ID;constraint:OnDelete:CASCADE,-:save,-:update"`
	CreatedAt  string          `json:"created_at" gorm:"not null"`
	UpdatedAt  string          `json:"updated_at" gorm:"not null"`
}

type PromptTemplate struct {
	ID        string `json:"id" gorm:"primaryKey;type:varchar(36)"`
	Title     string `json:"title" gorm:"type:varchar(255);not null;index"`
	Content   string `json:"content" gorm:"type:text;not null"`
	CreatedAt string `json:"created_at" gorm:"not null"`
	UpdatedAt string `json:"updated_at" gorm:"not null"`
}

type Channels struct {
	ID        string `json:"id" gorm:"primaryKey;type:varchar(36)"`
	Name      string `json:"name" gorm:"type:varchar(100);not null;uniqueIndex"`
	CreatedAt string `json:"created_at" gorm:"not null"`
	UpdatedAt string `json:"updated_at" gorm:"not null"`
}

type Integrations struct {
	ID        string `json:"id" gorm:"primaryKey;type:varchar(36)"`
	Name      string `json:"name" gorm:"type:varchar(100);not null;uniqueIndex"`
	CreatedAt string `json:"created_at" gorm:"not null"`
	UpdatedAt string `json:"updated_at" gorm:"not null"`
}

type ApiFunctions struct {
	ID          string            `json:"id" gorm:"primaryKey;type:varchar(36)"`
	Name        string            `json:"name" gorm:"type:varchar(255);not null;index"`
	Description string            `json:"description" gorm:"type:text;not null"`
	Method      string            `json:"method" gorm:"type:varchar(10);not null;index"`
	Url         string            `json:"url" gorm:"type:text;not null"`
	AuthType    string            `json:"auth_type" gorm:"type:varchar(50);not null"`
	Headers     map[string]string `json:"headers" gorm:"type:jsonb"`
	Params      map[string]string `json:"params" gorm:"type:jsonb"`
	Body        string            `json:"body" gorm:"type:text"`
	Retries     int               `json:"retries" gorm:"type:int;default:3"`
	RetryDelay  int               `json:"retry_delay" gorm:"type:int;default:1000"`
	Timeout     int               `json:"timeout" gorm:"type:int;default:30000"`
	FallbackMsg *string           `json:"fallback_message" gorm:"type:text"`
	CreatedAt   string            `json:"created_at" gorm:"not null"`
	UpdatedAt   string            `json:"updated_at" gorm:"not null"`
}
