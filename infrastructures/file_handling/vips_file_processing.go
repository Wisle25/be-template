﻿package file_handling

import (
	"github.com/davidbyttow/govips/v2/vips"
	"github.com/wisle25/be-template/applications/file_handling"
	"github.com/wisle25/be-template/commons"
	"os"
	"path/filepath"
)

type VipsFileProcessing struct /* implements FileProcessing */ {

}

func NewVipsFileProcessing() file_handling.FileProcessing {
	vips.Startup(nil)
	return &VipsFileProcessing{}
}

func (v *VipsFileProcessing) CompressImage(buffer []byte, to file_handling.ConvertTo) ([]byte, string) {
	if buffer == nil {
		return nil, ""
	}

	var result []byte
	var extension string
	var err error

	// Open Image
	image, err := vips.NewImageFromBuffer(buffer)
	if err != nil {
		commons.ThrowServerError("vips: new compress image", err)
	}

	// Compressing depends on what the option
	switch to {
	case file_handling.JPG:
		options := vips.NewJpegExportParams()
		options.Quality = 40
		options.StripMetadata = true
		result, _, err = image.ExportJpeg(options)
		if err != nil {
			commons.ThrowServerError("vips: compress image jpeg", err)
		}

		extension = ".jpg"
	case file_handling.WEBP:
		options := vips.NewWebpExportParams()
		options.Quality = 40
		options.StripMetadata = true
		options.Lossless = false

		result, _, err = image.ExportWebp(options)
		if err != nil {
			commons.ThrowServerError("vips: compress image webp", err)
		}

		extension = ".webp"
	}

	return result, extension
}

func (v *VipsFileProcessing) AddWatermark(buffer []byte) []byte {
	// Open original image
	originalImage, err := vips.NewImageFromBuffer(buffer)
	if err != nil {
		commons.ThrowServerError("add_watermark_err: opening original image", err)
	}
	defer originalImage.Close()

	// Open watermark image
	rootDir, _ := os.Getwd()
	watermarkImage, err := vips.NewImageFromFile(filepath.Join(rootDir, "resources", "watermark.png"))
	if err != nil {
		commons.ThrowServerError("add_watermark_err: opening watermark image", err)
	}
	defer watermarkImage.Close()

	// Resize watermark image to fit the original image
	err = watermarkImage.ResizeWithVScale(
		float64(originalImage.Width())/float64(watermarkImage.Width()),
		float64(originalImage.Height())/float64(watermarkImage.Height()),
		vips.KernelLanczos3,
	)
	if err != nil {
		commons.ThrowServerError("add_watermark_err: resizing watermark", err)
	}

	// Composite
	err = originalImage.Composite(watermarkImage, vips.BlendModeAdd, 0, 0)
	if err != nil {
		commons.ThrowServerError("add_watermark_err: compositing watermark", err)
	}

	// Get the buffer of the result
	resultBuffer, _, err := originalImage.ExportNative()
	if err != nil {
		commons.ThrowServerError("add_watermark_err: export native", err)
	}

	return resultBuffer
}
