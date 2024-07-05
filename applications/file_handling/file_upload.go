package file_handling

// FileUpload handling file_handling uploading, manipulating (like adding watermark), removing, etc
type FileUpload interface {
	// UploadFile before uploading file_handling, it will generate a new name.
	// Receiving the buffer and extension file_handling that want to be uploaded.
	// Returning uploaded link.
	UploadFile(buffer []byte, extension string) string

	// GetFile Getting the object file_handling as buffer
	GetFile(fileName string) []byte

	// RemoveFile deleting specified file_handling by its link
	// Do nothing if it's really not existed
	RemoveFile(oldFileLink string)
}
