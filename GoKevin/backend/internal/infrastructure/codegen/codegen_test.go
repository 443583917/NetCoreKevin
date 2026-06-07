package codegen

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCodeGenerator_New(t *testing.T) {
	gen := NewCodeGenerator(CodeGenConfig{
		ModuleName:  "user",
		EntityName:  "User",
		TableName:   "t_user",
		OutputDir:   "/tmp/test_output",
		PackageName: "github.com/test/app",
	})

	assert.NotNil(t, gen)
	assert.Equal(t, "user", gen.config.ModuleName)
	assert.Equal(t, "User", gen.config.EntityName)
	assert.Equal(t, "t_user", gen.config.TableName)
}

func TestCodeGenerator_Generate(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "codegen_test_*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	gen := NewCodeGenerator(CodeGenConfig{
		ModuleName:  "order",
		EntityName:  "Order",
		TableName:   "t_order",
		OutputDir:   tmpDir,
		PackageName: "github.com/kevin-ai/go-kevin",
	})

	err = gen.Generate()
	assert.NoError(t, err)

	// Verify entity file was created
	entityPath := filepath.Join(tmpDir, "domain", "entity", "order.go")
	assert.FileExists(t, entityPath)

	entityContent, err := os.ReadFile(entityPath)
	require.NoError(t, err)
	assert.Contains(t, string(entityContent), "type Order struct")
	assert.Contains(t, string(entityContent), "package entity")

	// Verify repository file was created
	repoPath := filepath.Join(tmpDir, "domain", "repository", "order_repo.go")
	assert.FileExists(t, repoPath)

	repoContent, err := os.ReadFile(repoPath)
	require.NoError(t, err)
	assert.Contains(t, string(repoContent), "type OrderRepository interface")
	assert.Contains(t, string(repoContent), "Create(ctx context.Context")
	assert.Contains(t, string(repoContent), "GetByID(ctx context.Context")

	// Verify service file was created
	servicePath := filepath.Join(tmpDir, "application", "service", "order_service.go")
	assert.FileExists(t, servicePath)

	serviceContent, err := os.ReadFile(servicePath)
	require.NoError(t, err)
	assert.Contains(t, string(serviceContent), "type OrderService struct")
	assert.Contains(t, string(serviceContent), "func NewOrderService")

	// Verify handler file was created
	handlerPath := filepath.Join(tmpDir, "interfaces", "http", "handler", "order_handler.go")
	assert.FileExists(t, handlerPath)

	handlerContent, err := os.ReadFile(handlerPath)
	require.NoError(t, err)
	assert.Contains(t, string(handlerContent), "type OrderHandler struct")
	assert.Contains(t, string(handlerContent), "func NewOrderHandler")
	assert.Contains(t, string(handlerContent), "func (h *OrderHandler) GetByID")
}

func TestCodeGenerator_GenerateFilesHaveCorrectImports(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "codegen_import_test_*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	gen := NewCodeGenerator(CodeGenConfig{
		ModuleName:  "product",
		EntityName:  "Product",
		TableName:   "t_product",
		OutputDir:   tmpDir,
		PackageName: "github.com/example/myapp",
	})

	err = gen.Generate()
	require.NoError(t, err)

	// Check repository imports
	repoContent, err := os.ReadFile(filepath.Join(tmpDir, "domain", "repository", "product_repo.go"))
	require.NoError(t, err)
	assert.Contains(t, string(repoContent), "github.com/example/myapp/internal/domain/entity")

	// Check service imports
	serviceContent, err := os.ReadFile(filepath.Join(tmpDir, "application", "service", "product_service.go"))
	require.NoError(t, err)
	assert.Contains(t, string(serviceContent), "github.com/example/myapp/internal/domain/entity")
	assert.Contains(t, string(serviceContent), "github.com/example/myapp/internal/domain/repository")
}

func TestCodeGenerator_GenerateCreatesDirectories(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "codegen_dir_test_*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	gen := NewCodeGenerator(CodeGenConfig{
		ModuleName:  "category",
		EntityName:  "Category",
		TableName:   "t_category",
		OutputDir:   tmpDir,
		PackageName: "github.com/test/app",
	})

	err = gen.Generate()
	require.NoError(t, err)

	// Verify all directories were created
	dirs := []string{
		filepath.Join(tmpDir, "domain", "entity"),
		filepath.Join(tmpDir, "domain", "repository"),
		filepath.Join(tmpDir, "application", "service"),
		filepath.Join(tmpDir, "interfaces", "http", "handler"),
	}

	for _, dir := range dirs {
		info, err := os.Stat(dir)
		assert.NoError(t, err, "directory should exist: %s", dir)
		assert.True(t, info.IsDir(), "path should be a directory: %s", dir)
	}
}
