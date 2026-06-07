package persistence

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDatabaseConnection(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	db, err := InitMySQL("root:admin123@tcp(127.0.0.1:3306)/kevin_app?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		t.Skip("Database not available:", err)
	}

	assert.NotNil(t, db)

	sqlDB, err := db.DB()
	assert.NoError(t, err)
	assert.NoError(t, sqlDB.Ping())
}
