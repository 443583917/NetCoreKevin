package storage

import (
	"context"
	"io"
)

// Storage is the interface for file storage
type Storage interface {
	// Upload uploads a file
	Upload(ctx context.Context, path string, reader io.Reader) error

	// Download downloads a file
	Download(ctx context.Context, path string) (io.ReadCloser, error)

	// Delete deletes a file
	Delete(ctx context.Context, path string) error

	// GetURL returns the URL for a file
	GetURL(ctx context.Context, path string) (string, error)

	// Exists checks if a file exists
	Exists(ctx context.Context, path string) (bool, error)
}

// File represents a file
type File struct {
	Name    string
	Size    int64
	Content []byte
	URL     string
}
