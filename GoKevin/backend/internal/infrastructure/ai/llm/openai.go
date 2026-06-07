package llm

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	defaultBaseURL = "https://api.openai.com/v1"
	defaultModel   = "gpt-4"
)

// OpenAIProvider implements the Provider interface for OpenAI
type OpenAIProvider struct {
	apiKey  string
	baseURL string
	model   string
	client  *http.Client
}

// OpenAIConfig represents OpenAI provider configuration
type OpenAIConfig struct {
	APIKey  string
	BaseURL string
	Model   string
}

// NewOpenAIProvider creates a new OpenAI provider
func NewOpenAIProvider(cfg OpenAIConfig) *OpenAIProvider {
	baseURL := cfg.BaseURL
	if baseURL == "" {
		baseURL = defaultBaseURL
	}

	model := cfg.Model
	if model == "" {
		model = defaultModel
	}

	return &OpenAIProvider{
		apiKey:  cfg.APIKey,
		baseURL: baseURL,
		model:   model,
		client:  &http.Client{},
	}
}

// GetModel returns the model name
func (p *OpenAIProvider) GetModel() string {
	return p.model
}

// Chat sends a chat completion request
func (p *OpenAIProvider) Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	if req.Model == "" {
		req.Model = p.model
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", p.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	var result openAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	if len(result.Choices) == 0 {
		return nil, fmt.Errorf("no choices in response")
	}

	choice := result.Choices[0]
	return &ChatResponse{
		ID:      result.ID,
		Content: choice.Message.Content,
		Role:    choice.Message.Role,
		ToolCalls: convertToolCalls(choice.Message.ToolCalls),
		Usage: Usage{
			PromptTokens:     result.Usage.PromptTokens,
			CompletionTokens: result.Usage.CompletionTokens,
			TotalTokens:      result.Usage.TotalTokens,
		},
	}, nil
}

// StreamChat sends a streaming chat completion request
func (p *OpenAIProvider) StreamChat(ctx context.Context, req *ChatRequest) (<-chan *ChatChunk, error) {
	req.Stream = true
	if req.Model == "" {
		req.Model = p.model
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", p.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	chunks := make(chan *ChatChunk, 100)
	go func() {
		defer close(chunks)
		defer resp.Body.Close()

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := scanner.Text()
			if !strings.HasPrefix(line, "data: ") {
				continue
			}

			data := strings.TrimPrefix(line, "data: ")
			if data == "[DONE]" {
				chunks <- &ChatChunk{Done: true}
				return
			}

			var chunk openAIStreamChunk
			if err := json.Unmarshal([]byte(data), &chunk); err != nil {
				continue
			}

			if len(chunk.Choices) > 0 {
				chunks <- &ChatChunk{
					ID:      chunk.ID,
					Content: chunk.Choices[0].Delta.Content,
					Done:    false,
				}
			}
		}
	}()

	return chunks, nil
}

// OpenAI API response types
type openAIResponse struct {
	ID      string `json:"id"`
	Choices []struct {
		Message struct {
			Role      string          `json:"role"`
			Content   string          `json:"content"`
			ToolCalls []openAIToolCall `json:"tool_calls,omitempty"`
		} `json:"message"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

type openAIToolCall struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Function struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"`
	} `json:"function"`
}

type openAIStreamChunk struct {
	ID      string `json:"id"`
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
	} `json:"choices"`
}

func convertToolCalls(calls []openAIToolCall) []ToolCall {
	result := make([]ToolCall, len(calls))
	for i, call := range calls {
		result[i] = ToolCall{
			ID:   call.ID,
			Type: call.Type,
			Function: FunctionCall{
				Name:      call.Function.Name,
				Arguments: call.Function.Arguments,
			},
		}
	}
	return result
}
