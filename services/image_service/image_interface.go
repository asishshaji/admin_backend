package image_service

import (
	"context"
	"mime/multipart"
)

type IImageService interface {
	UploadImage(ctx context.Context, file multipart.File) (string, error)
}
