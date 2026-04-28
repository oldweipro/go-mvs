//go:build windows && amd64

package mvsdk

import (
	"fmt"
	"time"
	"unsafe"
)

func (c *Camera) Info() DeviceInfo {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.info
}

func (c *Camera) IsConnected() (bool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.open || c.handle == 0 {
		return false, ErrCameraClosed
	}
	return c.sdk.driver.isDeviceConnected(c.handle), nil
}

func (c *Camera) ConfigureOptimalPacketSize() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.open || c.handle == 0 {
		return ErrCameraClosed
	}
	if c.info.TransportLayer != DeviceTypeGigE && c.info.TransportLayer != DeviceTypeGentlGigE {
		return nil
	}

	packetSize := c.sdk.driver.getOptimalPacketSize(c.handle)
	switch {
	case packetSize < 0:
		return SDKError{Op: "MV_CC_GetOptimalPacketSize", Code: uint32(packetSize)}
	case packetSize == 0:
		return fmt.Errorf("MV_CC_GetOptimalPacketSize returned 0")
	default:
		return c.sdk.driver.setIntValueEx(c.handle, NodeGevSCPSPacketSize, int64(packetSize))
	}
}

func (c *Camera) StartGrabbing() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.open || c.handle == 0 {
		return ErrCameraClosed
	}
	if c.grabbing {
		return nil
	}
	if err := c.sdk.driver.startGrabbing(c.handle); err != nil {
		return err
	}
	c.grabbing = true
	return nil
}

func (c *Camera) StopGrabbing() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.open || c.handle == 0 {
		return nil
	}
	if !c.grabbing {
		return nil
	}
	if err := c.sdk.driver.stopGrabbing(c.handle); err != nil {
		return err
	}
	c.grabbing = false
	return nil
}

func (c *Camera) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.open && c.handle == 0 {
		return nil
	}

	var firstErr error
	if c.grabbing {
		if err := c.sdk.driver.stopGrabbing(c.handle); err != nil && firstErr == nil {
			firstErr = err
		}
		c.grabbing = false
	}

	if c.handle != 0 {
		if err := c.sdk.driver.closeDevice(c.handle); err != nil && firstErr == nil {
			firstErr = err
		}
		if err := c.sdk.driver.destroyHandle(c.handle); err != nil && firstErr == nil {
			firstErr = err
		}
	}

	c.handle = 0
	c.open = false
	return firstErr
}

func (c *Camera) GetFrame(timeout time.Duration) (*Frame, error) {
	c.mu.Lock()
	if !c.open || c.handle == 0 {
		c.mu.Unlock()
		return nil, ErrCameraClosed
	}
	if !c.grabbing {
		c.mu.Unlock()
		return nil, ErrCameraNotGrabbing
	}
	handle := c.handle
	c.mu.Unlock()

	var rawFrame mvFrameOut
	if err := c.sdk.driver.getImageBuffer(handle, &rawFrame, timeoutMilliseconds(timeout)); err != nil {
		return nil, err
	}
	defer c.sdk.driver.freeImageBuffer(handle, &rawFrame)

	length := int(rawFrame.StFrameInfo.FrameLen)
	if rawFrame.BufAddr == nil && length > 0 {
		return nil, ErrNilFrameBuffer
	}

	data := make([]byte, length)
	copy(data, unsafe.Slice(rawFrame.BufAddr, length))

	width := uint32(rawFrame.StFrameInfo.Width)
	height := uint32(rawFrame.StFrameInfo.Height)
	if rawFrame.StFrameInfo.ExtendWidth > 0 {
		width = rawFrame.StFrameInfo.ExtendWidth
	}
	if rawFrame.StFrameInfo.ExtendHeight > 0 {
		height = rawFrame.StFrameInfo.ExtendHeight
	}

	frameLength := uint64(rawFrame.StFrameInfo.FrameLen)
	if rawFrame.StFrameInfo.FrameLenEx > 0 {
		frameLength = rawFrame.StFrameInfo.FrameLenEx
	}

	return &Frame{
		Width:           width,
		Height:          height,
		PixelType:       rawFrame.StFrameInfo.PixelType,
		FrameNum:        rawFrame.StFrameInfo.FrameNum,
		DataLength:      frameLength,
		DeviceTimestamp: uint64(rawFrame.StFrameInfo.DevTimestampHigh)<<32 | uint64(rawFrame.StFrameInfo.DevTimestampLow),
		HostTimestamp:   rawFrame.StFrameInfo.HostTimestamp,
		ExposureTime:    rawFrame.StFrameInfo.ExposureTime,
		Gain:            rawFrame.StFrameInfo.Gain,
		LostPacketCount: rawFrame.StFrameInfo.LostPacket,
		Data:            data,
	}, nil
}

