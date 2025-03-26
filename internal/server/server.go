package server

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router   *gin.Engine
	ApiGroup *gin.RouterGroup
}

func New() *Server {
	r := gin.Default()
	server := &Server{
		router:   r,
		ApiGroup: r.Group("api"),
	}

	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	s.ApiGroup.GET("/hc", func(c *gin.Context) {
		c.String(http.StatusOK, "Ok")
	})

	s.router.POST("/upload", s.handleImageUpload)
}

func (s *Server) Run(addr string) error {
	return s.router.Run(addr)
}

func (s *Server) handleImageUpload(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Файл не найден в запросе",
		})
		return
	}

	fmt.Printf("Загрузка файла: %s, Content-Type: %s\n", file.Filename, file.Header.Get("Content-Type"))

	if !isImageFile(file) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Загруженный файл не является изображением",
		})
		return
	}

	if err := os.MkdirAll("uploads", 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Ошибка при создании директории: %s", err.Error()),
		})
		return
	}

	filename := filepath.Join("uploads", file.Filename)
	if err := c.SaveUploadedFile(file, filename); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Ошибка при сохранении файла: %s", err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Файл успешно загружен",
		"filename": file.Filename,
	})
}

func isImageFile(file *multipart.FileHeader) bool {
	ext := filepath.Ext(file.Filename)
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp":
		return true
	default:
		contentType := file.Header.Get("Content-Type")
		switch contentType {
		case "image/jpeg", "image/png", "image/gif", "image/bmp", "image/webp":
			return true
		}
		return false
	}
}
