package repository

import (
	"context"
	"image-upload-api/internal/domain/image"

	"gorm.io/gorm"
)

type ImageRepository struct {
	db *gorm.DB
}

func NewImageRepository(db *gorm.DB) *ImageRepository {
	return &ImageRepository{db: db}
}

func (r *ImageRepository) Save(ctx context.Context, image *image.Image) error {
	return r.db.WithContext(ctx).Create(image).Error
}
