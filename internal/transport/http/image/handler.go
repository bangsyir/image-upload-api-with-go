package handler

import (
	"context"
	"fmt"
	"image-upload-api/internal/domain/image"
	"image-upload-api/internal/infrastructure/storage"
	"io"
	"net/http"

	"github.com/google/uuid"
)

type ImageHandler struct {
	repo    image.ImageRepository
	storage *storage.LocaStorage
}

func NewImageHandler(repo image.ImageRepository, storage *storage.LocaStorage) *ImageHandler {
	return &ImageHandler{repo: repo, storage: storage}
}

func (h *ImageHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
	// parse multipart/form for (10mb) max

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Failed to get image file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	//read file data
	data, err := io.ReadAll(file)

	if err != nil {
		http.Error(w, "Failed to image file", http.StatusInternalServerError)
		return
	}

	//generate unique filename
	filename := fmt.Sprintf("%s-%s", uuid.New().String(), header.Filename)

	// save file to localstorage
	path, err := h.storage.SaveFile(filename, data)
	if err != nil {
		http.Error(w, "failed to save image", http.StatusInternalServerError)
		return
	}

	//save image metadata to database
	image := &image.Image{
		ID:       uuid.New(),
		Filename: filename,
		Path:     path,
	}
	if err := h.repo.Save(context.Background(), image); err != nil {
		http.Error(w, "failed to save image metadata", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Image upload successfully: %s", filename)
}
