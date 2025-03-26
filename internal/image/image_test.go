package image

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInit(t *testing.T) {
	// Проверяем создание директории uploads
	Init()
	if _, err := os.Stat("uploads"); os.IsNotExist(err) {
		t.Error("Директория uploads не была создана")
	}
}

func TestGetTestImagePath(t *testing.T) {
	path := GetTestImagePath()

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("Тестовое изображение не найдено по пути: %s", path)
	}

	expectedPath := filepath.Join("internal", "image", "test_data", "rect.jpg")
	if !strings.Contains(path, expectedPath) {
		t.Errorf("Путь должен содержать %s, получено %s", expectedPath, path)
	}
}

func TestCalculateCompressedSize(t *testing.T) {
	img := &Image{
		Path: GetTestImagePath(),
		Name: "test_image",
	}

	size, err := img.CalculateCompressedSize(200)
	if err != nil {
		t.Fatalf("Ошибка при расчете размера: %v", err)
	}

	if size.Width != 200 {
		t.Errorf("Ожидалась ширина 200, получено %d", size.Width)
	}
}

func TestGetDirToSave(t *testing.T) {
	img := &Image{
		Name: "test_image",
	}

	path := img.GetDirToSave()
	expectedPath := filepath.Join("results", "test_image_compressed.jpg")
	if path != expectedPath {
		t.Errorf("Ожидался путь %s, получено %s", expectedPath, path)
	}

	if _, err := os.Stat("results"); os.IsNotExist(err) {
		t.Error("Директория results не была создана")
	}
}

func TestCompress(t *testing.T) {
	img := &Image{
		Path: GetTestImagePath(),
		Name: "test_image",
	}

	err := img.Compress(200)
	if err != nil {
		t.Fatalf("Ошибка при сжатии изображения: %v", err)
	}

	outputPath := img.GetDirToSave()
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Errorf("Сжатое изображение не было создано по пути: %s", outputPath)
	}

	err = img.Compress(MAX_W + 1)
	if err == nil {
		t.Error("Ожидалась ошибка при превышении максимальной ширины")
	}

	os.RemoveAll("results")
}
