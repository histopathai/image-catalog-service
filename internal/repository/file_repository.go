package repository

import (
	"context"
	"io"
	"mime/multipart"
)

type FileRepository interface {
	Create(ctx context.Context, file multipart.File, filepath string) error
	Read(ctx context.Context, filepath string) (io.ReadCloser, error)
	Update(ctx context.Context, file multipart.File, filepath string) error
	Delete(ctx context.Context, filepath string) error
}
