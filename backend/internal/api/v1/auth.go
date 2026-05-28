package v1

import (
	"github.com/gin-gonic/gin"
	"menu-recommend/internal/pkg/errcode"
	"menu-recommend/internal/pkg/response"
	"menu-recommend/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParam, "参数错误")
		return
	}

	user, err := h.authService.Register(req.Username, req.Password, req.Email, req.Phone)
	if err != nil {
		response.Error(c, errcode.ErrDuplicate, "用户名已存在")
		return
	}

	token, err := h.authService.GenerateToken(user.ID, "")
	if err != nil {
		response.Error(c, errcode.ErrServer, "生成token失败")
		return
	}

	response.Success(c, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"nickname": user.Nickname,
		},
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParam, "参数错误")
		return
	}

	user, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		response.Error(c, errcode.ErrAuth, "用户名或密码错误")
		return
	}

	token, err := h.authService.GenerateToken(user.ID, "")
	if err != nil {
		response.Error(c, errcode.ErrServer, "生成token失败")
		return
	}

	response.Success(c, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"nickname": user.Nickname,
			"avatar":   user.Avatar,
		},
	})
}
