package router

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRegisterRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()

	// Register routes with nil handlers (just testing route registration)
	RegisterRoutes(r, nil, nil)

	routes := r.Routes()

	routeMap := make(map[string]bool)
	for _, route := range routes {
		key := route.Method + ":" + route.Path
		routeMap[key] = true
	}

	// Check critical routes exist
	assert.True(t, routeMap["POST:/api/v1/auth/login"])
	assert.True(t, routeMap["POST:/api/v1/auth/register"])
	assert.True(t, routeMap["POST:/api/v1/auth/refresh"])
	assert.True(t, routeMap["GET:/api/v1/user"])
	assert.True(t, routeMap["GET:/api/v1/user/:id"])
	assert.True(t, routeMap["POST:/api/v1/user"])
	assert.True(t, routeMap["PUT:/api/v1/user/:id"])
	assert.True(t, routeMap["DELETE:/api/v1/user/:id"])

	// Check placeholder routes exist
	assert.True(t, routeMap["GET:/api/v1/aiapps"])
	assert.True(t, routeMap["GET:/api/v1/aiapps/:id"])
	assert.True(t, routeMap["POST:/api/v1/aiapps"])
	assert.True(t, routeMap["PUT:/api/v1/aiapps/:id"])
	assert.True(t, routeMap["DELETE:/api/v1/aiapps/:id"])

	assert.True(t, routeMap["POST:/api/v1/aichat/sessions"])
	assert.True(t, routeMap["GET:/api/v1/aichat/sessions"])
	assert.True(t, routeMap["GET:/api/v1/aichat/sessions/:id"])
	assert.True(t, routeMap["POST:/api/v1/aichat/sessions/:id/messages"])
	assert.True(t, routeMap["GET:/api/v1/aichat/sessions/:id/messages"])

	assert.True(t, routeMap["GET:/api/v1/aikmss"])
	assert.True(t, routeMap["POST:/api/v1/aikmss"])
	assert.True(t, routeMap["GET:/api/v1/aikmss/:id"])
	assert.True(t, routeMap["DELETE:/api/v1/aikmss/:id"])
	assert.True(t, routeMap["POST:/api/v1/aikmss/:id/documents"])
	assert.True(t, routeMap["POST:/api/v1/aikmss/:id/query"])

	assert.True(t, routeMap["POST:/api/v1/file/upload"])
	assert.True(t, routeMap["GET:/api/v1/file/:id/download"])
}
