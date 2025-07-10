package adapter

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"

	"cloud.google.com/go/storage"
)

type GCSFileRepository struct {
	client *storage.Client
	bucket string
}

func NewGCSFileRepository(client *storage.Client, bucketName string) *GCSFileRepository {
	return &GCSFileRepository{
		client: client,
		bucket: bucketName,
	}
}

func (g *GCSFileRepository) Create(ctx context.Context, file multipart.File, filepath string) error {
	wc := g.client.Bucket(g.bucket).Object(filepath).NewWriter(ctx)
	if _, err := io.Copy(wc, file); err != nil {
		return fmt.Errorf("failed to write file to GCS: %w", err)
	}

	if err := wc.Close(); err != nil {
		return fmt.Errorf("failed to close GCS writer: %w", err)
	}

	return nil
}

func (g *GCSFileRepository) Read(ctx context.Context, filepath string) (io.ReadCloser, error) {
	rc, err := g.client.Bucket(g.bucket).Object(filepath).NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to read file from GCS: %w", err)
	}

	return rc, nil
}

func (g *GCSFileRepository) Update(ctx context.Context, file multipart.File, filepath string) error {
	wc := g.client.Bucket(g.bucket).Object(filepath).NewWriter(ctx)
	defer wc.Close()
	if _, err := io.Copy(wc, file); err != nil {
		return fmt.Errorf("failed to update file in GCS: %w", err)
	}

	if err := wc.Close(); err != nil {
		return fmt.Errorf("failed to close GCS writer: %w", err)
	}

	return nil
}

func (g *GCSFileRepository) Delete(ctx context.Context, filepath string) error {
	if err := g.client.Bucket(g.bucket).Object(filepath).Delete(ctx); err != nil {
		return fmt.Errorf("failed to delete file from GCS: %w", err)
	}

	return nil
}
