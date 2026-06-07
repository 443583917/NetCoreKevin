package repository

import (
	"context"

	"github.com/kevin-ai/go-kevin/internal/domain/entity"
)

type AIAppRepository interface {
	Create(ctx context.Context, app *entity.AIApp) error
	GetByID(ctx context.Context, id int64) (*entity.AIApp, error)
	Update(ctx context.Context, app *entity.AIApp) error
	Delete(ctx context.Context, id int64) error
	ListByTenantID(ctx context.Context, tenantID int64, page, pageSize int) ([]*entity.AIApp, int64, error)
	GetWithSkills(ctx context.Context, id int64) (*entity.AIApp, error)
}

type ChatSessionRepository interface {
	Create(ctx context.Context, session *entity.ChatSession) error
	GetByID(ctx context.Context, id int64) (*entity.ChatSession, error)
	ListByUserID(ctx context.Context, userID int64, page, pageSize int) ([]*entity.ChatSession, int64, error)
	UpdateTitle(ctx context.Context, id int64, title string) error
	Delete(ctx context.Context, id int64) error
}

type ChatMessageRepository interface {
	Create(ctx context.Context, message *entity.ChatMessage) error
	ListBySessionID(ctx context.Context, sessionID int64, limit int) ([]*entity.ChatMessage, error)
	CountBySessionID(ctx context.Context, sessionID int64) (int64, error)
}

type KnowledgeBaseRepository interface {
	Create(ctx context.Context, kb *entity.KnowledgeBase) error
	GetByID(ctx context.Context, id int64) (*entity.KnowledgeBase, error)
	Update(ctx context.Context, kb *entity.KnowledgeBase) error
	Delete(ctx context.Context, id int64) error
	ListByTenantID(ctx context.Context, tenantID int64, page, pageSize int) ([]*entity.KnowledgeBase, int64, error)
}

type DocumentRepository interface {
	Create(ctx context.Context, doc *entity.Document) error
	GetByID(ctx context.Context, id int64) (*entity.Document, error)
	Update(ctx context.Context, doc *entity.Document) error
	Delete(ctx context.Context, id int64) error
	ListByKBID(ctx context.Context, kbID int64, page, pageSize int) ([]*entity.Document, int64, error)
	UpdateStatus(ctx context.Context, id int64, status string) error
}
