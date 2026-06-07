package websocket

import (
	"log"
	"net/http"
	"sync"
	"time"

	ws "github.com/gorilla/websocket"
)

var upgrader = ws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development
	},
}

// Client represents a WebSocket client
type Client struct {
	ID     string
	UserID int64
	Conn   *ws.Conn
	Hub    *Hub
	Send   chan []byte
	mu     sync.Mutex
}

// NewClient creates a new WebSocket client
func NewClient(id string, userID int64, conn *ws.Conn, hub *Hub) *Client {
	return &Client{
		ID:     id,
		UserID: userID,
		Conn:   conn,
		Hub:    hub,
		Send:   make(chan []byte, 256),
	}
}

// ReadPump reads messages from the WebSocket connection
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(512)
	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if ws.IsUnexpectedCloseError(err, ws.CloseGoingAway, ws.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}
		c.Hub.Broadcast <- message
	}
}

// WritePump writes messages to the WebSocket connection
func (c *Client) WritePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.Conn.WriteMessage(ws.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(ws.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(ws.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// SendMessage sends a message to the client
func (c *Client) SendMessage(message []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	select {
	case c.Send <- message:
	default:
		// Channel full, close connection
		close(c.Send)
	}
}

// Upgrade upgrades an HTTP connection to WebSocket
var Upgrade = upgrader.Upgrade
