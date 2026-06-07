package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSuccessResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	Success(c, map[string]string{"name": "test"})

	assert.Equal(t, http.StatusOK, w.Code)

	var resp Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.Code)
	assert.Equal(t, "操作成功", resp.Message)
	assert.NotNil(t, resp.Data)
}

func TestErrorResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	Error(c, http.StatusBadRequest, "参数错误")

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.Code)
	assert.Equal(t, "参数错误", resp.Message)
}

func TestPageResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	items := []string{"a", "b", "c"}
	Page(c, items, 100, 1, 10)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.Code)

	data := resp.Data.(map[string]interface{})
	assert.Equal(t, float64(100), data["total"])
	assert.Equal(t, float64(1), data["page"])
}
