package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/alpinesboltltd/boltz-ai/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for MVP
	},
	HandshakeTimeout: time.Duration(time.Second * 30),
}

func (h *AgentHandler) HandleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer conn.Close()

	for {
		// Read message
		var msg struct {
			Message string `json:"message"`
			AgentID string `json:"agent_id"`
			APIKey  string `json:"api_key"`
		}
		
		if err := conn.ReadJSON(&msg); err != nil {
			break
		}

		// Process message
		response, err := h.chatService.ProcessMessage(msg.AgentID, msg.Message, msg.APIKey)
		if err != nil {
			conn.WriteJSON(gin.H{"error": err.Error()})
			continue
		}

		// Send response
		if err := conn.WriteJSON(gin.H{
			"type": "response",
			"content": response,
			"agent_id": msg.AgentID,
		}); err != nil {
			break
		}
	}
}
