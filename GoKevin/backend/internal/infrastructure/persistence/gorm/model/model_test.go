package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserGORMTableName(t *testing.T) {
	user := UserGORM{}
	assert.Equal(t, "t_user", user.TableName())
}

func TestRoleGORMTableName(t *testing.T) {
	role := RoleGORM{}
	assert.Equal(t, "t_role", role.TableName())
}

func TestAIAppGORMTableName(t *testing.T) {
	app := AIAppGORM{}
	assert.Equal(t, "t_ai_apps", app.TableName())
}

func TestChatSessionGORMTableName(t *testing.T) {
	session := ChatSessionGORM{}
	assert.Equal(t, "t_ai_chat_sessions", session.TableName())
}
