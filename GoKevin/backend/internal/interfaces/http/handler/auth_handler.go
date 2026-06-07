package handler

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"github.com/kevin-ai/go-kevin/internal/domain/repository"
	"github.com/kevin-ai/go-kevin/pkg/auth"
	"github.com/kevin-ai/go-kevin/pkg/response"
)

type AuthHandler struct {
	userRepo   repository.UserRepository
	jwtSecret  string
	expireHour int
}

func NewAuthHandler(userRepo repository.UserRepository, jwtSecret string, expireHour int) *AuthHandler {
	return &AuthHandler{
		userRepo:   userRepo,
		jwtSecret:  jwtSecret,
		expireHour: expireHour,
	}
}

type LoginRequest struct {
	UserName string `json:"userName" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token     string `json:"token"`
	ExpiresIn int    `json:"expiresIn"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	user, err := h.userRepo.GetByUserName(c.Request.Context(), req.UserName)
	if err != nil {
		response.InternalError(c, "系统错误")
		return
	}
	if user == nil {
		response.Unauthorized(c, "用户名或密码错误")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		response.Unauthorized(c, "用户名或密码错误")
		return
	}

	token, err := auth.GenerateToken(h.jwtSecret, h.expireHour, user.ID, user.UserName, user.TenantID)
	if err != nil {
		response.InternalError(c, "生成 Token 失败")
		return
	}

	response.Success(c, LoginResponse{
		Token:     token,
		ExpiresIn: h.expireHour * 3600,
	})
}

func (h *AuthHandler) Register(c *gin.Context) {
	// Placeholder for register
	response.Success(c, nil)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		response.Unauthorized(c, "未提供 Token")
		return
	}

	tokenString := authHeader[7:]

	newToken, err := auth.RefreshToken(h.jwtSecret, h.expireHour, tokenString)
	if err != nil {
		response.Unauthorized(c, "Token 刷新失败: "+err.Error())
		return
	}

	response.Success(c, LoginResponse{
		Token:     newToken,
		ExpiresIn: h.expireHour * 3600,
	})
}
