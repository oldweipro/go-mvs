//go:build windows && amd64

package mvsdk

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"unsafe"
)

const (
	defaultDLLName     = "MvCameraControl.dll"
	infoMaxBufferSize  = 64
	mvMaxDeviceNum     = 256
	specialInfoRawSize = 604

	expectedMVCCDeviceInfoSize     = 636
	expectedMVCCDeviceInfoListSize = 2056
	expectedMVFrameOutInfoExSize   = 256
	expectedMVFrameOutSize         = 328
)

type driver struct {
	dllPath string
	dll     *syscall.LazyDLL

	procInitialize       *syscall.LazyProc
	procFinalize         *syscall.LazyProc
	procGetSDKVersion    *syscall.LazyProc
	procEnumDevices      *syscall.LazyProc
	procCreateHandle     *syscall.LazyProc
	procDestroyHandle    *syscall.LazyProc
	procOpenDevice       *syscall.LazyProc
	procCloseDevice      *syscall.LazyProc
	procIsDeviceConnected *syscall.LazyProc
	procStartGrabbing    *syscall.LazyProc
	procStopGrabbing     *syscall.LazyProc
	procGetImageBuffer   *syscall.LazyProc
	procFreeImageBuffer  *syscall.LazyProc
	procGetOptimalPacketSize *syscall.LazyProc
	procGetFloatValue    *syscall.LazyProc
	procSetFloatValue    *syscall.LazyProc
	procGetIntValueEx    *syscall.LazyProc
	procSetIntValueEx    *syscall.LazyProc
	procGetEnumValue     *syscall.LazyProc
	procSetEnumValue     *syscall.LazyProc
	procSetCommandValue  *syscall.LazyProc
}

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
	InterfaceID       [infoMaxBufferSize]byte
	VendorName        [infoMaxBufferSize]byte
	ModelName         [infoMaxBufferSize]byte
	ManufacturerInfo  [infoMaxBufferSize]byte
	DeviceVersion     [infoMaxBufferSize]byte
	SerialNumber      [infoMaxBufferSize]byte
	UserDefinedName   [infoMaxBufferSize]byte
	DeviceID          [infoMaxBufferSize]byte
	Reserved          [7]uint32
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
	Width             uint16
	Height            uint16
	PixelType         uint32
	FrameNum          uint32
	DevTimestampHigh  uint32
	DevTimestampLow   uint32
	Reserved0         uint32
	HostTimestamp     int64
	FrameLen          uint32
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
	ChunkWidth        uint16
	ChunkHeight       uint16
	LostPacket        uint32
	UnparsedChunkNum  uint32
	_                 uint32
	UnparsedChunkList [8]byte
	ExtendWidth       uint32
	ExtendHeight      uint32
	FrameLenEx        uint64
	ExtraType         uint32
	SubImageNum       uint32
	SubImageList      [8]byte
	UserPtr           [8]byte
	FirstLineEncoderCount uint32
	LastLineEncoderCount  uint32
	Reserved          [24]uint32
}

type mvFrameOut struct {
	BufAddr     *byte
	StFrameInfo mvFrameOutInfoEx
	Reserved    [16]uint32
}

type mvCCFloatValue struct {
	Current  float32
	Max      float32
	Min      float32
	Reserved [4]uint32
}

type mvCCIntValueEx struct {
	Current   int64
	Max       int64
	Min       int64
	Increment int64
	Reserved  [16]uint32
}

type mvCCEnumValue struct {
	Current      uint32
	SupportedNum uint32
	SupportValue [64]uint32
	Reserved     [4]uint32
}

var (
	_ [expectedMVCCDeviceInfoSize-int(unsafe.Sizeof(mvCCDeviceInfo{}))]byte
	_ [int(unsafe.Sizeof(mvCCDeviceInfo{}))-expectedMVCCDeviceInfoSize]byte

	_ [expectedMVCCDeviceInfoListSize-int(unsafe.Sizeof(mvCCDeviceInfoList{}))]byte
	_ [int(unsafe.Sizeof(mvCCDeviceInfoList{}))-expectedMVCCDeviceInfoListSize]byte

	_ [expectedMVFrameOutInfoExSize-int(unsafe.Sizeof(mvFrameOutInfoEx{}))]byte
	_ [int(unsafe.Sizeof(mvFrameOutInfoEx{}))-expectedMVFrameOutInfoExSize]byte

	_ [expectedMVFrameOutSize-int(unsafe.Sizeof(mvFrameOut{}))]byte
	_ [int(unsafe.Sizeof(mvFrameOut{}))-expectedMVFrameOutSize]byte
)

