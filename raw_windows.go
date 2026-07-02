//go:build windows && amd64

package mvs

import (
	"fmt"
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
	procCreateHandle         *syscall.LazyProc
	procDestroyHandle        *syscall.LazyProc
	procIsDeviceAccessible   *syscall.LazyProc
	procOpenDevice           *syscall.LazyProc
	procCloseDevice          *syscall.LazyProc
	procIsDeviceConnected    *syscall.LazyProc
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
	procSaveImageToFileEx    *syscall.LazyProc
	procConvertPixelTypeEx   *syscall.LazyProc
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
		procCreateHandle:         dll.NewProc("MV_CC_CreateHandle"),
		procDestroyHandle:        dll.NewProc("MV_CC_DestroyHandle"),
		procIsDeviceAccessible:   dll.NewProc("MV_CC_IsDeviceAccessible"),
		procOpenDevice:           dll.NewProc("MV_CC_OpenDevice"),
		procCloseDevice:          dll.NewProc("MV_CC_CloseDevice"),
		procIsDeviceConnected:    dll.NewProc("MV_CC_IsDeviceConnected"),
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
		procSaveImageToFileEx:    dll.NewProc("MV_CC_SaveImageToFileEx"),
		procConvertPixelTypeEx:   dll.NewProc("MV_CC_ConvertPixelTypeEx"),
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

func (d *driver) saveImageToFileEx(handle uintptr, param *mvSaveImageToFileParamEx) error {
	ret, _, _ := d.procSaveImageToFileEx.Call(
		handle,
		uintptr(unsafe.Pointer(param)),
	)
	return newSDKError("MV_CC_SaveImageToFileEx", ret)
}

func (d *driver) convertPixelTypeEx(handle uintptr, param *mvCCPixelConvertParamEx) error {
	ret, _, _ := d.procConvertPixelTypeEx.Call(
		handle,
		uintptr(unsafe.Pointer(param)),
	)
	return newSDKError("MV_CC_ConvertPixelTypeEx", ret)
}
