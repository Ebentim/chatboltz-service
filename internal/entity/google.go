package entity

type GoogleTokens struct {
	AccessToken  string `json:"access_token" gorm:"type:text;not null"`
	RefreshToken string `json:"refresh_token" gorm:"type:text;not null"`
	ExpiresIn    string `json:"expires_in" gorm:"type:varchar(50);not null"`
}