func (c *Camera) GetFloat(name string) (FloatValue, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.open || c.handle == 0 {
		return FloatValue{}, ErrCameraClosed
	}

	var raw mvCCFloatValue
	if err := c.sdk.driver.getFloatValue(c.handle, name, &raw); err != nil {
		return FloatValue{}, err
	}
	return FloatValue{
		Current: raw.Current,
		Max:     raw.Max,
		Min:     raw.Min,
	}, nil
}

func (c *Camera) SetFloat(name string, value float32) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.open || c.handle == 0 {
		return ErrCameraClosed
	}
	return c.sdk.driver.setFloatValue(c.handle, name, value)
}

func (c *Camera) GetInt(name string) (IntValue, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.open || c.handle == 0 {
		return IntValue{}, ErrCameraClosed
	}

	var raw mvCCIntValueEx
	if err := c.sdk.driver.getIntValueEx(c.handle, name, &raw); err != nil {
		return IntValue{}, err
	}
	return IntValue{
		Current:   raw.Current,
		Max:       raw.Max,
		Min:       raw.Min,
		Increment: raw.Increment,
	}, nil
}

func (c *Camera) SetInt(name string, value int64) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.open || c.handle == 0 {
		return ErrCameraClosed
	}
	return c.sdk.driver.setIntValueEx(c.handle, name, value)
}

func (c *Camera) GetEnum(name string) (EnumValue, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.open || c.handle == 0 {
		return EnumValue{}, ErrCameraClosed
	}

	var raw mvCCEnumValue
	if err := c.sdk.driver.getEnumValue(c.handle, name, &raw); err != nil {
		return EnumValue{}, err
	}

	count := int(raw.SupportedNum)
	if count > len(raw.SupportValue) {
		count = len(raw.SupportValue)
	}
	supported := make([]uint32, count)
	copy(supported, raw.SupportValue[:count])

	return EnumValue{
		Current:   raw.Current,
		Supported: supported,
	}, nil
}

func (c *Camera) SetEnum(name string, value uint32) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.open || c.handle == 0 {
		return ErrCameraClosed
	}
	return c.sdk.driver.setEnumValue(c.handle, name, value)
}

func (c *Camera) Command(name string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.open || c.handle == 0 {
		return ErrCameraClosed
	}
	return c.sdk.driver.setCommandValue(c.handle, name)
}

func (c *Camera) SetTriggerMode(enabled bool) error {
	if enabled {
		if err := c.SetEnum(NodeTriggerMode, TriggerModeOn); err != nil {
			return err
		}
		return c.SetEnum(NodeTriggerSource, TriggerSourceSoftware)
	}
	return c.SetEnum(NodeTriggerMode, TriggerModeOff)
}

func (c *Camera) TriggerOnce() error {
	return c.Command(NodeTriggerSoftware)
}

func timeoutMilliseconds(timeout time.Duration) uint32 {
	if timeout <= 0 {
		return 0
	}
	ms := timeout / time.Millisecond
	if ms == 0 {
		return 1
	}
	if ms > time.Duration(^uint32(0)) {
		return ^uint32(0)
	}
	return uint32(ms)
}
