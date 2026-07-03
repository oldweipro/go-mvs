//go:build integration && windows && amd64

package mvs

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestCameraIntegration(t *testing.T) {
	sdk, err := New(Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := sdk.Initialize(); err != nil {
		t.Fatal(err)
	}
	defer sdk.Finalize()

	devices, err := sdk.EnumerateDefaultDevices()
	if err != nil {
		t.Fatal(err)
	}
	if len(devices) == 0 {
		skipOrFailNoDevices(t)
	}

	assertInterfaceIntegration(t, sdk)
	assertCameraLinkSerialPortList(t, sdk)
	assertGenTLIntegration(t, sdk)

	device := devices[0]
	if strings.TrimSpace(device.SerialNumber) == "" {
		t.Fatalf("device serial number is empty: %+v", device)
	}
	if !sdk.IsDeviceAccessible(device, AccessExclusive) {
		skipOrFailDeviceInaccessible(t, device)
	}

	camera, err := sdk.OpenDeviceBySerial(device.SerialNumber, AccessExclusive)
	if err != nil {
		t.Fatal(err)
	}
	defer camera.Close()

	if err := camera.ConfigureOptimalPacketSize(); err != nil {
		t.Fatal(err)
	}
	if err := camera.SetImageNodeNum(3); err != nil {
		t.Fatal(err)
	}
	if err := camera.SetGrabStrategy(GrabStrategyOneByOne); err != nil {
		t.Fatal(err)
	}
	if err := camera.SetEnumByString(NodeTriggerMode, "Off"); err != nil {
		t.Fatal(err)
	}

	entries, err := camera.GetEnumEntries(NodeTriggerMode)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) == 0 {
		t.Fatal("TriggerMode enum entries are empty")
	}

	width, err := camera.GetInt(NodeWidth)
	if err != nil {
		t.Fatal(err)
	}
	height, err := camera.GetInt(NodeHeight)
	if err != nil {
		t.Fatal(err)
	}
	if width.Current <= 0 || height.Current <= 0 {
		t.Fatalf("invalid camera size: %d x %d", width.Current, height.Current)
	}
	bufferSize, err := camera.CurrentFrameBufferSize()
	if err != nil {
		t.Fatal(err)
	}
	if bufferSize <= 0 {
		t.Fatalf("invalid frame buffer size: %d", bufferSize)
	}

	featurePath := filepath.Join(t.TempDir(), "camera-features.mfs")
	if err := camera.FeatureSave(featurePath); err != nil {
		t.Fatal(err)
	}
	if info, err := os.Stat(featurePath); err != nil {
		t.Fatal(err)
	} else if info.Size() == 0 {
		t.Fatal("feature save produced an empty file")
	}
	if os.Getenv("MVS_TEST_FEATURE_LOAD") == "1" {
		if err := camera.FeatureLoad(featurePath); err != nil {
			t.Fatal(err)
		}
	}

	if err := camera.StartGrabbing(); err != nil {
		t.Fatal(err)
	}
	defer camera.StopGrabbing()
	if err := camera.ClearImageBuffer(); err != nil {
		t.Fatal(err)
	}

	frame, err := camera.GetFrame(2 * time.Second)
	if err != nil {
		t.Fatal(err)
	}
	if frame.Width == 0 || frame.Height == 0 || len(frame.Data) == 0 {
		t.Fatalf("invalid frame: %+v bytes=%d", frame, len(frame.Data))
	}
	if _, err := frame.Image(); err != nil {
		t.Fatalf("frame image conversion failed for %s: %v", frame.PixelTypeName(), err)
	}

	oneFrameBuffer := make([]byte, bufferSize)
	oneFrame, err := camera.GetOneFrameInto(oneFrameBuffer, 2*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	if len(oneFrame.Data) == 0 {
		t.Fatal("GetOneFrameInto returned empty frame data")
	}

	converted, err := camera.ConvertPixelType(frame, PixelTypeBGR8Packed)
	if err != nil {
		t.Fatal(err)
	}
	if converted.PixelType != PixelTypeBGR8Packed || len(converted.Data) == 0 {
		t.Fatalf("invalid converted frame: pixel=%s bytes=%d", converted.PixelTypeName(), len(converted.Data))
	}

	assertImageProcessingIntegration(t, camera, frame)
	assertEventIntegration(t, camera)
	assertRecordIntegration(t, camera)
	assertDeviceFileIntegration(t, camera)
	assertSerialIntegration(t, camera)

	imagePath := filepath.Join(t.TempDir(), "frame.jpg")
	if err := camera.SaveFrameToFile(frame, imagePath, ImageSaveOptions{Type: ImageTypeJPEG, Quality: 90}); err != nil {
		t.Fatal(err)
	}
	if info, err := os.Stat(imagePath); err != nil {
		t.Fatal(err)
	} else if info.Size() == 0 {
		t.Fatal("image save produced an empty file")
	}
}

func assertInterfaceIntegration(t *testing.T, sdk *SDK) {
	t.Helper()

	interfaces, err := sdk.EnumerateInterfaces(DefaultInterfaceTypeMask)
	if err != nil {
		if IsSDKErrorCode(err, ErrSupport) {
			t.Logf("interface enumeration unsupported: %v", err)
			return
		}
		t.Fatal(err)
	}
	if len(interfaces) == 0 {
		t.Log("no frame-grabber/interface entries found")
		return
	}

	iface, err := sdk.OpenInterface(interfaces[0])
	if err != nil {
		if IsSDKErrorCode(err, ErrSupport) || IsSDKErrorCode(err, ErrResource) || IsSDKErrorCode(err, ErrBusy) {
			t.Logf("interface open skipped: %v", err)
			return
		}
		t.Fatal(err)
	}
	defer iface.Close()

	devices, err := iface.EnumerateDevices()
	if err != nil {
		t.Fatal(err)
	}
	if len(devices) == 0 {
		t.Logf("interface %q has no devices", interfaces[0].InterfaceID)
	}
}

func assertCameraLinkSerialPortList(t *testing.T, sdk *SDK) {
	t.Helper()

	if _, err := sdk.GetCameraLinkSerialPortList(); err != nil {
		if IsSDKErrorCode(err, ErrSupport) || IsSDKErrorCode(err, ErrLoadLibrary) {
			t.Logf("CameraLink serial port list skipped: %v", err)
			return
		}
		t.Fatal(err)
	}
}

func assertGenTLIntegration(t *testing.T, sdk *SDK) {
	t.Helper()

	ctiPath := os.Getenv("MVS_TEST_GENTL_CTI")
	if ctiPath == "" {
		return
	}

	interfaces, err := sdk.EnumerateGenTLInterfaces(ctiPath)
	if err != nil {
		t.Fatal(err)
	}
	defer sdk.UnloadGenTLLibrary(ctiPath)
	if len(interfaces) == 0 {
		t.Logf("no GenTL interfaces from %s", ctiPath)
		return
	}
	if _, err := sdk.EnumerateGenTLDevices(interfaces[0]); err != nil {
		t.Fatal(err)
	}
}

func assertImageProcessingIntegration(t *testing.T, camera *Camera, frame *Frame) {
	t.Helper()

	switch frame.PixelType {
	case PixelTypeMono8, PixelTypeRGB8Packed, PixelTypeBGR8Packed:
		rotated, err := camera.RotateFrame(frame, RotationAngle180)
		if err != nil {
			t.Fatal(err)
		}
		if rotated.Width != frame.Width || rotated.Height != frame.Height || len(rotated.Data) == 0 {
			t.Fatalf("invalid rotated frame: %+v bytes=%d", rotated, len(rotated.Data))
		}

		flipped, err := camera.FlipFrame(frame, FlipHorizontal)
		if err != nil {
			t.Fatal(err)
		}
		if flipped.Width != frame.Width || flipped.Height != frame.Height || len(flipped.Data) == 0 {
			t.Fatalf("invalid flipped frame: %+v bytes=%d", flipped, len(flipped.Data))
		}
	default:
		t.Logf("rotate/flip skipped for %s", frame.PixelTypeName())
	}

	contrasted, err := camera.AdjustContrast(frame, ContrastOptions{Factor: 1000})
	if err != nil {
		if IsSDKErrorCode(err, ErrSupport) || IsSDKErrorCode(err, ErrParameter) {
			t.Logf("contrast skipped for %s: %v", frame.PixelTypeName(), err)
		} else {
			t.Fatal(err)
		}
	} else if len(contrasted.Data) == 0 {
		t.Fatal("contrast returned empty frame")
	}

	if frame.Height%2 == 0 && !IsCompressedPixelType(frame.PixelType) {
		frames, err := camera.ReconstructImage(frame, ReconstructImageOptions{ExposureNum: 2})
		if err != nil {
			if IsSDKErrorCode(err, ErrSupport) || IsSDKErrorCode(err, ErrParameter) {
				t.Logf("reconstruct skipped: %v", err)
			} else {
				t.Fatal(err)
			}
		} else if len(frames) != 2 || len(frames[0].Data) == 0 || len(frames[1].Data) == 0 {
			t.Fatalf("invalid reconstructed frames: %d", len(frames))
		}
	}
}

func assertEventIntegration(t *testing.T, camera *Camera) {
	t.Helper()

	if err := camera.RegisterAllEventCallback(func(EventInfo) {}); err != nil {
		if IsSDKErrorCode(err, ErrSupport) {
			t.Logf("event callback skipped: %v", err)
			return
		}
		t.Fatal(err)
	}

	eventName := os.Getenv("MVS_TEST_EVENT_NAME")
	if eventName == "" {
		return
	}
	if err := camera.RegisterEventCallback(eventName, func(EventInfo) {}); err != nil {
		t.Fatal(err)
	}
	if err := camera.EventNotificationOn(eventName); err != nil {
		t.Fatal(err)
	}
	defer camera.EventNotificationOff(eventName)
}

func assertRecordIntegration(t *testing.T, camera *Camera) {
	t.Helper()

	if os.Getenv("MVS_TEST_RECORD") != "1" {
		return
	}

	if err := camera.StopGrabbing(); err != nil {
		t.Fatal(err)
	}

	if err := camera.SetEnum(NodeImageCompressionMode, 0); err != nil {
		t.Logf("set ImageCompressionMode off skipped: %v", err)
	}

	width, err := camera.GetInt(NodeWidth)
	if err != nil {
		t.Fatal(err)
	}
	height, err := camera.GetInt(NodeHeight)
	if err != nil {
		t.Fatal(err)
	}
	pixel, err := camera.GetEnum(NodePixelFormat)
	if err != nil {
		t.Fatal(err)
	}
	frameRate, err := camera.GetFloat(NodeResultingFrameRate)
	if err != nil {
		t.Fatal(err)
	}

	recordPath := "go-mvs-record-integration.avi"
	_ = os.Remove(recordPath)
	t.Cleanup(func() {
		_ = os.Remove(recordPath)
	})

	t.Logf(
		"record params width=%d height=%d pixel=%s fps=%f bitrate=%d path=%s",
		width.Current,
		height.Current,
		PixelTypeName(pixel.Current),
		frameRate.Current,
		1000,
		recordPath,
	)
	if err := camera.StartRecord(RecordOptions{
		Path:        recordPath,
		PixelType:   pixel.Current,
		Width:       uint32(width.Current),
		Height:      uint32(height.Current),
		FrameRate:   frameRate.Current,
		BitRateKbps: 1000,
		Format:      RecordFormatAVI,
	}); err != nil {
		if IsSDKErrorCode(err, ErrSupport) || IsSDKErrorCode(err, ErrParameter) {
			t.Skipf("record skipped: %v", err)
		}
		t.Fatal(err)
	}

	if err := camera.StartGrabbing(); err != nil {
		_ = camera.StopRecord()
		t.Fatal(err)
	}

	for i := 0; i < 5; i++ {
		frame, err := camera.GetFrame(2 * time.Second)
		if err != nil {
			_ = camera.StopGrabbing()
			_ = camera.StopRecord()
			t.Fatal(err)
		}
		if err := camera.InputRecordFrame(frame); err != nil {
			_ = camera.StopGrabbing()
			_ = camera.StopRecord()
			t.Fatal(err)
		}
	}
	if err := camera.StopGrabbing(); err != nil {
		_ = camera.StopRecord()
		t.Fatal(err)
	}
	if err := camera.StopRecord(); err != nil {
		t.Fatal(err)
	}
	if info, err := os.Stat(recordPath); err != nil {
		t.Fatal(err)
	} else if info.Size() == 0 {
		t.Fatal("record produced an empty file")
	}
}

func assertDeviceFileIntegration(t *testing.T, camera *Camera) {
	t.Helper()

	deviceFile := os.Getenv("MVS_TEST_READ_DEVICE_FILE")
	if deviceFile != "" {
		if _, _, err := camera.ReadDeviceFile(deviceFile, 1024*1024); err != nil {
			t.Fatal(err)
		}
	}

	writeFile := os.Getenv("MVS_TEST_WRITE_DEVICE_FILE")
	if writeFile != "" {
		if _, err := camera.WriteDeviceFile(writeFile, []byte(os.Getenv("MVS_TEST_WRITE_DEVICE_FILE_DATA"))); err != nil {
			t.Fatal(err)
		}
	}
}

func assertSerialIntegration(t *testing.T, camera *Camera) {
	t.Helper()

	if os.Getenv("MVS_TEST_SERIAL") != "1" {
		return
	}
	if err := camera.OpenSerialPort(); err != nil {
		if IsSDKErrorCode(err, ErrSupport) {
			t.Skipf("serial skipped: %v", err)
		}
		t.Fatal(err)
	}
	defer camera.CloseSerialPort()
	if err := camera.ClearSerialPort(); err != nil {
		t.Fatal(err)
	}
	payload := []byte(os.Getenv("MVS_TEST_SERIAL_WRITE"))
	if len(payload) > 0 {
		if _, err := camera.WriteSerialPort(payload); err != nil {
			t.Fatal(err)
		}
	}
}

func TestCameraCallbackIntegration(t *testing.T) {
	sdk, err := New(Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := sdk.Initialize(); err != nil {
		t.Fatal(err)
	}
	defer sdk.Finalize()

	devices, err := sdk.EnumerateDefaultDevices()
	if err != nil {
		t.Fatal(err)
	}
	if len(devices) == 0 {
		skipOrFailNoDevices(t)
	}

	device := devices[0]
	if !sdk.IsDeviceAccessible(device, AccessExclusive) {
		skipOrFailDeviceInaccessible(t, device)
	}

	camera, err := sdk.OpenDeviceBySerial(device.SerialNumber, AccessExclusive)
	if err != nil {
		t.Fatal(err)
	}
	defer camera.Close()

	if err := camera.ConfigureOptimalPacketSize(); err != nil {
		t.Fatal(err)
	}
	if err := camera.SetTriggerMode(false); err != nil {
		t.Fatal(err)
	}

	frames := make(chan *Frame, 1)
	if err := camera.RegisterFrameCallback(func(frame *Frame) {
		select {
		case frames <- frame:
		default:
		}
	}); err != nil {
		t.Fatal(err)
	}

	if err := camera.StartGrabbing(); err != nil {
		t.Fatal(err)
	}
	defer camera.StopGrabbing()

	select {
	case frame := <-frames:
		if frame.Width == 0 || frame.Height == 0 || len(frame.Data) == 0 {
			t.Fatalf("invalid callback frame: %+v bytes=%d", frame, len(frame.Data))
		}
	case <-time.After(3 * time.Second):
		t.Fatal("timed out waiting for callback frame")
	}
}

func skipOrFailNoDevices(t *testing.T) {
	t.Helper()

	if os.Getenv("MVS_TEST_REQUIRE_CAMERA") == "1" {
		t.Fatal("no devices found; MVS_TEST_REQUIRE_CAMERA=1 requires at least one connected Hikrobot camera")
	}
	t.Skip("no devices found")
}

func skipOrFailDeviceInaccessible(t *testing.T, device DeviceInfo) {
	t.Helper()

	message := "device " + device.SerialNumber + " is not accessible in exclusive mode"
	if os.Getenv("MVS_TEST_REQUIRE_CAMERA") == "1" {
		t.Fatal(message)
	}
	t.Skip(message)
}
