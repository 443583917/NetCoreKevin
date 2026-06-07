package demo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kevin-ai/go-kevin/internal/infrastructure/ai/embedding"
	"github.com/kevin-ai/go-kevin/internal/infrastructure/ai/llm"
	"github.com/kevin-ai/go-kevin/internal/infrastructure/ai/vectorstore"
)

// TestQdrantClient_New tests Qdrant client creation
func TestQdrantClient_New(t *testing.T) {
	client := vectorstore.QdrantClient{}
	assert.NotNil(t, client)
}

// TestQdrantConfig tests Qdrant configuration
func TestQdrantConfig(t *testing.T) {
	cfg := vectorstore.QdrantConfig{
		Host:       "localhost",
		Port:       6333,
		Collection: "test_collection",
	}

	assert.Equal(t, "localhost", cfg.Host)
	assert.Equal(t, 6333, cfg.Port)
	assert.Equal(t, "test_collection", cfg.Collection)
}

// TestPoint_Struct tests Point struct
func TestPoint_Struct(t *testing.T) {
	point := vectorstore.Point{
		ID:     "doc1",
		Vector: []float32{0.1, 0.2, 0.3},
		Payload: map[string]interface{}{
			"content": "Test content",
		},
	}

	assert.Equal(t, "doc1", point.ID)
	assert.Len(t, point.Vector, 3)
	assert.Equal(t, "Test content", point.Payload["content"])
}

// TestSearchResult_Struct tests SearchResult struct
func TestSearchResult_Struct(t *testing.T) {
	result := vectorstore.SearchResult{
		ID:    "doc1",
		Score: 0.95,
		Payload: map[string]interface{}{
			"content": "Test content",
		},
	}

	assert.Equal(t, "doc1", result.ID)
	assert.Equal(t, 0.95, result.Score)
}

// TestMockEmbedding tests mock embedding service
func TestMockEmbedding(t *testing.T) {
	mock := embedding.NewMockEmbedding(128)

	ctx := context.Background()
	embedding, err := mock.Embed(ctx, "test text")

	assert.NoError(t, err)
	assert.Len(t, embedding, 128)
}

// TestMockEmbedding_Batch tests batch embedding
func TestMockEmbedding_Batch(t *testing.T) {
	mock := embedding.NewMockEmbedding(64)

	ctx := context.Background()
	texts := []string{"text1", "text2", "text3"}
	embeddings, err := mock.BatchEmbed(ctx, texts)

	assert.NoError(t, err)
	assert.Len(t, embeddings, 3)
	assert.Len(t, embeddings[0], 64)
}

// TestOpenAIEmbedding_New tests OpenAI embedding creation
func TestOpenAIEmbedding_New(t *testing.T) {
	emb := embedding.NewOpenAIEmbedding(embedding.OpenAIEmbeddingConfig{
		APIKey: "test-key",
		Model:  "text-embedding-ada-002",
	})

	assert.NotNil(t, emb)
}

// TestEmbeddingInterface tests embedding interface compliance
func TestEmbeddingInterface(t *testing.T) {
	var _ embedding.EmbeddingService = (*embedding.OpenAIEmbedding)(nil)
	var _ embedding.EmbeddingService = (*embedding.MockEmbedding)(nil)
}

// TestQdrantRAGWorkflow tests the complete RAG workflow
func TestQdrantRAGWorkflow(t *testing.T) {
	// This test demonstrates the RAG workflow without actual Qdrant connection

	// 1. Create mock embedding service
	mockEmbedding := embedding.NewMockEmbedding(128)

	// 2. Create mock LLM provider
	mockLLM := llm.NewMockProvider()
	mockLLM.SetDefaultResponse(llm.StringResponse("Go是一种编程语言。"))

	// 3. Test embedding generation
	ctx := context.Background()
	vector, err := mockEmbedding.Embed(ctx, "Go语言入门教程")
	assert.NoError(t, err)
	assert.Len(t, vector, 128)

	// 4. Test LLM response
	req := &llm.ChatRequest{
		Messages: []llm.Message{
			{Role: "user", Content: "什么是Go语言？"},
		},
	}
	resp, err := mockLLM.Chat(ctx, req)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.Content)

	t.Logf("Embedding dimension: %d", len(vector))
	t.Logf("LLM response: %s", resp.Content)
}
