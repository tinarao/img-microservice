package image

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
)

type Image struct {
	Path string
	Name string
}

type Size struct {
	Width  int
	Height int
}

const MAX_W = 400

func Init() {
	if err := os.MkdirAll("uploads", 0755); err != nil {
		log.Fatalf("Ошибка при создании директории uploads: %v", err)
	}
}

func GetTestImagePath() string {
	_, currentFile, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(currentFile)
	return filepath.Join(currentDir, "test_data", "rect.jpg")
}

func (i *Image) CalculateCompressedSize(toWidth int) (*Size, error) {
	img, err := imgio.Open(i.Path)
	if err != nil {
		return nil, err
	}

	size := &Size{
		Width:  toWidth,
		Height: img.Bounds().Dy() / (img.Bounds().Dx() / toWidth),
	}

	return size, nil
}

func (i *Image) GetDirToSave() string {
	resultsDir := filepath.Join("results")

	if err := os.MkdirAll(resultsDir, 0755); err != nil {
		log.Printf("Ошибка при создании директории %s: %v", resultsDir, err)
	}

	fp := filepath.Join("results", fmt.Sprintf("%s_compressed.jpg", i.Name))

	return fp
}

func (i *Image) Compress(toWidth uint) error {
	if toWidth > MAX_W {
		return fmt.Errorf("failed to compress: max width is %d", MAX_W)
	}

	img, err := imgio.Open(i.Path)
	if err != nil {
		return fmt.Errorf("failed to open image: %v", err)
	}

	size, err := i.CalculateCompressedSize(int(toWidth))
	if err != nil {
		return fmt.Errorf("failed to calculate size: %v", err)
	}

	if err := os.MkdirAll("results", 0755); err != nil {
		return fmt.Errorf("failed to create results directory: %v", err)
	}

	result := transform.Resize(img, size.Width, size.Height, transform.Linear)
	err = imgio.Save(i.GetDirToSave(), result, imgio.JPEGEncoder(80))
	if err != nil {
		return fmt.Errorf("failed to save compressed image: %v", err)
	}

	fmt.Printf("File saved successfully\n")
	return nil
}
