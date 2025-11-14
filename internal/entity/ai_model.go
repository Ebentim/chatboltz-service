package entity

func (AiModel) TableName() string {
	return "ai_models"
}

type AiModel struct {
	ID             string `json:"id" gorm:"primaryKey;type:varchar(36)"`
	Name           string `json:"name" gorm:"type:varchar(255);not null;uniqueIndex"`
	Provider       string `json:"provider" gorm:"type:varchar(50);not null;index"`
	CreditsPer1k   int    `json:"credits_per_1k" gorm:"type:int;not null"`
	SupportsText   bool   `json:"supports_text" gorm:"type:boolean;default:true"`
	SupportsVision bool   `json:"supports_vision" gorm:"type:boolean;default:false"`
	SupportsVoice  bool   `json:"supports_voice" gorm:"type:boolean;default:false"`
	IsReasoning    bool   `json:"is_reasoning" gorm:"type:boolean;default:false"`
	CreatedAt      string `json:"created_at" gorm:"not null"`
	UpdatedAt      string `json:"updated_at" gorm:"not null"`
}
