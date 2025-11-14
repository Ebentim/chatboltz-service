package handler

import (
	"net/http"

	"github.com/alpinesboltltd/boltz-ai/internal/entity"
	appErrors "github.com/alpinesboltltd/boltz-ai/internal/errors"
	"github.com/alpinesboltltd/boltz-ai/internal/usecase"
	"github.com/gin-gonic/gin"
)

type AiModelHandler struct {
	aiModelUsecase *usecase.AiModelUsecase
}

func NewAiModelHandler(aiModelUsecase *usecase.AiModelUsecase) *AiModelHandler {
	return &AiModelHandler{
		aiModelUsecase: aiModelUsecase,
	}
}

func (h *AiModelHandler) CreateAiModel(c *gin.Context) {
	userRole := c.GetString("role")

	if userRole != string(entity.SuperAdmin) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only superadmin can create AI models"})
		return
	}

	var req entity.AiModel
	if err := c.ShouldBindJSON(&req); err != nil {
		appErrors.HandleError(c, appErrors.NewValidationError("Invalid request format"), "CreateAiModel - JSON binding")
		return
	}

	model, err := h.aiModelUsecase.CreateAiModel(req.Name, req.Provider, req.CreditsPer1k, req.SupportsText, req.SupportsVision, req.SupportsVoice, req.IsReasoning)
	if err != nil {
		appErrors.HandleError(c, err, "CreateAiModel")
		return
	}

	c.JSON(http.StatusCreated, gin.H{"ai_model": model})
}

func (h *AiModelHandler) GetAiModel(c *gin.Context) {
	modelId := c.Param("modelId")
	model, err := h.aiModelUsecase.GetAiModel(modelId)
	if err != nil {
		appErrors.HandleError(c, err, "GetAiModel")
		return
	}
	c.JSON(http.StatusOK, gin.H{"ai_model": model})
}

func (h *AiModelHandler) ListAiModels(c *gin.Context) {
	provider := c.Query("provider")

	var models *[]entity.AiModel
	var err error

	if provider != "" {
		models, err = h.aiModelUsecase.ListAiModelsByProvider(provider)
	} else {
		models, err = h.aiModelUsecase.ListAiModels()
	}

	if err != nil {
		appErrors.HandleError(c, err, "ListAiModels")
		return
	}

	c.JSON(http.StatusOK, gin.H{"ai_models": models})
}

func (h *AiModelHandler) UpdateAiModel(c *gin.Context) {
	userRole := c.GetString("role")

	if userRole != string(entity.SuperAdmin) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only superadmin can update AI models"})
		return
	}

	modelId := c.Param("modelId")

	var req struct {
		Name           *string `json:"name,omitempty"`
		Provider       *string `json:"provider,omitempty"`
		CreditsPer1k   *int    `json:"credits_per_1k,omitempty"`
		SupportsText   *bool   `json:"supports_text,omitempty"`
		SupportsVision *bool   `json:"supports_vision,omitempty"`
		SupportsVoice  *bool   `json:"supports_voice,omitempty"`
		IsReasoning    *bool   `json:"is_reasoning,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		appErrors.HandleError(c, appErrors.NewValidationError("Invalid request format"), "UpdateAiModel - JSON binding")
		return
	}

	model, err := h.aiModelUsecase.UpdateAiModel(modelId, req.Name, req.Provider, req.CreditsPer1k, req.SupportsText, req.SupportsVision, req.SupportsVoice, req.IsReasoning)
	if err != nil {
		appErrors.HandleError(c, err, "UpdateAiModel")
		return
	}

	c.JSON(http.StatusOK, gin.H{"ai_model": model})
}

func (h *AiModelHandler) DeleteAiModel(c *gin.Context) {
	userRole := c.GetString("role")

	if userRole != string(entity.SuperAdmin) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only superadmin can delete AI models"})
		return
	}

	modelId := c.Param("modelId")
	err := h.aiModelUsecase.DeleteAiModel(modelId)
	if err != nil {
		appErrors.HandleError(c, err, "DeleteAiModel")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "AI model deleted successfully"})
}
