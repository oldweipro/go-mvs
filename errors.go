//go:build windows && amd64

package mvsdk

import (
	"errors"
	"fmt"
)

var (
	ErrSDKNotInitialized  = errors.New("mvsdk: SDK is not initialized")
	ErrCameraClosed       = errors.New("mvsdk: camera is closed")
	ErrCameraNotGrabbing  = errors.New("mvsdk: camera is not grabbing")
	ErrNilFrameBuffer     = errors.New("mvsdk: frame buffer pointer is nil")
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
	MVOK:            "MV_OK",
	ErrHandle:       "MV_E_HANDLE",
	ErrSupport:      "MV_E_SUPPORT",
	ErrBufferOver:   "MV_E_BUFOVER",
	ErrCallOrder:    "MV_E_CALLORDER",
	ErrParameter:    "MV_E_PARAMETER",
	ErrResource:     "MV_E_RESOURCE",
	ErrNoData:       "MV_E_NODATA",
	ErrPrecondition: "MV_E_PRECONDITION",
	ErrVersion:      "MV_E_VERSION",
	ErrNoEnoughBuf:  "MV_E_NOENOUGH_BUF",
	ErrAbnormalImg:  "MV_E_ABNORMAL_IMAGE",
	ErrLoadLibrary:  "MV_E_LOAD_LIBRARY",
	ErrNoOutBuf:     "MV_E_NOOUTBUF",
	ErrBusy:         "MV_E_BUSY",
}
