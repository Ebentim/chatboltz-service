package handler

import (
	"net/http"
	"strings"

	"github.com/alpinesboltltd/boltz-ai/internal/entity"
	appErrors "github.com/alpinesboltltd/boltz-ai/internal/errors"
	"github.com/alpinesboltltd/boltz-ai/internal/usecase"
	"github.com/gin-gonic/gin"
)

type GoogleHandler struct {
	agentUsecase *usecase.AgentUsecase
}

func NewGoogleHandler(agentUsecase *usecase.AgentUsecase) *GoogleHandler {
	return &GoogleHandler{
		agentUsecase: agentUsecase,
	}
}

// ConnectService enables a specific Google service for an agent
func (h *GoogleHandler) ConnectService(c *gin.Context) {
	agentId := c.Param("agentId")
	service := c.Param("service") // drive, calendar, mail

	// Normalize service name
	service = strings.ToLower(service)
	validServices := map[string]bool{
		"drive":     true,
		"calendar":  true,
		"mail":      true,
		"classroom": true,
		"slides":    true,
		"sheets":    true,
	}
	if !validServices[service] {
		appErrors.HandleError(c, appErrors.NewValidationError("Invalid service name"), "ConnectService")
		return
	}

	// 1. Get Integration
	integration, err := h.agentUsecase.GetAgentIntegration(agentId)
	if err != nil {
		// If not found, create one? GetAgentIntegration returns error if not found?
		// Usually we create empty one.
		// For now assume agent creation creates empty integration or handle error.
		// If NewNotFoundError, create it.
		// But let's assume it exists or use standard error handling.
		appErrors.HandleError(c, err, "ConnectService - GetIntegration")
		return
	}

	// 2. Add service if not present
	serviceId := "google_" + service
	exists := false
	for _, id := range integration.IntegrationId {
		if id == serviceId {
			exists = true
			break
		}
	}

	if !exists {
		newIds := append(integration.IntegrationId, serviceId)
		_, err := h.agentUsecase.UpdateAgentIntegration(agentId, entity.AgentIntegration{
			IntegrationId: newIds,
			IsActive:      true, // activate integration entry
		})
		if err != nil {
			appErrors.HandleError(c, err, "ConnectService - Update")
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Service connected", "service": service})
}

// DisconnectService disables a specific Google service
func (h *GoogleHandler) DisconnectService(c *gin.Context) {
	agentId := c.Param("agentId")
	service := c.Param("service")

	service = strings.ToLower(service)
	serviceId := "google_" + service

	integration, err := h.agentUsecase.GetAgentIntegration(agentId)
	if err != nil {
		appErrors.HandleError(c, err, "DisconnectService - GetIntegration")
		return
	}

	newIds := []string{}
	found := false
	for _, id := range integration.IntegrationId {
		if id == serviceId {
			found = true
			continue
		}
		newIds = append(newIds, id)
	}

	if found {
		_, err := h.agentUsecase.UpdateAgentIntegration(agentId, entity.AgentIntegration{
			IntegrationId: newIds,
		})
		if err != nil {
			appErrors.HandleError(c, err, "DisconnectService - Update")
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Service disconnected", "service": service})
}

// GetStatus returns the connection status of Google services
func (h *GoogleHandler) GetStatus(c *gin.Context) {
	agentId := c.Param("agentId")

	integration, err := h.agentUsecase.GetAgentIntegration(agentId)
	if err != nil {
		appErrors.HandleError(c, err, "GetStatus")
		return
	}

	status := map[string]bool{
		"drive":     false,
		"calendar":  false,
		"mail":      false,
		"classroom": false,
		"slides":    false,
		"sheets":    false,
	}

	for _, id := range integration.IntegrationId {
		if id == "google_drive" {
			status["drive"] = true
		}
		if id == "google_calendar" {
			status["calendar"] = true
		}
		if id == "google_mail" {
			status["mail"] = true
		}
		if id == "google_classroom" {
			status["classroom"] = true
		}
		if id == "google_slides" {
			status["slides"] = true
		}
		if id == "google_sheets" {
			status["sheets"] = true
		}
	}

	c.JSON(http.StatusOK, gin.H{"services": status})
}
