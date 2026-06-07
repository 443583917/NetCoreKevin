package demo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kevin-ai/go-kevin/internal/infrastructure/ai/chat"
	"github.com/kevin-ai/go-kevin/internal/infrastructure/ai/llm"
)

// TestChatSessionDemo demonstrates chat session management
func TestChatSessionDemo(t *testing.T) {
	// Create mock provider
	provider := llm.NewMockProvider()
	provider.SetDefaultResponse(llm.StringResponse("你好！有什么可以帮助你的吗？"))

	// Create session manager
	manager := chat.NewSessionManager(provider)

	// Create a session
	session := manager.CreateSession(1001, 1, "测试对话")

	assert.NotEmpty(t, session.ID)
	assert.Equal(t, int64(1001), session.UserID)
	assert.Equal(t, int64(1), session.AppID)
	assert.Equal(t, "测试对话", session.Title)
	assert.Empty(t, session.Messages)

	t.Logf("Created session: %s", session.ID)

	// Send a message
	ctx := context.Background()
	response, err := manager.SendMessage(ctx, session.ID, "你好")

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "assistant", response.Role)
	assert.NotEmpty(t, response.Content)

	t.Logf("Response: %s", response.Content)

	// Get history
	history, err := manager.GetHistory(session.ID)

	assert.NoError(t, err)
	assert.Len(t, history, 2) // user message + assistant response
	assert.Equal(t, "user", history[0].Role)
	assert.Equal(t, "你好", history[0].Content)
	assert.Equal(t, "assistant", history[1].Role)

	t.Logf("History length: %d", len(history))
}

// TestMultipleSessionsDemo demonstrates multiple sessions
func TestMultipleSessionsDemo(t *testing.T) {
	provider := llm.NewMockProvider()
	provider.SetDefaultResponse(llm.StringResponse("收到"))

	manager := chat.NewSessionManager(provider)

	// Create multiple sessions
	session1 := manager.CreateSession(1001, 1, "对话1")
	session2 := manager.CreateSession(1001, 2, "对话2")
	session3 := manager.CreateSession(1002, 1, "对话3")

	// List sessions for user 1001
	sessions := manager.ListSessions(1001)
	assert.Len(t, sessions, 2)

	// List sessions for user 1002
	sessions = manager.ListSessions(1002)
	assert.Len(t, sessions, 1)

	t.Logf("User 1001 has %d sessions", 2)
	t.Logf("User 1002 has %d sessions", 1)

	// Send messages to different sessions
	ctx := context.Background()

	_, err := manager.SendMessage(ctx, session1.ID, "消息1")
	assert.NoError(t, err)

	_, err = manager.SendMessage(ctx, session2.ID, "消息2")
	assert.NoError(t, err)

	_, err = manager.SendMessage(ctx, session3.ID, "消息3")
	assert.NoError(t, err)

	// Verify history
	history1, _ := manager.GetHistory(session1.ID)
	history2, _ := manager.GetHistory(session2.ID)
	history3, _ := manager.GetHistory(session3.ID)

	assert.Len(t, history1, 2)
	assert.Len(t, history2, 2)
	assert.Len(t, history3, 2)
}

// TestSessionDeleteDemo demonstrates session deletion
func TestSessionDeleteDemo(t *testing.T) {
	provider := llm.NewMockProvider()
	manager := chat.NewSessionManager(provider)

	// Create and delete session
	session := manager.CreateSession(1001, 1, "临时对话")

	_, ok := manager.GetSession(session.ID)
	assert.True(t, ok)

	// Delete session
	err := manager.DeleteSession(session.ID)
	assert.NoError(t, err)

	_, ok = manager.GetSession(session.ID)
	assert.False(t, ok)

	// Try to delete non-existent session
	err = manager.DeleteSession("non-existent")
	assert.Error(t, err)
}

// TestSessionNotFoundDemo demonstrates error handling
func TestSessionNotFoundDemo(t *testing.T) {
	provider := llm.NewMockProvider()
	manager := chat.NewSessionManager(provider)

	ctx := context.Background()

	// Try to send message to non-existent session
	_, err := manager.SendMessage(ctx, "non-existent", "hello")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "session not found")

	// Try to get history of non-existent session
	_, err = manager.GetHistory("non-existent")
	assert.Error(t, err)
}
