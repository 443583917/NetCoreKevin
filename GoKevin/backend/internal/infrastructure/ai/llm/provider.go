package llm

import (
	"context"
)

// Message represents a chat message
type Message struct {
	Role       string      `json:"role"` // system, user, assistant, tool
	Content    string      `json:"content"`
	ToolCalls  []ToolCall  `json:"tool_calls,omitempty"`
	ToolCallID string      `json:"tool_call_id,omitempty"`
	Name       string      `json:"name,omitempty"`
}

// ToolCall represents a tool call from the model
type ToolCall struct {
	ID       string       `json:"id"`
	Type     string       `json:"type"` // function
	Function FunctionCall `json:"function"`
}

// FunctionCall represents a function call
type FunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

// Tool represents a tool definition
type Tool struct {
	Type     string       `json:"type"` // function
	Function ToolFunction `json:"function"`
}

// ToolFunction represents a function tool
type ToolFunction struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Parameters  interface{} `json:"parameters"`
}

// ChatRequest represents a chat completion request
type ChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Tools       []Tool    `json:"tools,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Stream      bool      `json:"stream,omitempty"`
}

// ChatResponse represents a chat completion response
type ChatResponse struct {
	ID      string   `json:"id"`
	Content string   `json:"content"`
	Role    string   `json:"role"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
	Usage   Usage    `json:"usage"`
}

// Usage represents token usage
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// ChatChunk represents a streaming chat chunk
type ChatChunk struct {
	ID      string `json:"id"`
	Content string `json:"content"`
	Done    bool   `json:"done"`
}

// Provider is the interface for LLM providers
type Provider interface {
	// Chat sends a chat completion request
	Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error)

	// StreamChat sends a streaming chat completion request
	StreamChat(ctx context.Context, req *ChatRequest) (<-chan *ChatChunk, error)

	// GetModel returns the model name
	GetModel() string
}
