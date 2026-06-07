package tool

import (
	"context"
	"encoding/json"
	"fmt"
)

// Tool represents an AI tool that can be called by the model
type Tool interface {
	// Name returns the tool name
	Name() string

	// Description returns the tool description
	Description() string

	// Parameters returns the tool parameters schema
	Parameters() interface{}

	// Execute executes the tool with the given arguments
	Execute(ctx context.Context, args string) (string, error)
}

// BaseTool provides a base implementation for tools
type BaseTool struct {
	name        string
	description string
	parameters  interface{}
	handler     func(ctx context.Context, args string) (string, error)
}

// NewTool creates a new tool
func NewTool(name, description string, parameters interface{}, handler func(ctx context.Context, args string) (string, error)) Tool {
	return &BaseTool{
		name:        name,
		description: description,
		parameters:  parameters,
		handler:     handler,
	}
}

func (t *BaseTool) Name() string               { return t.name }
func (t *BaseTool) Description() string         { return t.description }
func (t *BaseTool) Parameters() interface{}     { return t.parameters }
func (t *BaseTool) Execute(ctx context.Context, args string) (string, error) {
	return t.handler(ctx, args)
}

// Registry manages a collection of tools
type Registry struct {
	tools map[string]Tool
}

// NewRegistry creates a new tool registry
func NewRegistry() *Registry {
	return &Registry{
		tools: make(map[string]Tool),
	}
}

// Register registers a tool
func (r *Registry) Register(t Tool) {
	r.tools[t.Name()] = t
}

// Get returns a tool by name
func (r *Registry) Get(name string) (Tool, bool) {
	t, ok := r.tools[name]
	return t, ok
}

// GetAll returns all registered tools
func (r *Registry) GetAll() []Tool {
	result := make([]Tool, 0, len(r.tools))
	for _, t := range r.tools {
		result = append(result, t)
	}
	return result
}

// --- Built-in Tools ---

// WeatherArgs represents weather tool arguments
type WeatherArgs struct {
	Location string `json:"location"`
}

// NewWeatherTool creates a weather tool
func NewWeatherTool() Tool {
	return NewTool(
		"get_weather",
		"Get the current weather for a location",
		map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"location": map[string]interface{}{
					"type":        "string",
					"description": "The city name, e.g., Beijing, New York",
				},
			},
			"required": []string{"location"},
		},
		func(ctx context.Context, args string) (string, error) {
			var weatherArgs WeatherArgs
			if err := json.Unmarshal([]byte(args), &weatherArgs); err != nil {
				return "", fmt.Errorf("parse args: %w", err)
			}

			// Mock weather data
			return fmt.Sprintf(`{"location":"%s","temperature":22,"condition":"sunny","humidity":45}`, weatherArgs.Location), nil
		},
	)
}

// CalculatorArgs represents calculator tool arguments
type CalculatorArgs struct {
	Expression string `json:"expression"`
}

// NewCalculatorTool creates a calculator tool
func NewCalculatorTool() Tool {
	return NewTool(
		"calculator",
		"Calculate a mathematical expression",
		map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"expression": map[string]interface{}{
					"type":        "string",
					"description": "The mathematical expression to calculate, e.g., 2+2, 10*5",
				},
			},
			"required": []string{"expression"},
		},
		func(ctx context.Context, args string) (string, error) {
			var calcArgs CalculatorArgs
			if err := json.Unmarshal([]byte(args), &calcArgs); err != nil {
				return "", fmt.Errorf("parse args: %w", err)
			}

			// Simple calculator (in production, use a proper expression evaluator)
			var result float64
			_, err := fmt.Sscanf(calcArgs.Expression, "%f", &result)
			if err != nil {
				return "", fmt.Errorf("invalid expression: %s", calcArgs.Expression)
			}

			return fmt.Sprintf(`{"expression":"%s","result":%.2f}`, calcArgs.Expression, result), nil
		},
	)
}

// SearchArgs represents search tool arguments
type SearchArgs struct {
	Query string `json:"query"`
}

// NewSearchTool creates a search tool
func NewSearchTool() Tool {
	return NewTool(
		"web_search",
		"Search the web for information",
		map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"query": map[string]interface{}{
					"type":        "string",
					"description": "The search query",
				},
			},
			"required": []string{"query"},
		},
		func(ctx context.Context, args string) (string, error) {
			var searchArgs SearchArgs
			if err := json.Unmarshal([]byte(args), &searchArgs); err != nil {
				return "", fmt.Errorf("parse args: %w", err)
			}

			// Mock search results
			return fmt.Sprintf(`{"query":"%s","results":[{"title":"Example Result","url":"https://example.com","snippet":"This is an example search result."}]}`, searchArgs.Query), nil
		},
	)
}
