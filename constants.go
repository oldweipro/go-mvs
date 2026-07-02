package mvs

const (
	MVOK = 0x00000000

	ErrHandle       = 0x80000000
	ErrSupport      = 0x80000001
	ErrBufferOver   = 0x80000002
	ErrCallOrder    = 0x80000003
	ErrParameter    = 0x80000004
	ErrResource     = 0x80000006
	ErrNoData       = 0x80000007
	ErrPrecondition = 0x80000008
	ErrVersion      = 0x80000009
	ErrNoEnoughBuf  = 0x8000000A
	ErrAbnormalImg  = 0x8000000B
	ErrLoadLibrary  = 0x8000000C
	ErrNoOutBuf     = 0x8000000D
	ErrBusy         = 0x80000204
)

const (
	DeviceTypeUnknown           uint32 = 0x00000000
	DeviceTypeGigE              uint32 = 0x00000001
	DeviceType1394              uint32 = 0x00000002
	DeviceTypeUSB               uint32 = 0x00000004
	DeviceTypeCameraLink        uint32 = 0x00000008
	DeviceTypeVirtualGigE       uint32 = 0x00000010
	DeviceTypeVirtualUSB        uint32 = 0x00000020
	DeviceTypeGentlGigE         uint32 = 0x00000040
	DeviceTypeGentlCameraLink   uint32 = 0x00000080
	DeviceTypeGentlCXP          uint32 = 0x00000100
	DeviceTypeGentlXOF          uint32 = 0x00000200
	DeviceTypeGentlVirtual      uint32 = 0x00000800
	DefaultDeviceTransportLayer uint32 = DeviceTypeGigE | DeviceTypeUSB | DeviceTypeGentlCameraLink | DeviceTypeGentlCXP | DeviceTypeGentlXOF
)

const (
	AccessExclusive              uint32 = 1
	AccessExclusiveWithSwitch    uint32 = 2
	AccessControl                uint32 = 3
	AccessControlWithSwitch      uint32 = 4
	AccessControlSwitchEnable    uint32 = 5
	AccessControlSwitchEnableKey uint32 = 6
	AccessMonitor                uint32 = 7
)

const (
	TriggerModeOff uint32 = 0
	TriggerModeOn  uint32 = 1

	TriggerSourceSoftware uint32 = 7
)

const (
	ImageTypeUndefined ImageType = 0
	ImageTypeBMP       ImageType = 1
	ImageTypeJPEG      ImageType = 2
	ImageTypePNG       ImageType = 3
	ImageTypeTIFF      ImageType = 4
)

const (
	InterpolationFast        InterpolationMethod = 0
	InterpolationBalanced    InterpolationMethod = 1
	InterpolationOptimal     InterpolationMethod = 2
	InterpolationOptimalPlus InterpolationMethod = 3
)

const (
	GrabStrategyOneByOne         GrabStrategy = 0
	GrabStrategyLatestImagesOnly GrabStrategy = 1
	GrabStrategyLatestImages     GrabStrategy = 2
	GrabStrategyUpcomingImage    GrabStrategy = 3
)

const (
	NodeTriggerMode          = "TriggerMode"
	NodeTriggerSource        = "TriggerSource"
	NodeTriggerSoftware      = "TriggerSoftware"
	NodeAcquisitionMode      = "AcquisitionMode"
	NodePixelFormat          = "PixelFormat"
	NodeWidth                = "Width"
	NodeHeight               = "Height"
	NodeGevSCPSPacketSize    = "GevSCPSPacketSize"
	NodeAcquisitionFrameRate = "AcquisitionFrameRate"
	NodeExposureTime         = "ExposureTime"
	NodeGain                 = "Gain"
	NodeDeviceUserID         = "DeviceUserID"
)

