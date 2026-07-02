package mvs

import (
	"errors"
	"fmt"
)

var (
	ErrSDKNotInitialized       = errors.New("mvs: SDK is not initialized")
	ErrCameraClosed            = errors.New("mvs: camera is closed")
	ErrCameraNotGrabbing       = errors.New("mvs: camera is not grabbing")
	ErrNilFrameBuffer          = errors.New("mvs: frame buffer pointer is nil")
	ErrDeviceNotFound          = errors.New("mvs: device not found")
	ErrInvalidFrameData        = errors.New("mvs: invalid frame data")
	ErrUnsupportedPixel        = errors.New("mvs: unsupported pixel type")
	ErrInvalidArgument         = errors.New("mvs: invalid argument")
	ErrCallbackRegistered      = errors.New("mvs: frame callback is already registered")
	ErrAcquisitionModeConflict = errors.New("mvs: acquisition mode conflict")
)

type SDKError struct {
	Op   string
	Code uint32
}

func (e SDKError) Error() string {
	if name, ok := sdkErrorNames[e.Code]; ok {
		return fmt.Sprintf("%s failed with 0x%08X (%s)", e.Op, e.Code, name)
	}
	return fmt.Sprintf("%s failed with 0x%08X", e.Op, e.Code)
}

func IsSDKErrorCode(err error, code uint32) bool {
	var sdkErr SDKError
	if errors.As(err, &sdkErr) {
		return sdkErr.Code == code
	}
	return false
}

func newSDKError(op string, ret uintptr) error {
	code := uint32(ret)
	if code == MVOK {
		return nil
	}
	return SDKError{Op: op, Code: code}
}

var sdkErrorNames = map[uint32]string{
	MVOK:             "MV_OK",
	ErrHandle:        "MV_E_HANDLE",
	ErrSupport:       "MV_E_SUPPORT",
	ErrBufferOver:    "MV_E_BUFOVER",
	ErrCallOrder:     "MV_E_CALLORDER",
	ErrParameter:     "MV_E_PARAMETER",
	ErrResource:      "MV_E_RESOURCE",
	ErrNoData:        "MV_E_NODATA",
	ErrPrecondition:  "MV_E_PRECONDITION",
	ErrVersion:       "MV_E_VERSION",
	ErrNoEnoughBuf:   "MV_E_NOENOUGH_BUF",
	ErrAbnormalImg:   "MV_E_ABNORMAL_IMAGE",
	ErrLoadLibrary:   "MV_E_LOAD_LIBRARY",
	ErrNoOutBuf:      "MV_E_NOOUTBUF",
	ErrGCGeneric:     "MV_E_GC_GENERIC",
	ErrGCArgument:    "MV_E_GC_ARGUMENT",
	ErrGCRange:       "MV_E_GC_RANGE",
	ErrGCProperty:    "MV_E_GC_PROPERTY",
	ErrGCRuntime:     "MV_E_GC_RUNTIME",
	ErrGCLogical:     "MV_E_GC_LOGICAL",
	ErrGCAccess:      "MV_E_GC_ACCESS",
	ErrGCTimeout:     "MV_E_GC_TIMEOUT",
	ErrGCDynamicCast: "MV_E_GC_DYNAMICCAST",
	ErrGCUnknown:     "MV_E_GC_UNKNOW",
	ErrAccessDenied:  "MV_E_ACCESS_DENIED",
	ErrBusy:          "MV_E_BUSY",
}
