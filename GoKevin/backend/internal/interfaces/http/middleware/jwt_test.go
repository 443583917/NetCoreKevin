package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/kevin-ai/go-kevin/pkg/auth"
)

func TestJWTAuth_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	secret := "test-secret"
	middleware := JWTAuth(secret)

	token, _ := auth.GenerateToken(secret, 24, 1, "testuser", 1000)

	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)
	r.Use(middleware)
	r.GET("/test", func(c *gin.Context) {
		userID := c.GetInt64("userId")
		c.JSON(200, gin.H{"userId": userID})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	c.Request = req

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestJWTAuth_MissingToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	secret := "test-secret"
	middleware := JWTAuth(secret)

	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)
	r.Use(middleware)
	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	c.Request = req

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestJWTAuth_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	secret := "test-secret"
	middleware := JWTAuth(secret)

	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)
	r.Use(middleware)
	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	c.Request = req

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
