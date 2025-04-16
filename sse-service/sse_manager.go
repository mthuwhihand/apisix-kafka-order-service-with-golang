package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type SSEManager struct {
	clients map[string]chan string
	mu      sync.Mutex
}

func NewSSEManager() *SSEManager {
	return &SSEManager{
		clients: make(map[string]chan string),
	}
}

func (s *SSEManager) AddClient(clientID string, ch chan string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.clients[clientID] = ch
}

func (s *SSEManager) RemoveClient(clientID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if ch, ok := s.clients[clientID]; ok {
		close(ch)
		delete(s.clients, clientID)
	}
}

func (s *SSEManager) SendToClient(clientID string, message string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if ch, ok := s.clients[clientID]; ok {
		ch <- message
	}
}

func (s *SSEManager) HandleOrderCreatedResponseSSE(c *gin.Context) {
	clientID := c.Query("clientId")
	if clientID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing clientId"})
		return
	}

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")

	messageChan := make(chan string)
	s.AddClient(clientID, messageChan)
	defer s.RemoveClient(clientID)

	notify := c.Request.Context().Done()

	for {
		select {
		case msg := <-messageChan:
			fmt.Fprintf(c.Writer, "data: %s\n\n", msg)
			c.Writer.Flush()
		case <-notify:
			return
		}
	}
}
