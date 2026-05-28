package admin

import (
	"github.com/gin-gonic/gin"
	"menu-recommend/internal/model"
	"menu-recommend/internal/pkg/errcode"
	"menu-recommend/internal/pkg/response"
	"menu-recommend/internal/service"

	"gorm.io/gorm"
)

type AuthHandler struct {
	authService *service.AuthService
	db          *gorm.DB
}

func NewAuthHandler(authService *service.AuthService, db *gorm.DB) *AuthHandler {
	return &AuthHandler{authService: authService, db: db}
}

type AdminLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req AdminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParam, "参数错误")
		return
	}

	var admin model.AdminUser
	if err := h.db.Where("username = ?", req.Username).First(&admin).Error; err != nil {
		response.Error(c, errcode.ErrAuth, "用户名或密码错误")
		return
	}

	if !h.authService.CheckPassword(req.Password, admin.PasswordHash) {
		response.Error(c, errcode.ErrAuth, "用户名或密码错误")
		return
	}

	token, err := h.authService.GenerateToken(admin.ID, "admin")
	if err != nil {
		response.Error(c, errcode.ErrServer, "生成token失败")
		return
	}

	response.Success(c, gin.H{
		"token": token,
		"user": gin.H{
			"id":       admin.ID,
			"username": admin.Username,
			"role":     admin.Role,
		},
	})
}
