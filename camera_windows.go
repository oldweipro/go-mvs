//go:build windows && amd64

package mvs

import (
	"encoding/binary"
	"fmt"
	"runtime"
	"syscall"
	"time"
	"unsafe"
)

const (
	maxSDKPathBytes = 260
	maxUint32Value  = uint64(^uint32(0))
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

func (c *Camera) SetImageNodeNum(num uint32) error {
	if num == 0 {
		return fmt.Errorf("%w: image node number must be greater than 0", ErrInvalidArgument)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.open || c.handle == 0 {
		return ErrCameraClosed
	}
	return c.sdk.driver.setImageNodeNum(c.handle, num)
}

func (c *Camera) SetGrabStrategy(strategy GrabStrategy) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.open || c.handle == 0 {
		return ErrCameraClosed
	}
	return c.sdk.driver.setGrabStrategy(c.handle, strategy)
}

func (c *Camera) SetOutputQueueSize(size uint32) error {
	if size == 0 {
		return fmt.Errorf("%w: output queue size must be greater than 0", ErrInvalidArgument)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.open || c.handle == 0 {
		return ErrCameraClosed
	}
	return c.sdk.driver.setOutputQueueSize(c.handle, size)
}

func (c *Camera) ClearImageBuffer() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.open || c.handle == 0 {
		return ErrCameraClosed
	}
	return c.sdk.driver.clearImageBuffer(c.handle)
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

	if c.recording {
		if err := c.sdk.driver.stopRecord(c.handle); err != nil && firstErr == nil {
			firstErr = err
		}
		c.recording = false
	}

	if c.serialOpen {
		if err := c.sdk.driver.serialPortClose(c.handle); err != nil && firstErr == nil {
			firstErr = err
		}
		c.serialOpen = false
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
	if c.callbackPtr != 0 {
		c.mu.Unlock()
		return nil, fmt.Errorf("%w: frame callback is registered", ErrAcquisitionModeConflict)
	}
	handle := c.handle
	c.mu.Unlock()

	var rawFrame mvFrameOut
	if err := c.sdk.driver.getImageBuffer(handle, &rawFrame, timeoutMilliseconds(timeout)); err != nil {
		return nil, err
	}
	defer c.sdk.driver.freeImageBuffer(handle, &rawFrame)

	length, err := frameDataLength(&rawFrame.StFrameInfo)
	if err != nil {
		return nil, err
	}
	if rawFrame.BufAddr == nil && length > 0 {
		return nil, ErrNilFrameBuffer
	}

	data := make([]byte, length)
	copy(data, unsafe.Slice(rawFrame.BufAddr, length))

	return frameFromInfo(&rawFrame.StFrameInfo, data)
}

func (c *Camera) GetOneFrameTimeout(timeout time.Duration) (*Frame, error) {
	bufferSize, err := c.CurrentFrameBufferSize()
	if err != nil {
		return nil, err
	}
	buffer := make([]byte, bufferSize)
	return c.GetOneFrameInto(buffer, timeout)
}

func (c *Camera) GetOneFrameInto(buffer []byte, timeout time.Duration) (*Frame, error) {
	if len(buffer) == 0 {
		return nil, fmt.Errorf("%w: frame buffer is empty", ErrInvalidArgument)
	}
	if uint64(len(buffer)) > maxUint32Value {
		return nil, fmt.Errorf("%w: frame buffer is larger than UINT_MAX", ErrInvalidArgument)
	}

	c.mu.Lock()
	if !c.open || c.handle == 0 {
		c.mu.Unlock()
		return nil, ErrCameraClosed
	}
	if !c.grabbing {
		c.mu.Unlock()
		return nil, ErrCameraNotGrabbing
	}
	if c.callbackPtr != 0 {
		c.mu.Unlock()
		return nil, fmt.Errorf("%w: frame callback is registered", ErrAcquisitionModeConflict)
	}
	handle := c.handle
	c.mu.Unlock()

	var frameInfo mvFrameOutInfoEx
	if err := c.sdk.driver.getOneFrameTimeout(handle, buffer, &frameInfo, timeoutMilliseconds(timeout)); err != nil {
		return nil, err
	}

	length, err := frameDataLength(&frameInfo)
	if err != nil {
		return nil, err
	}
	if length > len(buffer) {
		return nil, fmt.Errorf("%w: SDK returned %d bytes into %d-byte buffer", ErrInvalidFrameData, length, len(buffer))
	}
	return frameFromInfo(&frameInfo, buffer[:length])
}

func (c *Camera) CurrentFrameBufferSize() (int, error) {
	width, err := c.GetInt(NodeWidth)
	if err != nil {
		return 0, err
	}
	height, err := c.GetInt(NodeHeight)
	if err != nil {
		return 0, err
	}
	pixel, err := c.GetEnum(NodePixelFormat)
	if err != nil {
		return 0, err
	}
	if width.Current <= 0 || height.Current <= 0 || uint64(width.Current) > maxUint32Value || uint64(height.Current) > maxUint32Value {
		return 0, fmt.Errorf("%w: invalid camera size %dx%d", ErrInvalidFrameData, width.Current, height.Current)
	}
	return ExpectedFrameDataLength(uint32(width.Current), uint32(height.Current), pixel.Current)
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

func (c *Camera) SetEnumByString(name string, value string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.open || c.handle == 0 {
		return ErrCameraClosed
	}
	return c.sdk.driver.setEnumValueByString(c.handle, name, value)
}

func (c *Camera) GetEnumEntrySymbolic(name string, value uint32) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.open || c.handle == 0 {
		return "", ErrCameraClosed
	}

	entry := mvCCEnumEntry{Value: value}
	if err := c.sdk.driver.getEnumEntrySymbolic(c.handle, name, &entry); err != nil {
		return "", err
	}
	return byteString(entry.Symbolic[:]), nil
}

func (c *Camera) GetEnumEntries(name string) ([]EnumEntry, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.open || c.handle == 0 {
		return nil, ErrCameraClosed
	}

	var raw mvCCEnumValue
	if err := c.sdk.driver.getEnumValue(c.handle, name, &raw); err != nil {
		return nil, err
	}

	count := int(raw.SupportedNum)
	if count > len(raw.SupportValue) {
		count = len(raw.SupportValue)
	}

	entries := make([]EnumEntry, 0, count)
	for i := 0; i < count; i++ {
		entry := mvCCEnumEntry{Value: raw.SupportValue[i]}
		if err := c.sdk.driver.getEnumEntrySymbolic(c.handle, name, &entry); err != nil {
			return nil, err
		}
		entries = append(entries, EnumEntry{
			Value:    raw.SupportValue[i],
			Symbolic: byteString(entry.Symbolic[:]),
		})
	}
	return entries, nil
}

func (c *Camera) GetBool(name string) (bool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.open || c.handle == 0 {
		return false, ErrCameraClosed
	}
	return c.sdk.driver.getBoolValue(c.handle, name)
}

func (c *Camera) SetBool(name string, value bool) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.open || c.handle == 0 {
		return ErrCameraClosed
	}
	return c.sdk.driver.setBoolValue(c.handle, name, value)
}

func (c *Camera) GetString(name string) (StringValue, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.open || c.handle == 0 {
		return StringValue{}, ErrCameraClosed
	}

	var raw mvCCStringValue
	if err := c.sdk.driver.getStringValue(c.handle, name, &raw); err != nil {
		return StringValue{}, err
	}
	return StringValue{
		Current:   byteString(raw.Current[:]),
		MaxLength: raw.MaxLength,
	}, nil
}

func (c *Camera) SetString(name string, value string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.open || c.handle == 0 {
		return ErrCameraClosed
	}
	return c.sdk.driver.setStringValue(c.handle, name, value)
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

func (c *Camera) FeatureSave(path string) error {
	if err := validateSDKPath(path); err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.open || c.handle == 0 {
		return ErrCameraClosed
	}
	return c.sdk.driver.featureSave(c.handle, path)
}

func (c *Camera) FeatureLoad(path string) error {
	if err := validateSDKPath(path); err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.open || c.handle == 0 {
		return ErrCameraClosed
	}
	return c.sdk.driver.featureLoad(c.handle, path)
}

func (c *Camera) SaveFrameToFile(frame *Frame, path string, options ImageSaveOptions) error {
	if err := validateFrame(frame); err != nil {
		return err
	}
	if err := validateSDKPath(path); err != nil {
		return err
	}
	if uint64(len(frame.Data)) > maxUint32Value {
		return fmt.Errorf("%w: frame data is larger than UINT_MAX", ErrInvalidFrameData)
	}

	imageType := options.Type
	if imageType == ImageTypeUndefined {
		var ok bool
		imageType, ok = ImageTypeFromExtension(path)
		if !ok {
			return fmt.Errorf("%w: image type is not set and cannot be inferred from %q", ErrInvalidArgument, path)
		}
	}

	quality := options.Quality
	if quality == 0 {
		quality = 90
	}
	if imageType == ImageTypeJPEG && (quality <= 50 || quality > 99) {
		return fmt.Errorf("%w: JPEG quality must be in (50, 99]", ErrInvalidArgument)
	}

	pathPtr, err := syscall.BytePtrFromString(path)
	if err != nil {
		return fmt.Errorf("build image path %q: %w", path, err)
	}

	c.mu.Lock()
	if !c.open || c.handle == 0 {
		c.mu.Unlock()
		return ErrCameraClosed
	}
	handle := c.handle
	c.mu.Unlock()

	param := mvSaveImageToFileParamEx{
		Width:       frame.Width,
		Height:      frame.Height,
		PixelType:   frame.PixelType,
		Data:        &frame.Data[0],
		DataLen:     uint32(len(frame.Data)),
		ImageType:   uint32(imageType),
		ImagePath:   pathPtr,
		Quality:     quality,
		MethodValue: int32(options.MethodValue),
	}
	err = c.sdk.driver.saveImageToFileEx(handle, &param)
	runtime.KeepAlive(frame)
	runtime.KeepAlive(pathPtr)
	return err
}

func (c *Camera) ConvertPixelType(frame *Frame, dstPixelType uint32) (*Frame, error) {
	return c.ConvertFrame(frame, PixelConvertOptions{DstPixelType: dstPixelType})
}

func (c *Camera) ConvertFrame(frame *Frame, options PixelConvertOptions) (*Frame, error) {
	if err := validateFrame(frame); err != nil {
		return nil, err
	}
	if options.DstPixelType == 0 || options.DstPixelType == PixelTypeUndefined {
		return nil, fmt.Errorf("%w: destination pixel type is invalid", ErrInvalidArgument)
	}
	if uint64(len(frame.Data)) > maxUint32Value {
		return nil, fmt.Errorf("%w: frame data is larger than UINT_MAX", ErrInvalidFrameData)
	}

	dstSize := options.DstBufferSize
	if dstSize <= 0 {
		var err error
		dstSize, err = ExpectedFrameDataLength(frame.Width, frame.Height, options.DstPixelType)
		if err != nil {
			return nil, err
		}
	}
	if dstSize <= 0 || uint64(dstSize) > maxUint32Value {
		return nil, fmt.Errorf("%w: destination buffer size is invalid: %d", ErrInvalidArgument, dstSize)
	}

	dst := make([]byte, dstSize)

	c.mu.Lock()
	if !c.open || c.handle == 0 {
		c.mu.Unlock()
		return nil, ErrCameraClosed
	}
	handle := c.handle
	c.mu.Unlock()

	param := mvCCPixelConvertParamEx{
		Width:         frame.Width,
		Height:        frame.Height,
		SrcPixelType:  frame.PixelType,
		SrcData:       &frame.Data[0],
		SrcDataLen:    uint32(len(frame.Data)),
		DstPixelType:  options.DstPixelType,
		DstBuffer:     &dst[0],
		DstBufferSize: uint32(len(dst)),
	}
	if err := c.sdk.driver.convertPixelTypeEx(handle, &param); err != nil {
		runtime.KeepAlive(frame)
		runtime.KeepAlive(dst)
		return nil, err
	}
	runtime.KeepAlive(frame)
	runtime.KeepAlive(dst)

	if int(param.DstLen) > len(dst) {
		return nil, fmt.Errorf("%w: SDK returned %d bytes into %d-byte destination buffer", ErrInvalidFrameData, param.DstLen, len(dst))
	}

	converted := *frame
	converted.PixelType = options.DstPixelType
	converted.Data = dst[:param.DstLen]
	converted.DataLength = uint64(param.DstLen)
	return &converted, nil
}

func (c *Camera) RegisterFrameCallback(callback FrameCallback) error {
	if callback == nil {
		return fmt.Errorf("%w: frame callback is nil", ErrInvalidArgument)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.open || c.handle == 0 {
		return ErrCameraClosed
	}
	if c.grabbing {
		return fmt.Errorf("%w: register frame callback before StartGrabbing", ErrAcquisitionModeConflict)
	}
	if c.callbackPtr != 0 {
		return ErrCallbackRegistered
	}

	callbackPtr := syscall.NewCallback(func(dataPtr *byte, info *mvFrameOutInfoEx, user uintptr) uintptr {
		defer func() {
			_ = recover()
		}()
		if dataPtr == nil || info == nil {
			return 0
		}

		length, err := frameDataLength(info)
		if err != nil || length == 0 {
			return 0
		}

		data := make([]byte, length)
		copy(data, unsafe.Slice(dataPtr, length))

		frame, err := frameFromInfo(info, data)
		if err != nil {
			return 0
		}
		callback(frame)
		return 0
	})

	if err := c.sdk.driver.registerImageCallbackEx(c.handle, callbackPtr, 0); err != nil {
		return err
	}
	c.callbackPtr = callbackPtr
	c.callback = callback
	return nil
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

func validateFrame(frame *Frame) error {
	if frame == nil {
		return fmt.Errorf("%w: frame is nil", ErrInvalidArgument)
	}
	if frame.Width == 0 || frame.Height == 0 {
		return fmt.Errorf("%w: empty image size %dx%d", ErrInvalidFrameData, frame.Width, frame.Height)
	}
	if len(frame.Data) == 0 {
		return fmt.Errorf("%w: frame data is empty", ErrInvalidFrameData)
	}
	if uint64(len(frame.Data)) > maxUint32Value {
		return fmt.Errorf("%w: frame data is larger than UINT_MAX", ErrInvalidFrameData)
	}
	return nil
}

func validateSDKPath(path string) error {
	if path == "" {
		return fmt.Errorf("%w: path is empty", ErrInvalidArgument)
	}
	if len(path)+1 > maxSDKPathBytes {
		return fmt.Errorf("%w: SDK path must be shorter than %d bytes", ErrInvalidArgument, maxSDKPathBytes)
	}
	return nil
}

func frameDataLength(info *mvFrameOutInfoEx) (int, error) {
	if info == nil {
		return 0, fmt.Errorf("%w: frame info is nil", ErrInvalidFrameData)
	}
	length := uint64(info.FrameLen)
	if info.FrameLenEx > 0 {
		length = info.FrameLenEx
	}
	maxInt := uint64(int(^uint(0) >> 1))
	if length > maxInt {
		return 0, fmt.Errorf("%w: frame data is too large", ErrInvalidFrameData)
	}
	return int(length), nil
}

func frameFromInfo(info *mvFrameOutInfoEx, data []byte) (*Frame, error) {
	length, err := frameDataLength(info)
	if err != nil {
		return nil, err
	}
	if len(data) < length {
		return nil, fmt.Errorf("%w: need %d bytes, got %d", ErrInvalidFrameData, length, len(data))
	}

	width := uint32(info.Width)
	height := uint32(info.Height)
	if info.ExtendWidth > 0 {
		width = info.ExtendWidth
	}
	if info.ExtendHeight > 0 {
		height = info.ExtendHeight
	}

	frameLength := uint64(info.FrameLen)
	if info.FrameLenEx > 0 {
		frameLength = info.FrameLenEx
	}
	parts := framePartsFromInfo(info)

	return &Frame{
		Width:           width,
		Height:          height,
		PixelType:       info.PixelType,
		FrameNum:        info.FrameNum,
		DataLength:      frameLength,
		DeviceTimestamp: uint64(info.DevTimestampHigh)<<32 | uint64(info.DevTimestampLow),
		HostTimestamp:   info.HostTimestamp,
		ExposureTime:    info.ExposureTime,
		Gain:            info.Gain,
		LostPacketCount: info.LostPacket,
		ExtraType:       info.ExtraType,
		SubImageNum:     info.SubImageNum,
		Parts:           parts,
		Data:            data[:length],
	}, nil
}

func framePartsFromInfo(info *mvFrameOutInfoEx) []FramePart {
	if info == nil || info.SubImageNum == 0 || info.SubImageList == nil {
		return nil
	}

	count := int(info.SubImageNum)
	if count <= 0 {
		return nil
	}
	if count > mvMaxSplitNum && info.ExtraType != FrameExtraMultiParts {
		count = mvMaxSplitNum
	}

	switch info.ExtraType {
	case FrameExtraSubImages:
		return frameSubImagesFromInfo(info.SubImageList, count)
	case FrameExtraMultiParts:
		return frameMultiPartsFromInfo(info.SubImageList, count)
	default:
		return nil
	}
}

func frameSubImagesFromInfo(ptr unsafe.Pointer, count int) []FramePart {
	rawImages := unsafe.Slice((*mvCCImage)(ptr), count)
	parts := make([]FramePart, 0, count)
	for _, raw := range rawImages {
		length, ok := uint64ToInt(raw.ImageLen)
		if !ok || raw.ImageBuf == nil || length == 0 {
			continue
		}
		data := make([]byte, length)
		copy(data, unsafe.Slice(raw.ImageBuf, length))
		parts = append(parts, FramePart{
			DataFormat: uint32(raw.PixelType),
			Width:      raw.Width,
			Height:     raw.Height,
			PixelType:  raw.PixelType,
			Length:     raw.ImageLen,
			Data:       data,
		})
	}
	return parts
}

func frameMultiPartsFromInfo(ptr unsafe.Pointer, count int) []FramePart {
	rawParts := unsafe.Slice((*mvGigeMultiPartInfo)(ptr), count)
	parts := make([]FramePart, 0, count)
	for _, raw := range rawParts {
		length, ok := uint64ToInt(raw.Length)
		if !ok || raw.PartAddr == nil || length == 0 {
			continue
		}
		data := make([]byte, length)
		copy(data, unsafe.Slice(raw.PartAddr, length))

		part := FramePart{
			DataType:         MultiPartDataType(raw.DataType),
			DataFormat:       raw.DataFormat,
			PixelType:        raw.DataFormat,
			SourceID:         raw.SourceID,
			RegionID:         raw.RegionID,
			DataPurposeID:    raw.DataPurposeID,
			Zones:            raw.Zones,
			Length:           raw.Length,
			DataTypeSpecific: raw.DataTypeSpecific.Data,
			Data:             data,
		}
		if raw.DataType <= uint32(MultiPartDataConfidenceMap) {
			part.Width = binary.LittleEndian.Uint32(raw.DataTypeSpecific.Data[0:4])
			part.Height = binary.LittleEndian.Uint32(raw.DataTypeSpecific.Data[4:8])
		}
		parts = append(parts, part)
	}
	return parts
}

func uint64ToInt(value uint64) (int, bool) {
	maxInt := uint64(int(^uint(0) >> 1))
	if value > maxInt {
		return 0, false
	}
	return int(value), true
}
