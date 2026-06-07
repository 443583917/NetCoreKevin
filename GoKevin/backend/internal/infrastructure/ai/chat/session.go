package chat

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/kevin-ai/go-kevin/internal/infrastructure/ai/llm"
)

// Message represents a chat message
type Message struct {
	ID        string    `json:"id"`
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

var messageCounter int64

// Session represents a chat session
type Session struct {
	ID        string    `json:"id"`
	UserID    int64     `json:"userId"`
	AppID     int64     `json:"appId"`
	Title     string    `json:"title"`
	Messages  []Message `json:"messages"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// SessionManager manages chat sessions
type SessionManager struct {
	sessions map[string]*Session
	provider llm.Provider
	mu       sync.RWMutex
	counter  int64
}

// NewSessionManager creates a new session manager
func NewSessionManager(provider llm.Provider) *SessionManager {
	return &SessionManager{
		sessions: make(map[string]*Session),
		provider: provider,
	}
}

// CreateSession creates a new chat session
func (m *SessionManager) CreateSession(userID, appID int64, title string) *Session {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.counter++
	session := &Session{
		ID:        fmt.Sprintf("session_%d_%d_%d", userID, time.Now().UnixNano(), m.counter),
		UserID:    userID,
		AppID:     appID,
		Title:     title,
		Messages:  []Message{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	m.sessions[session.ID] = session
	return session
}

// GetSession returns a session by ID
func (m *SessionManager) GetSession(sessionID string) (*Session, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	session, ok := m.sessions[sessionID]
	return session, ok
}

// ListSessions returns all sessions for a user
func (m *SessionManager) ListSessions(userID int64) []*Session {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var sessions []*Session
	for _, session := range m.sessions {
		if session.UserID == userID {
			sessions = append(sessions, session)
		}
	}
	return sessions
}

// SendMessage sends a message to a session and gets a response
func (m *SessionManager) SendMessage(ctx context.Context, sessionID, content string) (*Message, error) {
	m.mu.Lock()
	session, ok := m.sessions[sessionID]
	if !ok {
		m.mu.Unlock()
		return nil, fmt.Errorf("session not found: %s", sessionID)
	}

	// Add user message
	msgID := atomic.AddInt64(&messageCounter, 1)
	userMsg := Message{
		ID:        fmt.Sprintf("msg_%d_%d", time.Now().UnixNano(), msgID),
		Role:      "user",
		Content:   content,
		Timestamp: time.Now(),
	}
	session.Messages = append(session.Messages, userMsg)
	m.mu.Unlock()

	// Build messages for LLM
	var messages []llm.Message
	for _, msg := range session.Messages {
		messages = append(messages, llm.Message{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	// Get response from LLM
	req := &llm.ChatRequest{
		Messages: messages,
	}

	resp, err := m.provider.Chat(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("get response: %w", err)
	}

	// Add assistant message
	m.mu.Lock()
	msgID = atomic.AddInt64(&messageCounter, 1)
	assistantMsg := Message{
		ID:        fmt.Sprintf("msg_%d_%d", time.Now().UnixNano(), msgID),
		Role:      "assistant",
		Content:   resp.Content,
		Timestamp: time.Now(),
	}
	session.Messages = append(session.Messages, assistantMsg)
	session.UpdatedAt = time.Now()
	m.mu.Unlock()

	return &assistantMsg, nil
}

// DeleteSession deletes a session
func (m *SessionManager) DeleteSession(sessionID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.sessions[sessionID]; !ok {
		return fmt.Errorf("session not found: %s", sessionID)
	}

	delete(m.sessions, sessionID)
	return nil
}

// GetHistory returns the message history of a session
func (m *SessionManager) GetHistory(sessionID string) ([]Message, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	session, ok := m.sessions[sessionID]
	if !ok {
		return nil, fmt.Errorf("session not found: %s", sessionID)
	}

	return session.Messages, nil
}
