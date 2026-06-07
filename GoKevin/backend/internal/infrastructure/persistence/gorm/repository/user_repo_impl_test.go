package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kevin-ai/go-kevin/internal/domain/entity"
	"github.com/kevin-ai/go-kevin/internal/infrastructure/persistence/gorm/model"
)

func TestToGORMModel(t *testing.T) {
	repo := &UserRepositoryImpl{}

	user := &entity.User{
		ID:       1,
		UserName: "testuser",
		Password: "hashed",
		RealName: "测试",
		Email:    "test@example.com",
		TenantID: 1000,
	}

	gormModel := repo.toGORMModel(user)

	assert.Equal(t, int64(1), gormModel.ID)
	assert.Equal(t, "testuser", gormModel.UserName)
	assert.Equal(t, "hashed", gormModel.Password)
	assert.Equal(t, int64(1000), gormModel.TenantID)
}

func TestToEntity(t *testing.T) {
	repo := &UserRepositoryImpl{}

	gormModel := &model.UserGORM{
		ID:       1,
		UserName: "testuser",
		RealName: "测试",
		Email:    "test@example.com",
		TenantID: 1000,
	}

	entity := repo.toEntity(gormModel)

	assert.Equal(t, int64(1), entity.ID)
	assert.Equal(t, "testuser", entity.UserName)
	assert.Equal(t, "测试", entity.RealName)
}
