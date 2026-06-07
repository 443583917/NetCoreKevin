package storage

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLocalStorage_New(t *testing.T) {
	storage := NewLocalStorage("/tmp/storage", "http://localhost:8080/files")
	assert.NotNil(t, storage)
}

func TestLocalStorage_Upload(t *testing.T) {
	dir := t.TempDir()
	storage := NewLocalStorage(dir, "http://localhost:8080/files")
	ctx := context.Background()

	reader := strings.NewReader("hello world")
	err := storage.Upload(ctx, "test/file.txt", reader)
	require.NoError(t, err)

	// Verify file exists on disk
	fullPath := filepath.Join(dir, "test/file.txt")
	data, err := os.ReadFile(fullPath)
	require.NoError(t, err)
	assert.Equal(t, "hello world", string(data))
}

func TestLocalStorage_Download(t *testing.T) {
	dir := t.TempDir()
	storage := NewLocalStorage(dir, "http://localhost:8080/files")
	ctx := context.Background()

	// Upload first
	reader := strings.NewReader("download content")
	err := storage.Upload(ctx, "download.txt", reader)
	require.NoError(t, err)

	// Download
	rc, err := storage.Download(ctx, "download.txt")
	require.NoError(t, err)
	defer rc.Close()

	data, err := io.ReadAll(rc)
	require.NoError(t, err)
	assert.Equal(t, "download content", string(data))
}

func TestLocalStorage_Download_NotFound(t *testing.T) {
	dir := t.TempDir()
	storage := NewLocalStorage(dir, "http://localhost:8080/files")
	ctx := context.Background()

	_, err := storage.Download(ctx, "nonexistent.txt")
	assert.Error(t, err)
}

func TestLocalStorage_Delete(t *testing.T) {
	dir := t.TempDir()
	storage := NewLocalStorage(dir, "http://localhost:8080/files")
	ctx := context.Background()

	// Upload then delete
	reader := strings.NewReader("to be deleted")
	err := storage.Upload(ctx, "delete_me.txt", reader)
	require.NoError(t, err)

	err = storage.Delete(ctx, "delete_me.txt")
	require.NoError(t, err)

	exists, err := storage.Exists(ctx, "delete_me.txt")
	require.NoError(t, err)
	assert.False(t, exists)
}

func TestLocalStorage_Delete_NotFound(t *testing.T) {
	dir := t.TempDir()
	storage := NewLocalStorage(dir, "http://localhost:8080/files")
	ctx := context.Background()

	err := storage.Delete(ctx, "nonexistent.txt")
	assert.Error(t, err)
}

func TestLocalStorage_GetURL(t *testing.T) {
	storage := NewLocalStorage("/tmp/storage", "http://localhost:8080/files")
	ctx := context.Background()

	url, err := storage.GetURL(ctx, "images/photo.jpg")
	require.NoError(t, err)
	assert.Equal(t, "http://localhost:8080/files/images/photo.jpg", url)
}

func TestLocalStorage_Exists(t *testing.T) {
	dir := t.TempDir()
	storage := NewLocalStorage(dir, "http://localhost:8080/files")
	ctx := context.Background()

	// Does not exist yet
	exists, err := storage.Exists(ctx, "exists.txt")
	require.NoError(t, err)
	assert.False(t, exists)

	// Upload
	reader := strings.NewReader("exists")
	err = storage.Upload(ctx, "exists.txt", reader)
	require.NoError(t, err)

	// Now exists
	exists, err = storage.Exists(ctx, "exists.txt")
	require.NoError(t, err)
	assert.True(t, exists)
}

func TestLocalStorage_NestedDirectories(t *testing.T) {
	dir := t.TempDir()
	storage := NewLocalStorage(dir, "http://localhost:8080/files")
	ctx := context.Background()

	reader := strings.NewReader("nested content")
	err := storage.Upload(ctx, "a/b/c/deep.txt", reader)
	require.NoError(t, err)

	exists, err := storage.Exists(ctx, "a/b/c/deep.txt")
	require.NoError(t, err)
	assert.True(t, exists)
}
