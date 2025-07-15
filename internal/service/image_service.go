package service

import (
	"context"
	"fmt"
	"time"

	"github.com/histopathai/image-catalog-service/config"
	"github.com/histopathai/image-catalog-service/internal/models"
	"github.com/histopathai/image-catalog-service/internal/repository"
)

// ImageService provides methods to manage images in the catalog.
type ImageService struct {
	repo repository.ImageRepository
	cfg  *config.Config
}

// NewImageService creates a new ImageService instance.
func NewImageService(repo repository.ImageRepository, cfg *config.Config) *ImageService {
	return &ImageService{
		repo: repo,
		cfg:  cfg,
	}
}

func (s *ImageService) GetImage(ctx context.Context, imageID string) (*models.Image, error) {
	image, err := s.repo.Read(ctx, imageID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve image: %w", err)
	}
	return image, nil
}

// UpdateImage updates an existing image record.
func (s *ImageService) UpdateImage(ctx context.Context, imageID string, updateRequest *models.ImageUpdateRequest) (*models.Image, error) {
	image, err := s.repo.Read(ctx, imageID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve image: %w", err)
	}

	if updateRequest.DatasetName != nil {
		image.DatasetName = *updateRequest.DatasetName
	}
	if updateRequest.OrganType != nil {
		image.OrganType = *updateRequest.OrganType
	}
	if updateRequest.DiseaseType != nil {
		image.DiseaseType = updateRequest.DiseaseType
	}
	if updateRequest.Classification != nil {
		image.Classification = updateRequest.Classification
	}
	if updateRequest.SubType != nil {
		image.SubType = updateRequest.SubType
	}
	if updateRequest.Grade != nil {
		image.Grade = updateRequest.Grade
	}

	image.UpdatedAt = time.Now()
	err = s.repo.Update(ctx, image)
	if err != nil {
		return nil, fmt.Errorf("failed to update image: %w", err)
	}
	return image, nil
}

// DeleteImage deletes an image record and its associated files.
func (s *ImageService) DeleteImage(ctx context.Context, imageID string) error {

	// Delete the image record
	if err := s.repo.Delete(ctx, imageID); err != nil {
		return fmt.Errorf("failed to delete image record: %w", err)
	}

	return nil
}

// ListImages retrieves all images with optional filtering.
func (s *ImageService) ListImages(ctx context.Context, filter *models.ImageFilter) ([]*models.Image, error) {
	images, err := s.repo.Filter(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list images: %w", err)
	}
	return images, nil
}
