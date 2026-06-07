package command

import (
	"context"

	"github.com/kevin-ai/go-kevin/internal/domain/repository"
)

type UpdateUserCommand struct {
	ID       int64
	RealName string
	Email    string
	Phone    string
}

type UpdateUserHandler struct {
	userRepo repository.UserRepository
}

func NewUpdateUserHandler(userRepo repository.UserRepository) *UpdateUserHandler {
	return &UpdateUserHandler{userRepo: userRepo}
}

func (h *UpdateUserHandler) Handle(ctx context.Context, cmd *UpdateUserCommand) error {
	user, err := h.userRepo.GetByID(ctx, cmd.ID)
	if err != nil {
		return err
	}
	if user == nil {
		return nil
	}

	user.RealName = cmd.RealName
	user.Email = cmd.Email
	user.Phone = cmd.Phone

	return h.userRepo.Update(ctx, user)
}
