package repository

import (
	"context"
	"mime/multipart"
)

type FileRepository interface {
	Create(ctx context.Context, file multipart.File, filepath string) error
	Read(ctx context.Context, filepath string) (multipart.File, error)
	Update(ctx context.Context, file multipart.File, filepath string) error
	Delete(ctx context.Context, filepath string) error
}
