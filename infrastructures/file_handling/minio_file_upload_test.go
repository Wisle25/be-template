package file_handling_test

import (
	"context"
	"testing"

	"github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/assert"
	"github.com/wisle25/be-template/commons"
	"github.com/wisle25/be-template/infrastructures/file_handling"
	"github.com/wisle25/be-template/infrastructures/generator"
	"github.com/wisle25/be-template/infrastructures/services"
)

func TestMinioFileUpload_RealService(t *testing.T) {
	// Load the configuration
	config := commons.LoadConfig("../..")

	// Initialize Minio client
	minioClient, bucketName := services.NewMinio(config)
	idGenerator := generator.NewUUIDGenerator()
	minioUploader := file_handling.NewMinioFileUpload(minioClient, idGenerator, bucketName)

	t.Run("UploadFile", func(t *testing.T) {
		t.Run("SuccessfulUpload", func(t *testing.T) {
			// Arrange
			buffer := []byte("file content")
			extension := ".jpg"
			expectedName := idGenerator.Generate() + extension

			// Action
			fileName := minioUploader.UploadFile(buffer, extension)

			// Assert
			assert.NotEmpty(t, fileName)
			assert.Equal(t, expectedName, fileName)

			// Cleanup
			_ = minioClient.RemoveObject(context.Background(), bucketName, fileName, minio.RemoveObjectOptions{})
		})

		t.Run("UploadFileWithEmptyBuffer", func(t *testing.T) {
			// Arrange
			var buffer []byte
			extension := ".jpg"

			// Action
			fileName := minioUploader.UploadFile(buffer, extension)

			// Assert
			assert.Empty(t, fileName)
		})

		t.Run("UploadFileError", func(t *testing.T) {
			// Arrange
			buffer := []byte("file content")
			extension := ".jpg"

			// Simulate error by using an invalid bucket name
			invalidBucketName := "invalid-bucket"
			minioUploader := file_handling.NewMinioFileUpload(minioClient, idGenerator, invalidBucketName)

			// Action and Assert
			assert.Panics(t, func() {
				minioUploader.UploadFile(buffer, extension)
			})
		})
	})

	t.Run("GetFile", func(t *testing.T) {
		t.Run("SuccessfulGet", func(t *testing.T) {
			// Arrange
			buffer := []byte("file content")
			extension := ".jpg"
			fileName := minioUploader.UploadFile(buffer, extension)

			// Action
			fileContent := minioUploader.GetFile(fileName)

			// Assert
			assert.NotEmpty(t, fileContent)
			assert.Equal(t, buffer, fileContent)

			// Cleanup
			_ = minioClient.RemoveObject(context.Background(), bucketName, fileName, minio.RemoveObjectOptions{})
		})

		t.Run("GetFileError", func(t *testing.T) {
			// Arrange
			filename := "non-existent-file.jpg"

			// Action and Assert
			assert.Panics(t, func() {
				minioUploader.GetFile(filename)
			})
		})
	})

	t.Run("RemoveFile", func(t *testing.T) {
		t.Run("SuccessfulRemove", func(t *testing.T) {
			// Arrange
			buffer := []byte("file content")
			extension := ".jpg"
			fileName := minioUploader.UploadFile(buffer, extension)

			// Action
			minioUploader.RemoveFile(fileName)

			// Assert
			_, err := minioClient.GetObject(context.Background(), bucketName, fileName, minio.GetObjectOptions{})
			assert.Error(t, err)
		})

		t.Run("RemoveFileError", func(t *testing.T) {
			// Arrange
			filename := "non-existent-file.jpg"

			// Action and Assert
			assert.Panics(t, func() {
				minioUploader.RemoveFile(filename)
			})
		})
	})
}
