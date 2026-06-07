package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/kevin-ai/go-kevin/internal/application/command"
	"github.com/kevin-ai/go-kevin/internal/application/query"
)

type MockCreateUserHandler struct {
	mock.Mock
}

func (m *MockCreateUserHandler) Handle(ctx context.Context, cmd *command.CreateUserCommand) (*command.CreateUserResult, error) {
	args := m.Called(ctx, cmd)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*command.CreateUserResult), args.Error(1)
}

type MockGetUserHandler struct {
	mock.Mock
}

func (m *MockGetUserHandler) Handle(ctx context.Context, q *query.GetUserQuery) (*query.GetUserResult, error) {
	args := m.Called(ctx, q)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*query.GetUserResult), args.Error(1)
}

type MockListUsersHandler struct {
	mock.Mock
}

func (m *MockListUsersHandler) Handle(ctx context.Context, q *query.ListUsersQuery) (*query.ListUsersResult, error) {
	args := m.Called(ctx, q)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*query.ListUsersResult), args.Error(1)
}

type MockUpdateUserHandler struct {
	mock.Mock
}

func (m *MockUpdateUserHandler) Handle(ctx context.Context, cmd *command.UpdateUserCommand) error {
	args := m.Called(ctx, cmd)
	return args.Error(0)
}

type MockDeleteUserHandler struct {
	mock.Mock
}

func (m *MockDeleteUserHandler) Handle(ctx context.Context, cmd *command.DeleteUserCommand) error {
	args := m.Called(ctx, cmd)
	return args.Error(0)
}

func TestUserHandler_GetByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockGet := new(MockGetUserHandler)
	handler := NewUserHandler(nil, mockGet, nil, nil, nil)

	result := &query.GetUserResult{
		ID:       1,
		UserName: "testuser",
		RealName: "测试用户",
	}

	mockGet.On("Handle", mock.Anything, &query.GetUserQuery{ID: int64(1)}).Return(result, nil)

	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)
	r.GET("/user/:id", handler.GetByID)

	req := httptest.NewRequest(http.MethodGet, "/user/1", nil)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockGet.AssertExpectations(t)
}

func TestUserHandler_GetByID_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockGet := new(MockGetUserHandler)
	handler := NewUserHandler(nil, mockGet, nil, nil, nil)

	mockGet.On("Handle", mock.Anything, &query.GetUserQuery{ID: int64(999)}).Return(nil, nil)

	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)
	r.GET("/user/:id", handler.GetByID)

	req := httptest.NewRequest(http.MethodGet, "/user/999", nil)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "999"}}

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockGet.AssertExpectations(t)
}

func TestUserHandler_GetByID_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewUserHandler(nil, nil, nil, nil, nil)

	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)
	r.GET("/user/:id", handler.GetByID)

	req := httptest.NewRequest(http.MethodGet, "/user/abc", nil)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "abc"}}

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockCreate := new(MockCreateUserHandler)
	handler := NewUserHandler(mockCreate, nil, nil, nil, nil)

	createResult := &command.CreateUserResult{
		ID:       1,
		UserName: "newuser",
	}

	mockCreate.On("Handle", mock.Anything, mock.AnythingOfType("*command.CreateUserCommand")).Return(createResult, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/user", nil)
	c.Set("tenantId", int64(1))

	handler.Create(c)

	// Without proper body binding, this should return 400
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_List(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockList := new(MockListUsersHandler)
	handler := NewUserHandler(nil, nil, mockList, nil, nil)

	listResult := &query.ListUsersResult{
		Users: []*query.GetUserResult{
			{ID: 1, UserName: "user1"},
			{ID: 2, UserName: "user2"},
		},
		Total: 2,
	}

	mockList.On("Handle", mock.Anything, mock.AnythingOfType("*query.ListUsersQuery")).Return(listResult, nil)

	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)
	r.GET("/users", handler.List)

	req := httptest.NewRequest(http.MethodGet, "/users?page=1&pageSize=10", nil)
	c.Request = req
	c.Set("tenantId", int64(1))

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockList.AssertExpectations(t)
}

func TestUserHandler_Delete(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockDelete := new(MockDeleteUserHandler)
	handler := NewUserHandler(nil, nil, nil, nil, mockDelete)

	mockDelete.On("Handle", mock.Anything, &command.DeleteUserCommand{ID: int64(1)}).Return(nil)

	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)
	r.DELETE("/user/:id", handler.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/user/1", nil)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockDelete.AssertExpectations(t)
}