const (
	PixelTypeUndefined uint32 = 0xFFFFFFFF

	PixelTypeMono1p       uint32 = 0x01010037
	PixelTypeMono2p       uint32 = 0x01020038
	PixelTypeMono4p       uint32 = 0x01040039
	PixelTypeMono8        uint32 = 0x01080001
	PixelTypeMono8Signed  uint32 = 0x01080002
	PixelTypeMono10       uint32 = 0x01100003
	PixelTypeMono10Packed uint32 = 0x010C0004
	PixelTypeMono12       uint32 = 0x01100005
	PixelTypeMono12Packed uint32 = 0x010C0006
	PixelTypeMono14       uint32 = 0x01100025
	PixelTypeMono16       uint32 = 0x01100007

	PixelTypeBayerGR8        uint32 = 0x01080008
	PixelTypeBayerRG8        uint32 = 0x01080009
	PixelTypeBayerGB8        uint32 = 0x0108000A
	PixelTypeBayerBG8        uint32 = 0x0108000B
	PixelTypeBayerRBGG8      uint32 = 0x01080046
	PixelTypeBayerGR10       uint32 = 0x0110000C
	PixelTypeBayerRG10       uint32 = 0x0110000D
	PixelTypeBayerGB10       uint32 = 0x0110000E
	PixelTypeBayerBG10       uint32 = 0x0110000F
	PixelTypeBayerGR12       uint32 = 0x01100010
	PixelTypeBayerRG12       uint32 = 0x01100011
	PixelTypeBayerGB12       uint32 = 0x01100012
	PixelTypeBayerBG12       uint32 = 0x01100013
	PixelTypeBayerGR10Packed uint32 = 0x010C0026
	PixelTypeBayerRG10Packed uint32 = 0x010C0027
	PixelTypeBayerGB10Packed uint32 = 0x010C0028
	PixelTypeBayerBG10Packed uint32 = 0x010C0029
	PixelTypeBayerGR12Packed uint32 = 0x010C002A
	PixelTypeBayerRG12Packed uint32 = 0x010C002B
	PixelTypeBayerGB12Packed uint32 = 0x010C002C
	PixelTypeBayerBG12Packed uint32 = 0x010C002D
	PixelTypeBayerGR16       uint32 = 0x0110002E
	PixelTypeBayerRG16       uint32 = 0x0110002F
	PixelTypeBayerGB16       uint32 = 0x01100030
	PixelTypeBayerBG16       uint32 = 0x01100031

	PixelTypeRGB8Packed    uint32 = 0x02180014
	PixelTypeBGR8Packed    uint32 = 0x02180015
	PixelTypeRGBA8Packed   uint32 = 0x02200016
	PixelTypeBGRA8Packed   uint32 = 0x02200017
	PixelTypeRGB10Packed   uint32 = 0x02300018
	PixelTypeBGR10Packed   uint32 = 0x02300019
	PixelTypeRGB12Packed   uint32 = 0x0230001A
	PixelTypeBGR12Packed   uint32 = 0x0230001B
	PixelTypeRGB16Packed   uint32 = 0x02300033
	PixelTypeBGR16Packed   uint32 = 0x0230004B
	PixelTypeRGBA16Packed  uint32 = 0x02400064
	PixelTypeBGRA16Packed  uint32 = 0x02400051
	PixelTypeRGB10V1Packed uint32 = 0x0220001C
	PixelTypeRGB10V2Packed uint32 = 0x0220001D
	PixelTypeRGB12V1Packed uint32 = 0x02240034
	PixelTypeRGB565Packed  uint32 = 0x02100035
	PixelTypeBGR565Packed  uint32 = 0x02100036

	PixelTypeYUV411Packed         uint32 = 0x020C001E
	PixelTypeYUV422Packed         uint32 = 0x0210001F
	PixelTypeYUV422YUYVPacked     uint32 = 0x02100032
	PixelTypeYUV444Packed         uint32 = 0x02180020
	PixelTypeYCbCr8CbYCr          uint32 = 0x0218003A
	PixelTypeYCbCr4228            uint32 = 0x0210003B
	PixelTypeYCbCr4228CbYCrY      uint32 = 0x02100043
	PixelTypeYCbCr4118CbYYCrYY    uint32 = 0x020C003C
	PixelTypeYCbCr6018CbYCr       uint32 = 0x0218003D
	PixelTypeYCbCr6014228         uint32 = 0x0210003E
	PixelTypeYCbCr6014228CbYCrY   uint32 = 0x02100044
	PixelTypeYCbCr6014118CbYYCrYY uint32 = 0x020C003F
	PixelTypeYCbCr7098CbYCr       uint32 = 0x02180040
	PixelTypeYCbCr7094228         uint32 = 0x02100041
	PixelTypeYCbCr7094228CbYCrY   uint32 = 0x02100045
	PixelTypeYCbCr7094118CbYYCrYY uint32 = 0x020C0042
	PixelTypeYUV420SPNV12         uint32 = 0x020C8001
	PixelTypeYUV420SPNV21         uint32 = 0x020C8002

	PixelTypeRGB8Planar  uint32 = 0x02180021
	PixelTypeRGB10Planar uint32 = 0x02300022
	PixelTypeRGB12Planar uint32 = 0x02300023
	PixelTypeRGB16Planar uint32 = 0x02300024

	PixelTypeJPEG                    uint32 = 0x80180001
	PixelTypeCoord3DABC32f           uint32 = 0x026000C0
	PixelTypeCoord3DABC32fPlanar     uint32 = 0x026000C1
	PixelTypeCoord3DAC32f            uint32 = 0x022800C2
	PixelTypeCoord3DDepthPlusMask    uint32 = 0x821C0001
	PixelTypeCoord3DABC32            uint32 = 0x82603001
	PixelTypeCoord3DAB32f            uint32 = 0x82403002
	PixelTypeCoord3DAB32             uint32 = 0x82403003
	PixelTypeCoord3DAC32f64          uint32 = 0x024000C2
	PixelTypeCoord3DAC32fPlanar      uint32 = 0x024000C3
	PixelTypeCoord3DAC32             uint32 = 0x82403004
	PixelTypeCoord3DA32f             uint32 = 0x012000BD
	PixelTypeCoord3DA32              uint32 = 0x81203005
	PixelTypeCoord3DC32f             uint32 = 0x012000BF
	PixelTypeCoord3DC32              uint32 = 0x81203006
	PixelTypeCoord3DABC16            uint32 = 0x023000B9
	PixelTypeCoord3DC16              uint32 = 0x011000B8
	PixelTypeFloat32                 uint32 = 0x81200001
	PixelTypeHighBandwidthMono8      uint32 = 0x81080001
	PixelTypeHighBandwidthMono10     uint32 = 0x81100003
	PixelTypeHighBandwidthMono10Pack uint32 = 0x810C0004
	PixelTypeHighBandwidthMono12     uint32 = 0x81100005
	PixelTypeHighBandwidthMono12Pack uint32 = 0x810C0006
	PixelTypeHighBandwidthMono16     uint32 = 0x81100007

	PixelTypeHighBandwidthBayerGR8        uint32 = 0x81080008
	PixelTypeHighBandwidthBayerRG8        uint32 = 0x81080009
	PixelTypeHighBandwidthBayerGB8        uint32 = 0x8108000A
	PixelTypeHighBandwidthBayerBG8        uint32 = 0x8108000B
	PixelTypeHighBandwidthBayerRBGG8      uint32 = 0x81080046
	PixelTypeHighBandwidthBayerGR10       uint32 = 0x8110000C
	PixelTypeHighBandwidthBayerRG10       uint32 = 0x8110000D
	PixelTypeHighBandwidthBayerGB10       uint32 = 0x8110000E
	PixelTypeHighBandwidthBayerBG10       uint32 = 0x8110000F
	PixelTypeHighBandwidthBayerGR12       uint32 = 0x81100010
	PixelTypeHighBandwidthBayerRG12       uint32 = 0x81100011
	PixelTypeHighBandwidthBayerGB12       uint32 = 0x81100012
	PixelTypeHighBandwidthBayerBG12       uint32 = 0x81100013
	PixelTypeHighBandwidthBayerGR10Packed uint32 = 0x810C0026
	PixelTypeHighBandwidthBayerRG10Packed uint32 = 0x810C0027
	PixelTypeHighBandwidthBayerGB10Packed uint32 = 0x810C0028
	PixelTypeHighBandwidthBayerBG10Packed uint32 = 0x810C0029
	PixelTypeHighBandwidthBayerGR12Packed uint32 = 0x810C002A
	PixelTypeHighBandwidthBayerRG12Packed uint32 = 0x810C002B
	PixelTypeHighBandwidthBayerGB12Packed uint32 = 0x810C002C
	PixelTypeHighBandwidthBayerBG12Packed uint32 = 0x810C002D

	PixelTypeHighBandwidthYUV422Packed     uint32 = 0x8210001F
	PixelTypeHighBandwidthYUV422YUYVPacked uint32 = 0x82100032
	PixelTypeHighBandwidthRGB8Packed       uint32 = 0x82180014
	PixelTypeHighBandwidthBGR8Packed       uint32 = 0x82180015
	PixelTypeHighBandwidthRGBA8Packed      uint32 = 0x82200016
	PixelTypeHighBandwidthBGRA8Packed      uint32 = 0x82200017
	PixelTypeHighBandwidthRGB16Packed      uint32 = 0x82300033
	PixelTypeHighBandwidthBGR16Packed      uint32 = 0x8230004B
	PixelTypeHighBandwidthRGBA16Packed     uint32 = 0x82400064
	PixelTypeHighBandwidthBGRA16Packed     uint32 = 0x82400051
)
