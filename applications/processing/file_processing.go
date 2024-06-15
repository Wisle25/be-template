package processing

import "mime/multipart"

// FileProcessing handling file uploading, manipulating (like adding watermark), removing, etc
type FileProcessing interface {
	// UploadFile before uploading file, it will generate a new name.
	// Receiving the file that want to be uploaded.
	// Returning uploaded link.
	UploadFile(fileHeader *multipart.FileHeader) string

	// RemoveFile deleting specified file by its link
	// Do nothing if it's really not existed
	RemoveFile(oldFileLink string)
}
