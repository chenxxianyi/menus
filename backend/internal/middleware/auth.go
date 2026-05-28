package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"menu-recommend/internal/pkg/errcode"
	"menu-recommend/internal/pkg/response"
	"menu-recommend/internal/service"
)

func AuthMiddleware(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		token = strings.TrimPrefix(token, "Bearer ")
		claims, err := authService.ParseToken(token)
		if err != nil {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Next()
	}
}

func AdminMiddleware(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		token = strings.TrimPrefix(token, "Bearer ")
		claims, err := authService.ParseToken(token)
		if err != nil || claims.Role != "admin" {
			response.Forbidden(c)
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Next()
	}
}

func OptionalAuth(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.Next()
			return
		}

		token = strings.TrimPrefix(token, "Bearer ")
		claims, err := authService.ParseToken(token)
		if err == nil {
			c.Set("user_id", claims.UserID)
		}
		c.Next()
	}
}

func GetUserID(c *gin.Context) uint {
	if v, ok := c.Get("user_id"); ok {
		if id, ok := v.(uint); ok {
			return id
		}
	}
	return 0
}

func ResponseError(c *gin.Context, code int) {
	switch code {
	case errcode.ErrParam:
		response.Error(c, code, "参数错误")
	case errcode.ErrNotFound:
		response.Error(c, code, "资源不存在")
	case errcode.ErrDuplicate:
		response.Error(c, code, "数据已存在")
	default:
		response.Error(c, errcode.ErrServer, "服务器内部错误")
	}
}
