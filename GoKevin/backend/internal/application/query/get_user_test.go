package query

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/kevin-ai/go-kevin/internal/domain/entity"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetByID(ctx context.Context, id int64) (*entity.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) GetWithRoles(ctx context.Context, id int64) (*entity.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) Create(ctx context.Context, user *entity.User) error {
	return m.Called(ctx, user).Error(0)
}

func (m *MockUserRepository) Update(ctx context.Context, user *entity.User) error {
	return m.Called(ctx, user).Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id int64) error {
	return m.Called(ctx, id).Error(0)
}

func (m *MockUserRepository) GetByUserName(ctx context.Context, userName string) (*entity.User, error) {
	args := m.Called(ctx, userName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) ListByTenantID(ctx context.Context, tenantID int64, page, pageSize int) ([]*entity.User, int64, error) {
	args := m.Called(ctx, tenantID, page, pageSize)
	return args.Get(0).([]*entity.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserRepository) BindRoles(ctx context.Context, userID int64, roleIDs []int64) error {
	return m.Called(ctx, userID, roleIDs).Error(0)
}

func TestGetUser_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	handler := NewGetUserHandler(mockRepo)

	user := &entity.User{
		ID:       1,
		UserName: "testuser",
		RealName: "测试用户",
		Email:    "test@example.com",
		Roles: []*entity.Role{
			{ID: 1, RoleName: "管理员", RoleCode: "admin"},
		},
	}

	mockRepo.On("GetWithRoles", mock.Anything, int64(1)).Return(user, nil)

	result, err := handler.Handle(context.Background(), &GetUserQuery{ID: 1})

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(1), result.ID)
	assert.Equal(t, "testuser", result.UserName)
	assert.Len(t, result.Roles, 1)
	assert.Equal(t, "管理员", result.Roles[0].RoleName)
}

func TestGetUser_NotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	handler := NewGetUserHandler(mockRepo)

	mockRepo.On("GetWithRoles", mock.Anything, int64(999)).Return(nil, nil)

	result, err := handler.Handle(context.Background(), &GetUserQuery{ID: 999})

	assert.NoError(t, err)
	assert.Nil(t, result)
}
