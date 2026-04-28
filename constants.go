//go:build windows && amd64

package mvsdk

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
	AccessExclusive                uint32 = 1
	AccessExclusiveWithSwitch      uint32 = 2
	AccessControl                  uint32 = 3
	AccessControlWithSwitch        uint32 = 4
	AccessControlSwitchEnable      uint32 = 5
	AccessControlSwitchEnableKey   uint32 = 6
	AccessMonitor                  uint32 = 7
)

const (
	TriggerModeOff uint32 = 0
	TriggerModeOn  uint32 = 1

	TriggerSourceSoftware uint32 = 7
)

const (
	NodeTriggerMode           = "TriggerMode"
	NodeTriggerSource         = "TriggerSource"
	NodeTriggerSoftware       = "TriggerSoftware"
	NodeGevSCPSPacketSize     = "GevSCPSPacketSize"
	NodeAcquisitionFrameRate  = "AcquisitionFrameRate"
	NodeExposureTime          = "ExposureTime"
	NodeGain                  = "Gain"
)
