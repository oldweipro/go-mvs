package mvs

import (
	"fmt"
	"image"
	"image/color"
)

const (
	pixelColorMask               uint32 = 0xFF000000
	pixelMonoFlag                uint32 = 0x01000000
	pixelColorFlag               uint32 = 0x02000000
	pixelEffectivePixelSizeMask  uint32 = 0x00FF0000
	pixelEffectivePixelSizeShift        = 16
)

func PixelTypeName(pixelType uint32) string {
	if name, ok := pixelTypeNames[pixelType]; ok {
		return name
	}
	return fmt.Sprintf("0x%08X", pixelType)
}

var pixelTypeNames = map[uint32]string{
	PixelTypeUndefined: "Undefined",

	PixelTypeMono1p:       "Mono1p",
	PixelTypeMono2p:       "Mono2p",
	PixelTypeMono4p:       "Mono4p",
	PixelTypeMono8:        "Mono8",
	PixelTypeMono8Signed:  "Mono8Signed",
	PixelTypeMono10:       "Mono10",
	PixelTypeMono10Packed: "Mono10Packed",
	PixelTypeMono12:       "Mono12",
	PixelTypeMono12Packed: "Mono12Packed",
	PixelTypeMono14:       "Mono14",
	PixelTypeMono16:       "Mono16",

	PixelTypeBayerGR8:        "BayerGR8",
	PixelTypeBayerRG8:        "BayerRG8",
	PixelTypeBayerGB8:        "BayerGB8",
	PixelTypeBayerBG8:        "BayerBG8",
	PixelTypeBayerRBGG8:      "BayerRBGG8",
	PixelTypeBayerGR10:       "BayerGR10",
	PixelTypeBayerRG10:       "BayerRG10",
	PixelTypeBayerGB10:       "BayerGB10",
	PixelTypeBayerBG10:       "BayerBG10",
	PixelTypeBayerGR12:       "BayerGR12",
	PixelTypeBayerRG12:       "BayerRG12",
	PixelTypeBayerGB12:       "BayerGB12",
	PixelTypeBayerBG12:       "BayerBG12",
	PixelTypeBayerGR10Packed: "BayerGR10Packed",
	PixelTypeBayerRG10Packed: "BayerRG10Packed",
	PixelTypeBayerGB10Packed: "BayerGB10Packed",
	PixelTypeBayerBG10Packed: "BayerBG10Packed",
	PixelTypeBayerGR12Packed: "BayerGR12Packed",
	PixelTypeBayerRG12Packed: "BayerRG12Packed",
	PixelTypeBayerGB12Packed: "BayerGB12Packed",
	PixelTypeBayerBG12Packed: "BayerBG12Packed",
	PixelTypeBayerGR16:       "BayerGR16",
	PixelTypeBayerRG16:       "BayerRG16",
	PixelTypeBayerGB16:       "BayerGB16",
	PixelTypeBayerBG16:       "BayerBG16",

	PixelTypeRGB8Packed:    "RGB8Packed",
	PixelTypeBGR8Packed:    "BGR8Packed",
	PixelTypeRGBA8Packed:   "RGBA8Packed",
	PixelTypeBGRA8Packed:   "BGRA8Packed",
	PixelTypeRGB10Packed:   "RGB10Packed",
	PixelTypeBGR10Packed:   "BGR10Packed",
	PixelTypeRGB12Packed:   "RGB12Packed",
	PixelTypeBGR12Packed:   "BGR12Packed",
	PixelTypeRGB16Packed:   "RGB16Packed",
	PixelTypeBGR16Packed:   "BGR16Packed",
	PixelTypeRGBA16Packed:  "RGBA16Packed",
	PixelTypeBGRA16Packed:  "BGRA16Packed",
	PixelTypeRGB10V1Packed: "RGB10V1Packed",
	PixelTypeRGB10V2Packed: "RGB10V2Packed",
	PixelTypeRGB12V1Packed: "RGB12V1Packed",
	PixelTypeRGB565Packed:  "RGB565Packed",
	PixelTypeBGR565Packed:  "BGR565Packed",

	PixelTypeYUV411Packed:         "YUV411Packed",
	PixelTypeYUV422Packed:         "YUV422Packed",
	PixelTypeYUV422YUYVPacked:     "YUV422YUYVPacked",
	PixelTypeYUV444Packed:         "YUV444Packed",
	PixelTypeYCbCr8CbYCr:          "YCbCr8CbYCr",
	PixelTypeYCbCr4228:            "YCbCr4228",
	PixelTypeYCbCr4228CbYCrY:      "YCbCr4228CbYCrY",
	PixelTypeYCbCr4118CbYYCrYY:    "YCbCr4118CbYYCrYY",
	PixelTypeYCbCr6018CbYCr:       "YCbCr6018CbYCr",
	PixelTypeYCbCr6014228:         "YCbCr6014228",
	PixelTypeYCbCr6014228CbYCrY:   "YCbCr6014228CbYCrY",
	PixelTypeYCbCr6014118CbYYCrYY: "YCbCr6014118CbYYCrYY",
	PixelTypeYCbCr7098CbYCr:       "YCbCr7098CbYCr",
	PixelTypeYCbCr7094228:         "YCbCr7094228",
	PixelTypeYCbCr7094228CbYCrY:   "YCbCr7094228CbYCrY",
	PixelTypeYCbCr7094118CbYYCrYY: "YCbCr7094118CbYYCrYY",
	PixelTypeYUV420SPNV12:         "YUV420SPNV12",
	PixelTypeYUV420SPNV21:         "YUV420SPNV21",

	PixelTypeRGB8Planar:  "RGB8Planar",
	PixelTypeRGB10Planar: "RGB10Planar",
	PixelTypeRGB12Planar: "RGB12Planar",
	PixelTypeRGB16Planar: "RGB16Planar",

	PixelTypeJPEG:                 "JPEG",
	PixelTypeCoord3DABC32f:        "Coord3DABC32f",
	PixelTypeCoord3DABC32fPlanar:  "Coord3DABC32fPlanar",
	PixelTypeCoord3DAC32f:         "Coord3DAC32f",
	PixelTypeCoord3DDepthPlusMask: "Coord3DDepthPlusMask",
	PixelTypeCoord3DABC32:         "Coord3DABC32",
	PixelTypeCoord3DAB32f:         "Coord3DAB32f",
	PixelTypeCoord3DAB32:          "Coord3DAB32",
	PixelTypeCoord3DAC32f64:       "Coord3DAC32f64",
	PixelTypeCoord3DAC32fPlanar:   "Coord3DAC32fPlanar",
	PixelTypeCoord3DAC32:          "Coord3DAC32",
	PixelTypeCoord3DA32f:          "Coord3DA32f",
	PixelTypeCoord3DA32:           "Coord3DA32",
	PixelTypeCoord3DC32f:          "Coord3DC32f",
	PixelTypeCoord3DC32:           "Coord3DC32",
	PixelTypeCoord3DABC16:         "Coord3DABC16",
	PixelTypeCoord3DC16:           "Coord3DC16",
	PixelTypeFloat32:              "Float32",

	PixelTypeHighBandwidthMono8:      "HighBandwidthMono8",
	PixelTypeHighBandwidthMono10:     "HighBandwidthMono10",
	PixelTypeHighBandwidthMono10Pack: "HighBandwidthMono10Packed",
	PixelTypeHighBandwidthMono12:     "HighBandwidthMono12",
	PixelTypeHighBandwidthMono12Pack: "HighBandwidthMono12Packed",
	PixelTypeHighBandwidthMono16:     "HighBandwidthMono16",

	PixelTypeHighBandwidthBayerGR8:        "HighBandwidthBayerGR8",
	PixelTypeHighBandwidthBayerRG8:        "HighBandwidthBayerRG8",
	PixelTypeHighBandwidthBayerGB8:        "HighBandwidthBayerGB8",
	PixelTypeHighBandwidthBayerBG8:        "HighBandwidthBayerBG8",
	PixelTypeHighBandwidthBayerRBGG8:      "HighBandwidthBayerRBGG8",
	PixelTypeHighBandwidthBayerGR10:       "HighBandwidthBayerGR10",
	PixelTypeHighBandwidthBayerRG10:       "HighBandwidthBayerRG10",
	PixelTypeHighBandwidthBayerGB10:       "HighBandwidthBayerGB10",
	PixelTypeHighBandwidthBayerBG10:       "HighBandwidthBayerBG10",
	PixelTypeHighBandwidthBayerGR12:       "HighBandwidthBayerGR12",
	PixelTypeHighBandwidthBayerRG12:       "HighBandwidthBayerRG12",
	PixelTypeHighBandwidthBayerGB12:       "HighBandwidthBayerGB12",
	PixelTypeHighBandwidthBayerBG12:       "HighBandwidthBayerBG12",
	PixelTypeHighBandwidthBayerGR10Packed: "HighBandwidthBayerGR10Packed",
	PixelTypeHighBandwidthBayerRG10Packed: "HighBandwidthBayerRG10Packed",
	PixelTypeHighBandwidthBayerGB10Packed: "HighBandwidthBayerGB10Packed",
	PixelTypeHighBandwidthBayerBG10Packed: "HighBandwidthBayerBG10Packed",
	PixelTypeHighBandwidthBayerGR12Packed: "HighBandwidthBayerGR12Packed",
	PixelTypeHighBandwidthBayerRG12Packed: "HighBandwidthBayerRG12Packed",
	PixelTypeHighBandwidthBayerGB12Packed: "HighBandwidthBayerGB12Packed",
	PixelTypeHighBandwidthBayerBG12Packed: "HighBandwidthBayerBG12Packed",

	PixelTypeHighBandwidthYUV422Packed:     "HighBandwidthYUV422Packed",
	PixelTypeHighBandwidthYUV422YUYVPacked: "HighBandwidthYUV422YUYVPacked",
	PixelTypeHighBandwidthRGB8Packed:       "HighBandwidthRGB8Packed",
	PixelTypeHighBandwidthBGR8Packed:       "HighBandwidthBGR8Packed",
	PixelTypeHighBandwidthRGBA8Packed:      "HighBandwidthRGBA8Packed",
	PixelTypeHighBandwidthBGRA8Packed:      "HighBandwidthBGRA8Packed",
	PixelTypeHighBandwidthRGB16Packed:      "HighBandwidthRGB16Packed",
	PixelTypeHighBandwidthBGR16Packed:      "HighBandwidthBGR16Packed",
	PixelTypeHighBandwidthRGBA16Packed:     "HighBandwidthRGBA16Packed",
	PixelTypeHighBandwidthBGRA16Packed:     "HighBandwidthBGRA16Packed",
}

