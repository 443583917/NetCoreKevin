package swagger

import (
	"github.com/gin-gonic/gin"
)

// SwaggerInfo represents Swagger information
type SwaggerInfo struct {
	Title       string
	Description string
	Version     string
	Host        string
	BasePath    string
}

// DefaultSwaggerInfo returns default Swagger info
func DefaultSwaggerInfo() *SwaggerInfo {
	return &SwaggerInfo{
		Title:       "Go Kevin API",
		Description: "NetCoreKevin Go Backend API Documentation",
		Version:     "1.0",
		Host:        "localhost:9901",
		BasePath:    "/api/v1",
	}
}

// SetupSwagger sets up Swagger UI routes
func SetupSwagger(r *gin.Engine, info *SwaggerInfo) {
	if info == nil {
		info = DefaultSwaggerInfo()
	}

	// Swagger UI
	r.GET("/swagger/*any", func(c *gin.Context) {
		c.Header("Content-Type", "text/html")
		c.String(200, getSwaggerHTML(info))
	})

	// Swagger JSON
	r.GET("/swagger.json", func(c *gin.Context) {
		c.JSON(200, getSwaggerJSON(info))
	})
}

func getSwaggerHTML(info *SwaggerInfo) string {
	return `<!DOCTYPE html>
<html>
<head>
    <title>` + info.Title + `</title>
    <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
    <script>
        SwaggerUIBundle({
            url: "/swagger.json",
            dom_id: '#swagger-ui',
        })
    </script>
</body>
</html>`
}

func getSwaggerJSON(info *SwaggerInfo) map[string]interface{} {
	return map[string]interface{}{
		"openapi": "3.0.0",
		"info": map[string]interface{}{
			"title":       info.Title,
			"description": info.Description,
			"version":     info.Version,
		},
		"host":     info.Host,
		"basePath": info.BasePath,
		"paths": map[string]interface{}{
			"/auth/login": map[string]interface{}{
				"post": map[string]interface{}{
					"summary": "User login",
					"requestBody": map[string]interface{}{
						"content": map[string]interface{}{
							"application/json": map[string]interface{}{
								"schema": map[string]interface{}{
									"type": "object",
									"properties": map[string]interface{}{
										"userName": map[string]interface{}{"type": "string"},
										"password": map[string]interface{}{"type": "string"},
									},
								},
							},
						},
					},
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "Success",
						},
					},
				},
			},
		},
	}
}
