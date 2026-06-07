package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestLoggerInit(t *testing.T) {
	Init()

	logger := Get()
	assert.NotNil(t, logger)

	// Test log output doesn't panic
	Info("test info message", zap.String("key", "value"))
	Debug("test debug message")
	Warn("test warn message")
	Error("test error message")
}
