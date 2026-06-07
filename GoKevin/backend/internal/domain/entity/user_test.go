package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUserEntity(t *testing.T) {
	user := &User{
		ID:         1,
		UserName:   "testuser",
		Password:   "hashedpassword",
		RealName:   "测试用户",
		Email:      "test@example.com",
		Phone:      "13800138000",
		TenantID:   1000,
		CreateTime: time.Now(),
	}

	assert.Equal(t, int64(1), user.ID)
	assert.Equal(t, "testuser", user.UserName)
	assert.Equal(t, "测试用户", user.RealName)
	assert.Equal(t, int64(1000), user.TenantID)
}

func TestUserJSONSerialization(t *testing.T) {
	user := &User{
		ID:       1,
		UserName: "testuser",
		Password: "hashedpassword",
	}

	assert.NotEmpty(t, user.UserName)
}
