package rag

import (
	"context"
	"fmt"
	"strings"

	"github.com/kevin-ai/go-kevin/internal/infrastructure/ai/embedding"
	"github.com/kevin-ai/go-kevin/internal/infrastructure/ai/llm"
	"github.com/kevin-ai/go-kevin/internal/infrastructure/ai/vectorstore"
)

// QdrantRAGService provides RAG functionality using Qdrant
type QdrantRAGService struct {
	qdrant   *vectorstore.QdrantClient
	embedding embedding.EmbeddingService
	llm      llm.Provider
}

// NewQdrantRAGService creates a new Qdrant RAG service
func NewQdrantRAGService(
	qdrant *vectorstore.QdrantClient,
	embeddingService embedding.EmbeddingService,
	llmProvider llm.Provider,
) *QdrantRAGService {
	return &QdrantRAGService{
		qdrant:    qdrant,
		embedding: embeddingService,
		llm:       llmProvider,
	}
}

// InitializeCollection creates the collection if it doesn't exist
func (s *QdrantRAGService) InitializeCollection(ctx context.Context, vectorSize int) error {
	return s.qdrant.CreateCollection(ctx, vectorSize)
}

// AddDocument adds a document to the knowledge base
func (s *QdrantRAGService) AddDocument(ctx context.Context, docID string, content string, metadata map[string]interface{}) error {
	// Generate embedding
	vector, err := s.embedding.Embed(ctx, content)
	if err != nil {
		return fmt.Errorf("generate embedding: %w", err)
	}

	// Prepare payload
	payload := map[string]interface{}{
		"content": content,
	}
	for k, v := range metadata {
		payload[k] = v
	}

	// Upsert to Qdrant
	point := vectorstore.Point{
		ID:      docID,
		Vector:  vector,
		Payload: payload,
	}

	return s.qdrant.UpsertPoints(ctx, []vectorstore.Point{point})
}

// Query queries the knowledge base and generates an answer
func (s *QdrantRAGService) Query(ctx context.Context, question string, topK int) (string, error) {
	// Generate embedding for the question
	queryVector, err := s.embedding.Embed(ctx, question)
	if err != nil {
		return "", fmt.Errorf("generate query embedding: %w", err)
	}

	// Search Qdrant
	results, err := s.qdrant.SearchPoints(ctx, queryVector, topK)
	if err != nil {
		return "", fmt.Errorf("search points: %w", err)
	}

	if len(results) == 0 {
		return "没有找到相关的文档。", nil
	}

	// Build context from search results
	var contextParts []string
	for i, result := range results {
		if content, ok := result.Payload["content"]; ok {
			contextParts = append(contextParts, fmt.Sprintf("[%d] %s (相似度: %.2f)", i+1, content, result.Score))
		}
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

	resp, err := s.llm.Chat(ctx, req)
	if err != nil {
		return "", fmt.Errorf("generate answer: %w", err)
	}

	return resp.Content, nil
}

// QueryWithContext queries with additional context
func (s *QdrantRAGService) QueryWithContext(ctx context.Context, question string, topK int, systemPrompt string) (string, error) {
	queryVector, err := s.embedding.Embed(ctx, question)
	if err != nil {
		return "", fmt.Errorf("generate query embedding: %w", err)
	}

	results, err := s.qdrant.SearchPoints(ctx, queryVector, topK)
	if err != nil {
		return "", fmt.Errorf("search points: %w", err)
	}

	var contextParts []string
	for i, result := range results {
		if content, ok := result.Payload["content"]; ok {
			contextParts = append(contextParts, fmt.Sprintf("[%d] %s", i+1, content))
		}
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

	resp, err := s.llm.Chat(ctx, req)
	if err != nil {
		return "", fmt.Errorf("generate answer: %w", err)
	}

	return resp.Content, nil
}

// DeleteDocument deletes a document from the knowledge base
func (s *QdrantRAGService) DeleteDocument(ctx context.Context, docID string) error {
	return s.qdrant.DeletePoints(ctx, []string{docID})
}

// GetCollectionInfo gets collection information
func (s *QdrantRAGService) GetCollectionInfo(ctx context.Context) (*vectorstore.CollectionInfo, error) {
	return s.qdrant.GetCollectionInfo(ctx)
}
