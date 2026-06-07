package demo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kevin-ai/go-kevin/internal/infrastructure/ai/tool"
)

// TestWeatherToolDemo demonstrates the weather tool
func TestWeatherToolDemo(t *testing.T) {
	weatherTool := tool.NewWeatherTool()

	assert.Equal(t, "get_weather", weatherTool.Name())
	assert.NotEmpty(t, weatherTool.Description())
	assert.NotNil(t, weatherTool.Parameters())

	// Execute the tool
	ctx := context.Background()
	result, err := weatherTool.Execute(ctx, `{"location":"北京"}`)

	assert.NoError(t, err)
	assert.Contains(t, result, "北京")
	assert.Contains(t, result, "temperature")
	t.Logf("Weather result: %s", result)
}

// TestCalculatorToolDemo demonstrates the calculator tool
func TestCalculatorToolDemo(t *testing.T) {
	calcTool := tool.NewCalculatorTool()

	assert.Equal(t, "calculator", calcTool.Name())

	ctx := context.Background()
	result, err := calcTool.Execute(ctx, `{"expression":"42"}`)

	assert.NoError(t, err)
	assert.Contains(t, result, "42")
	t.Logf("Calculator result: %s", result)
}

// TestSearchToolDemo demonstrates the search tool
func TestSearchToolDemo(t *testing.T) {
	searchTool := tool.NewSearchTool()

	assert.Equal(t, "web_search", searchTool.Name())

	ctx := context.Background()
	result, err := searchTool.Execute(ctx, `{"query":"Go语言教程"}`)

	assert.NoError(t, err)
	assert.Contains(t, result, "Go语言教程")
	t.Logf("Search result: %s", result)
}

// TestCustomToolDemo demonstrates creating a custom tool
func TestCustomToolDemo(t *testing.T) {
	// Create a custom tool
	translateTool := tool.NewTool(
		"translate",
		"Translate text to another language",
		map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"text": map[string]interface{}{
					"type":        "string",
					"description": "The text to translate",
				},
				"target_lang": map[string]interface{}{
					"type":        "string",
					"description": "The target language",
				},
			},
			"required": []string{"text", "target_lang"},
		},
		func(ctx context.Context, args string) (string, error) {
			// Mock translation
			return `{"original":"Hello","translated":"你好","target_lang":"zh"}`, nil
		},
	)

	assert.Equal(t, "translate", translateTool.Name())

	ctx := context.Background()
	result, err := translateTool.Execute(ctx, `{"text":"Hello","target_lang":"zh"}`)

	assert.NoError(t, err)
	assert.Contains(t, result, "你好")
	t.Logf("Translate result: %s", result)
}

// TestToolRegistryDemo demonstrates the tool registry
func TestToolRegistryDemo(t *testing.T) {
	registry := tool.NewRegistry()

	// Register tools
	registry.Register(tool.NewWeatherTool())
	registry.Register(tool.NewCalculatorTool())
	registry.Register(tool.NewSearchTool())

	// Get all tools
	allTools := registry.GetAll()
	assert.Len(t, allTools, 3)

	// Get specific tool
	weatherTool, ok := registry.Get("get_weather")
	assert.True(t, ok)
	assert.NotNil(t, weatherTool)

	// Get non-existent tool
	_, ok = registry.Get("non_existent")
	assert.False(t, ok)
}
