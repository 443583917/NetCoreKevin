package command

import (
	"context"

	"github.com/kevin-ai/go-kevin/internal/domain/repository"
)

type DeleteUserCommand struct {
	ID int64
}

type DeleteUserHandler struct {
	userRepo repository.UserRepository
}

func NewDeleteUserHandler(userRepo repository.UserRepository) *DeleteUserHandler {
	return &DeleteUserHandler{userRepo: userRepo}
}

func (h *DeleteUserHandler) Handle(ctx context.Context, cmd *DeleteUserCommand) error {
	return h.userRepo.Delete(ctx, cmd.ID)
}
