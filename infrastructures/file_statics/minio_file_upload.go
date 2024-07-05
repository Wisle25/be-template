package file_statics

import (
	"bytes"
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/wisle25/be-template/applications/file"
	"github.com/wisle25/be-template/applications/generator"
	"github.com/wisle25/be-template/commons"
	"io"
)

type MinioFileUpload struct /* implements FileUpload */ {
	minio       *minio.Client
	idGenerator generator.IdGenerator
	bucketName  string
}

func NewMinioFileUpload(
	minio *minio.Client,
	idGenerator generator.IdGenerator,
	bucketName string,
) file.FileUpload {
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
		commons.ThrowServerError("minio: upload file err", err)
	}

	return newName
}

func (m *MinioFileUpload) GetFile(filename string) []byte {
	ctx := context.Background()

	// Get from minio
	object, err := m.minio.GetObject(ctx, m.bucketName, filename, minio.GetObjectOptions{})
	if err != nil {
		commons.ThrowServerError("minio: get file err", err)
	}

	// Convert to bytes buffer
	buffer := new(bytes.Buffer)
	_, err = io.Copy(buffer, object)
	if err != nil {
		commons.ThrowServerError("minio: copy file err", err)
	}

	return buffer.Bytes()
}

func (m *MinioFileUpload) RemoveFile(oldFileLink string) {
	ctx := context.Background()

	// Remove
	removeOpts := minio.RemoveObjectOptions{}
	err := m.minio.RemoveObject(ctx, m.bucketName, oldFileLink, removeOpts)
	if err != nil {
		commons.ThrowServerError("minio: remove file err", err)
	}
}
