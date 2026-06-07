package codegen

// Entity template for generating entity structs
const entityTemplate = `package entity

import "time"

// {{.EntityName}} represents {{.ModuleName}} entity
type {{.EntityName}} struct {
	ID         int64     ` + "`json:\"id\"`" + `
	TenantID   int64     ` + "`json:\"tenantId\"`" + `
	IsDelete   bool      ` + "`json:\"isDelete\"`" + `
	CreateTime time.Time ` + "`json:\"createTime\"`" + `
	UpdateTime time.Time ` + "`json:\"updateTime\"`" + `
}
`

// Repository template for generating repository interfaces
const repositoryTemplate = `package repository

import (
	"context"

	"{{.PackageName}}/internal/domain/entity"
)

// {{.EntityName}}Repository is the interface for {{.ModuleName}} repository
type {{.EntityName}}Repository interface {
	Create(ctx context.Context, entity *entity.{{.EntityName}}) error
	GetByID(ctx context.Context, id int64) (*entity.{{.EntityName}}, error)
	Update(ctx context.Context, entity *entity.{{.EntityName}}) error
	Delete(ctx context.Context, id int64) error
	ListByTenantID(ctx context.Context, tenantID int64, page, pageSize int) ([]*entity.{{.EntityName}}, int64, error)
}
`

// Service template for generating service implementations
const serviceTemplate = `package service

import (
	"context"

	"{{.PackageName}}/internal/domain/entity"
	"{{.PackageName}}/internal/domain/repository"
)

// {{.EntityName}}Service provides {{.ModuleName}} business logic
type {{.EntityName}}Service struct {
	repo repository.{{.EntityName}}Repository
}

// New{{.EntityName}}Service creates a new service
func New{{.EntityName}}Service(repo repository.{{.EntityName}}Repository) *{{.EntityName}}Service {
	return &{{.EntityName}}Service{repo: repo}
}

// Create creates a new {{.ModuleName}}
func (s *{{.EntityName}}Service) Create(ctx context.Context, entity *entity.{{.EntityName}}) error {
	return s.repo.Create(ctx, entity)
}

// GetByID gets {{.ModuleName}} by ID
func (s *{{.EntityName}}Service) GetByID(ctx context.Context, id int64) (*entity.{{.EntityName}}, error) {
	return s.repo.GetByID(ctx, id)
}

// Update updates {{.ModuleName}}
func (s *{{.EntityName}}Service) Update(ctx context.Context, entity *entity.{{.EntityName}}) error {
	return s.repo.Update(ctx, entity)
}

// Delete deletes {{.ModuleName}}
func (s *{{.EntityName}}Service) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
`

// Handler template for generating HTTP handler code
const handlerTemplate = `package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"{{.PackageName}}/internal/domain/service"
	"{{.PackageName}}/pkg/response"
)

// {{.EntityName}}Handler handles {{.ModuleName}} HTTP requests
type {{.EntityName}}Handler struct {
	service *service.{{.EntityName}}Service
}

// New{{.EntityName}}Handler creates a new handler
func New{{.EntityName}}Handler(service *service.{{.EntityName}}Service) *{{.EntityName}}Handler {
	return &{{.EntityName}}Handler{service: service}
}

// GetByID handles GET /:id
func (h *{{.EntityName}}Handler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid ID")
		return
	}

	entity, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	if entity == nil {
		response.NotFound(c, "{{.EntityName}} not found")
		return
	}

	response.Success(c, entity)
}
`
