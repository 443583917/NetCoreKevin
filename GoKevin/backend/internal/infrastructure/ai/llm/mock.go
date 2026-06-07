package llm

import (
	"context"
	"fmt"
)

// MockProvider is a mock LLM provider for testing
type MockProvider struct {
	responses map[string]*ChatResponse
	defaultResponse *ChatResponse
}

// NewMockProvider creates a new mock provider
func NewMockProvider() *MockProvider {
	return &MockProvider{
		responses: make(map[string]*ChatResponse),
		defaultResponse: &ChatResponse{
			ID:      "mock-response",
			Content: "This is a mock response.",
			Role:    "assistant",
			Usage: Usage{
				PromptTokens:     10,
				CompletionTokens: 20,
				TotalTokens:      30,
			},
		},
	}
}

// SetResponse sets a mock response for a specific input
func (m *MockProvider) SetResponse(input string, response *ChatResponse) {
	m.responses[input] = response
}

// SetDefaultResponse sets the default mock response
func (m *MockProvider) SetDefaultResponse(response *ChatResponse) {
	m.defaultResponse = response
}

// Chat sends a chat completion request
func (m *MockProvider) Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	if len(req.Messages) > 0 {
		lastMsg := req.Messages[len(req.Messages)-1]
		if resp, ok := m.responses[lastMsg.Content]; ok {
			return resp, nil
		}
	}
	return m.defaultResponse, nil
}

// StreamChat sends a streaming chat completion request
func (m *MockProvider) StreamChat(ctx context.Context, req *ChatRequest) (<-chan *ChatChunk, error) {
	chunks := make(chan *ChatChunk, 100)

	go func() {
		defer close(chunks)

		content := "This is a streaming mock response."
		words := splitWords(content)

		for _, word := range words {
			chunks <- &ChatChunk{
				ID:      "mock-chunk",
				Content: word + " ",
				Done:    false,
			}
		}

		chunks <- &ChatChunk{Done: true}
	}()

	return chunks, nil
}

// GetModel returns the model name
func (m *MockProvider) GetModel() string {
	return "mock-model"
}

func splitWords(s string) []string {
	words := []string{}
	word := ""
	for _, c := range s {
		if c == ' ' {
			if word != "" {
				words = append(words, word)
				word = ""
			}
		} else {
			word += string(c)
		}
	}
	if word != "" {
		words = append(words, word)
	}
	return words
}

// MockProviderWithToolCalls creates a mock provider that returns tool calls
func MockProviderWithToolCalls(toolCalls []ToolCall) *MockProvider {
	provider := NewMockProvider()
	provider.defaultResponse = &ChatResponse{
		ID:        "mock-tool-response",
		Content:   "",
		Role:      "assistant",
		ToolCalls: toolCalls,
		Usage: Usage{
			PromptTokens:     10,
			CompletionTokens: 5,
			TotalTokens:      15,
		},
	}
	return provider
}

// MockProviderWithError creates a mock provider that returns an error
func MockProviderWithError(err error) *ErrorMockProvider {
	return &ErrorMockProvider{err: err}
}

// ErrorMockProvider is a mock provider that always returns an error
type ErrorMockProvider struct {
	err error
}

func (m *ErrorMockProvider) Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	return nil, m.err
}

func (m *ErrorMockProvider) StreamChat(ctx context.Context, req *ChatRequest) (<-chan *ChatChunk, error) {
	return nil, m.err
}

func (m *ErrorMockProvider) GetModel() string {
	return "error-model"
}

// StringResponse creates a simple mock response with text content
func StringResponse(content string) *ChatResponse {
	return &ChatResponse{
		ID:      fmt.Sprintf("mock-%d", len(content)),
		Content: content,
		Role:    "assistant",
		Usage: Usage{
			PromptTokens:     10,
			CompletionTokens: len(content) / 4,
			TotalTokens:      10 + len(content)/4,
		},
	}
}
