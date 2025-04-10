package main

import (
	"fmt"
	"net/http"
	"sync"
)

type SSEManager struct {
	clients map[chan string]bool
	mu      sync.Mutex
}

func NewSSEManager() *SSEManager {
	return &SSEManager{
		clients: make(map[chan string]bool),
	}
}

func (s *SSEManager) AddClient(ch chan string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.clients[ch] = true
}

func (s *SSEManager) RemoveClient(ch chan string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.clients, ch)
	close(ch)
}

func (s *SSEManager) Broadcast(message string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for ch := range s.clients {
		ch <- message
	}
}

func (s *SSEManager) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Headers để giữ kết nối SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	messageChan := make(chan string)
	s.AddClient(messageChan)
	defer s.RemoveClient(messageChan)

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
