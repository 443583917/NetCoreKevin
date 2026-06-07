package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/kevin-ai/go-kevin/internal/interfaces/http/middleware"
	"github.com/kevin-ai/go-kevin/pkg/auth"
	"github.com/kevin-ai/go-kevin/pkg/response"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	v1 := r.Group("/api/v1")

	// Public routes
	public := v1.Group("")
	{
		public.POST("/auth/login", func(c *gin.Context) {
			response.Success(c, gin.H{"token": "test-token"})
		})
		public.POST("/auth/register", func(c *gin.Context) {
			response.Success(c, nil)
		})
	}

	// Protected routes
	protected := v1.Group("")
	protected.Use(middleware.JWTAuth("test-secret"))
	{
		protected.GET("/user", func(c *gin.Context) {
			response.Success(c, []gin.H{
				{"id": 1, "userName": "admin"},
			})
		})
		protected.POST("/user", func(c *gin.Context) {
			response.Success(c, gin.H{"id": 2, "userName": "newuser"})
		})
		protected.GET("/user/:id", func(c *gin.Context) {
			id := c.Param("id")
			if id == "1" {
				response.Success(c, gin.H{"id": 1, "userName": "admin"})
			} else {
				response.NotFound(c, "用户不存在")
			}
		})
	}

	return r
}

func generateTestToken(t *testing.T) string {
	token, err := auth.GenerateToken("test-secret", 24, 1, "admin", 1000)
	assert.NoError(t, err)
	return token
}

func TestLoginAPI(t *testing.T) {
	r := setupTestRouter()

	body := map[string]string{
		"userName": "admin",
		"password": "123456",
	}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp response.Response
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, 200, resp.Code)
}

func TestGetUserAPI_Success(t *testing.T) {
	r := setupTestRouter()
	token := generateTestToken(t)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/user/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp response.Response
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, 200, resp.Code)
}

func TestGetUserAPI_NotFound(t *testing.T) {
	r := setupTestRouter()
	token := generateTestToken(t)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/user/999", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetUserAPI_Unauthorized(t *testing.T) {
	r := setupTestRouter()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/user/1", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestListUsersAPI(t *testing.T) {
	r := setupTestRouter()
	token := generateTestToken(t)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/user", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreateUserAPI(t *testing.T) {
	r := setupTestRouter()
	token := generateTestToken(t)

	body := map[string]string{
		"userName": "newuser",
		"password": "password123",
	}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/user", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
