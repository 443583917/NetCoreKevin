package query

import (
	"context"

	"github.com/kevin-ai/go-kevin/internal/domain/repository"
)

type ListUsersQuery struct {
	TenantID int64
	Page     int
	PageSize int
}

type ListUsersResult struct {
	Users []*GetUserResult
	Total int64
}

type ListUsersHandler struct {
	userRepo repository.UserRepository
}

func NewListUsersHandler(userRepo repository.UserRepository) *ListUsersHandler {
	return &ListUsersHandler{userRepo: userRepo}
}

func (h *ListUsersHandler) Handle(ctx context.Context, query *ListUsersQuery) (*ListUsersResult, error) {
	users, total, err := h.userRepo.ListByTenantID(ctx, query.TenantID, query.Page, query.PageSize)
	if err != nil {
		return nil, err
	}

	result := make([]*GetUserResult, len(users))
	for i, user := range users {
		result[i] = &GetUserResult{
			ID:       user.ID,
			UserName: user.UserName,
			RealName: user.RealName,
			Email:    user.Email,
			Phone:    user.Phone,
			TenantID: user.TenantID,
		}
	}

	return &ListUsersResult{
		Users: result,
		Total: total,
	}, nil
}
