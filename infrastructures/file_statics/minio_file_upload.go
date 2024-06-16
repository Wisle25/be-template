﻿package file_statics

import (
	"bytes"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/wisle25/be-template/applications/file_statics"
	"github.com/wisle25/be-template/applications/generator"
)

type MinioFileUpload struct {
	minio       *minio.Client
	idGenerator generator.IdGenerator
	bucketName  string
}

func NewMinioFileUpload(
	minio *minio.Client,
	idGenerator generator.IdGenerator,
	bucketName string,
) file_statics.FileUpload {
	return &MinioFileUpload{
		minio,
		idGenerator,
		bucketName,
	}
}

func (m *MinioFileUpload) UploadFile(buffer []byte, extension string) string {
	if buffer == nil {
		return ""
	}

	ctx := context.Background()
	var err error

	// Create new name
	newName := m.idGenerator.Generate() + extension

	// Upload
	uploadOpts := minio.PutObjectOptions{
		ContentType: "image/" + extension[1:],
	}
	_, err = m.minio.PutObject(
		ctx,
		m.bucketName,
		newName,
		bytes.NewReader(buffer),
		int64(len(buffer)),
		uploadOpts,
	)
	if err != nil {
		panic(fmt.Errorf("minio: upload file err: %v", err))
	}

	return newName
}

func (m *MinioFileUpload) RemoveFile(oldFileLink string) {
	ctx := context.Background()

	// Remove
	removeOpts := minio.RemoveObjectOptions{}
	err := m.minio.RemoveObject(ctx, m.bucketName, oldFileLink, removeOpts)
	if err != nil {
		panic(fmt.Errorf("minio: remove file err: %v", err))
	}
}
