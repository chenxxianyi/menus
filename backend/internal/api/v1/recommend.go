package v1

import (
	"github.com/gin-gonic/gin"
	"menu-recommend/internal/pkg/errcode"
	"menu-recommend/internal/pkg/response"
	"menu-recommend/internal/service"
)

type RecommendHandler struct {
	recommendService *service.RecommendService
}

func NewRecommendHandler(recommendService *service.RecommendService) *RecommendHandler {
	return &RecommendHandler{recommendService: recommendService}
}

func (h *RecommendHandler) Menu(c *gin.Context) {
	var params service.RecommendParams
	if err := c.ShouldBindJSON(&params); err != nil {
		response.Error(c, errcode.ErrParam, "参数错误")
		return
	}

	result, err := h.recommendService.RecommendMenu(&params)
	if err != nil {
		response.Error(c, errcode.ErrServer, "推荐失败")
		return
	}
	response.Success(c, result)
}

type ByIngredientsRequest struct {
	Ingredients []string `json:"ingredients" binding:"required"`
}

func (h *RecommendHandler) ByIngredients(c *gin.Context) {
	var req ByIngredientsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParam, "参数错误")
		return
	}

	result, err := h.recommendService.RecommendByIngredients(req.Ingredients)
	if err != nil {
		response.Error(c, errcode.ErrServer, "推荐失败")
		return
	}
	response.Success(c, result)
}

func (h *RecommendHandler) WeekMenu(c *gin.Context) {
	var params service.RecommendParams
	if err := c.ShouldBindJSON(&params); err != nil {
		response.Error(c, errcode.ErrParam, "参数错误")
		return
	}

	result, err := h.recommendService.GenerateWeekMenu(&params)
	if err != nil {
		response.Error(c, errcode.ErrServer, "生成失败")
		return
	}
	response.Success(c, result)
}
