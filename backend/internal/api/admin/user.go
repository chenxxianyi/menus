package admin

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"menu-recommend/internal/pkg/errcode"
	"menu-recommend/internal/pkg/response"
	"menu-recommend/internal/service"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) List(c *gin.Context) {
	keyword := c.Query("keyword")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	users, total, err := h.userService.ListUsers(keyword, page, pageSize)
	if err != nil {
		response.Error(c, errcode.ErrServer, "查询失败")
		return
	}
	response.SuccessPage(c, users, total, page, pageSize)
}

func (h *UserHandler) Detail(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	user, err := h.userService.GetUserByID(uint(id))
	if err != nil {
		response.Error(c, errcode.ErrNotFound, "用户不存在")
		return
	}
	response.Success(c, user)
}
