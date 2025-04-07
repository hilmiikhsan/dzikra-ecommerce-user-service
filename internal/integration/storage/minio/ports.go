package minio

import (
	"context"
	"mime/multipart"
)

type MinioService interface {
	UploadFile(ctx context.Context, objectName string, file multipart.File, fileSize int64, contentType string) (string, error)
}
