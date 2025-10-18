package handler

import (
	"fmt"
	"net/http"

	"github.com/alpinesboltltd/boltz-ai/internal/entity"
	appErrors "github.com/alpinesboltltd/boltz-ai/internal/errors"
	"github.com/alpinesboltltd/boltz-ai/internal/usecase"
	"github.com/gin-gonic/gin"
)

type AgentHandler struct {
	agentUsecase *usecase.AgentUsecase
}

func NewAgentHandler(agentUsecase *usecase.AgentUsecase) *AgentHandler {
	return &AgentHandler{
		agentUsecase: agentUsecase,
	}
}

func (h *AgentHandler) CreateAgent(c *gin.Context) {
	// Get user info from JWT token
	userRole := c.GetString("role")
	userID := c.GetString("userID")

	fmt.Println(userRole, userID, "User Role and User Id")

	// Check if user is admin or superadmin
	if userRole != string(entity.Admin) && userRole != string(entity.SuperAdmin) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admin and superadmin can create agents"})
		return
	}

	var req entity.Agent
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println(err)
		appErrors.HandleError(c, appErrors.NewValidationError("Invalid request format"), "CreateAgent - JSON binding")
		return
	}

	// Use the authenticated user's ID if not provided
	if req.UserId == "" {
		req.UserId = userID
	}

	agent, err := h.agentUsecase.CreateNewAgent(req.UserId, req.Name, req.Description, req.AiModel, req.AiProvider, req.AgentType, req.CreditsPer1k, req.Status)
	if err != nil {
		appErrors.HandleError(c, err, "CreateAgent")
		return
	}
	c.JSON(http.StatusCreated, gin.H{"agent": agent})
}

func (h *AgentHandler) UpdateAgent(c *gin.Context) {
	agentId := c.Param("agentId")
	var req entity.AgentUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		appErrors.HandleError(c, appErrors.NewValidationError("Invalid request format"), "UpdateAgent - JSON binding")
		return
	}

	agent, err := h.agentUsecase.UpdateAgent(agentId, req)
	if err != nil {
		appErrors.HandleError(c, err, "UpdateAgent")
		return
	}

	c.JSON(http.StatusOK, gin.H{"agent": agent})
}

func (h *AgentHandler) GetAgent(c *gin.Context) {
	agentId := c.Param("agentId")
	agent, err := h.agentUsecase.GetAgent(agentId)
	if err != nil {
		appErrors.HandleError(c, err, "GetAgent")
		return
	}
	c.JSON(http.StatusOK, gin.H{"agent": agent})
}

