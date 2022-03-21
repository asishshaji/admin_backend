package file_service

import (
	"context"
	"mime/multipart"
)

type IFileService interface {
	UploadFile(ctx context.Context, file multipart.File) (string, error)
}
