package v1

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"menu-recommend/config"
	"menu-recommend/internal/pkg/errcode"
	"menu-recommend/internal/pkg/response"
)

type UploadHandler struct {
	cfg *config.UploadConfig
}

func NewUploadHandler(cfg *config.UploadConfig) *UploadHandler {
	return &UploadHandler{cfg: cfg}
}

func (h *UploadHandler) Avatar(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.Error(c, errcode.ErrParam, "请选择图片")
		return
	}
	if h.cfg.MaxSize > 0 && file.Size > h.cfg.MaxSize {
		response.Error(c, errcode.ErrParam, "图片过大")
		return
	}
	if !isAllowedImage(file) {
		response.Error(c, errcode.ErrParam, "仅支持 jpg、png、webp 图片")
		return
	}

	dir := filepath.Join(h.cfg.Dir, "avatar")
	if err := os.MkdirAll(dir, 0755); err != nil {
		response.Error(c, errcode.ErrServer, "创建上传目录失败")
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	name := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	dst := filepath.Join(dir, name)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		response.Error(c, errcode.ErrServer, "上传失败")
		return
	}

	response.Success(c, gin.H{
		"url": "/uploads/avatar/" + name,
	})
}

func isAllowedImage(file *multipart.FileHeader) bool {
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" {
		return false
	}

	contentType := strings.ToLower(file.Header.Get("Content-Type"))
	if contentType == "" {
		return true
	}
	return contentType == "image/jpeg" ||
		contentType == "image/png" ||
		contentType == "image/webp"
}