func newDriver(dllPath string) (*driver, error) {
	if dllPath == "" {
		dllPath = resolveDLLPath()
	}

	dll := syscall.NewLazyDLL(dllPath)
	d := &driver{
		dllPath: dllPath,
		dll:     dll,
		procInitialize:          dll.NewProc("MV_CC_Initialize"),
		procFinalize:            dll.NewProc("MV_CC_Finalize"),
		procGetSDKVersion:       dll.NewProc("MV_CC_GetSDKVersion"),
		procEnumDevices:         dll.NewProc("MV_CC_EnumDevices"),
		procCreateHandle:        dll.NewProc("MV_CC_CreateHandle"),
		procDestroyHandle:       dll.NewProc("MV_CC_DestroyHandle"),
		procOpenDevice:          dll.NewProc("MV_CC_OpenDevice"),
		procCloseDevice:         dll.NewProc("MV_CC_CloseDevice"),
		procIsDeviceConnected:   dll.NewProc("MV_CC_IsDeviceConnected"),
		procStartGrabbing:       dll.NewProc("MV_CC_StartGrabbing"),
		procStopGrabbing:        dll.NewProc("MV_CC_StopGrabbing"),
		procGetImageBuffer:      dll.NewProc("MV_CC_GetImageBuffer"),
		procFreeImageBuffer:     dll.NewProc("MV_CC_FreeImageBuffer"),
		procGetOptimalPacketSize: dll.NewProc("MV_CC_GetOptimalPacketSize"),
		procGetFloatValue:       dll.NewProc("MV_CC_GetFloatValue"),
		procSetFloatValue:       dll.NewProc("MV_CC_SetFloatValue"),
		procGetIntValueEx:       dll.NewProc("MV_CC_GetIntValueEx"),
		procSetIntValueEx:       dll.NewProc("MV_CC_SetIntValueEx"),
		procGetEnumValue:        dll.NewProc("MV_CC_GetEnumValue"),
		procSetEnumValue:        dll.NewProc("MV_CC_SetEnumValue"),
		procSetCommandValue:     dll.NewProc("MV_CC_SetCommandValue"),
	}

	if err := dll.Load(); err != nil {
		return nil, fmt.Errorf("load %q: %w", dllPath, err)
	}
	return d, nil
}

func resolveDLLPath() string {
	if explicit := os.Getenv("MVS_SDK_DLL"); explicit != "" {
		return explicit
	}

	candidates := []string{
		filepath.Join(os.Getenv("ProgramFiles(x86)"), "Common Files", "MVS", "Runtime", "Win64_x64", defaultDLLName),
		filepath.Join(os.Getenv("ProgramFiles"), "Common Files", "MVS", "Runtime", "Win64_x64", defaultDLLName),
		defaultDLLName,
	}

	for _, candidate := range candidates {
		if candidate == "" {
			continue
		}
		if candidate == defaultDLLName {
			return candidate
		}
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}

	return defaultDLLName
}

func (d *driver) initialize() error {
	ret, _, _ := d.procInitialize.Call()
	return newSDKError("MV_CC_Initialize", ret)
}

func (d *driver) finalize() error {
	ret, _, _ := d.procFinalize.Call()
	return newSDKError("MV_CC_Finalize", ret)
}

func (d *driver) getSDKVersion() uint32 {
	ret, _, _ := d.procGetSDKVersion.Call()
	return uint32(ret)
}

func (d *driver) enumDevices(layerType uint32, list *mvCCDeviceInfoList) error {
	ret, _, _ := d.procEnumDevices.Call(
		uintptr(layerType),
		uintptr(unsafe.Pointer(list)),
	)
	return newSDKError("MV_CC_EnumDevices", ret)
}

func (d *driver) createHandle(info *mvCCDeviceInfo) (uintptr, error) {
	var handle uintptr
	ret, _, _ := d.procCreateHandle.Call(
		uintptr(unsafe.Pointer(&handle)),
		uintptr(unsafe.Pointer(info)),
	)
	if err := newSDKError("MV_CC_CreateHandle", ret); err != nil {
		return 0, err
	}
	return handle, nil
}

func (d *driver) destroyHandle(handle uintptr) error {
	ret, _, _ := d.procDestroyHandle.Call(handle)
	return newSDKError("MV_CC_DestroyHandle", ret)
}

func (d *driver) openDevice(handle uintptr, accessMode uint32, switchoverKey uint16) error {
	ret, _, _ := d.procOpenDevice.Call(handle, uintptr(accessMode), uintptr(switchoverKey))
	return newSDKError("MV_CC_OpenDevice", ret)
}

func (d *driver) closeDevice(handle uintptr) error {
	ret, _, _ := d.procCloseDevice.Call(handle)
	return newSDKError("MV_CC_CloseDevice", ret)
}

func (d *driver) isDeviceConnected(handle uintptr) bool {
	ret, _, _ := d.procIsDeviceConnected.Call(handle)
	return ret != 0
}

