package mvs

import (
	"errors"
	"image"
	"testing"
)

func TestPixelTypeName(t *testing.T) {
	tests := []struct {
		pixelType uint32
		want      string
	}{
		{PixelTypeMono8, "Mono8"},
		{PixelTypeMono10Packed, "Mono10Packed"},
		{PixelTypeBayerRG8, "BayerRG8"},
		{PixelTypeRGB8Packed, "RGB8Packed"},
		{PixelTypeHighBandwidthBayerGR10Packed, "HighBandwidthBayerGR10Packed"},
		{0x12345678, "0x12345678"},
	}

	for _, tt := range tests {
		if got := PixelTypeName(tt.pixelType); got != tt.want {
			t.Fatalf("PixelTypeName(0x%08X) = %q, want %q", tt.pixelType, got, tt.want)
		}
	}
}

func TestExpectedFrameDataLength(t *testing.T) {
	got, err := ExpectedFrameDataLength(4, 3, PixelTypeMono8)
	if err != nil {
		t.Fatal(err)
	}
	if got != 12 {
		t.Fatalf("Mono8 length = %d, want 12", got)
	}

	got, err = ExpectedFrameDataLength(4, 3, PixelTypeRGB8Packed)
	if err != nil {
		t.Fatal(err)
	}
	if got != 36 {
		t.Fatalf("RGB8 length = %d, want 36", got)
	}

	if _, err := ExpectedFrameDataLength(4, 3, PixelTypeUndefined); !errors.Is(err, ErrUnsupportedPixel) {
		t.Fatalf("undefined pixel error = %v, want ErrUnsupportedPixel", err)
	}

	if _, err := ExpectedFrameDataLength(0, 3, PixelTypeMono8); !errors.Is(err, ErrInvalidFrameData) {
		t.Fatalf("empty size error = %v, want ErrInvalidFrameData", err)
	}

	if _, err := ExpectedFrameDataLength(4, 3, PixelTypeHighBandwidthMono8); !errors.Is(err, ErrUnsupportedPixel) {
		t.Fatalf("compressed pixel error = %v, want ErrUnsupportedPixel", err)
	}
}

func TestPixelTypePredicates(t *testing.T) {
	if !IsMonoPixelType(PixelTypeMono8) {
		t.Fatal("Mono8 should be mono")
	}
	if !IsColorPixelType(PixelTypeRGB8Packed) {
		t.Fatal("RGB8Packed should be color")
	}
	if !IsMonoPixelType(PixelTypeHighBandwidthMono8) {
		t.Fatal("HighBandwidthMono8 should keep mono flag")
	}
	if !IsCompressedPixelType(PixelTypeHighBandwidthBayerRG8) {
		t.Fatal("HighBandwidthBayerRG8 should be compressed")
	}
	if !IsBayerPixelType(PixelTypeBayerGR12Packed) {
		t.Fatal("BayerGR12Packed should be bayer")
	}
}

func TestImageTypeFromExtension(t *testing.T) {
	tests := []struct {
		path string
		want ImageType
		ok   bool
	}{
		{"frame.bmp", ImageTypeBMP, true},
		{"frame.JPG", ImageTypeJPEG, true},
		{"frame.png", ImageTypePNG, true},
		{"frame.tiff", ImageTypeTIFF, true},
		{"frame.raw", ImageTypeUndefined, false},
	}

	for _, tt := range tests {
		got, ok := ImageTypeFromExtension(tt.path)
		if got != tt.want || ok != tt.ok {
			t.Fatalf("ImageTypeFromExtension(%q) = %v,%v want %v,%v", tt.path, got, ok, tt.want, tt.ok)
		}
	}
}

func TestImageFromMono8(t *testing.T) {
	img, err := ImageFromBuffer([]byte{1, 2, 3, 4}, 2, 2, PixelTypeMono8)
	if err != nil {
		t.Fatal(err)
	}

	gray, ok := img.(*image.Gray)
	if !ok {
		t.Fatalf("image type = %T, want *image.Gray", img)
	}
	if got := gray.GrayAt(1, 1).Y; got != 4 {
		t.Fatalf("pixel(1,1) = %d, want 4", got)
	}
}

func TestImageFromRGBAndBGR(t *testing.T) {
	rgb, err := ImageFromBuffer([]byte{10, 20, 30}, 1, 1, PixelTypeRGB8Packed)
	if err != nil {
		t.Fatal(err)
	}
	r, g, b, a := rgb.At(0, 0).RGBA()
	if r>>8 != 10 || g>>8 != 20 || b>>8 != 30 || a>>8 != 255 {
		t.Fatalf("RGB pixel = %d,%d,%d,%d", r>>8, g>>8, b>>8, a>>8)
	}

	bgr, err := ImageFromBuffer([]byte{30, 20, 10}, 1, 1, PixelTypeBGR8Packed)
	if err != nil {
		t.Fatal(err)
	}
	r, g, b, a = bgr.At(0, 0).RGBA()
	if r>>8 != 10 || g>>8 != 20 || b>>8 != 30 || a>>8 != 255 {
		t.Fatalf("BGR pixel = %d,%d,%d,%d", r>>8, g>>8, b>>8, a>>8)
	}
}

func TestImageFromBayer8(t *testing.T) {
	img, err := ImageFromBuffer([]byte{
		255, 10,
		20, 30,
	}, 2, 2, PixelTypeBayerRG8)
	if err != nil {
		t.Fatal(err)
	}

	r, _, _, a := img.At(0, 0).RGBA()
	if r>>8 != 255 || a>>8 != 255 {
		t.Fatalf("Bayer pixel = r:%d a:%d", r>>8, a>>8)
	}
}

func TestImageFromBufferErrors(t *testing.T) {
	if _, err := ImageFromBuffer([]byte{1}, 2, 2, PixelTypeMono8); !errors.Is(err, ErrInvalidFrameData) {
		t.Fatalf("short data error = %v, want ErrInvalidFrameData", err)
	}
	if _, err := ImageFromBuffer([]byte{1, 2, 3, 4}, 2, 2, PixelTypeMono16); !errors.Is(err, ErrInvalidFrameData) {
		t.Fatalf("short Mono16 error = %v, want ErrInvalidFrameData", err)
	}
	if _, err := ImageFromBuffer(make([]byte, 8), 2, 2, PixelTypeMono16); !errors.Is(err, ErrUnsupportedPixel) {
		t.Fatalf("unsupported pixel error = %v, want ErrUnsupportedPixel", err)
	}
}
