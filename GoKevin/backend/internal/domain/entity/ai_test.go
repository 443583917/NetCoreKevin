package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAIAppEntity(t *testing.T) {
	app := &AIApp{
		ID:           1,
		AppName:      "智能客服",
		AppDesc:      "24小时在线客服机器人",
		ModelID:      1,
		SystemPrompt: "你是一个专业的客服助手",
		TenantID:     1000,
	}

	assert.Equal(t, int64(1), app.ID)
	assert.Equal(t, "智能客服", app.AppName)
	assert.Equal(t, int64(1000), app.TenantID)
}

func TestChatSessionEntity(t *testing.T) {
	session := &ChatSession{
		ID:     1,
		UserID: 100,
		AppID:  1,
		Title:  "新对话",
	}

	assert.Equal(t, int64(1), session.ID)
	assert.Equal(t, int64(100), session.UserID)
}

func TestKnowledgeBaseEntity(t *testing.T) {
	kb := &KnowledgeBase{
		ID:          1,
		Name:        "产品知识库",
		Description: "包含所有产品文档",
		VectorModel: "text-embedding-ada-002",
		TenantID:    1000,
	}

	assert.Equal(t, "产品知识库", kb.Name)
	assert.Equal(t, "text-embedding-ada-002", kb.VectorModel)
}
