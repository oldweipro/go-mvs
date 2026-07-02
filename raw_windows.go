//go:build windows && amd64

package mvs

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"syscall"
	"unsafe"
)

const defaultDLLName = "MvCameraControl.dll"

type driver struct {
	dllPath string
	dll     *syscall.LazyDLL

	procInitialize           *syscall.LazyProc
	procFinalize             *syscall.LazyProc
	procGetSDKVersion        *syscall.LazyProc
	procEnumDevices          *syscall.LazyProc
	procEnumDevicesByIF      *syscall.LazyProc
	procCreateHandle         *syscall.LazyProc
	procCreateHandleByGenTL  *syscall.LazyProc
	procDestroyHandle        *syscall.LazyProc
	procIsDeviceAccessible   *syscall.LazyProc
	procOpenDevice           *syscall.LazyProc
	procCloseDevice          *syscall.LazyProc
	procIsDeviceConnected    *syscall.LazyProc
	procEnumInterfaces       *syscall.LazyProc
	procCreateInterface      *syscall.LazyProc
	procCreateInterfaceByID  *syscall.LazyProc
	procOpenInterface        *syscall.LazyProc
	procCloseInterface       *syscall.LazyProc
	procDestroyInterface     *syscall.LazyProc
	procRegisterCallbackEx   *syscall.LazyProc
	procStartGrabbing        *syscall.LazyProc
	procStopGrabbing         *syscall.LazyProc
	procGetImageBuffer       *syscall.LazyProc
	procFreeImageBuffer      *syscall.LazyProc
	procGetOneFrameTimeout   *syscall.LazyProc
	procClearImageBuffer     *syscall.LazyProc
	procSetImageNodeNum      *syscall.LazyProc
	procSetGrabStrategy      *syscall.LazyProc
	procSetOutputQueueSize   *syscall.LazyProc
	procGetOptimalPacketSize *syscall.LazyProc
	procGetFloatValue        *syscall.LazyProc
	procSetFloatValue        *syscall.LazyProc
	procGetIntValueEx        *syscall.LazyProc
	procSetIntValueEx        *syscall.LazyProc
	procGetEnumValue         *syscall.LazyProc
	procGetEnumEntrySymbolic *syscall.LazyProc
	procSetEnumValue         *syscall.LazyProc
	procSetEnumValueByString *syscall.LazyProc
	procGetBoolValue         *syscall.LazyProc
	procSetBoolValue         *syscall.LazyProc
	procGetStringValue       *syscall.LazyProc
	procSetStringValue       *syscall.LazyProc
	procSetCommandValue      *syscall.LazyProc
	procFeatureLoad          *syscall.LazyProc
	procFeatureSave          *syscall.LazyProc
	procFileAccessRead       *syscall.LazyProc
	procFileAccessReadEx     *syscall.LazyProc
	procFileAccessWrite      *syscall.LazyProc
	procFileAccessWriteEx    *syscall.LazyProc
	procFileAccessProgress   *syscall.LazyProc
	procRegisterAllEventCB   *syscall.LazyProc
	procRegisterEventCBEx    *syscall.LazyProc
	procEventNotificationOn  *syscall.LazyProc
	procEventNotificationOff *syscall.LazyProc
	procCamlSerialPortList   *syscall.LazyProc
	procCamlSetEnumPorts     *syscall.LazyProc
	procCamlSetBaudrate      *syscall.LazyProc
	procCamlGetBaudrate      *syscall.LazyProc
	procCamlSupportBaudrates *syscall.LazyProc
	procCamlSetGenCPTimeout  *syscall.LazyProc
	procEnumInterfacesGenTL  *syscall.LazyProc
	procUnloadGenTLLibrary   *syscall.LazyProc
	procEnumDevicesGenTL     *syscall.LazyProc
	procSaveImageToFileEx    *syscall.LazyProc
	procRotateImage          *syscall.LazyProc
	procFlipImage            *syscall.LazyProc
	procConvertPixelTypeEx   *syscall.LazyProc
	procSetBayerCvtQuality   *syscall.LazyProc
	procSetBayerFilterEnable *syscall.LazyProc
	procSetBayerGammaValue   *syscall.LazyProc
	procSetGammaValue        *syscall.LazyProc
	procSetBayerGammaParam   *syscall.LazyProc
	procSetBayerCCMParam     *syscall.LazyProc
	procSetBayerCCMParamEx   *syscall.LazyProc
	procImageContrast        *syscall.LazyProc
	procPurpleFringing       *syscall.LazyProc
	procSetISPConfig         *syscall.LazyProc
	procISPProcess           *syscall.LazyProc
	procHBDecode             *syscall.LazyProc
	procStartRecord          *syscall.LazyProc
	procInputOneFrame        *syscall.LazyProc
	procStopRecord           *syscall.LazyProc
	procReconstructImage     *syscall.LazyProc
	procSerialPortOpen       *syscall.LazyProc
	procSerialPortWrite      *syscall.LazyProc
	procSerialPortRead       *syscall.LazyProc
	procSerialPortClear      *syscall.LazyProc
	procSerialPortClose      *syscall.LazyProc
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

type mvCCStringValue struct {
	Current   [256]byte
	MaxLength int64
	Reserved  [2]uint32
}

func newDriver(dllPath string) (*driver, error) {
	if dllPath == "" {
		dllPath = resolveDLLPath()
	}

	dll := syscall.NewLazyDLL(dllPath)
	d := &driver{
		dllPath:                  dllPath,
		dll:                      dll,
		procInitialize:           dll.NewProc("MV_CC_Initialize"),
		procFinalize:             dll.NewProc("MV_CC_Finalize"),
		procGetSDKVersion:        dll.NewProc("MV_CC_GetSDKVersion"),
		procEnumDevices:          dll.NewProc("MV_CC_EnumDevices"),
		procEnumDevicesByIF:      dll.NewProc("MV_CC_EnumDevicesByInterface"),
		procCreateHandle:         dll.NewProc("MV_CC_CreateHandle"),
		procCreateHandleByGenTL:  dll.NewProc("MV_CC_CreateHandleByGenTL"),
		procDestroyHandle:        dll.NewProc("MV_CC_DestroyHandle"),
		procIsDeviceAccessible:   dll.NewProc("MV_CC_IsDeviceAccessible"),
		procOpenDevice:           dll.NewProc("MV_CC_OpenDevice"),
		procCloseDevice:          dll.NewProc("MV_CC_CloseDevice"),
		procIsDeviceConnected:    dll.NewProc("MV_CC_IsDeviceConnected"),
		procEnumInterfaces:       dll.NewProc("MV_CC_EnumInterfaces"),
		procCreateInterface:      dll.NewProc("MV_CC_CreateInterface"),
		procCreateInterfaceByID:  dll.NewProc("MV_CC_CreateInterfaceByID"),
		procOpenInterface:        dll.NewProc("MV_CC_OpenInterface"),
		procCloseInterface:       dll.NewProc("MV_CC_CloseInterface"),
		procDestroyInterface:     dll.NewProc("MV_CC_DestroyInterface"),
		procRegisterCallbackEx:   dll.NewProc("MV_CC_RegisterImageCallBackEx"),
		procStartGrabbing:        dll.NewProc("MV_CC_StartGrabbing"),
		procStopGrabbing:         dll.NewProc("MV_CC_StopGrabbing"),
		procGetImageBuffer:       dll.NewProc("MV_CC_GetImageBuffer"),
		procFreeImageBuffer:      dll.NewProc("MV_CC_FreeImageBuffer"),
		procGetOneFrameTimeout:   dll.NewProc("MV_CC_GetOneFrameTimeout"),
		procClearImageBuffer:     dll.NewProc("MV_CC_ClearImageBuffer"),
		procSetImageNodeNum:      dll.NewProc("MV_CC_SetImageNodeNum"),
		procSetGrabStrategy:      dll.NewProc("MV_CC_SetGrabStrategy"),
		procSetOutputQueueSize:   dll.NewProc("MV_CC_SetOutputQueueSize"),
		procGetOptimalPacketSize: dll.NewProc("MV_CC_GetOptimalPacketSize"),
		procGetFloatValue:        dll.NewProc("MV_CC_GetFloatValue"),
		procSetFloatValue:        dll.NewProc("MV_CC_SetFloatValue"),
		procGetIntValueEx:        dll.NewProc("MV_CC_GetIntValueEx"),
		procSetIntValueEx:        dll.NewProc("MV_CC_SetIntValueEx"),
		procGetEnumValue:         dll.NewProc("MV_CC_GetEnumValue"),
		procGetEnumEntrySymbolic: dll.NewProc("MV_CC_GetEnumEntrySymbolic"),
		procSetEnumValue:         dll.NewProc("MV_CC_SetEnumValue"),
		procSetEnumValueByString: dll.NewProc("MV_CC_SetEnumValueByString"),
		procGetBoolValue:         dll.NewProc("MV_CC_GetBoolValue"),
		procSetBoolValue:         dll.NewProc("MV_CC_SetBoolValue"),
		procGetStringValue:       dll.NewProc("MV_CC_GetStringValue"),
		procSetStringValue:       dll.NewProc("MV_CC_SetStringValue"),
		procSetCommandValue:      dll.NewProc("MV_CC_SetCommandValue"),
		procFeatureLoad:          dll.NewProc("MV_CC_FeatureLoad"),
		procFeatureSave:          dll.NewProc("MV_CC_FeatureSave"),
		procFileAccessRead:       dll.NewProc("MV_CC_FileAccessRead"),
		procFileAccessReadEx:     dll.NewProc("MV_CC_FileAccessReadEx"),
		procFileAccessWrite:      dll.NewProc("MV_CC_FileAccessWrite"),
		procFileAccessWriteEx:    dll.NewProc("MV_CC_FileAccessWriteEx"),
		procFileAccessProgress:   dll.NewProc("MV_CC_GetFileAccessProgress"),
		procRegisterAllEventCB:   dll.NewProc("MV_CC_RegisterAllEventCallBack"),
		procRegisterEventCBEx:    dll.NewProc("MV_CC_RegisterEventCallBackEx"),
		procEventNotificationOn:  dll.NewProc("MV_CC_EventNotificationOn"),
		procEventNotificationOff: dll.NewProc("MV_CC_EventNotificationOff"),
		procCamlSerialPortList:   dll.NewProc("MV_CAML_GetSerialPortList"),
		procCamlSetEnumPorts:     dll.NewProc("MV_CAML_SetEnumSerialPorts"),
		procCamlSetBaudrate:      dll.NewProc("MV_CAML_SetDeviceBaudrate"),
		procCamlGetBaudrate:      dll.NewProc("MV_CAML_GetDeviceBaudrate"),
		procCamlSupportBaudrates: dll.NewProc("MV_CAML_GetSupportBaudrates"),
		procCamlSetGenCPTimeout:  dll.NewProc("MV_CAML_SetGenCPTimeOut"),
		procEnumInterfacesGenTL:  dll.NewProc("MV_CC_EnumInterfacesByGenTL"),
		procUnloadGenTLLibrary:   dll.NewProc("MV_CC_UnloadGenTLLibrary"),
		procEnumDevicesGenTL:     dll.NewProc("MV_CC_EnumDevicesByGenTL"),
		procSaveImageToFileEx:    dll.NewProc("MV_CC_SaveImageToFileEx"),
		procRotateImage:          dll.NewProc("MV_CC_RotateImage"),
		procFlipImage:            dll.NewProc("MV_CC_FlipImage"),
		procConvertPixelTypeEx:   dll.NewProc("MV_CC_ConvertPixelTypeEx"),
		procSetBayerCvtQuality:   dll.NewProc("MV_CC_SetBayerCvtQuality"),
		procSetBayerFilterEnable: dll.NewProc("MV_CC_SetBayerFilterEnable"),
		procSetBayerGammaValue:   dll.NewProc("MV_CC_SetBayerGammaValue"),
		procSetGammaValue:        dll.NewProc("MV_CC_SetGammaValue"),
		procSetBayerGammaParam:   dll.NewProc("MV_CC_SetBayerGammaParam"),
		procSetBayerCCMParam:     dll.NewProc("MV_CC_SetBayerCCMParam"),
		procSetBayerCCMParamEx:   dll.NewProc("MV_CC_SetBayerCCMParamEx"),
		procImageContrast:        dll.NewProc("MV_CC_ImageContrast"),
		procPurpleFringing:       dll.NewProc("MV_CC_PurpleFringing"),
		procSetISPConfig:         dll.NewProc("MV_CC_SetISPConfig"),
		procISPProcess:           dll.NewProc("MV_CC_ISPProcess"),
		procHBDecode:             dll.NewProc("MV_CC_HB_Decode"),
		procStartRecord:          dll.NewProc("MV_CC_StartRecord"),
		procInputOneFrame:        dll.NewProc("MV_CC_InputOneFrame"),
		procStopRecord:           dll.NewProc("MV_CC_StopRecord"),
		procReconstructImage:     dll.NewProc("MV_CC_ReconstructImage"),
		procSerialPortOpen:       dll.NewProc("MV_CC_SerialPort_Open"),
		procSerialPortWrite:      dll.NewProc("MV_CC_SerialPort_Write"),
		procSerialPortRead:       dll.NewProc("MV_CC_SerialPort_Read"),
		procSerialPortClear:      dll.NewProc("MV_CC_SerialPort_ClearBuffer"),
		procSerialPortClose:      dll.NewProc("MV_CC_SerialPort_Close"),
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

func (d *driver) enumDevicesByInterface(handle uintptr, list *mvCCDeviceInfoList) error {
	ret, _, _ := d.procEnumDevicesByIF.Call(
		handle,
		uintptr(unsafe.Pointer(list)),
	)
	return newSDKError("MV_CC_EnumDevicesByInterface", ret)
}

func (d *driver) enumInterfaces(layerType uint32, list *mvInterfaceInfoList) error {
	ret, _, _ := d.procEnumInterfaces.Call(
		uintptr(layerType),
		uintptr(unsafe.Pointer(list)),
	)
	return newSDKError("MV_CC_EnumInterfaces", ret)
}

func (d *driver) createInterface(info *mvInterfaceInfo) (uintptr, error) {
	var handle uintptr
	ret, _, _ := d.procCreateInterface.Call(
		uintptr(unsafe.Pointer(&handle)),
		uintptr(unsafe.Pointer(info)),
	)
	if err := newSDKError("MV_CC_CreateInterface", ret); err != nil {
		return 0, err
	}
	return handle, nil
}

func (d *driver) createInterfaceByID(interfaceID string) (uintptr, error) {
	idPtr, err := syscall.BytePtrFromString(interfaceID)
	if err != nil {
		return 0, fmt.Errorf("build interface id %q: %w", interfaceID, err)
	}
	var handle uintptr
	ret, _, _ := d.procCreateInterfaceByID.Call(
		uintptr(unsafe.Pointer(&handle)),
		uintptr(unsafe.Pointer(idPtr)),
	)
	if err := newSDKError("MV_CC_CreateInterfaceByID", ret); err != nil {
		return 0, err
	}
	return handle, nil
}

func (d *driver) openInterface(handle uintptr) error {
	ret, _, _ := d.procOpenInterface.Call(handle, 0)
	return newSDKError("MV_CC_OpenInterface", ret)
}

func (d *driver) closeInterface(handle uintptr) error {
	ret, _, _ := d.procCloseInterface.Call(handle)
	return newSDKError("MV_CC_CloseInterface", ret)
}

func (d *driver) destroyInterface(handle uintptr) error {
	ret, _, _ := d.procDestroyInterface.Call(handle)
	return newSDKError("MV_CC_DestroyInterface", ret)
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

func (d *driver) createHandleByGenTL(info *mvGenTLDevInfo) (uintptr, error) {
	var handle uintptr
	ret, _, _ := d.procCreateHandleByGenTL.Call(
		uintptr(unsafe.Pointer(&handle)),
		uintptr(unsafe.Pointer(info)),
	)
	if err := newSDKError("MV_CC_CreateHandleByGenTL", ret); err != nil {
		return 0, err
	}
	return handle, nil
}

func (d *driver) destroyHandle(handle uintptr) error {
	ret, _, _ := d.procDestroyHandle.Call(handle)
	return newSDKError("MV_CC_DestroyHandle", ret)
}

func (d *driver) isDeviceAccessible(info *mvCCDeviceInfo, accessMode uint32) bool {
	ret, _, _ := d.procIsDeviceAccessible.Call(
		uintptr(unsafe.Pointer(info)),
		uintptr(accessMode),
	)
	return ret != 0
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

func (d *driver) registerImageCallbackEx(handle uintptr, callback uintptr, user uintptr) error {
	ret, _, _ := d.procRegisterCallbackEx.Call(handle, callback, user)
	return newSDKError("MV_CC_RegisterImageCallBackEx", ret)
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

func (d *driver) getOneFrameTimeout(handle uintptr, data []byte, frameInfo *mvFrameOutInfoEx, timeoutMs uint32) error {
	var dataPtr uintptr
	if len(data) > 0 {
		dataPtr = uintptr(unsafe.Pointer(&data[0]))
	}
	ret, _, _ := d.procGetOneFrameTimeout.Call(
		handle,
		dataPtr,
		uintptr(uint32(len(data))),
		uintptr(unsafe.Pointer(frameInfo)),
		uintptr(timeoutMs),
	)
	return newSDKError("MV_CC_GetOneFrameTimeout", ret)
}

func (d *driver) clearImageBuffer(handle uintptr) error {
	ret, _, _ := d.procClearImageBuffer.Call(handle)
	return newSDKError("MV_CC_ClearImageBuffer", ret)
}

func (d *driver) setImageNodeNum(handle uintptr, num uint32) error {
	ret, _, _ := d.procSetImageNodeNum.Call(handle, uintptr(num))
	return newSDKError("MV_CC_SetImageNodeNum", ret)
}

func (d *driver) setGrabStrategy(handle uintptr, strategy GrabStrategy) error {
	ret, _, _ := d.procSetGrabStrategy.Call(handle, uintptr(strategy))
	return newSDKError("MV_CC_SetGrabStrategy", ret)
}

func (d *driver) setOutputQueueSize(handle uintptr, size uint32) error {
	ret, _, _ := d.procSetOutputQueueSize.Call(handle, uintptr(size))
	return newSDKError("MV_CC_SetOutputQueueSize", ret)
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

func (d *driver) getEnumEntrySymbolic(handle uintptr, key string, entry *mvCCEnumEntry) error {
	keyPtr, err := syscall.BytePtrFromString(key)
	if err != nil {
		return fmt.Errorf("build key %q: %w", key, err)
	}
	ret, _, _ := d.procGetEnumEntrySymbolic.Call(
		handle,
		uintptr(unsafe.Pointer(keyPtr)),
		uintptr(unsafe.Pointer(entry)),
	)
	return newSDKError("MV_CC_GetEnumEntrySymbolic", ret)
}

func (d *driver) setEnumValueByString(handle uintptr, key string, value string) error {
	keyPtr, err := syscall.BytePtrFromString(key)
	if err != nil {
		return fmt.Errorf("build key %q: %w", key, err)
	}
	valuePtr, err := syscall.BytePtrFromString(value)
	if err != nil {
		return fmt.Errorf("build enum value for key %q: %w", key, err)
	}
	ret, _, _ := d.procSetEnumValueByString.Call(
		handle,
		uintptr(unsafe.Pointer(keyPtr)),
		uintptr(unsafe.Pointer(valuePtr)),
	)
	return newSDKError("MV_CC_SetEnumValueByString", ret)
}

func (d *driver) getBoolValue(handle uintptr, key string) (bool, error) {
	keyPtr, err := syscall.BytePtrFromString(key)
	if err != nil {
		return false, fmt.Errorf("build key %q: %w", key, err)
	}

	var value byte
	ret, _, _ := d.procGetBoolValue.Call(
		handle,
		uintptr(unsafe.Pointer(keyPtr)),
		uintptr(unsafe.Pointer(&value)),
	)
	if err := newSDKError("MV_CC_GetBoolValue", ret); err != nil {
		return false, err
	}
	return value != 0, nil
}

func (d *driver) setBoolValue(handle uintptr, key string, value bool) error {
	keyPtr, err := syscall.BytePtrFromString(key)
	if err != nil {
		return fmt.Errorf("build key %q: %w", key, err)
	}

	var raw uintptr
	if value {
		raw = 1
	}
	ret, _, _ := d.procSetBoolValue.Call(
		handle,
		uintptr(unsafe.Pointer(keyPtr)),
		raw,
	)
	return newSDKError("MV_CC_SetBoolValue", ret)
}

func (d *driver) getStringValue(handle uintptr, key string, value *mvCCStringValue) error {
	keyPtr, err := syscall.BytePtrFromString(key)
	if err != nil {
		return fmt.Errorf("build key %q: %w", key, err)
	}
	ret, _, _ := d.procGetStringValue.Call(
		handle,
		uintptr(unsafe.Pointer(keyPtr)),
		uintptr(unsafe.Pointer(value)),
	)
	return newSDKError("MV_CC_GetStringValue", ret)
}

func (d *driver) setStringValue(handle uintptr, key string, value string) error {
	keyPtr, err := syscall.BytePtrFromString(key)
	if err != nil {
		return fmt.Errorf("build key %q: %w", key, err)
	}
	valuePtr, err := syscall.BytePtrFromString(value)
	if err != nil {
		return fmt.Errorf("build value for key %q: %w", key, err)
	}
	ret, _, _ := d.procSetStringValue.Call(
		handle,
		uintptr(unsafe.Pointer(keyPtr)),
		uintptr(unsafe.Pointer(valuePtr)),
	)
	return newSDKError("MV_CC_SetStringValue", ret)
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

func (d *driver) featureLoad(handle uintptr, path string) error {
	pathPtr, err := syscall.BytePtrFromString(path)
	if err != nil {
		return fmt.Errorf("build feature file path %q: %w", path, err)
	}
	ret, _, _ := d.procFeatureLoad.Call(
		handle,
		uintptr(unsafe.Pointer(pathPtr)),
	)
	return newSDKError("MV_CC_FeatureLoad", ret)
}

func (d *driver) featureSave(handle uintptr, path string) error {
	pathPtr, err := syscall.BytePtrFromString(path)
	if err != nil {
		return fmt.Errorf("build feature file path %q: %w", path, err)
	}
	ret, _, _ := d.procFeatureSave.Call(
		handle,
		uintptr(unsafe.Pointer(pathPtr)),
	)
	return newSDKError("MV_CC_FeatureSave", ret)
}

func (d *driver) fileAccessRead(handle uintptr, param *mvCCFileAccess) error {
	ret, _, _ := d.procFileAccessRead.Call(
		handle,
		uintptr(unsafe.Pointer(param)),
	)
	return newSDKError("MV_CC_FileAccessRead", ret)
}

func (d *driver) fileAccessReadEx(handle uintptr, param *mvCCFileAccessEx) error {
	ret, _, _ := d.procFileAccessReadEx.Call(
		handle,
		uintptr(unsafe.Pointer(param)),
	)
	return newSDKError("MV_CC_FileAccessReadEx", ret)
}

func (d *driver) fileAccessWrite(handle uintptr, param *mvCCFileAccess) error {
	ret, _, _ := d.procFileAccessWrite.Call(
		handle,
		uintptr(unsafe.Pointer(param)),
	)
	return newSDKError("MV_CC_FileAccessWrite", ret)
}

func (d *driver) fileAccessWriteEx(handle uintptr, param *mvCCFileAccessEx) error {
	ret, _, _ := d.procFileAccessWriteEx.Call(
		handle,
		uintptr(unsafe.Pointer(param)),
	)
	return newSDKError("MV_CC_FileAccessWriteEx", ret)
}

func (d *driver) getFileAccessProgress(handle uintptr, progress *mvCCFileAccessProgress) error {
	ret, _, _ := d.procFileAccessProgress.Call(
		handle,
		uintptr(unsafe.Pointer(progress)),
	)
	return newSDKError("MV_CC_GetFileAccessProgress", ret)
}

func (d *driver) registerAllEventCallback(handle uintptr, callback uintptr, user uintptr) error {
	ret, _, _ := d.procRegisterAllEventCB.Call(handle, callback, user)
	return newSDKError("MV_CC_RegisterAllEventCallBack", ret)
}

func (d *driver) registerEventCallbackEx(handle uintptr, eventName string, callback uintptr, user uintptr) error {
	namePtr, err := syscall.BytePtrFromString(eventName)
	if err != nil {
		return fmt.Errorf("build event name %q: %w", eventName, err)
	}
	ret, _, _ := d.procRegisterEventCBEx.Call(
		handle,
		uintptr(unsafe.Pointer(namePtr)),
		callback,
		user,
	)
	return newSDKError("MV_CC_RegisterEventCallBackEx", ret)
}

func (d *driver) eventNotificationOn(handle uintptr, eventName string) error {
	namePtr, err := syscall.BytePtrFromString(eventName)
	if err != nil {
		return fmt.Errorf("build event name %q: %w", eventName, err)
	}
	ret, _, _ := d.procEventNotificationOn.Call(
		handle,
		uintptr(unsafe.Pointer(namePtr)),
	)
	return newSDKError("MV_CC_EventNotificationOn", ret)
}

func (d *driver) eventNotificationOff(handle uintptr, eventName string) error {
	namePtr, err := syscall.BytePtrFromString(eventName)
	if err != nil {
		return fmt.Errorf("build event name %q: %w", eventName, err)
	}
	ret, _, _ := d.procEventNotificationOff.Call(
		handle,
		uintptr(unsafe.Pointer(namePtr)),
	)
	return newSDKError("MV_CC_EventNotificationOff", ret)
}

func (d *driver) camlGetSerialPortList(list *mvCamlSerialPortList) error {
	ret, _, _ := d.procCamlSerialPortList.Call(uintptr(unsafe.Pointer(list)))
	return newSDKError("MV_CAML_GetSerialPortList", ret)
}

func (d *driver) camlSetEnumSerialPorts(list *mvCamlSerialPortList) error {
	ret, _, _ := d.procCamlSetEnumPorts.Call(uintptr(unsafe.Pointer(list)))
	return newSDKError("MV_CAML_SetEnumSerialPorts", ret)
}

func (d *driver) camlSetDeviceBaudrate(handle uintptr, baudrate uint32) error {
	ret, _, _ := d.procCamlSetBaudrate.Call(handle, uintptr(baudrate))
	return newSDKError("MV_CAML_SetDeviceBaudrate", ret)
}

func (d *driver) camlGetDeviceBaudrate(handle uintptr) (uint32, error) {
	var baudrate uint32
	ret, _, _ := d.procCamlGetBaudrate.Call(
		handle,
		uintptr(unsafe.Pointer(&baudrate)),
	)
	if err := newSDKError("MV_CAML_GetDeviceBaudrate", ret); err != nil {
		return 0, err
	}
	return baudrate, nil
}

func (d *driver) camlGetSupportBaudrates(handle uintptr) (uint32, error) {
	var baudrates uint32
	ret, _, _ := d.procCamlSupportBaudrates.Call(
		handle,
		uintptr(unsafe.Pointer(&baudrates)),
	)
	if err := newSDKError("MV_CAML_GetSupportBaudrates", ret); err != nil {
		return 0, err
	}
	return baudrates, nil
}

func (d *driver) camlSetGenCPTimeout(handle uintptr, timeoutMs uint32) error {
	ret, _, _ := d.procCamlSetGenCPTimeout.Call(handle, uintptr(timeoutMs))
	return newSDKError("MV_CAML_SetGenCPTimeOut", ret)
}

func (d *driver) enumInterfacesByGenTL(path string, list *mvGenTLIFInfoList) error {
	pathPtr, err := syscall.BytePtrFromString(path)
	if err != nil {
		return fmt.Errorf("build gentl path %q: %w", path, err)
	}
	ret, _, _ := d.procEnumInterfacesGenTL.Call(
		uintptr(unsafe.Pointer(list)),
		uintptr(unsafe.Pointer(pathPtr)),
	)
	return newSDKError("MV_CC_EnumInterfacesByGenTL", ret)
}

func (d *driver) unloadGenTLLibrary(path string) error {
	pathPtr, err := syscall.BytePtrFromString(path)
	if err != nil {
		return fmt.Errorf("build gentl path %q: %w", path, err)
	}
	ret, _, _ := d.procUnloadGenTLLibrary.Call(uintptr(unsafe.Pointer(pathPtr)))
	return newSDKError("MV_CC_UnloadGenTLLibrary", ret)
}

func (d *driver) enumDevicesByGenTL(info *mvGenTLIFInfo, list *mvGenTLDevInfoList) error {
	ret, _, _ := d.procEnumDevicesGenTL.Call(
		uintptr(unsafe.Pointer(info)),
		uintptr(unsafe.Pointer(list)),
	)
	return newSDKError("MV_CC_EnumDevicesByGenTL", ret)
}

func (d *driver) saveImageToFileEx(handle uintptr, param *mvSaveImageToFileParamEx) error {
	ret, _, _ := d.procSaveImageToFileEx.Call(
		handle,
		uintptr(unsafe.Pointer(param)),
	)
	return newSDKError("MV_CC_SaveImageToFileEx", ret)
}

func (d *driver) rotateImage(handle uintptr, param *mvCCRotateImageParam) error {
	ret, _, _ := d.procRotateImage.Call(
		handle,
		uintptr(unsafe.Pointer(param)),
	)
	return newSDKError("MV_CC_RotateImage", ret)
}

func (d *driver) flipImage(handle uintptr, param *mvCCFlipImageParam) error {
	ret, _, _ := d.procFlipImage.Call(
		handle,
		uintptr(unsafe.Pointer(param)),
	)
	return newSDKError("MV_CC_FlipImage", ret)
}

func (d *driver) convertPixelTypeEx(handle uintptr, param *mvCCPixelConvertParamEx) error {
	ret, _, _ := d.procConvertPixelTypeEx.Call(
		handle,
		uintptr(unsafe.Pointer(param)),
	)
	return newSDKError("MV_CC_ConvertPixelTypeEx", ret)
}

func (d *driver) setBayerCvtQuality(handle uintptr, quality InterpolationMethod) error {
	ret, _, _ := d.procSetBayerCvtQuality.Call(handle, uintptr(uint32(quality)))
	return newSDKError("MV_CC_SetBayerCvtQuality", ret)
}

func (d *driver) setBayerFilterEnable(handle uintptr, enabled bool) error {
	var raw uintptr
	if enabled {
		raw = 1
	}
	ret, _, _ := d.procSetBayerFilterEnable.Call(handle, raw)
	return newSDKError("MV_CC_SetBayerFilterEnable", ret)
}

func (d *driver) setBayerGammaValue(handle uintptr, gamma float32) error {
	ret, _, _ := d.procSetBayerGammaValue.Call(handle, uintptr(math.Float32bits(gamma)))
	return newSDKError("MV_CC_SetBayerGammaValue", ret)
}

func (d *driver) setGammaValue(handle uintptr, pixelType uint32, gamma float32) error {
	ret, _, _ := d.procSetGammaValue.Call(
		handle,
		uintptr(pixelType),
		uintptr(math.Float32bits(gamma)),
	)
	return newSDKError("MV_CC_SetGammaValue", ret)
}

func (d *driver) setBayerGammaParam(handle uintptr, param *mvCCGammaParam) error {
	ret, _, _ := d.procSetBayerGammaParam.Call(
		handle,
		uintptr(unsafe.Pointer(param)),
	)
	return newSDKError("MV_CC_SetBayerGammaParam", ret)
}

func (d *driver) setBayerCCMParam(handle uintptr, param *mvCCCCMParam) error {
	ret, _, _ := d.procSetBayerCCMParam.Call(
		handle,
		uintptr(unsafe.Pointer(param)),
	)
	return newSDKError("MV_CC_SetBayerCCMParam", ret)
}

func (d *driver) setBayerCCMParamEx(handle uintptr, param *mvCCCCMParamEx) error {
	ret, _, _ := d.procSetBayerCCMParamEx.Call(
		handle,
		uintptr(unsafe.Pointer(param)),
	)
	return newSDKError("MV_CC_SetBayerCCMParamEx", ret)
}

func (d *driver) imageContrast(handle uintptr, param *mvCCContrastParam) error {
	ret, _, _ := d.procImageContrast.Call(
		handle,
		uintptr(unsafe.Pointer(param)),
	)
	return newSDKError("MV_CC_ImageContrast", ret)
}

func (d *driver) purpleFringing(handle uintptr, param *mvCCPurpleFringingParam) error {
	ret, _, _ := d.procPurpleFringing.Call(
		handle,
		uintptr(unsafe.Pointer(param)),
	)
	return newSDKError("MV_CC_PurpleFringing", ret)
}

func (d *driver) setISPConfig(handle uintptr, param *mvCCISPConfigParam) error {
	ret, _, _ := d.procSetISPConfig.Call(
		handle,
		uintptr(unsafe.Pointer(param)),
	)
	return newSDKError("MV_CC_SetISPConfig", ret)
}

func (d *driver) ispProcess(handle uintptr, input *mvCCImage, output *mvCCImage) error {
	ret, _, _ := d.procISPProcess.Call(
		handle,
		uintptr(unsafe.Pointer(input)),
		uintptr(unsafe.Pointer(output)),
	)
	return newSDKError("MV_CC_ISPProcess", ret)
}

func (d *driver) hbDecode(handle uintptr, param *mvCCHBDecodeParam) error {
	ret, _, _ := d.procHBDecode.Call(
		handle,
		uintptr(unsafe.Pointer(param)),
	)
	return newSDKError("MV_CC_HB_Decode", ret)
}

func (d *driver) startRecord(handle uintptr, param *mvCCRecordParam) error {
	ret, _, _ := d.procStartRecord.Call(
		handle,
		uintptr(unsafe.Pointer(param)),
	)
	return newSDKError("MV_CC_StartRecord", ret)
}

func (d *driver) inputOneFrame(handle uintptr, param *mvCCInputFrameInfo) error {
	ret, _, _ := d.procInputOneFrame.Call(
		handle,
		uintptr(unsafe.Pointer(param)),
	)
	return newSDKError("MV_CC_InputOneFrame", ret)
}

func (d *driver) stopRecord(handle uintptr) error {
	ret, _, _ := d.procStopRecord.Call(handle)
	return newSDKError("MV_CC_StopRecord", ret)
}

func (d *driver) reconstructImage(handle uintptr, param *mvReconstructImageParam) error {
	ret, _, _ := d.procReconstructImage.Call(
		handle,
		uintptr(unsafe.Pointer(param)),
	)
	return newSDKError("MV_CC_ReconstructImage", ret)
}

func (d *driver) serialPortOpen(handle uintptr) error {
	ret, _, _ := d.procSerialPortOpen.Call(handle)
	return newSDKError("MV_CC_SerialPort_Open", ret)
}

func (d *driver) serialPortWrite(handle uintptr, data []byte) (uint32, error) {
	var dataPtr uintptr
	if len(data) > 0 {
		dataPtr = uintptr(unsafe.Pointer(&data[0]))
	}
	var written uint32
	ret, _, _ := d.procSerialPortWrite.Call(
		handle,
		dataPtr,
		uintptr(uint32(len(data))),
		uintptr(unsafe.Pointer(&written)),
	)
	if err := newSDKError("MV_CC_SerialPort_Write", ret); err != nil {
		return 0, err
	}
	return written, nil
}

func (d *driver) serialPortRead(handle uintptr, buffer []byte, timeoutMs uint32) (uint32, error) {
	var bufferPtr uintptr
	if len(buffer) > 0 {
		bufferPtr = uintptr(unsafe.Pointer(&buffer[0]))
	}
	var read uint32
	ret, _, _ := d.procSerialPortRead.Call(
		handle,
		bufferPtr,
		uintptr(uint32(len(buffer))),
		uintptr(unsafe.Pointer(&read)),
		uintptr(timeoutMs),
	)
	if err := newSDKError("MV_CC_SerialPort_Read", ret); err != nil {
		return 0, err
	}
	return read, nil
}

func (d *driver) serialPortClearBuffer(handle uintptr) error {
	ret, _, _ := d.procSerialPortClear.Call(handle)
	return newSDKError("MV_CC_SerialPort_ClearBuffer", ret)
}

func (d *driver) serialPortClose(handle uintptr) error {
	ret, _, _ := d.procSerialPortClose.Call(handle)
	return newSDKError("MV_CC_SerialPort_Close", ret)
}
