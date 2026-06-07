package query

import (
	"context"

	"github.com/kevin-ai/go-kevin/internal/domain/repository"
)

type GetUserQuery struct {
	ID int64 `json:"id"`
}

type GetUserResult struct {
	ID       int64      `json:"id"`
	UserName string     `json:"userName"`
	RealName string     `json:"realName"`
	Email    string     `json:"email"`
	Phone    string     `json:"phone"`
	TenantID int64      `json:"tenantId"`
	Roles    []RoleInfo `json:"roles"`
}

type RoleInfo struct {
	ID       int64  `json:"id"`
	RoleName string `json:"roleName"`
	RoleCode string `json:"roleCode"`
}

type GetUserHandler struct {
	userRepo repository.UserRepository
}

func NewGetUserHandler(userRepo repository.UserRepository) *GetUserHandler {
	return &GetUserHandler{userRepo: userRepo}
}

func (h *GetUserHandler) Handle(ctx context.Context, query *GetUserQuery) (*GetUserResult, error) {
	user, err := h.userRepo.GetWithRoles(ctx, query.ID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, nil
	}

	roles := make([]RoleInfo, len(user.Roles))
	for i, role := range user.Roles {
		roles[i] = RoleInfo{
			ID:       role.ID,
			RoleName: role.RoleName,
			RoleCode: role.RoleCode,
		}
	}

	return &GetUserResult{
		ID:       user.ID,
		UserName: user.UserName,
		RealName: user.RealName,
		Email:    user.Email,
		Phone:    user.Phone,
		TenantID: user.TenantID,
		Roles:    roles,
	}, nil
}
