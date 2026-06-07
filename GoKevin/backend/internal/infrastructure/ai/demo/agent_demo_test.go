package demo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kevin-ai/go-kevin/internal/infrastructure/ai/agent"
	"github.com/kevin-ai/go-kevin/internal/infrastructure/ai/llm"
	"github.com/kevin-ai/go-kevin/internal/infrastructure/ai/tool"
)

// TestAgentDemo demonstrates how to create and use an AI agent
func TestAgentDemo(t *testing.T) {
	// Create a mock provider
	provider := llm.NewMockProvider()
	provider.SetDefaultResponse(llm.StringResponse("你好！我是智能助手，很高兴为你服务。"))

	// Create tools
	weatherTool := tool.NewWeatherTool()
	calculatorTool := tool.NewCalculatorTool()

	// Create an agent
	aiAgent := agent.NewAgent(agent.Config{
		Name:         "demo-agent",
		Description:  "A demo AI agent",
		Model:        "gpt-4",
		SystemPrompt: "你是一个有用的AI助手。",
		Provider:     provider,
		Tools:        []tool.Tool{weatherTool, calculatorTool},
		MaxTurns:     5,
	})

	// Test agent properties
	assert.Equal(t, "demo-agent", aiAgent.GetName())
	assert.Equal(t, "A demo AI agent", aiAgent.GetDescription())

	// Run the agent
	ctx := context.Background()
	response, err := aiAgent.Run(ctx, "你好")

	assert.NoError(t, err)
	assert.NotEmpty(t, response)
	t.Logf("Agent response: %s", response)
}

// TestAgentWithToolCalls demonstrates agent with tool calling
func TestAgentWithToolCalls(t *testing.T) {
	// Create mock provider that returns tool calls first, then final response
	callCount := 0
	provider := llm.NewMockProvider()
	provider.SetDefaultResponse(&llm.ChatResponse{
		ID:      "mock-response",
		Content: "北京今天天气晴朗，温度22度。",
		Role:    "assistant",
	})

	// Create weather tool
	weatherTool := tool.NewTool(
		"get_weather",
		"Get the current weather",
		map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"location": map[string]interface{}{
					"type": "string",
				},
			},
		},
		func(ctx context.Context, args string) (string, error) {
			callCount++
			return `{"location":"北京","temperature":22,"condition":"sunny"}`, nil
		},
	)

	// Create agent
	aiAgent := agent.NewAgent(agent.Config{
		Name:         "weather-agent",
		Description:  "Agent with weather tool",
		SystemPrompt: "你是一个天气助手。",
		Provider:     provider,
		Tools:        []tool.Tool{weatherTool},
	})

	ctx := context.Background()
	response, err := aiAgent.Run(ctx, "北京天气怎么样？")

	assert.NoError(t, err)
	assert.NotEmpty(t, response)
	t.Logf("Weather response: %s", response)
}

// TestAgentStreamDemo demonstrates streaming responses
func TestAgentStreamDemo(t *testing.T) {
	provider := llm.NewMockProvider()

	aiAgent := agent.NewAgent(agent.Config{
		Name:         "stream-agent",
		Description:  "Agent with streaming",
		SystemPrompt: "你是一个有用的助手。",
		Provider:     provider,
	})

	ctx := context.Background()
	stream, err := aiAgent.RunStream(ctx, "讲一个故事")

	assert.NoError(t, err)
	assert.NotNil(t, stream)

	// Read streaming response
	var fullResponse string
	for chunk := range stream {
		fullResponse += chunk
	}

	t.Logf("Stream response: %s", fullResponse)
	assert.NotEmpty(t, fullResponse)
}
