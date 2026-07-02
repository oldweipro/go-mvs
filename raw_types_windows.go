//go:build windows && amd64

package mvs

import "unsafe"

const (
	infoMaxBufferSize  = 64
	mvMaxSymbolicLen   = 64
	mvMaxDeviceNum     = 256
	specialInfoRawSize = 604

	expectedMVCCDeviceInfoSize           = 636
	expectedMVCCDeviceInfoListSize       = 2056
	expectedMVCCStringValueSize          = 272
	expectedMVCCEnumEntrySize            = 84
	expectedMVSaveImageToFileParamExSize = 80
	expectedMVCCPixelConvertParamExSize  = 64
	expectedMVFrameOutInfoExSize         = 256
	expectedMVFrameOutSize               = 328
)

type mvGigeDeviceInfo struct {
	IpCfgOption              uint32
	IpCfgCurrent             uint32
	CurrentIP                uint32
	CurrentSubNetMask        uint32
	DefaultGateway           uint32
	ManufacturerName         [32]byte
	ModelName                [32]byte
	DeviceVersion            [32]byte
	ManufacturerSpecificInfo [48]byte
	SerialNumber             [16]byte
	UserDefinedName          [16]byte
	NetExport                uint32
	Reserved                 [4]uint32
}

type mvUSB3DeviceInfo struct {
	CtrlInEndPoint   byte
	CtrlOutEndPoint  byte
	StreamEndPoint   byte
	EventEndPoint    byte
	VendorID         uint16
	ProductID        uint16
	DeviceNumber     uint32
	DeviceGUID       [infoMaxBufferSize]byte
	VendorName       [infoMaxBufferSize]byte
	ModelName        [infoMaxBufferSize]byte
	FamilyName       [infoMaxBufferSize]byte
	DeviceVersion    [infoMaxBufferSize]byte
	ManufacturerName [infoMaxBufferSize]byte
	SerialNumber     [infoMaxBufferSize]byte
	UserDefinedName  [infoMaxBufferSize]byte
	BCDUSB           uint32
	DeviceAddress    uint32
	Reserved         [2]uint32
}

type mvCamLDevInfo struct {
	PortID           [infoMaxBufferSize]byte
	ModelName        [infoMaxBufferSize]byte
	FamilyName       [infoMaxBufferSize]byte
	DeviceVersion    [infoMaxBufferSize]byte
	ManufacturerName [infoMaxBufferSize]byte
	SerialNumber     [infoMaxBufferSize]byte
	Reserved         [38]uint32
}

type mvCxpDeviceInfo struct {
	InterfaceID      [infoMaxBufferSize]byte
	VendorName       [infoMaxBufferSize]byte
	ModelName        [infoMaxBufferSize]byte
	ManufacturerInfo [infoMaxBufferSize]byte
	DeviceVersion    [infoMaxBufferSize]byte
	SerialNumber     [infoMaxBufferSize]byte
	UserDefinedName  [infoMaxBufferSize]byte
	DeviceID         [infoMaxBufferSize]byte
	Reserved         [7]uint32
}

type mvCmlDeviceInfo = mvCxpDeviceInfo
type mvXofDeviceInfo = mvCxpDeviceInfo

type mvGentlVirDeviceInfo struct {
	InterfaceID      [infoMaxBufferSize]byte
	VendorName       [infoMaxBufferSize]byte
	ModelName        [infoMaxBufferSize]byte
	ManufacturerInfo [infoMaxBufferSize]byte
	DeviceVersion    [infoMaxBufferSize]byte
	SerialNumber     [infoMaxBufferSize]byte
	UserDefinedName  [infoMaxBufferSize]byte
	DeviceID         [infoMaxBufferSize]byte
	TLType           [infoMaxBufferSize]byte
	Reserved         [7]uint32
}

type mvCCDeviceInfo struct {
	MajorVer    uint16
	MinorVer    uint16
	MacAddrHigh uint32
	MacAddrLow  uint32
	TLayerType  uint32
	DevTypeInfo uint32
	Reserved    [3]uint32
	SpecialInfo [specialInfoRawSize]byte
}

type mvCCDeviceInfoList struct {
	DeviceNum  uint32
	_          uint32
	DeviceInfo [mvMaxDeviceNum]*mvCCDeviceInfo
}

