package image

import "context"

type ImageRepository interface {
	Save(ctx context.Context, image *Image) error
}
