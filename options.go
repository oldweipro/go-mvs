package mvs

import (
	"path/filepath"
	"strings"
)

type ImageType uint32

type InterpolationMethod int32

type GrabStrategy uint32

type ImageSaveOptions struct {
	Type        ImageType
	Quality     uint32
	MethodValue InterpolationMethod
}

type PixelConvertOptions struct {
	DstPixelType  uint32
	DstBufferSize int
}

type EnumEntry struct {
	Value    uint32
	Symbolic string
}

type FrameCallback func(*Frame)

func ImageTypeFromExtension(path string) (ImageType, bool) {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".bmp":
		return ImageTypeBMP, true
	case ".jpg", ".jpeg":
		return ImageTypeJPEG, true
	case ".png":
		return ImageTypePNG, true
	case ".tif", ".tiff":
		return ImageTypeTIFF, true
	default:
		return ImageTypeUndefined, false
	}
}
