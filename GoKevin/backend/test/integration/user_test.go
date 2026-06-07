package integration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kevin-ai/go-kevin/internal/domain/entity"
	"github.com/kevin-ai/go-kevin/internal/infrastructure/persistence/gorm/repository"
)

func TestUserRepository_Create(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	if testDB == nil {
		t.Skip("Database not available")
	}

	repo := repository.NewUserRepositoryImpl(testDB)

	user := &entity.User{
		UserName: "integrationtest",
		Password: "hashedpassword",
		RealName: "集成测试",
		Email:    "integration@test.com",
		TenantID: 1000,
	}

	err := repo.Create(context.Background(), user)
	assert.NoError(t, err)
	assert.NotZero(t, user.ID)

	// Cleanup
	testDB.Delete(&entity.User{}, user.ID)
}

func TestUserRepository_GetByID(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	if testDB == nil {
		t.Skip("Database not available")
	}

	repo := repository.NewUserRepositoryImpl(testDB)

	// Create test data
	user := &entity.User{
		UserName: "gettest",
		Password: "hashedpassword",
		TenantID: 1000,
	}
	testDB.Create(user)

	// Query
	found, err := repo.GetByID(context.Background(), user.ID)
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, user.UserName, found.UserName)

	// Cleanup
	testDB.Delete(&entity.User{}, user.ID)
}

func TestUserRepository_GetByUserName(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	if testDB == nil {
		t.Skip("Database not available")
	}

	repo := repository.NewUserRepositoryImpl(testDB)

	// Create test data
	user := &entity.User{
		UserName: "usernametest",
		Password: "hashedpassword",
		TenantID: 1000,
	}
	testDB.Create(user)

	// Query
	found, err := repo.GetByUserName(context.Background(), "usernametest")
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, user.ID, found.ID)

	// Cleanup
	testDB.Delete(&entity.User{}, user.ID)
}
