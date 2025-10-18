package usecase

import (
	"github.com/alpinesboltltd/boltz-ai/internal/entity"
	appErrors "github.com/alpinesboltltd/boltz-ai/internal/errors"
	"github.com/alpinesboltltd/boltz-ai/internal/repository"
)

type SystemUsecase struct {
	System repository.SystemRepositoryInterface
}

func NewSystemUsecase(systemRepo repository.SystemRepositoryInterface) *SystemUsecase {
	return &SystemUsecase{
		System: systemRepo,
	}
}

func (u *SystemUsecase) CreateSystemInstruction(title, content, createdBy string, templateId *string) (*entity.SystemInstruction, error) {
	if title == "" || content == "" || createdBy == "" {
		return nil, appErrors.NewValidationError("Title, content, and created_by are required")
	}

	if templateId != nil && *templateId != "" {
		if _, err := u.System.GetPromptTemplate(*templateId); err != nil {
			return nil, appErrors.NewValidationError("Invalid template ID")
		}
	}

	return u.System.CreateSystemInstruction(title, content, createdBy, templateId)
}

func (u *SystemUsecase) GetSystemInstruction(id string) (*entity.SystemInstruction, error) {
	if id == "" {
		return nil, appErrors.NewValidationError("System instruction ID is required")
	}
	return u.System.GetSystemInstruction(id)
}

func (u *SystemUsecase) UpdateSystemInstruction(id, title, content string) (*entity.SystemInstruction, error) {
	if id == "" {
		return nil, appErrors.NewValidationError("System instruction ID is required")
	}

	instruction, err := u.System.GetSystemInstruction(id)
	if err != nil {
		return nil, err
	}

	if title != "" {
		instruction.Title = title
	}
	if content != "" {
		instruction.Content = content
	}

	if err := u.System.UpdateSystemInstruction(instruction); err != nil {
		return nil, err
	}
	return instruction, nil
}

func (u *SystemUsecase) DeleteSystemInstruction(id string) error {
	if id == "" {
		return appErrors.NewValidationError("System instruction ID is required")
	}
	return u.System.DeleteSystemInstruction(id)
}

func (u *SystemUsecase) ListSystemInstructions() (*[]entity.SystemInstruction, error) {
	return u.System.ListSystemInstructions()
}

func (u *SystemUsecase) CreatePromptTemplate(title, content string) (*entity.PromptTemplate, error) {
	if title == "" || content == "" {
		return nil, appErrors.NewValidationError("Title and content are required")
	}
	return u.System.CreatePromptTemplate(title, content)
}

func (u *SystemUsecase) GetPromptTemplate(id string) (*entity.PromptTemplate, error) {
	if id == "" {
		return nil, appErrors.NewValidationError("Prompt template ID is required")
	}
	return u.System.GetPromptTemplate(id)
}

func (u *SystemUsecase) ListPromptTemplates() (*[]entity.PromptTemplate, error) {
	return u.System.ListPromptTemplates()
}
