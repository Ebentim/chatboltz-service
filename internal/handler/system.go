package handler

import (
	"net/http"

	"github.com/alpinesboltltd/boltz-ai/internal/entity"
	appErrors "github.com/alpinesboltltd/boltz-ai/internal/errors"
	"github.com/alpinesboltltd/boltz-ai/internal/usecase"
	"github.com/gin-gonic/gin"
)

type SystemHandler struct {
	systemUsecase *usecase.SystemUsecase
}

func NewSystemHandler(systemUsecase *usecase.SystemUsecase) *SystemHandler {
	return &SystemHandler{
		systemUsecase: systemUsecase,
	}
}

func (h *SystemHandler) CreateSystemInstruction(c *gin.Context) {
	userRole := c.GetString("role")
	userID := c.GetString("userID")

	if userRole != string(entity.SuperAdmin) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only superadmin can create system instructions"})
		return
	}

	var req struct {
		Title      string  `json:"title" binding:"required"`
		Content    string  `json:"content" binding:"required"`
		TemplateId *string `json:"template_id,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		appErrors.HandleError(c, appErrors.NewValidationError("Invalid request format"), "CreateSystemInstruction - JSON binding")
		return
	}

	instruction, err := h.systemUsecase.CreateSystemInstruction(req.Title, req.Content, userID, req.TemplateId)
	if err != nil {
		appErrors.HandleError(c, err, "CreateSystemInstruction")
		return
	}

	c.JSON(http.StatusCreated, gin.H{"instruction": instruction})
}

func (h *SystemHandler) GetSystemInstruction(c *gin.Context) {
	id := c.Param("id")
	instruction, err := h.systemUsecase.GetSystemInstruction(id)
	if err != nil {
		appErrors.HandleError(c, err, "GetSystemInstruction")
		return
	}
	c.JSON(http.StatusOK, gin.H{"instruction": instruction})
}

func (h *SystemHandler) UpdateSystemInstruction(c *gin.Context) {
	userRole := c.GetString("role")
	if userRole != string(entity.SuperAdmin) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only superadmin can update system instructions"})
		return
	}

	id := c.Param("id")
	var req struct {
		Title   string `json:"title,omitempty"`
		Content string `json:"content,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		appErrors.HandleError(c, appErrors.NewValidationError("Invalid request format"), "UpdateSystemInstruction - JSON binding")
		return
	}

	instruction, err := h.systemUsecase.UpdateSystemInstruction(id, req.Title, req.Content)
	if err != nil {
		appErrors.HandleError(c, err, "UpdateSystemInstruction")
		return
	}

	c.JSON(http.StatusOK, gin.H{"instruction": instruction})
}

func (h *SystemHandler) DeleteSystemInstruction(c *gin.Context) {
	userRole := c.GetString("role")
	if userRole != string(entity.SuperAdmin) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only superadmin can delete system instructions"})
		return
	}

	id := c.Param("id")
	err := h.systemUsecase.DeleteSystemInstruction(id)
	if err != nil {
		appErrors.HandleError(c, err, "DeleteSystemInstruction")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "System instruction deleted successfully"})
}

func (h *SystemHandler) ListSystemInstructions(c *gin.Context) {
	instructions, err := h.systemUsecase.ListSystemInstructions()
	if err != nil {
		appErrors.HandleError(c, err, "ListSystemInstructions")
		return
	}
	c.JSON(http.StatusOK, gin.H{"instructions": instructions})
}

func (h *SystemHandler) CreatePromptTemplate(c *gin.Context) {
	userRole := c.GetString("role")
	if userRole != string(entity.SuperAdmin) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only superadmin can create prompt templates"})
		return
	}

	var req struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
		Role    string `json:"role"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		appErrors.HandleError(c, appErrors.NewValidationError("Invalid request format"), "CreatePromptTemplate - JSON binding")
		return
	}

	template, err := h.systemUsecase.CreatePromptTemplate(req.Title, req.Content, req.Role)
	if err != nil {
		appErrors.HandleError(c, err, "CreatePromptTemplate")
		return
	}

	c.JSON(http.StatusCreated, gin.H{"template": template})
}

func (h *SystemHandler) GetPromptTemplate(c *gin.Context) {
	id := c.Param("id")
	template, err := h.systemUsecase.GetPromptTemplate(id)
	if err != nil {
		appErrors.HandleError(c, err, "GetPromptTemplate")
		return
	}
	c.JSON(http.StatusOK, gin.H{"template": template})
}

func (h *SystemHandler) ListPromptTemplates(c *gin.Context) {
	role := c.Query("role")
	templates, err := h.systemUsecase.ListPromptTemplates(role)
	if err != nil {
		appErrors.HandleError(c, err, "ListPromptTemplates")
		return
	}
	c.JSON(http.StatusOK, gin.H{"templates": templates})
}

func (h *SystemHandler) UpdatePromptTemplate(c *gin.Context) {
	userRole := c.GetString("role")
	if userRole != string(entity.SuperAdmin) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only superadmin can update prompt templates"})
		return
	}

	id := c.Param("id")
	var req struct {
		Title   string `json:"title,omitempty"`
		Content string `json:"content,omitempty"`
		Role    string `json:"role,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		appErrors.HandleError(c, appErrors.NewValidationError("Invalid request format"), "UpdatePromptTemplate - JSON binding")
		return
	}

	template, err := h.systemUsecase.UpdatePromptTemplate(id, req.Title, req.Content, req.Role)
	if err != nil {
		appErrors.HandleError(c, err, "UpdatePromptTemplate")
		return
	}

	c.JSON(http.StatusOK, gin.H{"template": template})
}

func (h *SystemHandler) DeletePromptTemplate(c *gin.Context) {
	userRole := c.GetString("role")
	if userRole != string(entity.SuperAdmin) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only superadmin can delete prompt templates"})
		return
	}

	id := c.Param("id")
	err := h.systemUsecase.DeletePromptTemplate(id)
	if err != nil {
		appErrors.HandleError(c, err, "DeletePromptTemplate")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Prompt template deleted successfully"})
}
