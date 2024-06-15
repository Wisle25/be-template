package processing

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/wisle25/be-template/applications/generator"
	"mime/multipart"
	"path/filepath"
)

type MinioFileProcessing struct {
	minio       *minio.Client
	idGenerator generator.IdGenerator
	bucketName  string
}

func NewMinioFileProcessing(
	minio *minio.Client,
	idGenerator generator.IdGenerator,
	bucketName string,
) *MinioFileProcessing {
	return &MinioFileProcessing{
		minio,
		idGenerator,
		bucketName,
	}
}

func (m *MinioFileProcessing) UploadFile(fileHeader *multipart.FileHeader) string {
	if fileHeader == nil {
		return ""
	}

	ctx := context.Background()
	var err error

	// Create new name
	extension := filepath.Ext(fileHeader.Filename)
	newName := m.idGenerator.Generate() + extension

	// Get buffer from file
	buffer, err := fileHeader.Open()

	if err != nil {
		panic(fmt.Errorf("minio: get buffer file err: %v", err))
	}
	defer buffer.Close()

	// Upload
	uploadOpts := minio.PutObjectOptions{
		ContentType: fileHeader.Header["Content-Type"][0],
	}
	_, err = m.minio.PutObject(
		ctx,
		m.bucketName,
		newName,
		buffer,
		fileHeader.Size,
		uploadOpts,
	)
	if err != nil {
		panic(fmt.Errorf("minio: upload file err: %v", err))
	}

	return newName
}

func (m *MinioFileProcessing) RemoveFile(oldFileLink string) {
	ctx := context.Background()

	// Remove
	removeOpts := minio.RemoveObjectOptions{}
	err := m.minio.RemoveObject(ctx, m.bucketName, oldFileLink, removeOpts)
	if err != nil {
		panic(fmt.Errorf("minio: remove file err: %v", err))
	}
}
