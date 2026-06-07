package llm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpenAIProvider_GetModel(t *testing.T) {
	provider := NewOpenAIProvider(OpenAIConfig{
		APIKey: "test-key",
		Model:  "gpt-4",
	})

	assert.Equal(t, "gpt-4", provider.GetModel())
}

func TestOpenAIProvider_DefaultModel(t *testing.T) {
	provider := NewOpenAIProvider(OpenAIConfig{
		APIKey: "test-key",
	})

	assert.Equal(t, defaultModel, provider.GetModel())
}

func TestOpenAIProvider_DefaultBaseURL(t *testing.T) {
	provider := NewOpenAIProvider(OpenAIConfig{
		APIKey: "test-key",
	})

	assert.Equal(t, defaultBaseURL, provider.baseURL)
}

func TestMessage_Struct(t *testing.T) {
	msg := Message{
		Role:    "user",
		Content: "Hello, world!",
	}

	assert.Equal(t, "user", msg.Role)
	assert.Equal(t, "Hello, world!", msg.Content)
}

func TestChatRequest_Struct(t *testing.T) {
	req := ChatRequest{
		Model: "gpt-4",
		Messages: []Message{
			{Role: "system", Content: "You are a helpful assistant."},
			{Role: "user", Content: "Hello!"},
		},
		Temperature: 0.7,
		MaxTokens:   1000,
	}

	assert.Equal(t, "gpt-4", req.Model)
	assert.Len(t, req.Messages, 2)
	assert.Equal(t, 0.7, req.Temperature)
}

func TestTool_Struct(t *testing.T) {
	tool := Tool{
		Type: "function",
		Function: ToolFunction{
			Name:        "get_weather",
			Description: "Get the current weather",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"location": map[string]interface{}{
						"type":        "string",
						"description": "The city name",
					},
				},
			},
		},
	}

	assert.Equal(t, "function", tool.Type)
	assert.Equal(t, "get_weather", tool.Function.Name)
}

func TestConvertToolCalls(t *testing.T) {
	calls := []openAIToolCall{
		{
			ID:   "call_123",
			Type: "function",
			Function: struct {
				Name      string `json:"name"`
				Arguments string `json:"arguments"`
			}{
				Name:      "get_weather",
				Arguments: `{"location":"Beijing"}`,
			},
		},
	}

	result := convertToolCalls(calls)

	assert.Len(t, result, 1)
	assert.Equal(t, "call_123", result[0].ID)
	assert.Equal(t, "get_weather", result[0].Function.Name)
	assert.Equal(t, `{"location":"Beijing"}`, result[0].Function.Arguments)
}
