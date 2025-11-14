package usecase

import (
	"github.com/alpinesboltltd/boltz-ai/internal/entity"
	appErrors "github.com/alpinesboltltd/boltz-ai/internal/errors"
	"github.com/alpinesboltltd/boltz-ai/internal/repository"
)

type AiModelUsecase struct {
	AiModel repository.AiModelRepositoryInterface
}

func NewAiModelUseCase(aiModelRepo repository.AiModelRepositoryInterface) *AiModelUsecase {
	return &AiModelUsecase{
		AiModel: aiModelRepo,
	}
}

func (u *AiModelUsecase) CreateAiModel(name, provider string, creditsPer1k int, supportsText, supportsVision, supportsVoice, isReasoning bool) (*entity.AiModel, error) {
	if name == "" || provider == "" {
		return nil, appErrors.NewValidationError("Name and provider are required")
	}

	return u.AiModel.CreateAiModel(name, provider, creditsPer1k, supportsText, supportsVision, supportsVoice, isReasoning)
}

func (u *AiModelUsecase) GetAiModel(id string) (*entity.AiModel, error) {
	if id == "" {
		return nil, appErrors.NewValidationError("AI model ID is required")
	}
	return u.AiModel.GetAiModel(id)
}

func (u *AiModelUsecase) GetAiModelByName(name string) (*entity.AiModel, error) {
	if name == "" {
		return nil, appErrors.NewValidationError("AI model name is required")
	}
	return u.AiModel.GetAiModelByName(name)
}

func (u *AiModelUsecase) ListAiModels() (*[]entity.AiModel, error) {
	return u.AiModel.ListAiModels()
}

func (u *AiModelUsecase) ListAiModelsByProvider(provider string) (*[]entity.AiModel, error) {
	if provider == "" {
		return nil, appErrors.NewValidationError("Provider is required")
	}
	return u.AiModel.ListAiModelsByProvider(provider)
}

func (u *AiModelUsecase) UpdateAiModel(id string, name, provider *string, creditsPer1k *int, supportsText, supportsVision, supportsVoice, isReasoning *bool) (*entity.AiModel, error) {
	if id == "" {
		return nil, appErrors.NewValidationError("AI model ID is required")
	}

	model, err := u.AiModel.GetAiModel(id)
	if err != nil {
		return nil, err
	}

	if name != nil {
		model.Name = *name
	}
	if provider != nil {
		model.Provider = *provider
	}
	if creditsPer1k != nil {
		model.CreditsPer1k = *creditsPer1k
	}
	if supportsText != nil {
		model.SupportsText = *supportsText
	}
	if supportsVision != nil {
		model.SupportsVision = *supportsVision
	}
	if supportsVoice != nil {
		model.SupportsVoice = *supportsVoice
	}
	if isReasoning != nil {
		model.IsReasoning = *isReasoning
	}

	if err := u.AiModel.UpdateAiModel(model); err != nil {
		return nil, err
	}

	return model, nil
}

func (u *AiModelUsecase) DeleteAiModel(id string) error {
	if id == "" {
		return appErrors.NewValidationError("AI model ID is required")
	}
	return u.AiModel.DeleteAiModel(id)
}
