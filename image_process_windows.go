//go:build windows && amd64

package mvs

import (
	"fmt"
	"runtime"
	"syscall"
)

func (c *Camera) RotateFrame(frame *Frame, angle RotationAngle) (*Frame, error) {
	if err := validateFrame(frame); err != nil {
		return nil, err
	}
	if angle != RotationAngle90 && angle != RotationAngle180 && angle != RotationAngle270 {
		return nil, fmt.Errorf("%w: invalid rotation angle %d", ErrInvalidArgument, angle)
	}

	dst := make([]byte, len(frame.Data))
	param := mvCCRotateImageParam{
		PixelType:     frame.PixelType,
		Width:         frame.Width,
		Height:        frame.Height,
		SrcData:       &frame.Data[0],
		SrcDataLen:    uint32(len(frame.Data)),
		DstBuffer:     &dst[0],
		DstBufferSize: uint32(len(dst)),
		RotationAngle: uint32(angle),
	}

	handle, err := c.cameraHandle()
	if err != nil {
		return nil, err
	}
	if err := c.sdk.driver.rotateImage(handle, &param); err != nil {
		runtime.KeepAlive(frame)
		runtime.KeepAlive(dst)
		return nil, err
	}
	runtime.KeepAlive(frame)
	runtime.KeepAlive(dst)

	if int(param.DstLen) > len(dst) {
		return nil, fmt.Errorf("%w: SDK returned %d bytes into %d-byte destination buffer", ErrInvalidFrameData, param.DstLen, len(dst))
	}
	rotated := *frame
	rotated.Width = param.Width
	rotated.Height = param.Height
	rotated.Data = dst[:param.DstLen]
	rotated.DataLength = uint64(param.DstLen)
	return &rotated, nil
}

func (c *Camera) FlipFrame(frame *Frame, flip FlipType) (*Frame, error) {
	if err := validateFrame(frame); err != nil {
		return nil, err
	}
	if flip != FlipVertical && flip != FlipHorizontal {
		return nil, fmt.Errorf("%w: invalid flip type %d", ErrInvalidArgument, flip)
	}

	dst := make([]byte, len(frame.Data))
	param := mvCCFlipImageParam{
		PixelType:     frame.PixelType,
		Width:         frame.Width,
		Height:        frame.Height,
		SrcData:       &frame.Data[0],
		SrcDataLen:    uint32(len(frame.Data)),
		DstBuffer:     &dst[0],
		DstBufferSize: uint32(len(dst)),
		FlipType:      uint32(flip),
	}

	handle, err := c.cameraHandle()
	if err != nil {
		return nil, err
	}
	if err := c.sdk.driver.flipImage(handle, &param); err != nil {
		runtime.KeepAlive(frame)
		runtime.KeepAlive(dst)
		return nil, err
	}
	runtime.KeepAlive(frame)
	runtime.KeepAlive(dst)

	if int(param.DstLen) > len(dst) {
		return nil, fmt.Errorf("%w: SDK returned %d bytes into %d-byte destination buffer", ErrInvalidFrameData, param.DstLen, len(dst))
	}
	flipped := *frame
	flipped.Data = dst[:param.DstLen]
	flipped.DataLength = uint64(param.DstLen)
	return &flipped, nil
}

func (c *Camera) SetBayerCvtQuality(quality InterpolationMethod) error {
	handle, err := c.cameraHandle()
	if err != nil {
		return err
	}
	return c.sdk.driver.setBayerCvtQuality(handle, quality)
}

func (c *Camera) SetBayerFilterEnable(enabled bool) error {
	handle, err := c.cameraHandle()
	if err != nil {
		return err
	}
	return c.sdk.driver.setBayerFilterEnable(handle, enabled)
}

func (c *Camera) SetBayerGammaValue(gamma float32) error {
	if gamma < 0.1 || gamma > 4.0 {
		return fmt.Errorf("%w: gamma must be in [0.1, 4.0]", ErrInvalidArgument)
	}
	handle, err := c.cameraHandle()
	if err != nil {
		return err
	}
	return c.sdk.driver.setBayerGammaValue(handle, gamma)
}

func (c *Camera) SetGammaValue(pixelType uint32, gamma float32) error {
	if pixelType == 0 || pixelType == PixelTypeUndefined {
		return fmt.Errorf("%w: pixel type is invalid", ErrInvalidArgument)
	}
	if gamma < 0.1 || gamma > 4.0 {
		return fmt.Errorf("%w: gamma must be in [0.1, 4.0]", ErrInvalidArgument)
	}
	handle, err := c.cameraHandle()
	if err != nil {
		return err
	}
	return c.sdk.driver.setGammaValue(handle, pixelType, gamma)
}

