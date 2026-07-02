//go:build windows && amd64

package mvs

import "image"

func (f *Frame) Image() (image.Image, error) {
	if f == nil {
		return nil, ErrInvalidFrameData
	}
	return ImageFromBuffer(f.Data, f.Width, f.Height, f.PixelType)
}

func (f *Frame) PixelTypeName() string {
	if f == nil {
		return PixelTypeName(PixelTypeUndefined)
	}
	return PixelTypeName(f.PixelType)
}
