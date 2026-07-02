package mvs

import (
	"path/filepath"
	"strings"
)

type ImageType uint32

type InterpolationMethod int32

type GrabStrategy uint32

type RecordFormat uint32

type RotationAngle uint32

type FlipType uint32

type GammaType uint32

type MultiPartDataType uint32

type ImageReconstructMethod uint32

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

type RecordOptions struct {
	Path        string
	PixelType   uint32
	Width       uint32
	Height      uint32
	FrameRate   float32
	BitRateKbps uint32
	Format      RecordFormat
}

type EventInfo struct {
	Name          string
	EventID       uint16
	StreamChannel uint16
	BlockID       uint64
	Timestamp     uint64
	Data          []byte
}

type FileAccessProgress struct {
	Completed int64
	Total     int64
}

type CameraLinkSerialPort struct {
	Name string
}

type InterfaceInfo struct {
	Index              int
	TransportLayer     uint32
	TransportLayerName string
	PCIEInfo           uint32
	InterfaceID        string
	DisplayName        string
	SerialNumber       string
	ModelName          string
	Manufacturer       string
	DeviceVersion      string
	UserDefinedName    string
}

type GenTLInterfaceInfo struct {
	Index       int
	InterfaceID string
	TLType      string
	DisplayName string
	CtiIndex    uint32
}

type GenTLDeviceInfo struct {
	Index           int
	InterfaceID     string
	DeviceID        string
	VendorName      string
	ModelName       string
	TLType          string
	DisplayName     string
	UserDefinedName string
	SerialNumber    string
	DeviceVersion   string
	CtiIndex        uint32
}

type FramePart struct {
	DataType         MultiPartDataType
	DataFormat       uint32
	Width            uint32
	Height           uint32
	PixelType        uint32
	SourceID         uint32
	RegionID         uint32
	DataPurposeID    uint32
	Zones            uint32
	Length           uint64
	DataTypeSpecific [24]byte
	Data             []byte
}

type FrameSpecInfo struct {
	SecondCount       uint32
	CycleCount        uint32
	CycleOffset       uint32
	Gain              float32
	ExposureTime      float32
	AverageBrightness uint32
	Red               uint32
	Green             uint32
	Blue              uint32
	FrameCounter      uint32
	TriggerIndex      uint32
	Input             uint32
	Output            uint32
	OffsetX           uint16
	OffsetY           uint16
	FrameWidth        uint16
	FrameHeight       uint16
}

type GammaOptions struct {
	Type  GammaType
	Value float32
	Curve []byte
}

type CCMOptions struct {
	Enabled bool
	Matrix  [9]int32
}

type CCMOptionsEx struct {
	Enabled bool
	Matrix  [9]int32
	Scale   uint32
}

type ContrastOptions struct {
	Factor        uint32
	DstBufferSize int
}

type PurpleFringingOptions struct {
	KernelSize    uint32
	EdgeThreshold uint32
	DstBufferSize int
}

type HBDecodeOptions struct {
	DstBufferSize int
}

type ReconstructImageOptions struct {
	ExposureNum       uint32
	Method            ImageReconstructMethod
	DstBufferSizes    []int
	DefaultBufferSize int
}

type FrameCallback func(*Frame)

type EventCallback func(EventInfo)

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
