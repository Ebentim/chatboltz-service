package repository

import (
	"errors"
	"time"

	"github.com/alpinesboltltd/boltz-ai/internal/entity"
	appErrors "github.com/alpinesboltltd/boltz-ai/internal/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SystemRepository struct {
	db *gorm.DB
}

func NewSystemRepository(db *gorm.DB) SystemRepositoryInterface {
	return &SystemRepository{db: db}
}

func (r *SystemRepository) CreateSystemInstruction(title, content, createdBy string, templateId *string) (*entity.SystemInstruction, error) {
	instruction := &entity.SystemInstruction{
		ID:         uuid.New().String(),
		Title:      title,
		Content:    content,
		TemplateId: templateId,
		CreatedBy:  createdBy,
		CreatedAt:  time.Now().UTC().Format(time.RFC3339),
		UpdatedAt:  time.Now().UTC().Format(time.RFC3339),
	}

	if err := r.db.Create(instruction).Error; err != nil {
		return nil, appErrors.WrapDatabaseError(err, "create system instruction")
	}
	return instruction, nil
}

func (r *SystemRepository) GetSystemInstruction(id string) (*entity.SystemInstruction, error) {
	var instruction entity.SystemInstruction
	if err := r.db.Where("id = ?", id).First(&instruction).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.NewNotFoundError("System instruction not found")
		}
		return nil, appErrors.WrapDatabaseError(err, "get system instruction")
	}
	return &instruction, nil
}

func (r *SystemRepository) UpdateSystemInstruction(instruction *entity.SystemInstruction) error {
	instruction.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	return r.db.Save(instruction).Error
}

func (r *SystemRepository) DeleteSystemInstruction(id string) error {
	if err := r.db.Where("id = ?", id).Delete(&entity.SystemInstruction{}).Error; err != nil {
		return appErrors.WrapDatabaseError(err, "delete system instruction")
	}
	return nil
}

func (r *SystemRepository) ListSystemInstructions() (*[]entity.SystemInstruction, error) {
	var instructions []entity.SystemInstruction
	if err := r.db.Find(&instructions).Error; err != nil {
		return nil, appErrors.WrapDatabaseError(err, "list system instructions")
	}
	return &instructions, nil
}

func (r *SystemRepository) CreatePromptTemplate(title, content string) (*entity.PromptTemplate, error) {
	template := &entity.PromptTemplate{
		ID:        uuid.New().String(),
		Title:     title,
		Content:   content,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
	}

	if err := r.db.Create(template).Error; err != nil {
		return nil, appErrors.WrapDatabaseError(err, "create prompt template")
	}
	return template, nil
}

func (r *SystemRepository) GetPromptTemplate(id string) (*entity.PromptTemplate, error) {
	var template entity.PromptTemplate
	if err := r.db.Where("id = ?", id).First(&template).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.NewNotFoundError("Prompt template not found")
		}
		return nil, appErrors.WrapDatabaseError(err, "get prompt template")
	}
	return &template, nil
}

func (r *SystemRepository) ListPromptTemplates() (*[]entity.PromptTemplate, error) {
	var templates []entity.PromptTemplate
	if err := r.db.Find(&templates).Error; err != nil {
		return nil, appErrors.WrapDatabaseError(err, "list prompt templates")
	}
	return &templates, nil
}
