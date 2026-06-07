package agent

import (
	"context"
	"fmt"

	"github.com/kevin-ai/go-kevin/internal/infrastructure/ai/llm"
	"github.com/kevin-ai/go-kevin/internal/infrastructure/ai/tool"
)

// Agent represents an AI agent
type Agent struct {
	name         string
	description  string
	model        string
	systemPrompt string
	provider     llm.Provider
	tools        []tool.Tool
	maxTurns     int
}

// Config represents agent configuration
type Config struct {
	Name         string
	Description  string
	Model        string
	SystemPrompt string
	Provider     llm.Provider
	Tools        []tool.Tool
	MaxTurns     int
}

// NewAgent creates a new AI agent
func NewAgent(cfg Config) *Agent {
	maxTurns := cfg.MaxTurns
	if maxTurns <= 0 {
		maxTurns = 10
	}

	return &Agent{
		name:         cfg.Name,
		description:  cfg.Description,
		model:        cfg.Model,
		systemPrompt: cfg.SystemPrompt,
		provider:     cfg.Provider,
		tools:        cfg.Tools,
		maxTurns:     maxTurns,
	}
}

// Run executes the agent with a user message
func (a *Agent) Run(ctx context.Context, message string) (string, error) {
	messages := []llm.Message{}

	// Add system prompt if configured
	if a.systemPrompt != "" {
		messages = append(messages, llm.Message{
			Role:    "system",
			Content: a.systemPrompt,
		})
	}

	// Add user message
	messages = append(messages, llm.Message{
		Role:    "user",
		Content: message,
	})

	// Convert tools to LLM format
	llmTools := a.convertTools()

	// Run conversation loop
	for turn := 0; turn < a.maxTurns; turn++ {
		req := &llm.ChatRequest{
			Model:    a.model,
			Messages: messages,
			Tools:    llmTools,
		}

		resp, err := a.provider.Chat(ctx, req)
		if err != nil {
			return "", fmt.Errorf("chat error: %w", err)
		}

		// If no tool calls, return the response
		if len(resp.ToolCalls) == 0 {
			return resp.Content, nil
		}

		// Add assistant message with tool calls
		messages = append(messages, llm.Message{
			Role:      "assistant",
			Content:   resp.Content,
			ToolCalls: resp.ToolCalls,
		})

		// Execute tool calls
		for _, tc := range resp.ToolCalls {
			result, err := a.executeTool(ctx, tc)
			if err != nil {
				result = fmt.Sprintf("Error: %v", err)
			}

			messages = append(messages, llm.Message{
				Role:       "tool",
				Content:    result,
				ToolCallID: tc.ID,
			})
		}
	}

	return "", fmt.Errorf("max turns (%d) exceeded", a.maxTurns)
}

// RunStream executes the agent with streaming response
func (a *Agent) RunStream(ctx context.Context, message string) (<-chan string, error) {
	messages := []llm.Message{}

	if a.systemPrompt != "" {
		messages = append(messages, llm.Message{
			Role:    "system",
			Content: a.systemPrompt,
		})
	}

	messages = append(messages, llm.Message{
		Role:    "user",
		Content: message,
	})

	req := &llm.ChatRequest{
		Model:    a.model,
		Messages: messages,
		Tools:    a.convertTools(),
		Stream:   true,
	}

	chunks, err := a.provider.StreamChat(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("stream chat error: %w", err)
	}

	output := make(chan string, 100)
	go func() {
		defer close(output)
		for chunk := range chunks {
			if chunk.Done {
				return
			}
			if chunk.Content != "" {
				output <- chunk.Content
			}
		}
	}()

	return output, nil
}

// GetName returns the agent name
func (a *Agent) GetName() string {
	return a.name
}

// GetDescription returns the agent description
func (a *Agent) GetDescription() string {
	return a.description
}

// convertTools converts tool.Tool to llm.Tool
func (a *Agent) convertTools() []llm.Tool {
	result := make([]llm.Tool, len(a.tools))
	for i, t := range a.tools {
		result[i] = llm.Tool{
			Type: "function",
			Function: llm.ToolFunction{
				Name:        t.Name(),
				Description: t.Description(),
				Parameters:  t.Parameters(),
			},
		}
	}
	return result
}

// executeTool executes a tool call
func (a *Agent) executeTool(ctx context.Context, tc llm.ToolCall) (string, error) {
	for _, t := range a.tools {
		if t.Name() == tc.Function.Name {
			return t.Execute(ctx, tc.Function.Arguments)
		}
	}
	return "", fmt.Errorf("tool not found: %s", tc.Function.Name)
}
