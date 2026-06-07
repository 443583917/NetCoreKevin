package embedding

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// EmbeddingService is the interface for embedding services
type EmbeddingService interface {
	// Embed generates an embedding for the given text
	Embed(ctx context.Context, text string) ([]float32, error)

	// BatchEmbed generates embeddings for multiple texts
	BatchEmbed(ctx context.Context, texts []string) ([][]float32, error)
}

// OpenAIEmbedding implements EmbeddingService using OpenAI API
type OpenAIEmbedding struct {
	apiKey  string
	baseURL string
	model   string
	client  *http.Client
}

// OpenAIEmbeddingConfig represents OpenAI embedding configuration
type OpenAIEmbeddingConfig struct {
	APIKey  string
	BaseURL string
	Model   string
}

// NewOpenAIEmbedding creates a new OpenAI embedding service
func NewOpenAIEmbedding(cfg OpenAIEmbeddingConfig) *OpenAIEmbedding {
	baseURL := cfg.BaseURL
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}

	model := cfg.Model
	if model == "" {
		model = "text-embedding-ada-002"
	}

	return &OpenAIEmbedding{
		apiKey:  cfg.APIKey,
		baseURL: baseURL,
		model:   model,
		client:  &http.Client{},
	}
}

// Embed generates an embedding for the given text
func (e *OpenAIEmbedding) Embed(ctx context.Context, text string) ([]float32, error) {
	embeddings, err := e.BatchEmbed(ctx, []string{text})
	if err != nil {
		return nil, err
	}
	if len(embeddings) == 0 {
		return nil, fmt.Errorf("no embedding returned")
	}
	return embeddings[0], nil
}

// BatchEmbed generates embeddings for multiple texts
func (e *OpenAIEmbedding) BatchEmbed(ctx context.Context, texts []string) ([][]float32, error) {
	body := map[string]interface{}{
		"input": texts,
		"model": e.model,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", e.baseURL+"/embeddings", bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+e.apiKey)

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	var result openAIEmbeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	embeddings := make([][]float32, len(result.Data))
	for i, d := range result.Data {
		embeddings[i] = d.Embedding
	}

	return embeddings, nil
}

type openAIEmbeddingResponse struct {
	Data []struct {
		Embedding []float32 `json:"embedding"`
	} `json:"data"`
}

// MockEmbedding is a mock embedding service for testing
type MockEmbedding struct {
	dimension int
}

// NewMockEmbedding creates a new mock embedding service
func NewMockEmbedding(dimension int) *MockEmbedding {
	return &MockEmbedding{dimension: dimension}
}

// Embed generates a mock embedding
func (m *MockEmbedding) Embed(ctx context.Context, text string) ([]float32, error) {
	embedding := make([]float32, m.dimension)
	for i := range embedding {
		embedding[i] = float32(i) / float32(m.dimension)
	}
	return embedding, nil
}

// BatchEmbed generates mock embeddings
func (m *MockEmbedding) BatchEmbed(ctx context.Context, texts []string) ([][]float32, error) {
	embeddings := make([][]float32, len(texts))
	for i := range texts {
		embedding, err := m.Embed(ctx, texts[i])
		if err != nil {
			return nil, err
		}
		embeddings[i] = embedding
	}
	return embeddings, nil
}
