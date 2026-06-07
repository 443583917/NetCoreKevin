package codegen

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// CodeGenConfig represents code generation configuration
type CodeGenConfig struct {
	ModuleName  string // e.g., "user"
	EntityName  string // e.g., "User"
	TableName   string // e.g., "t_user"
	OutputDir   string
	PackageName string
}

// CodeGenerator generates CRUD code
type CodeGenerator struct {
	config CodeGenConfig
}

// NewCodeGenerator creates a new code generator
func NewCodeGenerator(config CodeGenConfig) *CodeGenerator {
	return &CodeGenerator{config: config}
}

// Generate generates all CRUD code
func (g *CodeGenerator) Generate() error {
	// Generate entity
	if err := g.generateEntity(); err != nil {
		return fmt.Errorf("generate entity: %w", err)
	}

	// Generate repository
	if err := g.generateRepository(); err != nil {
		return fmt.Errorf("generate repository: %w", err)
	}

	// Generate service
	if err := g.generateService(); err != nil {
		return fmt.Errorf("generate service: %w", err)
	}

	// Generate handler
	if err := g.generateHandler(); err != nil {
		return fmt.Errorf("generate handler: %w", err)
	}

	return nil
}

func (g *CodeGenerator) generateEntity() error {
	dir := filepath.Join(g.config.OutputDir, "domain", "entity")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	filePath := filepath.Join(dir, strings.ToLower(g.config.ModuleName)+".go")
	return g.writeFile(filePath, entityTemplate)
}

func (g *CodeGenerator) generateRepository() error {
	dir := filepath.Join(g.config.OutputDir, "domain", "repository")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	filePath := filepath.Join(dir, strings.ToLower(g.config.ModuleName)+"_repo.go")
	return g.writeFile(filePath, repositoryTemplate)
}

func (g *CodeGenerator) generateService() error {
	dir := filepath.Join(g.config.OutputDir, "application", "service")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	filePath := filepath.Join(dir, strings.ToLower(g.config.ModuleName)+"_service.go")
	return g.writeFile(filePath, serviceTemplate)
}

func (g *CodeGenerator) generateHandler() error {
	dir := filepath.Join(g.config.OutputDir, "interfaces", "http", "handler")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	filePath := filepath.Join(dir, strings.ToLower(g.config.ModuleName)+"_handler.go")
	return g.writeFile(filePath, handlerTemplate)
}

func (g *CodeGenerator) writeFile(filePath, tmpl string) error {
	t, err := template.New("code").Parse(tmpl)
	if err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	return t.Execute(file, g.config)
}
