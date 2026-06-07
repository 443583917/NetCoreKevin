package router

import (
	"github.com/gin-gonic/gin"

	"github.com/kevin-ai/go-kevin/internal/interfaces/http/handler"
	"github.com/kevin-ai/go-kevin/internal/interfaces/http/middleware"
)

func RegisterRoutes(
	r *gin.Engine,
	userHandler *handler.UserHandler,
	authHandler *handler.AuthHandler,
) {
	v1 := r.Group("/api/v1")

	// Public routes
	public := v1.Group("")
	{
		public.POST("/auth/login", func(c *gin.Context) {
			if authHandler != nil {
				authHandler.Login(c)
			}
		})
		public.POST("/auth/register", func(c *gin.Context) {
			if authHandler != nil {
				authHandler.Register(c)
			}
		})
		public.POST("/auth/refresh", func(c *gin.Context) {
			if authHandler != nil {
				authHandler.RefreshToken(c)
			}
		})
	}

	// Protected routes
	protected := v1.Group("")
	protected.Use(middleware.JWTAuth("your-secret-key"))
	protected.Use(middleware.TenantContext())
	{
		// User management
		users := protected.Group("/user")
		{
			users.GET("", func(c *gin.Context) {
				if userHandler != nil {
					userHandler.List(c)
				}
			})
			users.GET("/:id", func(c *gin.Context) {
				if userHandler != nil {
					userHandler.GetByID(c)
				}
			})
			users.POST("", func(c *gin.Context) {
				if userHandler != nil {
					userHandler.Create(c)
				}
			})
			users.PUT("/:id", func(c *gin.Context) {
				if userHandler != nil {
					userHandler.Update(c)
				}
			})
			users.DELETE("/:id", func(c *gin.Context) {
				if userHandler != nil {
					userHandler.Delete(c)
				}
			})
		}

		// AI Apps (placeholder)
		apps := protected.Group("/aiapps")
		{
			apps.GET("", func(c *gin.Context) {})
			apps.GET("/:id", func(c *gin.Context) {})
			apps.POST("", func(c *gin.Context) {})
			apps.PUT("/:id", func(c *gin.Context) {})
			apps.DELETE("/:id", func(c *gin.Context) {})
		}

		// Chat (placeholder)
		chat := protected.Group("/aichat")
		{
			chat.POST("/sessions", func(c *gin.Context) {})
			chat.GET("/sessions", func(c *gin.Context) {})
			chat.GET("/sessions/:id", func(c *gin.Context) {})
			chat.POST("/sessions/:id/messages", func(c *gin.Context) {})
			chat.GET("/sessions/:id/messages", func(c *gin.Context) {})
		}

		// Knowledge Base (placeholder)
		kb := protected.Group("/aikmss")
		{
			kb.GET("", func(c *gin.Context) {})
			kb.POST("", func(c *gin.Context) {})
			kb.GET("/:id", func(c *gin.Context) {})
			kb.DELETE("/:id", func(c *gin.Context) {})
			kb.POST("/:id/documents", func(c *gin.Context) {})
			kb.POST("/:id/query", func(c *gin.Context) {})
		}

		// File management (placeholder)
		files := protected.Group("/file")
		{
			files.POST("/upload", func(c *gin.Context) {})
			files.GET("/:id/download", func(c *gin.Context) {})
		}
	}
}
