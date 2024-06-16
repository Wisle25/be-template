package file_statics

import "mime/multipart"

type ConvertTo int8

const (
	WEBP = iota
	JPG
)

type FileProcessing interface {
	CompressImage(buffer []byte, to ConvertTo) ([]byte, string)
	ResizeImage(fileHeader *multipart.FileHeader)
}
