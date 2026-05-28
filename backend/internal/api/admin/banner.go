package admin

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"menu-recommend/internal/model"
	"menu-recommend/internal/pkg/errcode"
	"menu-recommend/internal/pkg/response"
	"menu-recommend/internal/service"
)

type BannerHandler struct {
	bannerService *service.BannerService
}

func NewBannerHandler(bannerService *service.BannerService) *BannerHandler {
	return &BannerHandler{bannerService: bannerService}
}

func (h *BannerHandler) List(c *gin.Context) {
	banners, err := h.bannerService.GetAll()
	if err != nil {
		response.Error(c, errcode.ErrServer, "查询失败")
		return
	}
	response.Success(c, banners)
}

func (h *BannerHandler) Create(c *gin.Context) {
	var banner model.Banner
	if err := c.ShouldBindJSON(&banner); err != nil {
		response.Error(c, errcode.ErrParam, "参数错误")
		return
	}
	if err := h.bannerService.Create(&banner); err != nil {
		response.Error(c, errcode.ErrServer, "创建失败")
		return
	}
	response.Success(c, banner)
}

func (h *BannerHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var banner model.Banner
	if err := c.ShouldBindJSON(&banner); err != nil {
		response.Error(c, errcode.ErrParam, "参数错误")
		return
	}
	banner.ID = uint(id)
	if err := h.bannerService.Update(&banner); err != nil {
		response.Error(c, errcode.ErrServer, "更新失败")
		return
	}
	response.Success(c, banner)
}

func (h *BannerHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.bannerService.Delete(uint(id)); err != nil {
		response.Error(c, errcode.ErrServer, "删除失败")
		return
	}
	response.Success(c, nil)
}