type mvFrameOutInfoEx struct {
	Width                 uint16
	Height                uint16
	PixelType             uint32
	FrameNum              uint32
	DevTimestampHigh      uint32
	DevTimestampLow       uint32
	Reserved0             uint32
	HostTimestamp         int64
	FrameLen              uint32
	SecondCount           uint32
	CycleCount            uint32
	CycleOffset           uint32
	Gain                  float32
	ExposureTime          float32
	AverageBrightness     uint32
	Red                   uint32
	Green                 uint32
	Blue                  uint32
	FrameCounter          uint32
	TriggerIndex          uint32
	Input                 uint32
	Output                uint32
	OffsetX               uint16
	OffsetY               uint16
	ChunkWidth            uint16
	ChunkHeight           uint16
	LostPacket            uint32
	UnparsedChunkNum      uint32
	UnparsedChunkList     [8]byte
	ExtendWidth           uint32
	ExtendHeight          uint32
	FrameLenEx            uint64
	ExtraType             uint32
	SubImageNum           uint32
	SubImageList          [8]byte
	UserPtr               [8]byte
	FirstLineEncoderCount uint32
	LastLineEncoderCount  uint32
	Reserved              [24]uint32
}

type mvFrameOut struct {
	BufAddr     *byte
	StFrameInfo mvFrameOutInfoEx
	Reserved    [16]uint32
}

type mvCCEnumEntry struct {
	Value    uint32
	Symbolic [mvMaxSymbolicLen]byte
	Reserved [4]uint32
}

type mvSaveImageToFileParamEx struct {
	Width       uint32
	Height      uint32
	PixelType   uint32
	Data        *byte
	DataLen     uint32
	ImageType   uint32
	ImagePath   *byte
	Quality     uint32
	MethodValue int32
	Reserved    [8]uint32
}

type mvCCPixelConvertParamEx struct {
	Width         uint32
	Height        uint32
	SrcPixelType  uint32
	SrcData       *byte
	SrcDataLen    uint32
	DstPixelType  uint32
	DstBuffer     *byte
	DstLen        uint32
	DstBufferSize uint32
	Reserved      [4]uint32
}

var (
	_ [expectedMVCCDeviceInfoSize - int(unsafe.Sizeof(mvCCDeviceInfo{}))]byte
	_ [int(unsafe.Sizeof(mvCCDeviceInfo{})) - expectedMVCCDeviceInfoSize]byte

	_ [expectedMVCCDeviceInfoListSize - int(unsafe.Sizeof(mvCCDeviceInfoList{}))]byte
	_ [int(unsafe.Sizeof(mvCCDeviceInfoList{})) - expectedMVCCDeviceInfoListSize]byte

	_ [expectedMVCCStringValueSize - int(unsafe.Sizeof(mvCCStringValue{}))]byte
	_ [int(unsafe.Sizeof(mvCCStringValue{})) - expectedMVCCStringValueSize]byte

	_ [expectedMVCCEnumEntrySize - int(unsafe.Sizeof(mvCCEnumEntry{}))]byte
	_ [int(unsafe.Sizeof(mvCCEnumEntry{})) - expectedMVCCEnumEntrySize]byte

	_ [expectedMVSaveImageToFileParamExSize - int(unsafe.Sizeof(mvSaveImageToFileParamEx{}))]byte
	_ [int(unsafe.Sizeof(mvSaveImageToFileParamEx{})) - expectedMVSaveImageToFileParamExSize]byte

	_ [expectedMVCCPixelConvertParamExSize - int(unsafe.Sizeof(mvCCPixelConvertParamEx{}))]byte
	_ [int(unsafe.Sizeof(mvCCPixelConvertParamEx{})) - expectedMVCCPixelConvertParamExSize]byte

	_ [expectedMVFrameOutInfoExSize - int(unsafe.Sizeof(mvFrameOutInfoEx{}))]byte
	_ [int(unsafe.Sizeof(mvFrameOutInfoEx{})) - expectedMVFrameOutInfoExSize]byte

	_ [expectedMVFrameOutSize - int(unsafe.Sizeof(mvFrameOut{}))]byte
	_ [int(unsafe.Sizeof(mvFrameOut{})) - expectedMVFrameOutSize]byte
)
