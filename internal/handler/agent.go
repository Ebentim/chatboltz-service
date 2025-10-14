package handler

import (
	"net/http"

	"github.com/alpinesboltltd/boltz-ai/internal/entity"
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
	var req entity.Agent
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	agent, err := h.agentUsecase.CreateNewAgent(req.UserId, req.Name, req.Description, req.AiModel, req.AgentType, req.Capabilities, req.CreditsPer1k, req.Status)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"agent": agent})
}

func (h *AgentHandler) UpdateAgent(c *gin.Context) {
	agentId := c.Param("agentId")
	var req entity.AgentUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	agent, err := h.agentUsecase.UpdateAgent(agentId, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"agent": agent})
}

func (h *AgentHandler) GetAgent(c *gin.Context) {
	agentId := c.Param("agentId")
	agent, err := h.agentUsecase.GetAgent(agentId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"agent": agent})
}

func (h *AgentHandler) GetAgentByUser(c *gin.Context) {
	userId := c.Param("userId")
	agents, err := h.agentUsecase.GetUserAgents(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"agents": agents})
}
