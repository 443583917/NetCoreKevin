package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/kevin-ai/go-kevin/internal/domain/entity"
	"github.com/kevin-ai/go-kevin/internal/domain/event"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	if args.Get(0) != nil {
		user.ID = 1
	}
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id int64) (*entity.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
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

func (m *MockUserRepository) GetWithRoles(ctx context.Context, id int64) (*entity.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) BindRoles(ctx context.Context, userID int64, roleIDs []int64) error {
	args := m.Called(ctx, userID, roleIDs)
	return args.Error(0)
}

type MockEventBus struct {
	mock.Mock
}

func (m *MockEventBus) Subscribe(eventName string, handler event.EventHandler) {
	m.Called(eventName, handler)
}

func (m *MockEventBus) Publish(event event.DomainEvent) {
	m.Called(event)
}

func TestCreateUser_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockBus := new(MockEventBus)

	handler := NewCreateUserHandler(mockRepo, mockBus)

	cmd := &CreateUserCommand{
		UserName: "testuser",
		Password: "password123",
		RealName: "测试用户",
		Email:    "test@example.com",
	}

	mockRepo.On("GetByUserName", mock.Anything, "testuser").Return(nil, nil)
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.User")).
		Run(func(args mock.Arguments) {
			user := args.Get(1).(*entity.User)
			user.ID = 1
		}).Return(nil)
	mockBus.On("Publish", mock.AnythingOfType("*event.UserCreatedEvent")).Return()

	result, err := handler.Handle(context.Background(), cmd)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(1), result.ID)
	assert.Equal(t, "testuser", result.UserName)

	mockRepo.AssertExpectations(t)
	mockBus.AssertExpectations(t)
}

func TestCreateUser_DuplicateUserName(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockBus := new(MockEventBus)

	handler := NewCreateUserHandler(mockRepo, mockBus)

	cmd := &CreateUserCommand{
		UserName: "existinguser",
		Password: "password123",
	}

	existingUser := &entity.User{ID: 1, UserName: "existinguser"}
	mockRepo.On("GetByUserName", mock.Anything, "existinguser").Return(existingUser, nil)

	result, err := handler.Handle(context.Background(), cmd)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "用户名已存在")
}
