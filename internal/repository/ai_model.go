package repository

import (
	"errors"
	"time"

	"github.com/alpinesboltltd/boltz-ai/internal/entity"
	appErrors "github.com/alpinesboltltd/boltz-ai/internal/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AiModelRepository struct {
	db *gorm.DB
}

func NewAiModelRepository(db *gorm.DB) AiModelRepositoryInterface {
	return &AiModelRepository{db: db}
}

func (r *AiModelRepository) CreateAiModel(name, provider string, creditsPer1k int, supportsText, supportsVision, supportsVoice, isReasoning bool) (*entity.AiModel, error) {
	model := &entity.AiModel{
		ID:             uuid.New().String(),
		Name:           name,
		Provider:       provider,
		CreditsPer1k:   creditsPer1k,
		SupportsText:   supportsText,
		SupportsVision: supportsVision,
		SupportsVoice:  supportsVoice,
		IsReasoning:    isReasoning,
		CreatedAt:      time.Now().UTC().Format(time.RFC3339),
		UpdatedAt:      time.Now().UTC().Format(time.RFC3339),
	}

	if err := r.db.Create(model).Error; err != nil {
		return nil, appErrors.WrapDatabaseError(err, "create ai model")
	}

	return model, nil
}

func (r *AiModelRepository) GetAiModel(id string) (*entity.AiModel, error) {
	var model entity.AiModel
	if err := r.db.Where("id = ?", id).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.NewNotFoundError("AI model not found")
		}
		return nil, appErrors.WrapDatabaseError(err, "get ai model")
	}
	return &model, nil
}

func (r *AiModelRepository) GetAiModelByName(name string) (*entity.AiModel, error) {
	var model entity.AiModel
	if err := r.db.Where("name = ?", name).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.NewNotFoundError("AI model not found")
		}
		return nil, appErrors.WrapDatabaseError(err, "get ai model by name")
	}
	return &model, nil
}

func (r *AiModelRepository) ListAiModels() (*[]entity.AiModel, error) {
	var models []entity.AiModel
	if err := r.db.Find(&models).Error; err != nil {
		return nil, appErrors.WrapDatabaseError(err, "list ai models")
	}
	return &models, nil
}

func (r *AiModelRepository) ListAiModelsByProvider(provider string) (*[]entity.AiModel, error) {
	var models []entity.AiModel
	if err := r.db.Where("provider = ?", provider).Find(&models).Error; err != nil {
		return nil, appErrors.WrapDatabaseError(err, "list ai models by provider")
	}
	return &models, nil
}

func (r *AiModelRepository) UpdateAiModel(model *entity.AiModel) error {
	model.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	return r.db.Save(model).Error
}

func (r *AiModelRepository) DeleteAiModel(id string) error {
	if err := r.db.Where("id = ?", id).Delete(&entity.AiModel{}).Error; err != nil {
		return appErrors.WrapDatabaseError(err, "delete ai model")
	}
	return nil
}