func (c *Camera) SetBayerGammaParam(options GammaOptions) error {
	if options.Type > GammaTypeSRGBToLRGB {
		return fmt.Errorf("%w: invalid gamma type %d", ErrInvalidArgument, options.Type)
	}
	if options.Type == GammaTypeValue && options.Value == 0 {
		return fmt.Errorf("%w: gamma value is required", ErrInvalidArgument)
	}
	if options.Value != 0 && (options.Value < 0.1 || options.Value > 4.0) {
		return fmt.Errorf("%w: gamma must be in [0.1, 4.0]", ErrInvalidArgument)
	}
	if options.Type == GammaTypeUserCurve && len(options.Curve) == 0 {
		return fmt.Errorf("%w: gamma curve is required", ErrInvalidArgument)
	}
	if uint64(len(options.Curve)) > maxUint32Value {
		return fmt.Errorf("%w: gamma curve is larger than UINT_MAX", ErrInvalidArgument)
	}

	param := mvCCGammaParam{
		GammaType:        uint32(options.Type),
		GammaValue:       options.Value,
		GammaCurveBufLen: uint32(len(options.Curve)),
	}
	if len(options.Curve) > 0 {
		param.GammaCurveBuf = &options.Curve[0]
	}

	handle, err := c.cameraHandle()
	if err != nil {
		return err
	}
	err = c.sdk.driver.setBayerGammaParam(handle, &param)
	runtime.KeepAlive(options.Curve)
	return err
}

func (c *Camera) SetBayerCCMParam(options CCMOptions) error {
	for _, value := range options.Matrix {
		if value < -8192 || value > 8192 {
			return fmt.Errorf("%w: CCM matrix values must be in [-8192, 8192]", ErrInvalidArgument)
		}
	}

	param := mvCCCCMParam{
		CCMat: options.Matrix,
	}
	if options.Enabled {
		param.CCMEnable = 1
	}

	handle, err := c.cameraHandle()
	if err != nil {
		return err
	}
	return c.sdk.driver.setBayerCCMParam(handle, &param)
}

func (c *Camera) SetBayerCCMParamEx(options CCMOptionsEx) error {
	for _, value := range options.Matrix {
		if value < -65536 || value > 65536 {
			return fmt.Errorf("%w: CCM matrix values must be in [-65536, 65536]", ErrInvalidArgument)
		}
	}
	if options.Scale > 65536 || (options.Scale != 0 && options.Scale&(options.Scale-1) != 0) {
		return fmt.Errorf("%w: CCM scale must be a power of two and <= 65536", ErrInvalidArgument)
	}

	param := mvCCCCMParamEx{
		CCMat:    options.Matrix,
		CCMScale: options.Scale,
	}
	if options.Enabled {
		param.CCMEnable = 1
	}

	handle, err := c.cameraHandle()
	if err != nil {
		return err
	}
	return c.sdk.driver.setBayerCCMParamEx(handle, &param)
}

func (c *Camera) AdjustContrast(frame *Frame, options ContrastOptions) (*Frame, error) {
	if err := validateFrame(frame); err != nil {
		return nil, err
	}
	if frame.Width < 8 || frame.Height < 8 {
		return nil, fmt.Errorf("%w: contrast image size must be at least 8x8", ErrInvalidArgument)
	}
	if options.Factor < 1 || options.Factor > 10000 {
		return nil, fmt.Errorf("%w: contrast factor must be in [1, 10000]", ErrInvalidArgument)
	}

	dst, err := makeOutputBuffer(frame, options.DstBufferSize)
	if err != nil {
		return nil, err
	}
	param := mvCCContrastParam{
		Width:          frame.Width,
		Height:         frame.Height,
		SrcBuffer:      &frame.Data[0],
		SrcBufferLen:   uint32(len(frame.Data)),
		PixelType:      frame.PixelType,
		DstBuffer:      &dst[0],
		DstBufferSize:  uint32(len(dst)),
		ContrastFactor: options.Factor,
	}

	handle, err := c.cameraHandle()
	if err != nil {
		return nil, err
	}
	if err := c.sdk.driver.imageContrast(handle, &param); err != nil {
		runtime.KeepAlive(frame)
		runtime.KeepAlive(dst)
		return nil, err
	}
	runtime.KeepAlive(frame)
	runtime.KeepAlive(dst)
	return processedFrameFromBuffer(frame, dst, param.DstBufferLen)
}

