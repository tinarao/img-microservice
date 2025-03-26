package main

import (
	"go-image-processor/internal/db"
	"go-image-processor/internal/image"
	"go-image-processor/internal/oauth"
	"go-image-processor/internal/server"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("failed to load .env file: %s", err.Error())
	}

	srv := server.New()

	db.Init()
	image.Init()
	oauth.Init(srv.ApiGroup)

	if err := srv.Run(":8080"); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