func PixelBitCount(pixelType uint32) int {
	if pixelType == PixelTypeUndefined {
		return 0
	}
	return int((pixelType & pixelEffectivePixelSizeMask) >> pixelEffectivePixelSizeShift)
}

func IsMonoPixelType(pixelType uint32) bool {
	return pixelType&pixelMonoFlag != 0 && pixelType&pixelColorFlag == 0
}

func IsColorPixelType(pixelType uint32) bool {
	return pixelType&pixelColorFlag != 0
}

func IsCompressedPixelType(pixelType uint32) bool {
	switch pixelType {
	case PixelTypeJPEG,
		PixelTypeHighBandwidthMono8,
		PixelTypeHighBandwidthMono10,
		PixelTypeHighBandwidthMono10Pack,
		PixelTypeHighBandwidthMono12,
		PixelTypeHighBandwidthMono12Pack,
		PixelTypeHighBandwidthMono16,
		PixelTypeHighBandwidthBayerGR8,
		PixelTypeHighBandwidthBayerRG8,
		PixelTypeHighBandwidthBayerGB8,
		PixelTypeHighBandwidthBayerBG8,
		PixelTypeHighBandwidthBayerRBGG8,
		PixelTypeHighBandwidthBayerGR10,
		PixelTypeHighBandwidthBayerRG10,
		PixelTypeHighBandwidthBayerGB10,
		PixelTypeHighBandwidthBayerBG10,
		PixelTypeHighBandwidthBayerGR12,
		PixelTypeHighBandwidthBayerRG12,
		PixelTypeHighBandwidthBayerGB12,
		PixelTypeHighBandwidthBayerBG12,
		PixelTypeHighBandwidthBayerGR10Packed,
		PixelTypeHighBandwidthBayerRG10Packed,
		PixelTypeHighBandwidthBayerGB10Packed,
		PixelTypeHighBandwidthBayerBG10Packed,
		PixelTypeHighBandwidthBayerGR12Packed,
		PixelTypeHighBandwidthBayerRG12Packed,
		PixelTypeHighBandwidthBayerGB12Packed,
		PixelTypeHighBandwidthBayerBG12Packed,
		PixelTypeHighBandwidthYUV422Packed,
		PixelTypeHighBandwidthYUV422YUYVPacked,
		PixelTypeHighBandwidthRGB8Packed,
		PixelTypeHighBandwidthBGR8Packed,
		PixelTypeHighBandwidthRGBA8Packed,
		PixelTypeHighBandwidthBGRA8Packed,
		PixelTypeHighBandwidthRGB16Packed,
		PixelTypeHighBandwidthBGR16Packed,
		PixelTypeHighBandwidthRGBA16Packed,
		PixelTypeHighBandwidthBGRA16Packed:
		return true
	default:
		return false
	}
}

