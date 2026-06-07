package integration

import (
	"os"
	"testing"

	"gorm.io/gorm"

	"github.com/kevin-ai/go-kevin/internal/infrastructure/persistence"
)

var testDB *gorm.DB

func TestMain(m *testing.M) {
	dsn := os.Getenv("TEST_DB_DSN")
	if dsn == "" {
		dsn = "root:admin123@tcp(127.0.0.1:3306)/kevin_test?charset=utf8mb4&parseTime=True&loc=Local"
	}

	var err error
	testDB, err = persistence.InitMySQL(dsn)
	if err != nil {
		// Skip integration tests if database is not available
		os.Exit(0)
	}

	os.Exit(m.Run())
}
