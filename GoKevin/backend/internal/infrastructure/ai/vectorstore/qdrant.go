package vectorstore

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// QdrantClient implements VectorStore using Qdrant
type QdrantClient struct {
	baseURL    string
	httpClient *http.Client
	collection string
}

// QdrantConfig represents Qdrant configuration
type QdrantConfig struct {
	Host       string
	Port       int
	Collection string
	APIKey     string
}

// NewQdrantClient creates a new Qdrant client
func NewQdrantClient(cfg QdrantConfig) *QdrantClient {
	baseURL := fmt.Sprintf("http://%s:%d", cfg.Host, cfg.Port)
	if cfg.Host == "" {
		baseURL = "http://localhost:6333"
	}

	return &QdrantClient{
		baseURL:    baseURL,
		httpClient: &http.Client{},
		collection: cfg.Collection,
	}
}

// CreateCollection creates a new collection
func (c *QdrantClient) CreateCollection(ctx context.Context, vectorSize int) error {
	url := fmt.Sprintf("%s/collections/%s", c.baseURL, c.collection)

	body := map[string]interface{}{
		"vectors": map[string]interface{}{
			"size":     vectorSize,
			"distance": "Cosine",
		},
	}

	return c.doRequest(ctx, "PUT", url, body, nil)
}

// UpsertPoints upserts points (vectors) to the collection
func (c *QdrantClient) UpsertPoints(ctx context.Context, points []Point) error {
	url := fmt.Sprintf("%s/collections/%s/points", c.baseURL, c.collection)

	body := map[string]interface{}{
		"points": points,
	}

	return c.doRequest(ctx, "PUT", url, body, nil)
}

// SearchPoints searches for similar points
func (c *QdrantClient) SearchPoints(ctx context.Context, vector []float32, topK int) ([]SearchResult, error) {
	url := fmt.Sprintf("%s/collections/%s/points/search", c.baseURL, c.collection)

	body := map[string]interface{}{
		"vector":       vector,
		"top":          topK,
		"with_payload": true,
	}

	var result searchResponse
	if err := c.doRequest(ctx, "POST", url, body, &result); err != nil {
		return nil, err
	}

	searchResults := make([]SearchResult, len(result.Result))
	for i, r := range result.Result {
		searchResults[i] = SearchResult{
			ID:      r.ID,
			Score:   r.Score,
			Payload: r.Payload,
		}
	}

	return searchResults, nil
}

// DeletePoints deletes points by IDs
func (c *QdrantClient) DeletePoints(ctx context.Context, ids []string) error {
	url := fmt.Sprintf("%s/collections/%s/points/delete", c.baseURL, c.collection)

	body := map[string]interface{}{
		"points": ids,
	}

	return c.doRequest(ctx, "POST", url, body, nil)
}

// GetCollectionInfo gets collection information
func (c *QdrantClient) GetCollectionInfo(ctx context.Context) (*CollectionInfo, error) {
	url := fmt.Sprintf("%s/collections/%s", c.baseURL, c.collection)

	var result collectionInfoResponse
	if err := c.doRequest(ctx, "GET", url, nil, &result); err != nil {
		return nil, err
	}

	return &CollectionInfo{
		Name:           result.Result.Name,
		VectorsCount:   result.Result.VectorsCount,
		PointsCount:    result.Result.PointsCount,
	}, nil
}

func (c *QdrantClient) doRequest(ctx context.Context, method, url string, body interface{}, result interface{}) error {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshal body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("decode response: %w", err)
		}
	}

	return nil
}

// Qdrant API types
type Point struct {
	ID      string                 `json:"id"`
	Vector  []float32              `json:"vector"`
	Payload map[string]interface{} `json:"payload,omitempty"`
}

type SearchResult struct {
	ID      string                 `json:"id"`
	Score   float64                `json:"score"`
	Payload map[string]interface{} `json:"payload,omitempty"`
}

type CollectionInfo struct {
	Name         string `json:"name"`
	VectorsCount int64  `json:"vectors_count"`
	PointsCount  int64  `json:"points_count"`
}

type searchResponse struct {
	Result []struct {
		ID      string                 `json:"id"`
		Score   float64                `json:"score"`
		Payload map[string]interface{} `json:"payload,omitempty"`
	} `json:"result"`
}

type collectionInfoResponse struct {
	Result struct {
		Name         string `json:"name"`
		VectorsCount int64  `json:"vectors_count"`
		PointsCount  int64  `json:"points_count"`
	} `json:"result"`
}
