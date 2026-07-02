//go:build windows && amd64

package mvs

import (
	"fmt"
	"runtime"
	"syscall"
)

const (
	minRecordSide     = 96
	maxRecordSide     = 8000
	minRecordFrameFPS = 1.0 / 16.0
	maxRecordFrameFPS = 1000.0
	minRecordBitRate  = 128
	maxRecordBitRate  = 16 * 1024
)

func (c *Camera) StartRecord(options RecordOptions) error {
	normalized, err := c.normalizeRecordOptions(options)
	if err != nil {
		return err
	}

	pathPtr, err := syscall.BytePtrFromString(normalized.Path)
	if err != nil {
		return fmt.Errorf("build record path %q: %w", normalized.Path, err)
	}

	param := mvCCRecordParam{
		PixelType:    normalized.PixelType,
		Width:        uint16(normalized.Width),
		Height:       uint16(normalized.Height),
		FrameRate:    normalized.FrameRate,
		BitRate:      normalized.BitRateKbps,
		RecordFormat: uint32(normalized.Format),
		FilePath:     pathPtr,
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.open || c.handle == 0 {
		return ErrCameraClosed
	}
	if c.recording {
		return nil
	}
	if err := c.sdk.driver.startRecord(c.handle, &param); err != nil {
		runtime.KeepAlive(pathPtr)
		return err
	}
	runtime.KeepAlive(pathPtr)
	c.recording = true
	return nil
}

func (c *Camera) InputRecordFrame(frame *Frame) error {
	if err := validateFrame(frame); err != nil {
		return err
	}
	return c.InputRecordData(frame.Data)
}

func (c *Camera) InputRecordData(data []byte) error {
	if len(data) == 0 {
		return fmt.Errorf("%w: record frame data is empty", ErrInvalidArgument)
	}
	if uint64(len(data)) > maxUint32Value {
		return fmt.Errorf("%w: record frame data is larger than UINT_MAX", ErrInvalidArgument)
	}

	c.mu.Lock()
	if !c.open || c.handle == 0 {
		c.mu.Unlock()
		return ErrCameraClosed
	}
	if !c.recording {
		c.mu.Unlock()
		return fmt.Errorf("%w: recording is not started", ErrAcquisitionModeConflict)
	}
	handle := c.handle
	c.mu.Unlock()

	param := mvCCInputFrameInfo{
		Data:    &data[0],
		DataLen: uint32(len(data)),
	}
	err := c.sdk.driver.inputOneFrame(handle, &param)
	runtime.KeepAlive(data)
	return err
}

func (c *Camera) StopRecord() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.open || c.handle == 0 {
		return nil
	}
	if !c.recording {
		return nil
	}
	if err := c.sdk.driver.stopRecord(c.handle); err != nil {
		return err
	}
	c.recording = false
	return nil
}

func (c *Camera) normalizeRecordOptions(options RecordOptions) (RecordOptions, error) {
	if err := validateSDKPath(options.Path); err != nil {
		return RecordOptions{}, err
	}

	if options.PixelType == 0 || options.PixelType == PixelTypeUndefined {
		pixel, err := c.GetEnum(NodePixelFormat)
		if err != nil {
			return RecordOptions{}, err
		}
		options.PixelType = pixel.Current
	}
	if options.Width == 0 {
		width, err := c.GetInt(NodeWidth)
		if err != nil {
			return RecordOptions{}, err
		}
		if width.Current <= 0 || uint64(width.Current) > maxUint32Value {
			return RecordOptions{}, fmt.Errorf("%w: invalid record width %d", ErrInvalidArgument, width.Current)
		}
		options.Width = uint32(width.Current)
	}
	if options.Height == 0 {
		height, err := c.GetInt(NodeHeight)
		if err != nil {
			return RecordOptions{}, err
		}
		if height.Current <= 0 || uint64(height.Current) > maxUint32Value {
			return RecordOptions{}, fmt.Errorf("%w: invalid record height %d", ErrInvalidArgument, height.Current)
		}
		options.Height = uint32(height.Current)
	}
	if options.FrameRate == 0 {
		frameRate, err := c.GetFloat(NodeResultingFrameRate)
		if err != nil {
			options.FrameRate = 25
		} else {
			options.FrameRate = frameRate.Current
		}
	}
	if options.BitRateKbps == 0 {
		options.BitRateKbps = 1000
	}
	if options.Format == RecordFormatUndefined {
		options.Format = RecordFormatAVI
	}

	if options.Width < minRecordSide || options.Width > maxRecordSide || options.Width%2 != 0 {
		return RecordOptions{}, fmt.Errorf("%w: record width must be even and in [%d, %d]", ErrInvalidArgument, minRecordSide, maxRecordSide)
	}
	if options.Height < minRecordSide || options.Height > maxRecordSide || options.Height%2 != 0 {
		return RecordOptions{}, fmt.Errorf("%w: record height must be even and in [%d, %d]", ErrInvalidArgument, minRecordSide, maxRecordSide)
	}
	if options.Width > uint32(^uint16(0)) || options.Height > uint32(^uint16(0)) {
		return RecordOptions{}, fmt.Errorf("%w: record size exceeds SDK ushort fields", ErrInvalidArgument)
	}
	if options.FrameRate < minRecordFrameFPS || options.FrameRate > maxRecordFrameFPS {
		return RecordOptions{}, fmt.Errorf("%w: record frame rate must be in [1/16, 1000]", ErrInvalidArgument)
	}
	if options.BitRateKbps < minRecordBitRate || options.BitRateKbps > maxRecordBitRate {
		return RecordOptions{}, fmt.Errorf("%w: record bitrate must be in [%d, %d] kbps", ErrInvalidArgument, minRecordBitRate, maxRecordBitRate)
	}
	if options.Format != RecordFormatAVI {
		return RecordOptions{}, fmt.Errorf("%w: unsupported record format %d", ErrInvalidArgument, options.Format)
	}
	return options, nil
}
