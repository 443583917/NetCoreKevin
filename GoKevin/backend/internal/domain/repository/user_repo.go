package repository

import (
	"context"

	"github.com/kevin-ai/go-kevin/internal/domain/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	GetByID(ctx context.Context, id int64) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id int64) error

	GetByUserName(ctx context.Context, userName string) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	ListByTenantID(ctx context.Context, tenantID int64, page, pageSize int) ([]*entity.User, int64, error)

	GetWithRoles(ctx context.Context, id int64) (*entity.User, error)
	BindRoles(ctx context.Context, userID int64, roleIDs []int64) error
}