func IsBayerPixelType(pixelType uint32) bool {
	switch pixelType {
	case PixelTypeBayerGR8,
		PixelTypeBayerRG8,
		PixelTypeBayerGB8,
		PixelTypeBayerBG8,
		PixelTypeBayerRBGG8,
		PixelTypeBayerGR10,
		PixelTypeBayerRG10,
		PixelTypeBayerGB10,
		PixelTypeBayerBG10,
		PixelTypeBayerGR12,
		PixelTypeBayerRG12,
		PixelTypeBayerGB12,
		PixelTypeBayerBG12,
		PixelTypeBayerGR10Packed,
		PixelTypeBayerRG10Packed,
		PixelTypeBayerGB10Packed,
		PixelTypeBayerBG10Packed,
		PixelTypeBayerGR12Packed,
		PixelTypeBayerRG12Packed,
		PixelTypeBayerGB12Packed,
		PixelTypeBayerBG12Packed,
		PixelTypeBayerGR16,
		PixelTypeBayerRG16,
		PixelTypeBayerGB16,
		PixelTypeBayerBG16,
		PixelTypeHighBandwidthBayerGR8,
		PixelTypeHighBandwidthBayerRG8,
		PixelTypeHighBandwidthBayerGB8,
		PixelTypeHighBandwidthBayerBG8,
		PixelTypeHighBandwidthBayerRBGG8,
		PixelTypeHighBandwidthBayerGR10,
		PixelTypeHighBandwidthBayerRG10,
		PixelTypeHighBandwidthBayerGB10,
		PixelTypeHighBandwidthBayerBG10,
		PixelTypeHighBandwidthBayerGR12,
		PixelTypeHighBandwidthBayerRG12,
		PixelTypeHighBandwidthBayerGB12,
		PixelTypeHighBandwidthBayerBG12,
		PixelTypeHighBandwidthBayerGR10Packed,
		PixelTypeHighBandwidthBayerRG10Packed,
		PixelTypeHighBandwidthBayerGB10Packed,
		PixelTypeHighBandwidthBayerBG10Packed,
		PixelTypeHighBandwidthBayerGR12Packed,
		PixelTypeHighBandwidthBayerRG12Packed,
		PixelTypeHighBandwidthBayerGB12Packed,
		PixelTypeHighBandwidthBayerBG12Packed:
		return true
	default:
		return false
	}
}