func (h *AgentHandler) GetAgentByUser(c *gin.Context) {
	// Get authenticated user's ID from JWT
	authUserID := c.GetString("userID")
	authUserRole := c.GetString("role")

	userId := c.Param("userId")

	// Only allow users to get their own agents unless they're admin/superadmin
	if userId != authUserID && authUserRole != string(entity.Admin) && authUserRole != string(entity.SuperAdmin) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	agents, err := h.agentUsecase.GetUserAgents(userId)
	if err != nil {
		appErrors.HandleError(c, err, "GetAgentByUser")
		return
	}

	// Convert to simplified response format
	var agentResponses []entity.AgentResponse
	for _, agent := range *agents {
		agentResponses = append(agentResponses, entity.AgentResponse{
			ID:           agent.ID,
			UserId:       agent.UserId,
			Name:         agent.Name,
			Description:  agent.Description,
			AgentType:    agent.AgentType,
			AiModel:      agent.AiModel,
			AiProvider:   agent.AiProvider,
			CreditsPer1k: agent.CreditsPer1k,
			Status:       agent.Status,
			CreatedAt:    agent.CreatedAt,
			UpdatedAt:    agent.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"agents": agentResponses})
}

func (h *AgentHandler) DeleteAgent(c *gin.Context) {
	agentId := c.Param("agentId")
	role := c.GetString("role")
	userId := c.GetString("userID")
	if role != string(entity.SuperAdmin) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only the super admin can delete agents"})
		return
	}

	err := h.agentUsecase.DeleteAgent(agentId, userId)
	if err != nil {
		appErrors.HandleError(c, err, "DeleteAgent")
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Agent deleted successfully"})
}

func (h *AgentHandler) CreateAgentAppearance(c *gin.Context) {
	var req entity.AgentAppearance
	if err := c.ShouldBindJSON(&req); err != nil {
		appErrors.HandleError(c, appErrors.NewValidationError("Invalid request format"), "CreateAgentAppearance - JSON binding")
		return
	}

	appearance, err := h.agentUsecase.CreateAgentAppearance(req)
	if err != nil {
		appErrors.HandleError(c, err, "CreateAgentAppearance")
		return
	}
	c.JSON(http.StatusCreated, gin.H{"appearance": appearance})
}

func (h *AgentHandler) CreateAgentBehavior(c *gin.Context) {
	var req entity.AgentBehavior
	if err := c.ShouldBindJSON(&req); err != nil {
		appErrors.HandleError(c, appErrors.NewValidationError("Invalid request format"), "CreateAgentBehavior - JSON binding")
		return
	}

	behavior, err := h.agentUsecase.CreateAgentBehavior(req)
	if err != nil {
		appErrors.HandleError(c, err, "CreateAgentBehavior")
		return
	}
	c.JSON(http.StatusCreated, gin.H{"behavior": behavior})
}

func (h *AgentHandler) CreateAgentChannel(c *gin.Context) {
	var req entity.AgentChannel
	if err := c.ShouldBindJSON(&req); err != nil {
		appErrors.HandleError(c, appErrors.NewValidationError("Invalid request format"), "CreateAgentChannel - JSON binding")
		return
	}

	channel, err := h.agentUsecase.CreateAgentChannel(req)
	if err != nil {
		appErrors.HandleError(c, err, "CreateAgentChannel")
		return
	}
	c.JSON(http.StatusCreated, gin.H{"channel": channel})
}

func (h *AgentHandler) CreateAgentIntegration(c *gin.Context) {
	var req entity.AgentIntegration
	if err := c.ShouldBindJSON(&req); err != nil {
		appErrors.HandleError(c, appErrors.NewValidationError("Invalid request format"), "CreateAgentIntegration - JSON binding")
		return
	}

	integration, err := h.agentUsecase.CreateAgentIntegration(req)
	if err != nil {
		appErrors.HandleError(c, err, "CreateAgentIntegration")
		return
	}
	c.JSON(http.StatusCreated, gin.H{"integration": integration})
}

func (h *AgentHandler) GetAgentAppearance(c *gin.Context) {
	agentId := c.Param("agentId")
	appearance, err := h.agentUsecase.GetAgentAppearance(agentId)
	if err != nil {
		appErrors.HandleError(c, err, "GetAgentAppearance")
		return
	}
	c.JSON(http.StatusOK, gin.H{"appearance": appearance})
}

func (h *AgentHandler) GetAgentBehavior(c *gin.Context) {
	agentId := c.Param("agentId")
	behavior, err := h.agentUsecase.GetAgentBehavior(agentId)
	if err != nil {
		appErrors.HandleError(c, err, "GetAgentBehavior")
		return
	}
	c.JSON(http.StatusOK, gin.H{"behavior": behavior})
}

func (h *AgentHandler) GetAgentChannel(c *gin.Context) {
	agentId := c.Param("agentId")
	channel, err := h.agentUsecase.GetAgentChannel(agentId)
	if err != nil {
		appErrors.HandleError(c, err, "GetAgentChannel")
		return
	}
	c.JSON(http.StatusOK, gin.H{"channel": channel})
}

func (h *AgentHandler) GetAgentStats(c *gin.Context) {
	agentId := c.Param("agentId")
	stats, err := h.agentUsecase.GetAgentStats(agentId)
	if err != nil {
		appErrors.HandleError(c, err, "GetAgentStats")
		return
	}
	c.JSON(http.StatusOK, gin.H{"stats": stats})
}

func (h *AgentHandler) GetAgentIntegration(c *gin.Context) {
	agentId := c.Param("agentId")
	integration, err := h.agentUsecase.GetAgentIntegration(agentId)
	if err != nil {
		appErrors.HandleError(c, err, "GetAgentIntegration")
		return
	}
	c.JSON(http.StatusOK, gin.H{"integration": integration})
}

func (h *AgentHandler) UpdateAgentAppearance(c *gin.Context) {
	agentId := c.Param("agentId")
	var req entity.AgentAppearance
	if err := c.ShouldBindJSON(&req); err != nil {
		appErrors.HandleError(c, appErrors.NewValidationError("Invalid request format"), "UpdateAgentAppearance - JSON binding")
		return
	}

	appearance, err := h.agentUsecase.UpdateAgentAppearance(agentId, req)
	if err != nil {
		appErrors.HandleError(c, err, "UpdateAgentAppearance")
		return
	}
	c.JSON(http.StatusOK, gin.H{"appearance": appearance})
}

func (h *AgentHandler) UpdateAgentBehavior(c *gin.Context) {
	agentId := c.Param("agentId")
	var req entity.AgentBehavior
	if err := c.ShouldBindJSON(&req); err != nil {
		appErrors.HandleError(c, appErrors.NewValidationError("Invalid request format"), "UpdateAgentBehavior - JSON binding")
		return
	}

	behavior, err := h.agentUsecase.UpdateAgentBehavior(agentId, req)
	if err != nil {
		appErrors.HandleError(c, err, "UpdateAgentBehavior")
		return
	}
	c.JSON(http.StatusOK, gin.H{"behavior": behavior})
}

func (h *AgentHandler) UpdateAgentChannel(c *gin.Context) {
	agentId := c.Param("agentId")
	var req entity.AgentChannel
	if err := c.ShouldBindJSON(&req); err != nil {
		appErrors.HandleError(c, appErrors.NewValidationError("Invalid request format"), "UpdateAgentChannel - JSON binding")
		return
	}

	channel, err := h.agentUsecase.UpdateAgentChannel(agentId, req)
	if err != nil {
		appErrors.HandleError(c, err, "UpdateAgentChannel")
		return
	}
	c.JSON(http.StatusOK, gin.H{"channel": channel})
}

func (h *AgentHandler) UpdateAgentIntegration(c *gin.Context) {
	agentId := c.Param("agentId")
	var req entity.AgentIntegration
	if err := c.ShouldBindJSON(&req); err != nil {
		appErrors.HandleError(c, appErrors.NewValidationError("Invalid request format"), "UpdateAgentIntegration - JSON binding")
		return
	}

	integration, err := h.agentUsecase.UpdateAgentIntegration(agentId, req)
	if err != nil {
		appErrors.HandleError(c, err, "UpdateAgentIntegration")
		return
	}
	c.JSON(http.StatusOK, gin.H{"integration": integration})
}

func (h *AgentHandler) DeleteAgentAppearance(c *gin.Context) {
	agentId := c.Param("agentId")
	err := h.agentUsecase.DeleteAgentAppearance(agentId)
	if err != nil {
		appErrors.HandleError(c, err, "DeleteAgentAppearance")
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Agent appearance deleted successfully"})
}

func (h *AgentHandler) DeleteAgentBehavior(c *gin.Context) {
	agentId := c.Param("agentId")
	err := h.agentUsecase.DeleteAgentBehavior(agentId)
	if err != nil {
		appErrors.HandleError(c, err, "DeleteAgentBehavior")
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Agent behavior deleted successfully"})
}

func (h *AgentHandler) DeleteAgentChannel(c *gin.Context) {
	agentId := c.Param("agentId")
	err := h.agentUsecase.DeleteAgentChannel(agentId)
	if err != nil {
		appErrors.HandleError(c, err, "DeleteAgentChannel")
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Agent channel deleted successfully"})
}

func (h *AgentHandler) DeleteAgentStats(c *gin.Context) {
	agentId := c.Param("agentId")
	err := h.agentUsecase.DeleteAgentStats(agentId)
	if err != nil {
		appErrors.HandleError(c, err, "DeleteAgentStats")
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Agent stats deleted successfully"})
}

func (h *AgentHandler) DeleteAgentIntegration(c *gin.Context) {
	agentId := c.Param("agentId")
	err := h.agentUsecase.DeleteAgentIntegration(agentId)
	if err != nil {
		appErrors.HandleError(c, err, "DeleteAgentIntegration")
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Agent integration deleted successfully"})
}
