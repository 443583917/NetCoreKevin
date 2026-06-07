package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"

	"github.com/kevin-ai/go-kevin/internal/domain/entity"
)

type MockUserRepoForAuth struct {
	mock.Mock
}

func (m *MockUserRepoForAuth) GetByUserName(ctx context.Context, userName string) (*entity.User, error) {
	args := m.Called(ctx, userName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepoForAuth) Create(ctx context.Context, user *entity.User) error {
	return m.Called(ctx, user).Error(0)
}

func (m *MockUserRepoForAuth) GetByID(ctx context.Context, id int64) (*entity.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepoForAuth) Update(ctx context.Context, user *entity.User) error {
	return m.Called(ctx, user).Error(0)
}

func (m *MockUserRepoForAuth) Delete(ctx context.Context, id int64) error {
	return m.Called(ctx, id).Error(0)
}

func (m *MockUserRepoForAuth) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepoForAuth) ListByTenantID(ctx context.Context, tenantID int64, page, pageSize int) ([]*entity.User, int64, error) {
	args := m.Called(ctx, tenantID, page, pageSize)
	return args.Get(0).([]*entity.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserRepoForAuth) GetWithRoles(ctx context.Context, id int64) (*entity.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepoForAuth) BindRoles(ctx context.Context, userID int64, roleIDs []int64) error {
	return m.Called(ctx, userID, roleIDs).Error(0)
}

func TestLogin_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := new(MockUserRepoForAuth)
	handler := NewAuthHandler(mockRepo, "test-secret", 24)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := &entity.User{
		ID:       1,
		UserName: "testuser",
		Password: string(hashedPassword),
		TenantID: 1000,
	}

	mockRepo.On("GetByUserName", mock.Anything, "testuser").Return(user, nil)

	body := LoginRequest{
		UserName: "testuser",
		Password: "password123",
	}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)
	r.POST("/login", handler.Login)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestLogin_WrongPassword(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := new(MockUserRepoForAuth)
	handler := NewAuthHandler(mockRepo, "test-secret", 24)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.DefaultCost)
	user := &entity.User{
		ID:       1,
		UserName: "testuser",
		Password: string(hashedPassword),
	}

	mockRepo.On("GetByUserName", mock.Anything, "testuser").Return(user, nil)

	body := LoginRequest{
		UserName: "testuser",
		Password: "wrong-password",
	}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)
	r.POST("/login", handler.Login)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
