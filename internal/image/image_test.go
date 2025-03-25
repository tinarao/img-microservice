package image

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// getTestImagePath возвращает абсолютный путь к тестовому изображению
func getTestImagePath() string {
	// Получаем текущую директорию
	pwd, _ := os.Getwd()
	// Формируем путь к тестовому изображению
	return filepath.Join(pwd, "test_data", "rect.jpg")
}

func TestCalculateImg(t *testing.T) {
	imagePath := getTestImagePath()
	t.Logf("Путь к файлу: %s", imagePath)

	// Проверяем существование файла
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		t.Fatal(err)
	}

	pic := &Image{
		Path: imagePath,
		Name: "Test",	// h: 600, w: 800 init
	}

	size, err := pic.CalculateCompressedSize()
	if err != nil {
		t.Fatal(err)
	}
	
	fmt.Printf("h: %d\tw: %d\n", size.Height, size.Width)
	if size.Width != MAX_W || size.Height != 300 {
		t.Fatalf("incorrect aspect ratio")
	}
}