func (c *Camera) CorrectPurpleFringing(frame *Frame, options PurpleFringingOptions) (*Frame, error) {
	if err := validateFrame(frame); err != nil {
		return nil, err
	}
	if frame.Width < 4 || frame.Height < 4 {
		return nil, fmt.Errorf("%w: purple-fringing image size must be at least 4x4", ErrInvalidArgument)
	}
	if options.KernelSize != 3 && options.KernelSize != 5 && options.KernelSize != 7 && options.KernelSize != 9 {
		return nil, fmt.Errorf("%w: kernel size must be one of 3, 5, 7, 9", ErrInvalidArgument)
	}
	if options.EdgeThreshold > 2040 {
		return nil, fmt.Errorf("%w: edge threshold must be in [0, 2040]", ErrInvalidArgument)
	}

	dst, err := makeOutputBuffer(frame, options.DstBufferSize)
	if err != nil {
		return nil, err
	}
	param := mvCCPurpleFringingParam{
		Width:         frame.Width,
		Height:        frame.Height,
		SrcBuffer:     &frame.Data[0],
		SrcBufferLen:  uint32(len(frame.Data)),
		PixelType:     frame.PixelType,
		DstBuffer:     &dst[0],
		DstBufferSize: uint32(len(dst)),
		KernelSize:    options.KernelSize,
		EdgeThreshold: options.EdgeThreshold,
	}

	handle, err := c.cameraHandle()
	if err != nil {
		return nil, err
	}
	if err := c.sdk.driver.purpleFringing(handle, &param); err != nil {
		runtime.KeepAlive(frame)
		runtime.KeepAlive(dst)
		return nil, err
	}
	runtime.KeepAlive(frame)
	runtime.KeepAlive(dst)
	return processedFrameFromBuffer(frame, dst, param.DstBufferLen)
}

func (c *Camera) SetISPConfig(path string) error {
	if err := validateSDKPath(path); err != nil {
		return err
	}
	pathPtr, err := syscall.BytePtrFromString(path)
	if err != nil {
		return fmt.Errorf("build ISP config path %q: %w", path, err)
	}
	param := mvCCISPConfigParam{ConfigPath: pathPtr}

	handle, err := c.cameraHandle()
	if err != nil {
		return err
	}
	err = c.sdk.driver.setISPConfig(handle, &param)
	runtime.KeepAlive(pathPtr)
	return err
}

func (c *Camera) ISPProcess(frame *Frame, options PixelConvertOptions) (*Frame, error) {
	if err := validateFrame(frame); err != nil {
		return nil, err
	}
	dstPixelType := options.DstPixelType
	if dstPixelType == 0 || dstPixelType == PixelTypeUndefined {
		dstPixelType = frame.PixelType
	}

	dstSize := options.DstBufferSize
	if dstSize <= 0 {
		var err error
		dstSize, err = ExpectedFrameDataLength(frame.Width, frame.Height, dstPixelType)
		if err != nil {
			return nil, err
		}
	}
	if dstSize <= 0 || uint64(dstSize) > maxUint32Value {
		return nil, fmt.Errorf("%w: invalid ISP output buffer size %d", ErrInvalidArgument, dstSize)
	}
	dst := make([]byte, dstSize)

	input := mvCCImage{
		Width:        frame.Width,
		Height:       frame.Height,
		PixelType:    frame.PixelType,
		ImageBuf:     &frame.Data[0],
		ImageBufSize: uint64(len(frame.Data)),
		ImageLen:     uint64(len(frame.Data)),
	}
	output := mvCCImage{
		Width:        frame.Width,
		Height:       frame.Height,
		PixelType:    dstPixelType,
		ImageBuf:     &dst[0],
		ImageBufSize: uint64(len(dst)),
	}

	handle, err := c.cameraHandle()
	if err != nil {
		return nil, err
	}
	if err := c.sdk.driver.ispProcess(handle, &input, &output); err != nil {
		runtime.KeepAlive(frame)
		runtime.KeepAlive(dst)
		return nil, err
	}
	runtime.KeepAlive(frame)
	runtime.KeepAlive(dst)

	length, ok := uint64ToInt(output.ImageLen)
	if !ok || length > len(dst) {
		return nil, fmt.Errorf("%w: SDK returned %d bytes into %d-byte ISP output buffer", ErrInvalidFrameData, output.ImageLen, len(dst))
	}
	processed := *frame
	processed.Width = output.Width
	processed.Height = output.Height
	processed.PixelType = output.PixelType
	processed.Data = dst[:length]
	processed.DataLength = output.ImageLen
	return &processed, nil
}

