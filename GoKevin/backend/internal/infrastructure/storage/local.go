package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// LocalStorage implements Storage using local filesystem
type LocalStorage struct {
	basePath string
	baseURL  string
}

// NewLocalStorage creates a new local storage
func NewLocalStorage(basePath, baseURL string) *LocalStorage {
	return &LocalStorage{
		basePath: basePath,
		baseURL:  baseURL,
	}
}

// Upload uploads a file
func (s *LocalStorage) Upload(ctx context.Context, path string, reader io.Reader) error {
	fullPath := filepath.Join(s.basePath, path)

	// Create directory if not exists
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create directory: %w", err)
	}

	// Create file
	file, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer file.Close()

	// Copy content
	_, err = io.Copy(file, reader)
	if err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}

// Download downloads a file
func (s *LocalStorage) Download(ctx context.Context, path string) (io.ReadCloser, error) {
	fullPath := filepath.Join(s.basePath, path)

	file, err := os.Open(fullPath)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}

	return file, nil
}

// Delete deletes a file
func (s *LocalStorage) Delete(ctx context.Context, path string) error {
	fullPath := filepath.Join(s.basePath, path)

	if err := os.Remove(fullPath); err != nil {
		return fmt.Errorf("delete file: %w", err)
	}

	return nil
}

// GetURL returns the URL for a file
func (s *LocalStorage) GetURL(ctx context.Context, path string) (string, error) {
	return s.baseURL + "/" + path, nil
}

// Exists checks if a file exists
func (s *LocalStorage) Exists(ctx context.Context, path string) (bool, error) {
	fullPath := filepath.Join(s.basePath, path)

	_, err := os.Stat(fullPath)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}
