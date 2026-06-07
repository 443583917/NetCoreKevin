package websocket

import (
	"encoding/json"
	"log"
	"sync"
)

// Hub maintains the set of active clients
type Hub struct {
	clients    map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
	mu         sync.RWMutex
}

// Message represents a WebSocket message
type Message struct {
	Type    string      `json:"type"`
	From    string      `json:"from,omitempty"`
	To      string      `json:"to,omitempty"`
	Content interface{} `json:"content"`
}

// NewHub creates a new Hub
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

// Run starts the hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Printf("Client connected: %s", client.ID)

		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
			}
			h.mu.Unlock()
			log.Printf("Client disconnected: %s", client.ID)

		case message := <-h.Broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// SendToUser sends a message to a specific user
func (h *Hub) SendToUser(userID int64, message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients {
		if client.UserID == userID {
			client.SendMessage(message)
		}
	}
}

// SendToClient sends a message to a specific client
func (h *Hub) SendToClient(clientID string, message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients {
		if client.ID == clientID {
			client.SendMessage(message)
			break
		}
	}
}

// BroadcastMessage broadcasts a message to all clients
func (h *Hub) BroadcastMessage(msgType string, content interface{}) {
	msg := Message{
		Type:    msgType,
		Content: content,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Marshal message error: %v", err)
		return
	}

	h.Broadcast <- data
}

// GetClients returns all connected clients
func (h *Hub) GetClients() []*Client {
	h.mu.RLock()
	defer h.mu.RUnlock()

	clients := make([]*Client, 0, len(h.clients))
	for client := range h.clients {
		clients = append(clients, client)
	}
	return clients
}

// GetClientCount returns the number of connected clients
func (h *Hub) GetClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}
