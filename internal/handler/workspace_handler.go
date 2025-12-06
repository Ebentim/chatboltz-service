package handler

import (
	"net/http"

	"github.com/alpinesboltltd/boltz-ai/internal/usecase"
	"github.com/gin-gonic/gin"
)

type WorkspaceHandler struct {
	usecase usecase.WorkspaceUsecase
}

func NewWorkspaceHandler(usecase usecase.WorkspaceUsecase) *WorkspaceHandler {
	return &WorkspaceHandler{usecase: usecase}
}

func (h *WorkspaceHandler) CreateWorkspace(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString("userId") // Assuming AuthMiddleware sets this
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	workspace, err := h.usecase.CreateWorkspace(req.Name, req.Description, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, workspace)
}

func (h *WorkspaceHandler) GetUserWorkspaces(c *gin.Context) {
	userID := c.GetString("userId")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	workspaces, err := h.usecase.GetUserWorkspaces(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, workspaces)
}

func (h *WorkspaceHandler) GetWorkspace(c *gin.Context) {
	id := c.Param("id")
	workspace, err := h.usecase.GetWorkspace(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "workspace not found"})
		return
	}
	c.JSON(http.StatusOK, workspace)
}
