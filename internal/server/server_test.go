package server

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func getTestImagePath() string {
	pwd, _ := os.Getwd()

	parent := filepath.Dir(pwd)
	return filepath.Join(parent, "image", "test_data", "rect.jpg")
}

func TestHandleImageUpload(t *testing.T) {
	gin.SetMode(gin.TestMode)

	uploadDir := "uploads"
	err := os.MkdirAll(uploadDir, 0755)
	assert.NoError(t, err)
	defer os.RemoveAll(uploadDir)

	imagePath := getTestImagePath()
	imageFile, err := os.Open(imagePath)
	if err != nil {
		t.Fatalf("Не удалось открыть тестовое изображение: %v", err)
	}
	defer imageFile.Close()

	tests := []struct {
		name           string
		setupRequest   func() (*http.Request, error)
		expectedStatus int
		expectedJSON   map[string]string
	}{
		{
			name: "Успешная загрузка JPEG",
			setupRequest: func() (*http.Request, error) {
				body := new(bytes.Buffer)
				writer := multipart.NewWriter(body)

				// Поле для файла
				part, err := writer.CreateFormFile("image", "rect.jpg")
				if err != nil {
					return nil, err
				}

				imageFile.Seek(0, 0)
				if _, err := io.Copy(part, imageFile); err != nil {
					return nil, err
				}

				if err := writer.Close(); err != nil {
					return nil, err
				}

				req := httptest.NewRequest("POST", "/upload", body)
				req.Header.Set("Content-Type", writer.FormDataContentType())
				return req, nil
			},
			expectedStatus: http.StatusOK,
			expectedJSON: map[string]string{
				"message":  "Файл успешно загружен",
				"filename": "rect.jpg",
			},
		},
		{
			name: "Неверный тип файла",
			setupRequest: func() (*http.Request, error) {
				body := new(bytes.Buffer)
				writer := multipart.NewWriter(body)

				part, err := writer.CreateFormFile("image", "test.txt")
				if err != nil {
					return nil, err
				}

				if _, err := part.Write([]byte("это не изображение")); err != nil {
					return nil, err
				}

				if err := writer.Close(); err != nil {
					return nil, err
				}

				req := httptest.NewRequest("POST", "/upload", body)
				req.Header.Set("Content-Type", writer.FormDataContentType())
				return req, nil
			},
			expectedStatus: http.StatusBadRequest,
			expectedJSON: map[string]string{
				"error": "Загруженный файл не является изображением",
			},
		},
		{
			name: "Отсутствует файл в запросе",
			setupRequest: func() (*http.Request, error) {
				req := httptest.NewRequest("POST", "/upload", nil)
				return req, nil
			},
			expectedStatus: http.StatusBadRequest,
			expectedJSON: map[string]string{
				"error": "Файл не найден в запросе",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := New()

			w := httptest.NewRecorder()

			req, err := tt.setupRequest()
			assert.NoError(t, err)

			server.router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]string
			err = json.NewDecoder(w.Body).Decode(&response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedJSON, response)

			if tt.expectedStatus == http.StatusOK {
				uploadedFile := filepath.Join("uploads", "rect.jpg")
				_, err := os.Stat(uploadedFile)
				assert.NoError(t, err, "Файл должен быть создан")
			}
		})
	}
}

func TestIsImageFile(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		expected    bool
	}{
		{"JPEG изображение", "image/jpeg", true},
		{"PNG изображение", "image/png", true},
		{"GIF изображение", "image/gif", true},
		{"BMP изображение", "image/bmp", true},
		{"WebP изображение", "image/webp", true},
		{"Текстовый файл", "text/plain", false},
		{"PDF файл", "application/pdf", false},
		{"Пустой тип", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := &multipart.FileHeader{
				Header: make(map[string][]string),
			}
			file.Header.Set("Content-Type", tt.contentType)

			result := isImageFile(file)
			assert.Equal(t, tt.expected, result)
		})
	}
}
