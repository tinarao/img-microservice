package main

import (
	"go-image-processor/internal/image"
	"go-image-processor/internal/server"
	"log"
)

func main() {
	image.Init()

	image := &image.Image{
		Path: image.GetTestImagePath(),
		Name: "BikeTest",
	}

	if err := image.Compress(180); err != nil {
		log.Fatal(err)
	}

	srv := server.New()
	if err := srv.Run(":8080"); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
