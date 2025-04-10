package main

import (
	"fmt"
	"net/http"
	"sync"
)

type SSEManager struct {
	clients map[string]chan string // clientId -> channel
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

func (s *SSEManager) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	clientID := r.URL.Query().Get("clientId")
	if clientID == "" {
		http.Error(w, "Missing clientId", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	messageChan := make(chan string)
	s.AddClient(clientID, messageChan)
	defer s.RemoveClient(clientID)

	notify := w.(http.CloseNotifier).CloseNotify()

	for {
		select {
		case msg := <-messageChan:
			fmt.Fprintf(w, "data: %s\n\n", msg)
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
		case <-notify:
			return
		}
	}
}
