package file_statics

import (
	"fmt"
	"github.com/davidbyttow/govips/v2/vips"
	"github.com/wisle25/be-template/applications/file_statics"
	"mime/multipart"
)

type VipsFileProcessing struct {
}

func NewVipsFileProcessing() file_statics.FileProcessing {
	vips.Startup(nil)
	return &VipsFileProcessing{}
}

func (v *VipsFileProcessing) CompressImage(buffer []byte, to file_statics.ConvertTo) ([]byte, string) {
	if buffer == nil {
		return nil, ""
	}

	var result []byte
	var extension string
	var err error

	image, err := vips.NewImageFromBuffer(buffer)
	if err != nil {
		panic(fmt.Errorf("vips: new compress image: %v", err))
	}

	switch to {
	case file_statics.JPG:
		options := vips.NewJpegExportParams()
		options.Quality = 40
		options.StripMetadata = true
		result, _, err = image.ExportJpeg(options)
		if err != nil {
			panic(fmt.Errorf("vips: compress image jpeg: %v", err))
		}

		extension = ".jpg"
	case file_statics.WEBP:
		options := vips.NewWebpExportParams()
		options.Quality = 40
		options.StripMetadata = true
		options.Lossless = false

		result, _, err = image.ExportWebp(options)
		if err != nil {
			panic(fmt.Errorf("vips: compress image webp: %v", err))
		}

		extension = ".webp"
	}

	return result, extension
}

func (v *VipsFileProcessing) ResizeImage(fileHeader *multipart.FileHeader) {
	//TODO implement me
	panic("implement me")
}