func (d *driver) startGrabbing(handle uintptr) error {
	ret, _, _ := d.procStartGrabbing.Call(handle)
	return newSDKError("MV_CC_StartGrabbing", ret)
}

func (d *driver) stopGrabbing(handle uintptr) error {
	ret, _, _ := d.procStopGrabbing.Call(handle)
	return newSDKError("MV_CC_StopGrabbing", ret)
}

func (d *driver) getImageBuffer(handle uintptr, frame *mvFrameOut, timeoutMs uint32) error {
	ret, _, _ := d.procGetImageBuffer.Call(
		handle,
		uintptr(unsafe.Pointer(frame)),
		uintptr(timeoutMs),
	)
	return newSDKError("MV_CC_GetImageBuffer", ret)
}

func (d *driver) freeImageBuffer(handle uintptr, frame *mvFrameOut) error {
	ret, _, _ := d.procFreeImageBuffer.Call(
		handle,
		uintptr(unsafe.Pointer(frame)),
	)
	return newSDKError("MV_CC_FreeImageBuffer", ret)
}

func (d *driver) getOptimalPacketSize(handle uintptr) int32 {
	ret, _, _ := d.procGetOptimalPacketSize.Call(handle)
	return int32(ret)
}

func (d *driver) getFloatValue(handle uintptr, key string, value *mvCCFloatValue) error {
	keyPtr, err := syscall.BytePtrFromString(key)
	if err != nil {
		return fmt.Errorf("build key %q: %w", key, err)
	}
	ret, _, _ := d.procGetFloatValue.Call(
		handle,
		uintptr(unsafe.Pointer(keyPtr)),
		uintptr(unsafe.Pointer(value)),
	)
	return newSDKError("MV_CC_GetFloatValue", ret)
}

func (d *driver) setFloatValue(handle uintptr, key string, value float32) error {
	keyPtr, err := syscall.BytePtrFromString(key)
	if err != nil {
		return fmt.Errorf("build key %q: %w", key, err)
	}
	ret, _, _ := d.procSetFloatValue.Call(
		handle,
		uintptr(unsafe.Pointer(keyPtr)),
		uintptr(*(*uint32)(unsafe.Pointer(&value))),
	)
	return newSDKError("MV_CC_SetFloatValue", ret)
}

func (d *driver) getIntValueEx(handle uintptr, key string, value *mvCCIntValueEx) error {
	keyPtr, err := syscall.BytePtrFromString(key)
	if err != nil {
		return fmt.Errorf("build key %q: %w", key, err)
	}
	ret, _, _ := d.procGetIntValueEx.Call(
		handle,
		uintptr(unsafe.Pointer(keyPtr)),
		uintptr(unsafe.Pointer(value)),
	)
	return newSDKError("MV_CC_GetIntValueEx", ret)
}

func (d *driver) setIntValueEx(handle uintptr, key string, value int64) error {
	keyPtr, err := syscall.BytePtrFromString(key)
	if err != nil {
		return fmt.Errorf("build key %q: %w", key, err)
	}
	ret, _, _ := d.procSetIntValueEx.Call(
		handle,
		uintptr(unsafe.Pointer(keyPtr)),
		uintptr(value),
	)
	return newSDKError("MV_CC_SetIntValueEx", ret)
}

func (d *driver) getEnumValue(handle uintptr, key string, value *mvCCEnumValue) error {
	keyPtr, err := syscall.BytePtrFromString(key)
	if err != nil {
		return fmt.Errorf("build key %q: %w", key, err)
	}
	ret, _, _ := d.procGetEnumValue.Call(
		handle,
		uintptr(unsafe.Pointer(keyPtr)),
		uintptr(unsafe.Pointer(value)),
	)
	return newSDKError("MV_CC_GetEnumValue", ret)
}

func (d *driver) setEnumValue(handle uintptr, key string, value uint32) error {
	keyPtr, err := syscall.BytePtrFromString(key)
	if err != nil {
		return fmt.Errorf("build key %q: %w", key, err)
	}
	ret, _, _ := d.procSetEnumValue.Call(
		handle,
		uintptr(unsafe.Pointer(keyPtr)),
		uintptr(value),
	)
	return newSDKError("MV_CC_SetEnumValue", ret)
}

func (d *driver) setCommandValue(handle uintptr, key string) error {
	keyPtr, err := syscall.BytePtrFromString(key)
	if err != nil {
		return fmt.Errorf("build key %q: %w", key, err)
	}
	ret, _, _ := d.procSetCommandValue.Call(
		handle,
		uintptr(unsafe.Pointer(keyPtr)),
	)
	return newSDKError("MV_CC_SetCommandValue", ret)
}

