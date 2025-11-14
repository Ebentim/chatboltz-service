package repository

import (
	"time"

	"github.com/alpinesboltltd/boltz-ai/internal/entity"
	appErrors "github.com/alpinesboltltd/boltz-ai/internal/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TokenRepository struct {
	db *gorm.DB
}

func NewUserToken(db *gorm.DB) TokenRepositoryInterface {
	return &TokenRepository{db: db}
}

func (r *TokenRepository) CreateToken(email, purpose, token string, expiresAt time.Time) error {
	t := &entity.Token{
		ID:        uuid.New().String(),
		Email:     email,
		Token:     token,
		Purpose:   purpose,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
	}

	if err := r.db.Create(t).Error; err != nil {
		return err
	}
	return nil
}

func (r *TokenRepository) GetToken(email, purpose string) (*entity.Token, error) {
	var token entity.Token
	err := r.db.Where("email = ? AND purpose = ?", email, purpose).First(&token).Error
	if err == gorm.ErrRecordNotFound {
		return nil, appErrors.WrapDatabaseError(err, "Get token")
	}

	if err != nil {
		return nil, appErrors.WrapDatabaseError(err, "Get token")
	}

	return &token, nil
}
