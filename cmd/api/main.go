package main

import (
	"fmt"
	"image-upload-api/internal/config"
	db "image-upload-api/internal/infrastructure/database"
	storage "image-upload-api/internal/infrastructure/storage"
	repository "image-upload-api/internal/repository/image"
	handler "image-upload-api/internal/transport/http/image"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg := config.NewConfig()

	// initialize database
	db, err := db.NewDatabase(cfg.DatabasePath)
	if err != nil {
		log.Fatal("failed to initialize database: %w", err)
	}

	// initialize storage
	imageStorage := storage.NewLocalStorage(cfg.UploadDir)
	// initialize repository
	imageRepo := repository.NewImageRepository(db)
	imageHandler := handler.NewImageHandler(imageRepo, imageStorage)

	// route
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/api/images", imageHandler.UploadImage)

	// start server
	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("Server starting on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
