package storage

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
)

type LocalStorage struct {
	uploadDir string
}

func NewLocalStorage(uploadDir string) *LocalStorage {
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		fmt.Printf("warning: failed to create upload directory : %v\n", err)
	}
	return &LocalStorage{uploadDir: uploadDir}
}

type ImageOptions struct {
	Width   int
	Height  int
	Quality int    // 1-100 for JPEG compression
	Format  string // "jpeg", "png"
}

func (s *LocalStorage) SaveFile(filename string, data []byte, options *ImageOptions) (string, error) {
	// decode the image
	img, err := imaging.Decode(bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %w", err)
	}

	// Apply resize if specified
	if options.Width > 0 || options.Height > 0 {
		img = imaging.Resize(img, options.Width, options.Height, imaging.Lanczos)
	}

	// adjust filename extension based on output format
	ext := strings.ToLower(options.Format)
	if ext == "" {
		ext = filepath.Ext(filename)[1:] // keep original extension if not specified
	}

	filename = strings.TrimSuffix(filename, filepath.Ext(filename)) + "." + ext

	// prepare buffer for processed image
	var buf bytes.Buffer

	switch ext {
	case "jpeg", "jpg":
		quality := 80
		if options.Quality > 0 && options.Quality <= 100 {
			quality = options.Quality
		}
		if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: quality}); err != nil {
			return "", fmt.Errorf("failed to encode JPEG: %w", err)
		}
	case "png":
		if err := png.Encode(&buf, img); err != nil {
			return "", fmt.Errorf("failed to encode PNG: %w", err)
		}
	default:
		return "", fmt.Errorf("Unsupported formatL %s", ext)
	}

	// save processed image
	path := filepath.Join(s.uploadDir, filename)
	if err := os.WriteFile(path, buf.Bytes(), 0644); err != nil {
		return "", fmt.Errorf("failed to save file : %w", err)
	}
	return path, nil
}
