package middleware

import (
	"net/http"

	"github.com/alpinesboltltd/boltz-ai/internal/usecase"
	"github.com/gin-gonic/gin"
)

func WorkspaceAuthMiddleware(workspaceUsecase usecase.WorkspaceUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		workspaceID := c.GetHeader("X-Workspace-ID")
		if workspaceID == "" {
			// Try query param for certain GET requests if header is missing (optional but helpful)
			workspaceID = c.Query("workspaceID")
		}

		if workspaceID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "X-Workspace-ID header is required"})
			c.Abort()
			return
		}

		userID := c.GetString("userID")
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
			c.Abort()
			return
		}

		role, err := workspaceUsecase.GetMemberRole(workspaceID, userID)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "User is not a member of this workspace"})
			c.Abort()
			return
		}

		// Set Context
		c.Set("workspaceID", workspaceID)
		c.Set("workspaceRole", role)

		c.Next()
	}
}
