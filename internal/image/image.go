package image

import (
	"fmt"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
)

type Image struct {
	Path string
	Name string
}

type Size struct {
	Width int
	Height int
}

const MAX_W = 400

func (i *Image) CalculateCompressedSize() (*Size, error) {
	img, err := imgio.Open(i.Path)
	if err != nil {
		return nil, err
	}

	newW := MAX_W
	originalW := img.Bounds().Dx()
	originalH := img.Bounds().Dy()
	
	// Сохраняем пропорции
	newH := (originalH * newW) / originalW

	size := &Size {
		Width: newW,
		Height: newH,
	}

	return size, nil
}

func (i *Image) Compress() error {
	img, err := imgio.Open(i.Path)
	if err != nil {
		return err
	}
	
	size, err := i.CalculateCompressedSize()
	if err != nil {
		fmt.Printf("failed to calculate new pic size: %s\n", err.Error())
	}

	transform.Resize(img, size.Width, size.Height, transform.Linear)
	err = imgio.Save(fmt.Sprintf("./%s-resized", i.Name), img, imgio.JPEGEncoder(80))
	if err != nil {
		return err
	}

	return nil
}
