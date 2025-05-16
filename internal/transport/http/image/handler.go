package handler

import (
	"context"
	"fmt"
	"image-upload-api/internal/domain/image"
	"image-upload-api/internal/infrastructure/storage"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

type ImageHandler struct {
	repo    image.ImageRepository
	storage *storage.LocalStorage
}

func NewImageHandler(repo image.ImageRepository, storage *storage.LocalStorage) *ImageHandler {
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

	// validate content type (only allow images)
	// uncomment code bellow if want to allow image only

	// contentType := header.Header.Get("Content-Type")
	// if !strings.HasPrefix(contentType, "image/") {
	//   http.Error(w, "Only image files are allowed", http.StatusBadRequest)
	//   return
	// }

	// read file data
	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read image file", http.StatusInternalServerError)
	}

	// parse image processing options from query params
	options := &storage.ImageOptions{}
	if width := r.URL.Query().Get("width"); width != "" {
		if w, err := strconv.Atoi(width); err == nil && w > 0 {
			options.Width = w
		}
	}
	if height := r.URL.Query().Get("height"); height != "" {
		if h, err := strconv.Atoi(height); err == nil && h > 0 {
			options.Height = h
		}
	}
	if quality := r.URL.Query().Get("quality"); quality != "" {
		if q, err := strconv.Atoi(quality); err == nil && q > 0 && q <= 100 {
			options.Quality = q
		}
	}
	if format := r.URL.Query().Get("format"); format != "" {
		switch strings.ToLower(format) {
		case "jpeg", "jpg", "png", "webp":
			options.Format = strings.ToLower(format)
		default:
			http.Error(w, "Unsupported format. Use jpeg, png, or webp", http.StatusBadRequest)
			return
		}
	}
	//generate unique filename
	filename := fmt.Sprintf("%s-%s", uuid.New().String(), header.Filename)

	// save file to localstorage
	path, err := h.storage.SaveFile(filename, data, options)
	if err != nil {
		http.Error(w, "failed to save image: "+err.Error(), http.StatusInternalServerError)
		return
	}

	//save image metadata to database
	image := &image.Image{
		ID:       uuid.New(),
		Filename: filepath.Base(path), // Use processed filename
		Path:     path,
	}
	if err := h.repo.Save(context.Background(), image); err != nil {
		http.Error(w, "failed to save image metadata", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, `{"message":"Image upload successfully","filename: %s"}`, filename)
}