func ExpectedFrameDataLength(width, height uint32, pixelType uint32) (int, error) {
	if width == 0 || height == 0 {
		return 0, fmt.Errorf("%w: empty image size %dx%d", ErrInvalidFrameData, width, height)
	}
	if IsCompressedPixelType(pixelType) {
		return 0, fmt.Errorf("%w: compressed pixel type %s has variable data length", ErrUnsupportedPixel, PixelTypeName(pixelType))
	}

	bits := PixelBitCount(pixelType)
	if bits <= 0 {
		return 0, fmt.Errorf("%w: %s", ErrUnsupportedPixel, PixelTypeName(pixelType))
	}

	pixels := uint64(width) * uint64(height)
	totalBits := pixels * uint64(bits)
	maxInt := int(^uint(0) >> 1)
	if totalBits > uint64(maxInt)*8 {
		return 0, fmt.Errorf("%w: image data is too large", ErrInvalidFrameData)
	}
	return int((totalBits + 7) / 8), nil
}

func ImageFromBuffer(data []byte, width, height uint32, pixelType uint32) (image.Image, error) {
	expected, err := ExpectedFrameDataLength(width, height, pixelType)
	if err != nil {
		return nil, err
	}
	if len(data) < expected {
		return nil, fmt.Errorf("%w: need %d bytes, got %d", ErrInvalidFrameData, expected, len(data))
	}

	switch pixelType {
	case PixelTypeMono8:
		return imageFromMono8(data, int(width), int(height)), nil
	case PixelTypeRGB8Packed:
		return imageFromRGB(data, int(width), int(height), false, false), nil
	case PixelTypeBGR8Packed:
		return imageFromRGB(data, int(width), int(height), true, false), nil
	case PixelTypeRGBA8Packed:
		return imageFromRGB(data, int(width), int(height), false, true), nil
	case PixelTypeBGRA8Packed:
		return imageFromRGB(data, int(width), int(height), true, true), nil
	case PixelTypeBayerGR8, PixelTypeBayerRG8, PixelTypeBayerGB8, PixelTypeBayerBG8:
		return imageFromBayer8(data, int(width), int(height), pixelType), nil
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedPixel, PixelTypeName(pixelType))
	}
}

