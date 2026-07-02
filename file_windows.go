//go:build windows && amd64

package mvs

import (
	"fmt"
	"runtime"
	"syscall"
)

func (c *Camera) ReadDeviceFileToFile(deviceFileName string, localPath string) error {
	param, keepAlive, err := newFileAccessParam(localPath, deviceFileName)
	if err != nil {
		return err
	}

	c.mu.Lock()
	if !c.open || c.handle == 0 {
		c.mu.Unlock()
		return ErrCameraClosed
	}
	handle := c.handle
	c.mu.Unlock()

	err = c.sdk.driver.fileAccessRead(handle, &param)
	for _, ptr := range keepAlive {
		runtime.KeepAlive(ptr)
	}
	return err
}

func (c *Camera) WriteDeviceFileFromFile(localPath string, deviceFileName string) error {
	param, keepAlive, err := newFileAccessParam(localPath, deviceFileName)
	if err != nil {
		return err
	}

	c.mu.Lock()
	if !c.open || c.handle == 0 {
		c.mu.Unlock()
		return ErrCameraClosed
	}
	handle := c.handle
	c.mu.Unlock()

	err = c.sdk.driver.fileAccessWrite(handle, &param)
	for _, ptr := range keepAlive {
		runtime.KeepAlive(ptr)
	}
	return err
}

func (c *Camera) ReadDeviceFile(deviceFileName string, bufferSize int) ([]byte, FileAccessProgress, error) {
	if deviceFileName == "" {
		return nil, FileAccessProgress{}, fmt.Errorf("%w: device file name is empty", ErrInvalidArgument)
	}
	if bufferSize <= 0 || uint64(bufferSize) > maxUint32Value {
		return nil, FileAccessProgress{}, fmt.Errorf("%w: invalid file buffer size %d", ErrInvalidArgument, bufferSize)
	}

	buffer := make([]byte, bufferSize)
	namePtr, err := syscall.BytePtrFromString(deviceFileName)
	if err != nil {
		return nil, FileAccessProgress{}, fmt.Errorf("build device file name %q: %w", deviceFileName, err)
	}
	param := mvCCFileAccessEx{
		UserFileBuffer: &buffer[0],
		FileBufferSize: uint32(len(buffer)),
		DevFileName:    namePtr,
	}

	c.mu.Lock()
	if !c.open || c.handle == 0 {
		c.mu.Unlock()
		return nil, FileAccessProgress{}, ErrCameraClosed
	}
	handle := c.handle
	c.mu.Unlock()

	if err := c.sdk.driver.fileAccessReadEx(handle, &param); err != nil {
		runtime.KeepAlive(buffer)
		runtime.KeepAlive(namePtr)
		return nil, FileAccessProgress{}, err
	}
	runtime.KeepAlive(buffer)
	runtime.KeepAlive(namePtr)

	if int(param.FileBufferLen) > len(buffer) {
		return nil, FileAccessProgress{}, fmt.Errorf("%w: SDK returned %d bytes into %d-byte buffer", ErrInvalidFrameData, param.FileBufferLen, len(buffer))
	}
	progress, err := c.GetFileAccessProgress()
	if err != nil {
		return buffer[:param.FileBufferLen], FileAccessProgress{}, err
	}
	return buffer[:param.FileBufferLen], progress, nil
}

func (c *Camera) WriteDeviceFile(deviceFileName string, data []byte) (FileAccessProgress, error) {
	if deviceFileName == "" {
		return FileAccessProgress{}, fmt.Errorf("%w: device file name is empty", ErrInvalidArgument)
	}
	if len(data) == 0 || uint64(len(data)) > maxUint32Value {
		return FileAccessProgress{}, fmt.Errorf("%w: invalid file data length %d", ErrInvalidArgument, len(data))
	}

	namePtr, err := syscall.BytePtrFromString(deviceFileName)
	if err != nil {
		return FileAccessProgress{}, fmt.Errorf("build device file name %q: %w", deviceFileName, err)
	}
	param := mvCCFileAccessEx{
		UserFileBuffer: &data[0],
		FileBufferSize: uint32(len(data)),
		FileBufferLen:  uint32(len(data)),
		DevFileName:    namePtr,
	}

	c.mu.Lock()
	if !c.open || c.handle == 0 {
		c.mu.Unlock()
		return FileAccessProgress{}, ErrCameraClosed
	}
	handle := c.handle
	c.mu.Unlock()

	if err := c.sdk.driver.fileAccessWriteEx(handle, &param); err != nil {
		runtime.KeepAlive(data)
		runtime.KeepAlive(namePtr)
		return FileAccessProgress{}, err
	}
	runtime.KeepAlive(data)
	runtime.KeepAlive(namePtr)
	return c.GetFileAccessProgress()
}

func (c *Camera) GetFileAccessProgress() (FileAccessProgress, error) {
	c.mu.Lock()
	if !c.open || c.handle == 0 {
		c.mu.Unlock()
		return FileAccessProgress{}, ErrCameraClosed
	}
	handle := c.handle
	c.mu.Unlock()

	var raw mvCCFileAccessProgress
	if err := c.sdk.driver.getFileAccessProgress(handle, &raw); err != nil {
		return FileAccessProgress{}, err
	}
	return FileAccessProgress{
		Completed: raw.Completed,
		Total:     raw.Total,
	}, nil
}

func newFileAccessParam(localPath string, deviceFileName string) (mvCCFileAccess, []*byte, error) {
	if localPath == "" {
		return mvCCFileAccess{}, nil, fmt.Errorf("%w: local path is empty", ErrInvalidArgument)
	}
	if deviceFileName == "" {
		return mvCCFileAccess{}, nil, fmt.Errorf("%w: device file name is empty", ErrInvalidArgument)
	}
	localPtr, err := syscall.BytePtrFromString(localPath)
	if err != nil {
		return mvCCFileAccess{}, nil, fmt.Errorf("build local path %q: %w", localPath, err)
	}
	devicePtr, err := syscall.BytePtrFromString(deviceFileName)
	if err != nil {
		return mvCCFileAccess{}, nil, fmt.Errorf("build device file name %q: %w", deviceFileName, err)
	}
	return mvCCFileAccess{
		UserFileName: localPtr,
		DevFileName:  devicePtr,
	}, []*byte{localPtr, devicePtr}, nil
}
