//go:build windows && amd64

package mvs

import "unsafe"

const (
	infoMaxBufferSize  = 64
	mvMaxSymbolicLen   = 64
	mvMaxDeviceNum     = 256
	mvMaxInterfaceNum  = 64
	mvMaxGenTLIFNum    = 256
	mvMaxGenTLDevNum   = 256
	mvMaxSerialPortNum = 64
	mvMaxEventNameSize = 128
	mvMaxSplitNum      = 8
	specialInfoRawSize = 604

	expectedMVCCDeviceInfoSize           = 636
	expectedMVCCDeviceInfoListSize       = 2056
	expectedMVCCStringValueSize          = 272
	expectedMVInterfaceInfoSize          = 712
	expectedMVInterfaceInfoListSize      = 520
	expectedMVGenTLIFInfoSize            = 228
	expectedMVGenTLIFInfoListSize        = 2056
	expectedMVGenTLDevInfoSize           = 612
	expectedMVGenTLDevInfoListSize       = 2056
	expectedMVCCEnumEntrySize            = 84
	expectedMVCCImageSize                = 56
	expectedMVSaveImageToFileParamExSize = 80
	expectedMVCCPixelConvertParamExSize  = 64
	expectedMVCCRotateImageParamSize     = 88
	expectedMVCCFlipImageParamSize       = 88
	expectedMVCCGammaParamSize           = 56
	expectedMVCCCCMParamSize             = 72
	expectedMVCCCCMParamExSize           = 76
	expectedMVCCContrastParamSize        = 80
	expectedMVCCFrameSpecInfoSize        = 124
	expectedMVCCPurpleFringingParamSize  = 80
	expectedMVCCISPConfigParamSize       = 72
	expectedMVCCHBDecodeParamSize        = 200
	expectedMVCCRecordParamSize          = 64
	expectedMVCCInputFrameInfoSize       = 48
	expectedMVEventOutInfoSize           = 232
	expectedMVCCFileAccessSize           = 144
	expectedMVCCFileAccessExSize         = 152
	expectedMVCCFileAccessProgressSize   = 48
	expectedMVGigeZoneInfoSize           = 48
	expectedMVGigeMultiPartInfoSize      = 104
	expectedMVOutputImageInfoSize        = 64
	expectedMVReconstructImageParamSize  = 568
	expectedMVCamlSerialPortSize         = 80
	expectedMVCamlSerialPortListSize     = 5140
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

type mvInterfaceInfo struct {
	TLayerType      uint32
	PCIEInfo        uint32
	InterfaceID     [infoMaxBufferSize]byte
	DisplayName     [infoMaxBufferSize]byte
	SerialNumber    [infoMaxBufferSize]byte
	ModelName       [infoMaxBufferSize]byte
	Manufacturer    [infoMaxBufferSize]byte
	DeviceVersion   [infoMaxBufferSize]byte
	UserDefinedName [infoMaxBufferSize]byte
	Reserved        [64]uint32
}

type mvInterfaceInfoList struct {
	InterfaceNum uint32
	_            uint32
	Interface    [mvMaxInterfaceNum]*mvInterfaceInfo
}

type mvGenTLIFInfo struct {
	InterfaceID [infoMaxBufferSize]byte
	TLType      [infoMaxBufferSize]byte
	DisplayName [infoMaxBufferSize]byte
	CtiIndex    uint32
	Reserved    [8]uint32
}

type mvGenTLIFInfoList struct {
	InterfaceNum uint32
	_            uint32
	Interface    [mvMaxGenTLIFNum]*mvGenTLIFInfo
}

type mvGenTLDevInfo struct {
	InterfaceID     [infoMaxBufferSize]byte
	DeviceID        [infoMaxBufferSize]byte
	VendorName      [infoMaxBufferSize]byte
	ModelName       [infoMaxBufferSize]byte
	TLType          [infoMaxBufferSize]byte
	DisplayName     [infoMaxBufferSize]byte
	UserDefinedName [infoMaxBufferSize]byte
	SerialNumber    [infoMaxBufferSize]byte
	DeviceVersion   [infoMaxBufferSize]byte
	CtiIndex        uint32
	Reserved        [8]uint32
}

type mvGenTLDevInfoList struct {
	DeviceNum uint32
	_         uint32
	Device    [mvMaxGenTLDevNum]*mvGenTLDevInfo
}

type mvGigeZoneInfo struct {
	Direction uint32
	_         uint32
	Zone      uintptr
	Length    uint64
	Reserved  [6]uint32
}

type mvGigePartDataInfo struct {
	Data [24]byte
}

type mvGigeMultiPartInfo struct {
	DataType         uint32
	DataFormat       uint32
	SourceID         uint32
	RegionID         uint32
	DataPurposeID    uint32
	Zones            uint32
	ZoneInfo         *mvGigeZoneInfo
	Length           uint64
	PartAddr         *byte
	DataTypeSpecific mvGigePartDataInfo
	Reserved         [8]uint32
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
	SubImageList          unsafe.Pointer
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

type mvCCImage struct {
	Width        uint32
	Height       uint32
	PixelType    uint32
	_            uint32
	ImageBuf     *byte
	ImageBufSize uint64
	ImageLen     uint64
	Reserved     [4]uint32
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

type mvCCRotateImageParam struct {
	PixelType     uint32
	Width         uint32
	Height        uint32
	_             uint32
	SrcData       *byte
	SrcDataLen    uint32
	_             uint32
	DstBuffer     *byte
	DstLen        uint32
	DstBufferSize uint32
	RotationAngle uint32
	Reserved      [8]uint32
}

type mvCCFlipImageParam struct {
	PixelType     uint32
	Width         uint32
	Height        uint32
	_             uint32
	SrcData       *byte
	SrcDataLen    uint32
	_             uint32
	DstBuffer     *byte
	DstLen        uint32
	DstBufferSize uint32
	FlipType      uint32
	Reserved      [8]uint32
}

type mvCCGammaParam struct {
	GammaType        uint32
	GammaValue       float32
	GammaCurveBuf    *byte
	GammaCurveBufLen uint32
	Reserved         [8]uint32
}

type mvCCCCMParam struct {
	CCMEnable byte
	_         [3]byte
	CCMat     [9]int32
	Reserved  [8]uint32
}

type mvCCCCMParamEx struct {
	CCMEnable byte
	_         [3]byte
	CCMat     [9]int32
	CCMScale  uint32
	Reserved  [8]uint32
}

type mvCCContrastParam struct {
	Width          uint32
	Height         uint32
	SrcBuffer      *byte
	SrcBufferLen   uint32
	PixelType      uint32
	DstBuffer      *byte
	DstBufferSize  uint32
	DstBufferLen   uint32
	ContrastFactor uint32
	Reserved       [8]uint32
}

type mvCCFrameSpecInfo struct {
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
	Reserved          [16]uint32
}

type mvCCPurpleFringingParam struct {
	Width         uint32
	Height        uint32
	SrcBuffer     *byte
	SrcBufferLen  uint32
	PixelType     uint32
	DstBuffer     *byte
	DstBufferSize uint32
	DstBufferLen  uint32
	KernelSize    uint32
	EdgeThreshold uint32
	Reserved      [8]uint32
}

type mvCCISPConfigParam struct {
	ConfigPath *byte
	Reserved   [16]uint32
}

type mvCCHBDecodeParam struct {
	SrcBuffer     *byte
	SrcLen        uint32
	Width         uint32
	Height        uint32
	_             uint32
	DstBuffer     *byte
	DstBufferSize uint32
	DstBufferLen  uint32
	DstPixelType  uint32
	FrameSpecInfo mvCCFrameSpecInfo
	Reserved      [8]uint32
}

type mvCCRecordParam struct {
	PixelType    uint32
	Width        uint16
	Height       uint16
	FrameRate    float32
	BitRate      uint32
	RecordFormat uint32
	_            uint32
	FilePath     *byte
	Reserved     [8]uint32
}

type mvCCInputFrameInfo struct {
	Data     *byte
	DataLen  uint32
	Reserved [8]uint32
}

type mvEventOutInfo struct {
	EventName     [mvMaxEventNameSize]byte
	EventID       uint16
	StreamChannel uint16
	BlockIDHigh   uint32
	BlockIDLow    uint32
	TimestampHigh uint32
	TimestampLow  uint32
	_             uint32
	EventData     *byte
	EventDataSize uint32
	Reserved      [16]uint32
}

type mvCCFileAccess struct {
	UserFileName *byte
	DevFileName  *byte
	Reserved     [32]uint32
}

type mvCCFileAccessEx struct {
	UserFileBuffer *byte
	FileBufferSize uint32
	FileBufferLen  uint32
	DevFileName    *byte
	Reserved       [32]uint32
}

type mvCCFileAccessProgress struct {
	Completed int64
	Total     int64
	Reserved  [8]uint32
}

type mvOutputImageInfo struct {
	Width      uint32
	Height     uint32
	PixelType  uint32
	_          uint32
	Buffer     *byte
	BufferLen  uint32
	BufferSize uint32
	Reserved   [8]uint32
}

type mvReconstructImageParam struct {
	Width             uint32
	Height            uint32
	PixelType         uint32
	_                 uint32
	SrcData           *byte
	SrcDataLen        uint32
	ExposureNum       uint32
	ReconstructMethod uint32
	_                 uint32
	DstBufferList     [mvMaxSplitNum]mvOutputImageInfo
	Reserved          [4]uint32
}

type mvCamlSerialPort struct {
	SerialPort [infoMaxBufferSize]byte
	Reserved   [4]uint32
}

type mvCamlSerialPortList struct {
	SerialPortNum uint32
	SerialPort    [mvMaxSerialPortNum]mvCamlSerialPort
	Reserved      [4]uint32
}

var (
	_ [expectedMVCCDeviceInfoSize - int(unsafe.Sizeof(mvCCDeviceInfo{}))]byte
	_ [int(unsafe.Sizeof(mvCCDeviceInfo{})) - expectedMVCCDeviceInfoSize]byte

	_ [expectedMVCCDeviceInfoListSize - int(unsafe.Sizeof(mvCCDeviceInfoList{}))]byte
	_ [int(unsafe.Sizeof(mvCCDeviceInfoList{})) - expectedMVCCDeviceInfoListSize]byte

	_ [expectedMVCCStringValueSize - int(unsafe.Sizeof(mvCCStringValue{}))]byte
	_ [int(unsafe.Sizeof(mvCCStringValue{})) - expectedMVCCStringValueSize]byte

	_ [expectedMVInterfaceInfoSize - int(unsafe.Sizeof(mvInterfaceInfo{}))]byte
	_ [int(unsafe.Sizeof(mvInterfaceInfo{})) - expectedMVInterfaceInfoSize]byte

	_ [expectedMVInterfaceInfoListSize - int(unsafe.Sizeof(mvInterfaceInfoList{}))]byte
	_ [int(unsafe.Sizeof(mvInterfaceInfoList{})) - expectedMVInterfaceInfoListSize]byte

	_ [expectedMVGenTLIFInfoSize - int(unsafe.Sizeof(mvGenTLIFInfo{}))]byte
	_ [int(unsafe.Sizeof(mvGenTLIFInfo{})) - expectedMVGenTLIFInfoSize]byte

	_ [expectedMVGenTLIFInfoListSize - int(unsafe.Sizeof(mvGenTLIFInfoList{}))]byte
	_ [int(unsafe.Sizeof(mvGenTLIFInfoList{})) - expectedMVGenTLIFInfoListSize]byte

	_ [expectedMVGenTLDevInfoSize - int(unsafe.Sizeof(mvGenTLDevInfo{}))]byte
	_ [int(unsafe.Sizeof(mvGenTLDevInfo{})) - expectedMVGenTLDevInfoSize]byte

	_ [expectedMVGenTLDevInfoListSize - int(unsafe.Sizeof(mvGenTLDevInfoList{}))]byte
	_ [int(unsafe.Sizeof(mvGenTLDevInfoList{})) - expectedMVGenTLDevInfoListSize]byte

	_ [expectedMVCCEnumEntrySize - int(unsafe.Sizeof(mvCCEnumEntry{}))]byte
	_ [int(unsafe.Sizeof(mvCCEnumEntry{})) - expectedMVCCEnumEntrySize]byte

	_ [expectedMVCCImageSize - int(unsafe.Sizeof(mvCCImage{}))]byte
	_ [int(unsafe.Sizeof(mvCCImage{})) - expectedMVCCImageSize]byte

	_ [expectedMVSaveImageToFileParamExSize - int(unsafe.Sizeof(mvSaveImageToFileParamEx{}))]byte
	_ [int(unsafe.Sizeof(mvSaveImageToFileParamEx{})) - expectedMVSaveImageToFileParamExSize]byte

	_ [expectedMVCCPixelConvertParamExSize - int(unsafe.Sizeof(mvCCPixelConvertParamEx{}))]byte
	_ [int(unsafe.Sizeof(mvCCPixelConvertParamEx{})) - expectedMVCCPixelConvertParamExSize]byte

	_ [expectedMVCCRotateImageParamSize - int(unsafe.Sizeof(mvCCRotateImageParam{}))]byte
	_ [int(unsafe.Sizeof(mvCCRotateImageParam{})) - expectedMVCCRotateImageParamSize]byte

	_ [expectedMVCCFlipImageParamSize - int(unsafe.Sizeof(mvCCFlipImageParam{}))]byte
	_ [int(unsafe.Sizeof(mvCCFlipImageParam{})) - expectedMVCCFlipImageParamSize]byte

	_ [expectedMVCCGammaParamSize - int(unsafe.Sizeof(mvCCGammaParam{}))]byte
	_ [int(unsafe.Sizeof(mvCCGammaParam{})) - expectedMVCCGammaParamSize]byte

	_ [expectedMVCCCCMParamSize - int(unsafe.Sizeof(mvCCCCMParam{}))]byte
	_ [int(unsafe.Sizeof(mvCCCCMParam{})) - expectedMVCCCCMParamSize]byte

	_ [expectedMVCCCCMParamExSize - int(unsafe.Sizeof(mvCCCCMParamEx{}))]byte
	_ [int(unsafe.Sizeof(mvCCCCMParamEx{})) - expectedMVCCCCMParamExSize]byte

	_ [expectedMVCCContrastParamSize - int(unsafe.Sizeof(mvCCContrastParam{}))]byte
	_ [int(unsafe.Sizeof(mvCCContrastParam{})) - expectedMVCCContrastParamSize]byte

	_ [expectedMVCCFrameSpecInfoSize - int(unsafe.Sizeof(mvCCFrameSpecInfo{}))]byte
	_ [int(unsafe.Sizeof(mvCCFrameSpecInfo{})) - expectedMVCCFrameSpecInfoSize]byte

	_ [expectedMVCCPurpleFringingParamSize - int(unsafe.Sizeof(mvCCPurpleFringingParam{}))]byte
	_ [int(unsafe.Sizeof(mvCCPurpleFringingParam{})) - expectedMVCCPurpleFringingParamSize]byte

	_ [expectedMVCCISPConfigParamSize - int(unsafe.Sizeof(mvCCISPConfigParam{}))]byte
	_ [int(unsafe.Sizeof(mvCCISPConfigParam{})) - expectedMVCCISPConfigParamSize]byte

	_ [expectedMVCCHBDecodeParamSize - int(unsafe.Sizeof(mvCCHBDecodeParam{}))]byte
	_ [int(unsafe.Sizeof(mvCCHBDecodeParam{})) - expectedMVCCHBDecodeParamSize]byte

	_ [expectedMVCCRecordParamSize - int(unsafe.Sizeof(mvCCRecordParam{}))]byte
	_ [int(unsafe.Sizeof(mvCCRecordParam{})) - expectedMVCCRecordParamSize]byte

	_ [expectedMVCCInputFrameInfoSize - int(unsafe.Sizeof(mvCCInputFrameInfo{}))]byte
	_ [int(unsafe.Sizeof(mvCCInputFrameInfo{})) - expectedMVCCInputFrameInfoSize]byte

	_ [expectedMVEventOutInfoSize - int(unsafe.Sizeof(mvEventOutInfo{}))]byte
	_ [int(unsafe.Sizeof(mvEventOutInfo{})) - expectedMVEventOutInfoSize]byte

	_ [expectedMVCCFileAccessSize - int(unsafe.Sizeof(mvCCFileAccess{}))]byte
	_ [int(unsafe.Sizeof(mvCCFileAccess{})) - expectedMVCCFileAccessSize]byte

	_ [expectedMVCCFileAccessExSize - int(unsafe.Sizeof(mvCCFileAccessEx{}))]byte
	_ [int(unsafe.Sizeof(mvCCFileAccessEx{})) - expectedMVCCFileAccessExSize]byte

	_ [expectedMVCCFileAccessProgressSize - int(unsafe.Sizeof(mvCCFileAccessProgress{}))]byte
	_ [int(unsafe.Sizeof(mvCCFileAccessProgress{})) - expectedMVCCFileAccessProgressSize]byte

	_ [expectedMVGigeZoneInfoSize - int(unsafe.Sizeof(mvGigeZoneInfo{}))]byte
	_ [int(unsafe.Sizeof(mvGigeZoneInfo{})) - expectedMVGigeZoneInfoSize]byte

	_ [expectedMVGigeMultiPartInfoSize - int(unsafe.Sizeof(mvGigeMultiPartInfo{}))]byte
	_ [int(unsafe.Sizeof(mvGigeMultiPartInfo{})) - expectedMVGigeMultiPartInfoSize]byte

	_ [expectedMVOutputImageInfoSize - int(unsafe.Sizeof(mvOutputImageInfo{}))]byte
	_ [int(unsafe.Sizeof(mvOutputImageInfo{})) - expectedMVOutputImageInfoSize]byte

	_ [expectedMVReconstructImageParamSize - int(unsafe.Sizeof(mvReconstructImageParam{}))]byte
	_ [int(unsafe.Sizeof(mvReconstructImageParam{})) - expectedMVReconstructImageParamSize]byte

	_ [expectedMVCamlSerialPortSize - int(unsafe.Sizeof(mvCamlSerialPort{}))]byte
	_ [int(unsafe.Sizeof(mvCamlSerialPort{})) - expectedMVCamlSerialPortSize]byte

	_ [expectedMVCamlSerialPortListSize - int(unsafe.Sizeof(mvCamlSerialPortList{}))]byte
	_ [int(unsafe.Sizeof(mvCamlSerialPortList{})) - expectedMVCamlSerialPortListSize]byte

	_ [expectedMVFrameOutInfoExSize - int(unsafe.Sizeof(mvFrameOutInfoEx{}))]byte
	_ [int(unsafe.Sizeof(mvFrameOutInfoEx{})) - expectedMVFrameOutInfoExSize]byte

	_ [expectedMVFrameOutSize - int(unsafe.Sizeof(mvFrameOut{}))]byte
	_ [int(unsafe.Sizeof(mvFrameOut{})) - expectedMVFrameOutSize]byte
)
