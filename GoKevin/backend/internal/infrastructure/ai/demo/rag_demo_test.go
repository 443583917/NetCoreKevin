package demo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kevin-ai/go-kevin/internal/infrastructure/ai/llm"
	"github.com/kevin-ai/go-kevin/internal/infrastructure/ai/rag"
)

// MockVectorStore is a mock vector store for testing
type MockVectorStore struct {
	documents map[string]rag.Document
	embeddings map[string][]float32
}

func NewMockVectorStore() *MockVectorStore {
	return &MockVectorStore{
		documents: make(map[string]rag.Document),
		embeddings: make(map[string][]float32),
	}
}

func (s *MockVectorStore) Store(ctx context.Context, doc rag.Document, embedding []float32) error {
	s.documents[doc.ID] = doc
	s.embeddings[doc.ID] = embedding
	return nil
}

func (s *MockVectorStore) Search(ctx context.Context, embedding []float32, topK int) ([]rag.SearchResult, error) {
	// Simple mock: return all documents
	var results []rag.SearchResult
	for _, doc := range s.documents {
		results = append(results, rag.SearchResult{
			Document: doc,
			Score:    0.95,
		})
	}
	return results, nil
}

func (s *MockVectorStore) Delete(ctx context.Context, docID string) error {
	delete(s.documents, docID)
	delete(s.embeddings, docID)
	return nil
}

// MockEmbeddingService is a mock embedding service for testing
type MockEmbeddingService struct{}

func (s *MockEmbeddingService) Embed(ctx context.Context, text string) ([]float32, error) {
	// Simple mock: return fixed-size embedding
	return make([]float32, 128), nil
}

func (s *MockEmbeddingService) BatchEmbed(ctx context.Context, texts []string) ([][]float32, error) {
	embeddings := make([][]float32, len(texts))
	for i := range texts {
		embeddings[i] = make([]float32, 128)
	}
	return embeddings, nil
}

// TestRAGServiceDemo demonstrates the RAG service
func TestRAGServiceDemo(t *testing.T) {
	// Create mock dependencies
	vectorStore := NewMockVectorStore()
	embeddingService := &MockEmbeddingService{}
	llmProvider := llm.NewMockProvider()
	llmProvider.SetDefaultResponse(llm.StringResponse("Go是一种开源编程语言，由Google开发。"))

	// Create RAG service
	ragService := rag.NewRAGService(vectorStore, embeddingService, llmProvider)

	ctx := context.Background()

	// Add documents to knowledge base
	documents := []rag.Document{
		{
			ID:      "doc1",
			Content: "Go是一种开源编程语言，由Google开发。",
			Metadata: map[string]interface{}{
				"source": "wikipedia",
			},
		},
		{
			ID:      "doc2",
			Content: "Go语言具有简洁、高效、并发等特点。",
			Metadata: map[string]interface{}{
				"source": "tutorial",
			},
		},
		{
			ID:      "doc3",
			Content: "Go广泛应用于云原生、微服务、容器等领域。",
			Metadata: map[string]interface{}{
				"source": "blog",
			},
		},
	}

	for _, doc := range documents {
		err := ragService.AddDocument(ctx, doc)
		assert.NoError(t, err)
	}

	// Query the knowledge base
	answer, err := ragService.Query(ctx, "Go语言是什么？", 3)

	assert.NoError(t, err)
	assert.NotEmpty(t, answer)
	t.Logf("RAG answer: %s", answer)
}

// TestRAGWithContextDemo demonstrates RAG with custom context
func TestRAGWithContextDemo(t *testing.T) {
	vectorStore := NewMockVectorStore()
	embeddingService := &MockEmbeddingService{}
	llmProvider := llm.NewMockProvider()
	llmProvider.SetDefaultResponse(llm.StringResponse("Go语言适合用于构建高性能的后端服务。"))

	ragService := rag.NewRAGService(vectorStore, embeddingService, llmProvider)

	ctx := context.Background()

	// Add a document
	err := ragService.AddDocument(ctx, rag.Document{
		ID:      "doc1",
		Content: "Go语言适合构建高性能后端服务。",
	})
	assert.NoError(t, err)

	// Query with custom system prompt
	answer, err := ragService.QueryWithContext(
		ctx,
		"Go语言适合做什么？",
		3,
		"你是一个编程语言专家，请用专业的方式回答问题。",
	)

	assert.NoError(t, err)
	assert.NotEmpty(t, answer)
	t.Logf("RAG answer with context: %s", answer)
}
