package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type PageData struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{Code: 0, Message: "success", Data: data})
}

func SuccessPage(c *gin.Context, list interface{}, total int64, page, pageSize int) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    PageData{List: list, Total: total, Page: page, PageSize: pageSize},
	})
}

func Error(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, Response{Code: code, Message: message, Data: nil})
}

func Unauthorized(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, Response{Code: 401, Message: "Unauthorized", Data: nil})
}

func Forbidden(c *gin.Context) {
	c.JSON(http.StatusForbidden, Response{Code: 403, Message: "Forbidden", Data: nil})
}
