package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/kevin-ai/go-kevin/pkg/auth"
	"github.com/kevin-ai/go-kevin/pkg/response"
)

func JWTAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "未提供认证令牌")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, "认证格式错误")
			c.Abort()
			return
		}

		claims, err := auth.ParseToken(secret, parts[1])
		if err != nil {
			response.Unauthorized(c, "无效的认证令牌")
			c.Abort()
			return
		}

		c.Set("userId", claims.UserID)
		c.Set("userName", claims.UserName)
		c.Set("tenantId", claims.TenantID)

		c.Next()
	}
}

func TenantContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.GetInt64("tenantId")
		if tenantID == 0 {
			response.Forbidden(c, "租户信息缺失")
			c.Abort()
			return
		}
		c.Next()
	}
}

func RBAC(requiredPermission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetInt64("userId")
		if userID == 0 {
			response.Forbidden(c, "权限不足")
			c.Abort()
			return
		}
		c.Next()
	}
}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
