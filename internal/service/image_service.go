package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/histopathai/image-catalog-service/config"
	"github.com/histopathai/image-catalog-service/internal/models"
	"github.com/histopathai/image-catalog-service/internal/repository"
)

// ImageService provides methods to manage images in the catalog.
type ImageService struct {
	repo repository.ImageRepository
	fs   repository.FileRepository
	mq   MessageQueue
	cfg  *config.Config
}

// NewImageService creates a new ImageService instance.
func NewImageService(repo repository.ImageRepository, fileRepo repository.FileRepository, mq MessageQueue, cfg *config.Config) *ImageService {
	return &ImageService{
		repo: repo,
		fs:   fileRepo,
		mq:   mq,
		cfg:  cfg,
	}
}

// UploadImage uploads an image file and creates a new image record in the catalog.
func (s *ImageService) UploadImage(ctx context.Context, file multipart.File, header *multipart.FileHeader, request *models.ImageUploadRequest) (*models.Image, error) {

	imageID := uuid.New().String()

	fName := strings.ToLower(header.Filename)
	ext := filepath.Ext(fName)
	if !s.cfg.IsSupportedFormat(ext) {
		return nil, fmt.Errorf("unsupported file format: %s", ext)
	}

	originalName := strings.TrimSuffix(fName, ext)

	originalGCSPATH := fmt.Sprintf("originals/%s/%s/%s/%s", request.DatasetName, request.OrganType, request.DatasetName, fName)

	if err := s.fs.Create(ctx, file, originalGCSPATH); err != nil {
		return nil, fmt.Errorf("failed to upload original image: %w", err)
	}

	image := &models.Image{
		ID:               imageID,
		OriginalName:     originalName,
		OriginalUID:      request.OriginalUID,
		DatasetName:      request.DatasetName,
		ImageType:        ext,
		OrganType:        request.OrganType,
		DiseaseType:      request.DiseaseType,
		Classification:   request.Classification,
		SubType:          request.Subtype,
		Grade:            request.Grade,
		OriginalGCSPath:  originalGCSPATH,
		FileSize:         header.Size,
		ProcessingStatus: models.ProcessingStatusPending,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),

		// Placeholders, will be updated after processing
		DZIGCSPath:       "",
		TilesGCSPath:     "",
		ThumbnailGCSPath: "",
		Width:            0,
		Height:           0,
	}

	err := s.repo.Create(ctx, image)
	if err != nil {
		//rollback file upload
		s.fs.Delete(ctx, originalGCSPATH)
		return nil, fmt.Errorf("failed to create image record: %w", err)
	}

	// Publish to message queue for processing
	pj := &models.ProcessingJob{
		ImageID:           image.ID,
		OriginalGCSPath:   image.OriginalGCSPath,
		FileType:          ext,
		DestinationBucket: s.cfg.GCS.DZIBucket,
		ThumbnailSize:     s.cfg.Processing.ThumbnailSize,
		DZITileSize:       s.cfg.Processing.DZITileSize,
	}

	if err := s.mq.PublishProcessingJob(ctx, pj); err != nil {

		image.ProcessingStatus = models.ProcessingStatusFailed
		image.UpdatedAt = time.Now()
		if updateErr := s.repo.Update(ctx, image); updateErr != nil {
			return nil, fmt.Errorf("failed to update image processing status: %w", updateErr)
		}
		return nil, fmt.Errorf("failed to publish processing job: %w", err)
	}

	return image, nil
}

// GetImage retrieves an image by its ID.
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
	if updateRequest.OriginalUID != nil {
		image.OriginalUID = *updateRequest.OriginalUID
	}
	if updateRequest.Subtype != nil {
		image.SubType = updateRequest.Subtype
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
	image, err := s.repo.Read(ctx, imageID)
	if err != nil {
		return fmt.Errorf("failed to retrieve image: %w", err)
	}

	// Delete associated files
	if err := s.fs.Delete(ctx, image.OriginalGCSPath); err != nil {
		return fmt.Errorf("failed to delete original image file: %w", err)
	}
	if image.DZIGCSPath != "" {
		if err := s.fs.Delete(ctx, image.DZIGCSPath); err != nil {
			return fmt.Errorf("failed to delete DZI file: %w", err)
		}
	}
	if image.TilesGCSPath != "" {
		if err := s.fs.Delete(ctx, image.TilesGCSPath); err != nil {
			return fmt.Errorf("failed to delete tiles file: %w", err)
		}
	}
	if image.ThumbnailGCSPath != "" {
		if err := s.fs.Delete(ctx, image.ThumbnailGCSPath); err != nil {
			return fmt.Errorf("failed to delete thumbnail file: %w", err)
		}
	}

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
