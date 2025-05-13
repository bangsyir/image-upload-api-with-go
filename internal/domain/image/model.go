package image

import "github.com/google/uuid"

type Image struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	Filename string    `gorm:"not null"`
	Path     string    `gorm:"not null"`
}