func imageFromMono8(data []byte, width, height int) *image.Gray {
	img := image.NewGray(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		copy(img.Pix[y*img.Stride:y*img.Stride+width], data[y*width:y*width+width])
	}
	return img
}

func imageFromRGB(data []byte, width, height int, bgr bool, hasAlpha bool) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	srcStep := 3
	if hasAlpha {
		srcStep = 4
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			src := (y*width + x) * srcStep
			dst := y*img.Stride + x*4
			if bgr {
				img.Pix[dst+0] = data[src+2]
				img.Pix[dst+1] = data[src+1]
				img.Pix[dst+2] = data[src+0]
			} else {
				img.Pix[dst+0] = data[src+0]
				img.Pix[dst+1] = data[src+1]
				img.Pix[dst+2] = data[src+2]
			}
			if hasAlpha {
				img.Pix[dst+3] = data[src+3]
			} else {
				img.Pix[dst+3] = 0xFF
			}
		}
	}
	return img
}

func imageFromBayer8(data []byte, width, height int, pixelType uint32) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.SetRGBA(x, y, color.RGBA{
				R: averageBayerChannel(data, width, height, x, y, pixelType, 'r'),
				G: averageBayerChannel(data, width, height, x, y, pixelType, 'g'),
				B: averageBayerChannel(data, width, height, x, y, pixelType, 'b'),
				A: 0xFF,
			})
		}
	}
	return img
}

func averageBayerChannel(data []byte, width, height, x, y int, pixelType uint32, channel byte) uint8 {
	if bayerChannel(pixelType, x, y) == channel {
		return data[y*width+x]
	}

	var sum int
	var count int
	for radius := 1; radius <= 2 && count == 0; radius++ {
		for yy := maxIntValue(0, y-radius); yy <= minIntValue(height-1, y+radius); yy++ {
			for xx := maxIntValue(0, x-radius); xx <= minIntValue(width-1, x+radius); xx++ {
				if bayerChannel(pixelType, xx, yy) == channel {
					sum += int(data[yy*width+xx])
					count++
				}
			}
		}
	}
	if count == 0 {
		return data[y*width+x]
	}
	return uint8(sum / count)
}

func bayerChannel(pixelType uint32, x, y int) byte {
	evenRow := y%2 == 0
	evenCol := x%2 == 0

	switch pixelType {
	case PixelTypeBayerRG8:
		if evenRow && evenCol {
			return 'r'
		}
		if !evenRow && !evenCol {
			return 'b'
		}
		return 'g'
	case PixelTypeBayerGR8:
		if evenRow && !evenCol {
			return 'r'
		}
		if !evenRow && evenCol {
			return 'b'
		}
		return 'g'
	case PixelTypeBayerGB8:
		if !evenRow && evenCol {
			return 'r'
		}
		if evenRow && !evenCol {
			return 'b'
		}
		return 'g'
	case PixelTypeBayerBG8:
		if !evenRow && !evenCol {
			return 'r'
		}
		if evenRow && evenCol {
			return 'b'
		}
		return 'g'
	default:
		return 'g'
	}
}

func minIntValue(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxIntValue(a, b int) int {
	if a > b {
		return a
	}
	return b
}
