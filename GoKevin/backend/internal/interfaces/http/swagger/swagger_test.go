package swagger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultSwaggerInfo(t *testing.T) {
	info := DefaultSwaggerInfo()
	assert.NotNil(t, info)
	assert.Equal(t, "Go Kevin API", info.Title)
	assert.Equal(t, "1.0", info.Version)
}

func TestGetSwaggerJSON(t *testing.T) {
	info := DefaultSwaggerInfo()
	json := getSwaggerJSON(info)

	assert.NotNil(t, json)
	assert.Equal(t, "3.0.0", json["openapi"])
}
