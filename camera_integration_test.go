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
		t.Skip("no devices found")
	}

	device := devices[0]
	if strings.TrimSpace(device.SerialNumber) == "" {
		t.Fatalf("device serial number is empty: %+v", device)
	}
	if !sdk.IsDeviceAccessible(device, AccessExclusive) {
		t.Skipf("device %s is not accessible in exclusive mode", device.SerialNumber)
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
		t.Skip("no devices found")
	}

	device := devices[0]
	if !sdk.IsDeviceAccessible(device, AccessExclusive) {
		t.Skipf("device %s is not accessible in exclusive mode", device.SerialNumber)
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
