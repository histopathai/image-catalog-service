package repository

import (
	"context"

	"github.com/histopathai/image-catalog-service/internal/models"
)

type ImageRepository interface {
	Read(ctx context.Context, imageID string) (*models.Image, error)
	Update(ctx context.Context, image *models.Image) error
	Delete(ctx context.Context, imageID string) error
	Filter(ctx context.Context, filter *models.ImageFilter) ([]*models.Image, error)
}
