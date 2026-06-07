package handler

import (
	"context"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/kevin-ai/go-kevin/internal/application/command"
	"github.com/kevin-ai/go-kevin/internal/application/query"
	"github.com/kevin-ai/go-kevin/pkg/response"
)

// CreateUserHandler defines the interface for creating users
type CreateUserHandler interface {
	Handle(ctx context.Context, cmd *command.CreateUserCommand) (*command.CreateUserResult, error)
}

// GetUserHandler defines the interface for getting a user by ID
type GetUserHandler interface {
	Handle(ctx context.Context, q *query.GetUserQuery) (*query.GetUserResult, error)
}

// ListUsersHandler defines the interface for listing users
type ListUsersHandler interface {
	Handle(ctx context.Context, q *query.ListUsersQuery) (*query.ListUsersResult, error)
}

// UpdateUserHandler defines the interface for updating a user
type UpdateUserHandler interface {
	Handle(ctx context.Context, cmd *command.UpdateUserCommand) error
}

// DeleteUserHandler defines the interface for deleting a user
type DeleteUserHandler interface {
	Handle(ctx context.Context, cmd *command.DeleteUserCommand) error
}

type UserHandler struct {
	createUser CreateUserHandler
	getUser    GetUserHandler
	listUsers  ListUsersHandler
	updateUser UpdateUserHandler
	deleteUser DeleteUserHandler
}

func NewUserHandler(
	createUser CreateUserHandler,
	getUser GetUserHandler,
	listUsers ListUsersHandler,
	updateUser UpdateUserHandler,
	deleteUser DeleteUserHandler,
) *UserHandler {
	return &UserHandler{
		createUser: createUser,
		getUser:    getUser,
		listUsers:  listUsers,
		updateUser: updateUser,
		deleteUser: deleteUser,
	}
}

func (h *UserHandler) Create(c *gin.Context) {
	var cmd command.CreateUserCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	tenantID := c.GetInt64("tenantId")
	cmd.TenantID = tenantID

	result, err := h.createUser.Handle(c.Request.Context(), &cmd)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, result)
}

func (h *UserHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	result, err := h.getUser.Handle(c.Request.Context(), &query.GetUserQuery{ID: id})
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	if result == nil {
		response.NotFound(c, "用户不存在")
		return
	}

	response.Success(c, result)
}

func (h *UserHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	tenantID := c.GetInt64("tenantId")

	result, err := h.listUsers.Handle(c.Request.Context(), &query.ListUsersQuery{
		TenantID: tenantID,
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Page(c, result.Users, result.Total, page, pageSize)
}

func (h *UserHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	var cmd command.UpdateUserCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}
	cmd.ID = id

	if err := h.updateUser.Handle(c.Request.Context(), &cmd); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *UserHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	if err := h.deleteUser.Handle(c.Request.Context(), &command.DeleteUserCommand{ID: id}); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, nil)
}

