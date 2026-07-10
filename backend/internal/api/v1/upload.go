package v1

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
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
	if h.cfg.MaxSize > 0 {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, h.cfg.MaxSize+1024)
	}
	file, err := c.FormFile("file")
	if err != nil {
		response.Error(c, errcode.ErrParam, "请选择图片")
		return
	}
	if h.cfg.MaxSize > 0 && file.Size > h.cfg.MaxSize {
		response.Error(c, errcode.ErrParam, "图片过大")
		return
	}
	imageType, err := inspectImage(file)
	if err != nil {
		response.Error(c, errcode.ErrParam, "仅支持 jpg、png、webp 图片")
		return
	}

	dir := filepath.Join(h.cfg.Dir, "avatar")
	if err := os.MkdirAll(dir, 0755); err != nil {
		response.Error(c, errcode.ErrServer, "创建上传目录失败")
		return
	}

	name := fmt.Sprintf("%d%s", time.Now().UnixNano(), imageType.extension)
	dst := filepath.Join(dir, name)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		response.Error(c, errcode.ErrServer, "上传失败")
		return
	}

	response.Success(c, gin.H{
		"url": "/uploads/avatar/" + name,
	})
}

type verifiedImageType struct {
	extension string
}

const maxImagePixels = 20_000_000

func inspectImage(file *multipart.FileHeader) (verifiedImageType, error) {
	if file == nil {
		return verifiedImageType{}, fmt.Errorf("file is nil")
	}
	source, err := file.Open()
	if err != nil {
		return verifiedImageType{}, err
	}
	defer source.Close()

	data := make([]byte, 512)
	read, err := source.Read(data)
	if err != nil && read == 0 {
		return verifiedImageType{}, err
	}
	data = data[:read]
	contentType := http.DetectContentType(data)
	if contentType == "image/jpeg" || contentType == "image/png" {
		if _, err := source.Seek(0, 0); err != nil {
			return verifiedImageType{}, err
		}
		config, _, err := image.DecodeConfig(source)
		if err != nil || config.Width <= 0 || config.Height <= 0 || config.Width*config.Height > maxImagePixels {
			return verifiedImageType{}, fmt.Errorf("invalid image content")
		}
		if contentType == "image/jpeg" {
			return verifiedImageType{extension: ".jpg"}, nil
		}
		return verifiedImageType{extension: ".png"}, nil
	}

	// WebP does not have a standard-library decoder. Validate its RIFF/WebP
	// signature and use a server-generated extension instead of the user name.
	if len(data) >= 12 && bytes.Equal(data[:4], []byte("RIFF")) && bytes.Equal(data[8:12], []byte("WEBP")) {
		return verifiedImageType{extension: ".webp"}, nil
	}
	return verifiedImageType{}, fmt.Errorf("unsupported image content")
}