func transportLayerName(layer uint32) string {
	switch layer {
	case DeviceTypeGigE:
		return "GigE"
	case DeviceTypeUSB:
		return "USB3"
	case DeviceTypeCameraLink:
		return "CameraLink"
	case DeviceTypeGentlGigE:
		return "GenTL GigE"
	case DeviceTypeGentlCameraLink:
		return "GenTL CameraLink"
	case DeviceTypeGentlCXP:
		return "GenTL CoaXPress"
	case DeviceTypeGentlXOF:
		return "GenTL XoF"
	case DeviceTypeGentlVirtual:
		return "GenTL Virtual"
	case DeviceTypeVirtualGigE:
		return "Virtual GigE"
	case DeviceTypeVirtualUSB:
		return "Virtual USB"
	default:
		return fmt.Sprintf("Unknown(0x%08X)", layer)
	}
}

func deviceInfoFromRaw(index int, raw mvCCDeviceInfo) DeviceInfo {
	info := DeviceInfo{
		Index:              index,
		TransportLayer:     raw.TLayerType,
		TransportLayerName: transportLayerName(raw.TLayerType),
		raw:                raw,
	}

	switch raw.TLayerType {
	case DeviceTypeGigE, DeviceTypeGentlGigE:
		spec := (*mvGigeDeviceInfo)(unsafe.Pointer(&raw.SpecialInfo[0]))
		info.ModelName = byteString(spec.ModelName[:])
		info.SerialNumber = byteString(spec.SerialNumber[:])
		info.UserDefinedName = byteString(spec.UserDefinedName[:])
		info.ManufacturerName = byteString(spec.ManufacturerName[:])
		info.CurrentIP = ipv4String(spec.CurrentIP)
	case DeviceTypeUSB:
		spec := (*mvUSB3DeviceInfo)(unsafe.Pointer(&raw.SpecialInfo[0]))
		info.ModelName = byteString(spec.ModelName[:])
		info.SerialNumber = byteString(spec.SerialNumber[:])
		info.UserDefinedName = byteString(spec.UserDefinedName[:])
		info.ManufacturerName = byteString(spec.ManufacturerName[:])
	case DeviceTypeCameraLink:
		spec := (*mvCamLDevInfo)(unsafe.Pointer(&raw.SpecialInfo[0]))
		info.ModelName = byteString(spec.ModelName[:])
		info.SerialNumber = byteString(spec.SerialNumber[:])
		info.ManufacturerName = byteString(spec.ManufacturerName[:])
	case DeviceTypeGentlCameraLink:
		spec := (*mvCmlDeviceInfo)(unsafe.Pointer(&raw.SpecialInfo[0]))
		info.InterfaceID = byteString(spec.InterfaceID[:])
		info.ModelName = byteString(spec.ModelName[:])
		info.SerialNumber = byteString(spec.SerialNumber[:])
		info.UserDefinedName = byteString(spec.UserDefinedName[:])
		info.DeviceID = byteString(spec.DeviceID[:])
		info.ManufacturerName = byteString(spec.VendorName[:])
	case DeviceTypeGentlCXP:
		spec := (*mvCxpDeviceInfo)(unsafe.Pointer(&raw.SpecialInfo[0]))
		info.InterfaceID = byteString(spec.InterfaceID[:])
		info.ModelName = byteString(spec.ModelName[:])
		info.SerialNumber = byteString(spec.SerialNumber[:])
		info.UserDefinedName = byteString(spec.UserDefinedName[:])
		info.DeviceID = byteString(spec.DeviceID[:])
		info.ManufacturerName = byteString(spec.VendorName[:])
	case DeviceTypeGentlXOF:
		spec := (*mvXofDeviceInfo)(unsafe.Pointer(&raw.SpecialInfo[0]))
		info.InterfaceID = byteString(spec.InterfaceID[:])
		info.ModelName = byteString(spec.ModelName[:])
		info.SerialNumber = byteString(spec.SerialNumber[:])
		info.UserDefinedName = byteString(spec.UserDefinedName[:])
		info.DeviceID = byteString(spec.DeviceID[:])
		info.ManufacturerName = byteString(spec.VendorName[:])
	case DeviceTypeGentlVirtual:
		spec := (*mvGentlVirDeviceInfo)(unsafe.Pointer(&raw.SpecialInfo[0]))
		info.InterfaceID = byteString(spec.InterfaceID[:])
		info.ModelName = byteString(spec.ModelName[:])
		info.SerialNumber = byteString(spec.SerialNumber[:])
		info.UserDefinedName = byteString(spec.UserDefinedName[:])
		info.DeviceID = byteString(spec.DeviceID[:])
		info.ManufacturerName = byteString(spec.VendorName[:])
	}

	return info
}

func byteString(buf []byte) string {
	buf = bytes.TrimRight(buf, "\x00")
	return string(buf)
}

func ipv4String(ip uint32) string {
	return fmt.Sprintf("%d.%d.%d.%d", byte(ip>>24), byte(ip>>16), byte(ip>>8), byte(ip))
}
