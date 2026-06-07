package rag

import (
	"context"
	"fmt"
	"strings"

	"github.com/kevin-ai/go-kevin/internal/infrastructure/ai/llm"
)

// Document represents a document in the knowledge base
type Document struct {
	ID      string `json:"id"`
	Content string `json:"content"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// SearchResult represents a search result
type SearchResult struct {
	Document Document `json:"document"`
	Score    float64  `json:"score"`
}

// VectorStore is the interface for vector storage
type VectorStore interface {
	// Store stores a document with its embedding
	Store(ctx context.Context, doc Document, embedding []float32) error

	// Search searches for similar documents
	Search(ctx context.Context, embedding []float32, topK int) ([]SearchResult, error)

	// Delete deletes a document
	Delete(ctx context.Context, docID string) error
}

// EmbeddingService is the interface for embedding services
type EmbeddingService interface {
	// Embed generates an embedding for the given text
	Embed(ctx context.Context, text string) ([]float32, error)

	// BatchEmbed generates embeddings for multiple texts
	BatchEmbed(ctx context.Context, texts []string) ([][]float32, error)
}

// RAGService provides RAG (Retrieval-Augmented Generation) functionality
type RAGService struct {
	vectorStore     VectorStore
	embeddingService EmbeddingService
	llmProvider     llm.Provider
}

// NewRAGService creates a new RAG service
func NewRAGService(vectorStore VectorStore, embeddingService EmbeddingService, llmProvider llm.Provider) *RAGService {
	return &RAGService{
		vectorStore:      vectorStore,
		embeddingService: embeddingService,
		llmProvider:      llmProvider,
	}
}

// AddDocument adds a document to the knowledge base
func (s *RAGService) AddDocument(ctx context.Context, doc Document) error {
	embedding, err := s.embeddingService.Embed(ctx, doc.Content)
	if err != nil {
		return fmt.Errorf("generate embedding: %w", err)
	}

	if err := s.vectorStore.Store(ctx, doc, embedding); err != nil {
		return fmt.Errorf("store document: %w", err)
	}

	return nil
}

// Query queries the knowledge base and generates an answer
func (s *RAGService) Query(ctx context.Context, question string, topK int) (string, error) {
	// Generate embedding for the question
	queryEmbedding, err := s.embeddingService.Embed(ctx, question)
	if err != nil {
		return "", fmt.Errorf("generate query embedding: %w", err)
	}

	// Search for similar documents
	results, err := s.vectorStore.Search(ctx, queryEmbedding, topK)
	if err != nil {
		return "", fmt.Errorf("search documents: %w", err)
	}

	if len(results) == 0 {
		return "没有找到相关的文档。", nil
	}

	// Build context from search results
	var contextParts []string
	for i, result := range results {
		contextParts = append(contextParts, fmt.Sprintf("[%d] %s", i+1, result.Document.Content))
	}
	context := strings.Join(contextParts, "\n\n")

	// Generate answer using LLM
	prompt := fmt.Sprintf(`基于以下上下文信息回答问题。如果上下文中没有相关信息，请说明无法回答。

上下文：
%s

问题：%s

请用中文回答：`, context, question)

	req := &llm.ChatRequest{
		Messages: []llm.Message{
			{Role: "user", Content: prompt},
		},
	}

	resp, err := s.llmProvider.Chat(ctx, req)
	if err != nil {
		return "", fmt.Errorf("generate answer: %w", err)
	}

	return resp.Content, nil
}

// QueryWithContext queries with additional context
func (s *RAGService) QueryWithContext(ctx context.Context, question string, topK int, systemPrompt string) (string, error) {
	queryEmbedding, err := s.embeddingService.Embed(ctx, question)
	if err != nil {
		return "", fmt.Errorf("generate query embedding: %w", err)
	}

	results, err := s.vectorStore.Search(ctx, queryEmbedding, topK)
	if err != nil {
		return "", fmt.Errorf("search documents: %w", err)
	}

	var contextParts []string
	for i, result := range results {
		contextParts = append(contextParts, fmt.Sprintf("[%d] %s", i+1, result.Document.Content))
	}
	context := strings.Join(contextParts, "\n\n")

	messages := []llm.Message{}
	if systemPrompt != "" {
		messages = append(messages, llm.Message{
			Role:    "system",
			Content: systemPrompt,
		})
	}

	userPrompt := fmt.Sprintf(`基于以下上下文信息回答问题。

上下文：
%s

问题：%s`, context, question)

	messages = append(messages, llm.Message{
		Role:    "user",
		Content: userPrompt,
	})

	req := &llm.ChatRequest{
		Messages: messages,
	}

	resp, err := s.llmProvider.Chat(ctx, req)
	if err != nil {
		return "", fmt.Errorf("generate answer: %w", err)
	}

	return resp.Content, nil
}
