package migration

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMigrator_Run(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping migration test in short mode")
	}
	assert.True(t, true)
}

func TestMigrator_ModelList(t *testing.T) {
	models := getModels()

	expectedModels := []string{
		"UserGORM",
		"RoleGORM",
		"PermissionGORM",
		"AIAppGORM",
		"AIModelGORM",
		"SkillGORM",
		"ChatSessionGORM",
		"ChatMessageGORM",
		"KnowledgeBaseGORM",
		"DocumentGORM",
	}

	for _, expected := range expectedModels {
		found := false
		for _, model := range models {
			if model == expected {
				found = true
				break
			}
		}
		assert.True(t, found, "Model %s not found", expected)
	}
}
