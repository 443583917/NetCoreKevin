package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitConfig(t *testing.T) {
	content := `
server:
  port: 9901
  mode: debug
database:
  host: 127.0.0.1
  port: 3306
  user: root
  password: test123
  dbname: test_db
`
	tmpFile, err := os.CreateTemp("", "config-*.yaml")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(content)
	assert.NoError(t, err)
	tmpFile.Close()

	err = Init(tmpFile.Name())
	assert.NoError(t, err)

	cfg := Get()
	assert.NotNil(t, cfg)
	assert.Equal(t, 9901, cfg.Server.Port)
	assert.Equal(t, "debug", cfg.Server.Mode)
	assert.Equal(t, "127.0.0.1", cfg.Database.Host)
	assert.Equal(t, "test_db", cfg.Database.DBName)
}