func (c *Camera) DecodeHighBandwidthData(data []byte, options HBDecodeOptions) (*Frame, FrameSpecInfo, error) {
	if len(data) == 0 {
		return nil, FrameSpecInfo{}, fmt.Errorf("%w: high-bandwidth data is empty", ErrInvalidArgument)
	}
	if uint64(len(data)) > maxUint32Value {
		return nil, FrameSpecInfo{}, fmt.Errorf("%w: high-bandwidth data is larger than UINT_MAX", ErrInvalidArgument)
	}
	if options.DstBufferSize <= 0 || uint64(options.DstBufferSize) > maxUint32Value {
		return nil, FrameSpecInfo{}, fmt.Errorf("%w: invalid high-bandwidth output buffer size %d", ErrInvalidArgument, options.DstBufferSize)
	}

	dst := make([]byte, options.DstBufferSize)
	param := mvCCHBDecodeParam{
		SrcBuffer:     &data[0],
		SrcLen:        uint32(len(data)),
		DstBuffer:     &dst[0],
		DstBufferSize: uint32(len(dst)),
	}

	handle, err := c.cameraHandle()
	if err != nil {
		return nil, FrameSpecInfo{}, err
	}
	if err := c.sdk.driver.hbDecode(handle, &param); err != nil {
		runtime.KeepAlive(data)
		runtime.KeepAlive(dst)
		return nil, FrameSpecInfo{}, err
	}
	runtime.KeepAlive(data)
	runtime.KeepAlive(dst)

	if int(param.DstBufferLen) > len(dst) {
		return nil, FrameSpecInfo{}, fmt.Errorf("%w: SDK returned %d bytes into %d-byte high-bandwidth output buffer", ErrInvalidFrameData, param.DstBufferLen, len(dst))
	}
	frame := &Frame{
		Width:      param.Width,
		Height:     param.Height,
		PixelType:  param.DstPixelType,
		DataLength: uint64(param.DstBufferLen),
		Data:       dst[:param.DstBufferLen],
	}
	return frame, frameSpecInfoFromRaw(param.FrameSpecInfo), nil
}

func (c *Camera) HBDecode(data []byte, options HBDecodeOptions) (*Frame, FrameSpecInfo, error) {
	return c.DecodeHighBandwidthData(data, options)
}

func (c *Camera) DecodeHighBandwidthFrame(frame *Frame, options HBDecodeOptions) (*Frame, FrameSpecInfo, error) {
	if err := validateFrame(frame); err != nil {
		return nil, FrameSpecInfo{}, err
	}
	decoded, spec, err := c.DecodeHighBandwidthData(frame.Data, options)
	if err != nil {
		return nil, FrameSpecInfo{}, err
	}
	decoded.FrameNum = frame.FrameNum
	decoded.DeviceTimestamp = frame.DeviceTimestamp
	decoded.HostTimestamp = frame.HostTimestamp
	return decoded, spec, nil
}

func (c *Camera) HBDecodeFrame(frame *Frame, options HBDecodeOptions) (*Frame, FrameSpecInfo, error) {
	return c.DecodeHighBandwidthFrame(frame, options)
}

