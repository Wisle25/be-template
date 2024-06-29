package file_statics

// ConvertTo represents the type to which an image can be converted.
type ConvertTo int8

const (
	// WEBP represents the WEBP image format.
	WEBP ConvertTo = iota
	// JPG represents the JPEG image format.
	JPG
)

// FileProcessing defines the interface for file processing operations.
// Any struct that implements this interface can be used to process files.
type FileProcessing interface {
	// CompressImage compresses the given image buffer to the specified format.
	// The buffer parameter is the image data in bytes.
	// The "to" parameter specifies the target image format.
	// Returns the compressed image data in bytes and the MIME type as a string.
	CompressImage(buffer []byte, to ConvertTo) ([]byte, string)

	// AddWatermark adds a watermark to the given image buffer.
	// The buffer parameter is the image data in bytes.
	// Returns the image data with the watermark added.
	AddWatermark(buffer []byte) []byte
}
