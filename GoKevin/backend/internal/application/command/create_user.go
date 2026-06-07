package command

import (
	"context"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/kevin-ai/go-kevin/internal/domain/entity"
	"github.com/kevin-ai/go-kevin/internal/domain/event"
	"github.com/kevin-ai/go-kevin/internal/domain/repository"
)

type CreateUserCommand struct {
	UserName string `json:"userName" binding:"required"`
	Password string `json:"password" binding:"required"`
	RealName string `json:"realName"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	TenantID int64  `json:"tenantId"`
}

type CreateUserResult struct {
	ID       int64  `json:"id"`
	UserName string `json:"userName"`
}

type CreateUserHandler struct {
	userRepo repository.UserRepository
	eventBus event.EventBus
}

func NewCreateUserHandler(
	userRepo repository.UserRepository,
	eventBus event.EventBus,
) *CreateUserHandler {
	return &CreateUserHandler{
		userRepo: userRepo,
		eventBus: eventBus,
	}
}

func (h *CreateUserHandler) Handle(ctx context.Context, cmd *CreateUserCommand) (*CreateUserResult, error) {
	existingUser, err := h.userRepo.GetByUserName(ctx, cmd.UserName)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("用户名已存在")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(cmd.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		UserName: cmd.UserName,
		Password: string(hashedPassword),
		RealName: cmd.RealName,
		Email:    cmd.Email,
		Phone:    cmd.Phone,
		TenantID: cmd.TenantID,
	}

	if err := h.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	h.eventBus.Publish(&event.UserCreatedEvent{
		UserID:    user.ID,
		UserName:  user.UserName,
		Timestamp: time.Now(),
	})

	return &CreateUserResult{
		ID:       user.ID,
		UserName: user.UserName,
	}, nil
}
