package usecase

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// A websocket messeage
type WSMessage struct {
	Type      string                 `json:"type"`
	TimeStamp time.Time              `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

// AgentClient represents a connected websocket client
type AgentClient struct {
	ID   string
	Conn *websocket.Conn
	// Send is a channel that receives messages to be sent to the client
	Send chan *WSMessage
	// Close is a channel that receives a message when the client is closed
	Close  chan bool
	Hub    *AgentHub
	Mu     sync.RWMutex
	Ctx    context.Context
	Cancel context.CancelFunc
}

// AgentHub manages all connected clients and agent coordination
type AgentHub struct {
	Clients    map[*AgentClient]bool
	Broadcast  chan *WSMessage
	Register   chan *AgentClient
	Unregister chan *AgentClient
	Mu         sync.RWMutex
}

// Creates a new hub
func NewAgentHub() *AgentHub {
	return &AgentHub{
		Broadcast:  make(chan *WSMessage),
		Register:   make(chan *AgentClient),
		Unregister: make(chan *AgentClient),
		Clients:    make(map[*AgentClient]bool),
	}
}

// Run starts the hub's main loop
func (h *AgentHub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Mu.Lock()
			h.Clients[client] = true
			h.Mu.Unlock()
			log.Printf("Client registered: %s (TotalL %d) \n", client.ID, len(h.Clients))

		case client := <-h.Unregister:
			h.Mu.Lock()
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
			}
			h.Mu.Unlock()
			log.Printf("Client unregistered: %s (Total: %d) \n", client.ID, len(h.Clients))

		case msg := <-h.Broadcast:
			h.Mu.RLock()
			for client := range h.Clients {
				select {
				case client.Send <- msg:
				default:
					close(client.Send)
					delete(h.Clients, client)
				}
			}
			h.Mu.RUnlock()
		}

	}
}

func (c *AgentClient) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
		c.Cancel()
	}()

	c.Conn.SetReadDeadline(time.Now().Add(time.Second * 60))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(time.Second * 60))
		return nil
	})

	for {
		select {
		case <-c.Ctx.Done():
			return
		default:
			_, msgBytes, err := c.Conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Panicf("Websocket error: %v", err)
				}
				return
			}

			var msg WSMessage
			if err := json.Unmarshal(msgBytes, &msg); err != nil {
				log.Printf("JSON unmarshal error: %v", err)
				continue
			}

			c.handleMessage(msg)
		}
	}
}

func (c *AgentClient) handleMessage(msg WSMessage) {
	switch msg.Type {
	// handle different message type here
	}
}