func (c *Camera) ReconstructImage(frame *Frame, options ReconstructImageOptions) ([]*Frame, error) {
	if err := validateFrame(frame); err != nil {
		return nil, err
	}
	if options.ExposureNum == 0 || options.ExposureNum > mvMaxSplitNum {
		return nil, fmt.Errorf("%w: exposure num must be in [1, %d]", ErrInvalidArgument, mvMaxSplitNum)
	}
	method := options.Method
	if method == 0 {
		method = ImageReconstructSplitByLine
	}
	if method != ImageReconstructSplitByLine {
		return nil, fmt.Errorf("%w: unsupported image reconstruct method %d", ErrInvalidArgument, method)
	}

	buffers, err := makeReconstructBuffers(frame, options)
	if err != nil {
		return nil, err
	}
	param := mvReconstructImageParam{
		Width:             frame.Width,
		Height:            frame.Height,
		PixelType:         frame.PixelType,
		SrcData:           &frame.Data[0],
		SrcDataLen:        uint32(len(frame.Data)),
		ExposureNum:       options.ExposureNum,
		ReconstructMethod: uint32(method),
	}
	for i := 0; i < int(options.ExposureNum); i++ {
		param.DstBufferList[i].Buffer = &buffers[i][0]
		param.DstBufferList[i].BufferSize = uint32(len(buffers[i]))
	}

	handle, err := c.cameraHandle()
	if err != nil {
		return nil, err
	}
	if err := c.sdk.driver.reconstructImage(handle, &param); err != nil {
		runtime.KeepAlive(frame)
		runtime.KeepAlive(buffers)
		return nil, err
	}
	runtime.KeepAlive(frame)
	runtime.KeepAlive(buffers)

	frames := make([]*Frame, 0, int(options.ExposureNum))
	for i := 0; i < int(options.ExposureNum); i++ {
		out := param.DstBufferList[i]
		length := int(out.BufferLen)
		if length > len(buffers[i]) {
			return nil, fmt.Errorf("%w: SDK returned %d bytes into %d-byte reconstruct output buffer", ErrInvalidFrameData, out.BufferLen, len(buffers[i]))
		}
		width := out.Width
		if width == 0 {
			width = frame.Width
		}
		height := out.Height
		if height == 0 {
			height = frame.Height / options.ExposureNum
		}
		pixelType := out.PixelType
		if pixelType == 0 {
			pixelType = frame.PixelType
		}
		reconstructed := *frame
		reconstructed.Width = width
		reconstructed.Height = height
		reconstructed.PixelType = pixelType
		reconstructed.Data = buffers[i][:length]
		reconstructed.DataLength = uint64(length)
		frames = append(frames, &reconstructed)
	}
	return frames, nil
}

func makeOutputBuffer(frame *Frame, requested int) ([]byte, error) {
	size := requested
	if size <= 0 {
		size = len(frame.Data)
	}
	if size <= 0 || uint64(size) > maxUint32Value {
		return nil, fmt.Errorf("%w: invalid output buffer size %d", ErrInvalidArgument, size)
	}
	return make([]byte, size), nil
}

func processedFrameFromBuffer(source *Frame, buffer []byte, length uint32) (*Frame, error) {
	if int(length) > len(buffer) {
		return nil, fmt.Errorf("%w: SDK returned %d bytes into %d-byte output buffer", ErrInvalidFrameData, length, len(buffer))
	}
	processed := *source
	processed.Data = buffer[:length]
	processed.DataLength = uint64(length)
	return &processed, nil
}

func makeReconstructBuffers(frame *Frame, options ReconstructImageOptions) ([][]byte, error) {
	count := int(options.ExposureNum)
	if len(options.DstBufferSizes) > 0 && len(options.DstBufferSizes) < count {
		return nil, fmt.Errorf("%w: destination buffer sizes must cover every exposure", ErrInvalidArgument)
	}

	defaultSize := options.DefaultBufferSize
	if defaultSize <= 0 {
		outHeight := frame.Height / options.ExposureNum
		if outHeight == 0 {
			return nil, fmt.Errorf("%w: reconstructed image height is zero", ErrInvalidArgument)
		}
		var err error
		defaultSize, err = ExpectedFrameDataLength(frame.Width, outHeight, frame.PixelType)
		if err != nil {
			defaultSize = (len(frame.Data) + count - 1) / count
		}
	}

	buffers := make([][]byte, count)
	for i := 0; i < count; i++ {
		size := defaultSize
		if len(options.DstBufferSizes) > 0 {
			size = options.DstBufferSizes[i]
		}
		if size <= 0 || uint64(size) > maxUint32Value {
			return nil, fmt.Errorf("%w: invalid reconstruct output buffer size %d", ErrInvalidArgument, size)
		}
		buffers[i] = make([]byte, size)
	}
	return buffers, nil
}

func frameSpecInfoFromRaw(raw mvCCFrameSpecInfo) FrameSpecInfo {
	return FrameSpecInfo{
		SecondCount:       raw.SecondCount,
		CycleCount:        raw.CycleCount,
		CycleOffset:       raw.CycleOffset,
		Gain:              raw.Gain,
		ExposureTime:      raw.ExposureTime,
		AverageBrightness: raw.AverageBrightness,
		Red:               raw.Red,
		Green:             raw.Green,
		Blue:              raw.Blue,
		FrameCounter:      raw.FrameCounter,
		TriggerIndex:      raw.TriggerIndex,
		Input:             raw.Input,
		Output:            raw.Output,
		OffsetX:           raw.OffsetX,
		OffsetY:           raw.OffsetY,
		FrameWidth:        raw.FrameWidth,
		FrameHeight:       raw.FrameHeight,
	}
}

func (c *Camera) cameraHandle() (uintptr, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.open || c.handle == 0 {
		return 0, ErrCameraClosed
	}
	return c.handle, nil
}
