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
		origin := r.Header.Get("Origin")
		if origin != "http://"+r.Host && origin != "https://"+r.Host {
			return false
		}
		return true
	},
	HandshakeTimeout: time.Duration(time.Second * 30),
}

func WebSocketHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer conn.Close()
	// ctx, cancel := context.WithCancel(context.Background())
	_, cancel := context.WithCancel(context.Background())
	client := &usecase.AgentClient{
		ID: fmt.Sprintf("boltz-agent"),
	}

	for {
		MessageType, msg, err := conn.ReadMessage()
		if err != nil {
			cancel()

			break
		}
		if err := conn.WriteMessage(MessageType, msg); err != nil {
			break
		}
		client.ReadPump()

	}
}
