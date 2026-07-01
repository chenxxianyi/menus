package v1

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"menu-recommend/internal/middleware"
	"menu-recommend/internal/model"
	"menu-recommend/internal/pkg/errcode"
	"menu-recommend/internal/pkg/response"
	"menu-recommend/internal/service"
)

type ShoppingHandler struct {
	shoppingService *service.ShoppingService
}

func NewShoppingHandler(shoppingService *service.ShoppingService) *ShoppingHandler {
	return &ShoppingHandler{shoppingService: shoppingService}
}

func (h *ShoppingHandler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)
	lists, err := h.shoppingService.GetLists(userID)
	if err != nil {
		response.Error(c, errcode.ErrServer, "查询失败")
		return
	}
	response.Success(c, lists)
}

type CreateShoppingRequest struct {
	Name      string     `json:"name"`
	ItemsJSON model.JSON `json:"items_json"`
}

func (h *ShoppingHandler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req CreateShoppingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParam, "参数错误")
		return
	}

	list := &model.ShoppingList{
		UserID:    userID,
		Name:      req.Name,
		ItemsJSON: req.ItemsJSON,
	}

	if err := h.shoppingService.CreateList(list); err != nil {
		response.Error(c, errcode.ErrServer, "创建失败")
		return
	}
	response.Success(c, list)
}

type UpdateShoppingRequest struct {
	Name         string     `json:"name"`
	ItemsJSON    model.JSON `json:"items_json"`
	hasItemsJSON bool
}

func (r *UpdateShoppingRequest) UnmarshalJSON(data []byte) error {
	var raw struct {
		Name      string           `json:"name"`
		ItemsJSON *json.RawMessage `json:"items_json"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	r.Name = raw.Name
	if raw.ItemsJSON != nil {
		r.hasItemsJSON = true
		r.ItemsJSON = append(r.ItemsJSON[:0], (*raw.ItemsJSON)...)
	}
	return nil
}

func (h *ShoppingHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrParam, "无效的ID")
		return
	}

	existing, err := h.shoppingService.GetListByID(uint(id))
	if err != nil {
		response.Error(c, errcode.ErrNotFound, "清单不存在")
		return
	}

	var req UpdateShoppingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParam, "参数错误")
		return
	}

	if req.Name != "" {
		existing.Name = req.Name
	}
	if req.hasItemsJSON {
		existing.ItemsJSON = req.ItemsJSON
	}

	if err := h.shoppingService.UpdateList(existing); err != nil {
		response.Error(c, errcode.ErrServer, "更新失败")
		return
	}
	response.Success(c, existing)
}

func (h *ShoppingHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrParam, "无效的ID")
		return
	}

	if err := h.shoppingService.DeleteList(uint(id)); err != nil {
		response.Error(c, errcode.ErrServer, "删除失败")
		return
	}
	response.Success(c, nil)
}

type DeleteShoppingItemsRequest struct {
	Indices []int `json:"indices" binding:"required,min=1"`
}

type DeleteShoppingItemsResponse struct {
	ListID       uint       `json:"list_id"`
	DeletedCount int        `json:"deleted_count"`
	ItemsJSON    model.JSON `json:"items_json"`
}

type GenerateShoppingByDishRequest struct {
	DishName string `json:"dish_name" binding:"required"`
	Preview  bool   `json:"preview"`
}

func (h *ShoppingHandler) DeleteItems(c *gin.Context) {
	listID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || listID == 0 {
		response.Error(c, errcode.ErrParam, "无效的清单ID")
		return
	}

	var req DeleteShoppingItemsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParam, "请选择要删除的食材")
		return
	}

	userID := middleware.GetUserID(c)
	list, deletedCount, err := h.shoppingService.DeleteItems(userID, uint(listID), req.Indices)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			response.Error(c, errcode.ErrNotFound, "购物清单不存在")
		case errors.Is(err, service.ErrInvalidShoppingItemIndices):
			response.Error(c, errcode.ErrParam, "删除的食材位置无效")
		default:
			response.Error(c, errcode.ErrServer, "删除食材失败")
		}
		return
	}

	response.Success(c, DeleteShoppingItemsResponse{
		ListID:       list.ID,
		DeletedCount: deletedCount,
		ItemsJSON:    list.ItemsJSON,
	})
}

func (h *ShoppingHandler) GenerateByDish(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req GenerateShoppingByDishRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParam, "请输入想吃的菜品")
		return
	}

	result, err := h.shoppingService.GenerateFromDish(userID, req.DishName, req.Preview)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			response.Error(c, errcode.ErrNotFound, "没有匹配到真实菜谱，请先在后台维护该菜品")
		case errors.Is(err, service.ErrRecipeIngredientsEmpty):
			response.Error(c, errcode.ErrParam, "匹配到的菜谱还没有维护食材数据")
		default:
			response.Error(c, errcode.ErrServer, "生成采购清单失败")
		}
		return
	}

	response.Success(c, gin.H{
		"list":    result.List,
		"recipe":  result.Recipe,
		"items":   result.Items,
		"preview": req.Preview,
	})
}

func (h *ShoppingHandler) GenerateByAI(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req GenerateShoppingByDishRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParam, "请输入想吃的菜品")
		return
	}

	result, err := h.shoppingService.GenerateFromDishByAI(c.Request.Context(), userID, req.DishName, req.Preview)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrAIConfigMissing):
			response.Error(c, errcode.ErrParam, "AI 未配置，请先配置 AI_BASE_URL、AI_API_KEY 和 AI_MODEL")
		case errors.Is(err, service.ErrAIInvalidResponse):
			response.Error(c, errcode.ErrServer, "AI 返回的采购清单格式无效")
		default:
			response.Error(c, errcode.ErrServer, "AI 生成采购清单失败")
		}
		return
	}

	response.Success(c, gin.H{
		"list":    result.List,
		"recipe":  result.Recipe,
		"items":   result.Items,
		"source":  "ai",
		"preview": req.Preview,
	})
}
