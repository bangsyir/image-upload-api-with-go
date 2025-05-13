package storage

import (
	"fmt"
	"os"
	"path/filepath"
)

type LocaStorage struct {
	uploadDir string
}

func NewLocalStorage(uploadDir string) *LocaStorage {
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		fmt.Printf("warning: failed to create upload directory : %v\n", err)
	}
	return &LocaStorage{uploadDir: uploadDir}
}

func (s *LocaStorage) SaveFile(fileName string, data []byte) (string, error) {
	path := filepath.Join(s.uploadDir, fileName)
	if err := os.WriteFile(path, data, 0644); err != nil {
		return "", fmt.Errorf("failed to save file : %w", err)
	}
	return path, nil
}
