# go-mvs

`go-mvs` is a Go wrapper around Hikrobot MVS `MvCameraControl.dll`.

Current tagged release: `v0.1.0-beta.1`. The package is suitable for application integration and hardware validation on Windows amd64. This beta stabilizes the common camera-management and acquisition path before a formal `v0.1.0` release.

This package is not a complete Hikrobot MVS SDK binding and should not be treated as a stable v1 API yet.

## Stable Core Target

- Windows amd64
- No cgo dependency
- Runtime DLL binding through Go's Windows syscall layer
- Device enumeration, accessibility checks, and open helpers
- Frame acquisition through SDK buffer pull mode, caller-provided buffer mode, and image callback mode
- Common GenICam node access: int, float, enum, bool, string, command
- Enum symbolic lookup and enum-by-string writes
- GigE optimal packet size configuration, stream buffer node count, grab strategy, output queue size, and image buffer clearing
- Feature file save and opt-in feature file load through Hikrobot SDK
- SDK image export through `MV_CC_SaveImageToFileEx`
- SDK pixel conversion through `MV_CC_ConvertPixelTypeEx`
- Pixel helpers for Mono8, RGB8, BGR8, RGBA8, BGRA8, Bayer8 preview conversion, and MVS pixel type metadata

## Experimental Modules

- Interface/frame-grabber enumeration and GenTL CTI enumeration/open helpers through the `MV_CC_*` interface APIs exposed by `MvCameraControl.h`
- Camera event callback registration and event notification toggles
- Camera file access through SDK file read/write APIs
- CameraLink local serial-port enumeration, CameraLink baudrate helpers, and camera serial-port read/write APIs
- SDK recording bindings with Go-side parameter validation
- SDK rotate, flip, Bayer interpolation/Gamma/CCM, contrast, purple-fringing, ISP config/process, high-bandwidth decode, and image reconstruction bindings
- MultiPart/SubImage frame metadata extraction, including 3D image part metadata

Experimental modules are available for integration work, but they remain outside the stable release promise until they have positive validation on suitable hardware and documented failure behavior.

## Requirements

Install Hikrobot Machine Vision Industrial Camera SDK runtime on the host. The package has been validated locally with:

- MVS SDK: `4.8.0.3`
- Runtime DLL: `C:\Program Files (x86)\Common Files\MVS\Runtime\Win64_x64\MvCameraControl.dll`
- Camera: `MV-CS200-10GM`

## DLL Lookup

The runtime DLL is resolved in this order:

1. `MVS_SDK_DLL` environment variable
2. `C:\Program Files (x86)\Common Files\MVS\Runtime\Win64_x64\MvCameraControl.dll`
3. `C:\Program Files\Common Files\MVS\Runtime\Win64_x64\MvCameraControl.dll`
4. `MvCameraControl.dll` from `PATH`

## Install

```powershell
go get github.com/oldweipro/go-mvs@v0.1.0-beta.1
```

For local development:

```powershell
go mod edit -replace github.com/oldweipro/go-mvs=../go-mvs
```

## Example

```go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/oldweipro/go-mvs"
)

func main() {
	sdk, err := mvs.New(mvs.Config{})
	if err != nil {
		log.Fatal(err)
	}
	if err := sdk.Initialize(); err != nil {
		log.Fatal(err)
	}
	defer sdk.Finalize()

	devices, err := sdk.EnumerateDefaultDevices()
	if err != nil {
		log.Fatal(err)
	}
	if len(devices) == 0 {
		log.Fatal("no devices found")
	}

	cam, err := sdk.OpenDeviceBySerial(devices[0].SerialNumber, mvs.AccessExclusive)
	if err != nil {
		log.Fatal(err)
	}
	defer cam.Close()

	if err := cam.ConfigureOptimalPacketSize(); err != nil {
		log.Printf("packet size setup skipped: %v", err)
	}
	if err := cam.SetTriggerMode(false); err != nil {
		log.Fatal(err)
	}
	if err := cam.StartGrabbing(); err != nil {
		log.Fatal(err)
	}
	defer cam.StopGrabbing()

	frame, err := cam.GetFrame(time.Second)
	if err != nil {
		log.Fatal(err)
	}
	if err := cam.SaveFrameToFile(frame, "frame.jpg", mvs.ImageSaveOptions{
		Type:    mvs.ImageTypeJPEG,
		Quality: 90,
	}); err != nil {
		log.Fatal(err)
	}

	fmt.Printf(
		"frame=%d size=%dx%d pixel=%s bytes=%d\n",
		frame.FrameNum,
		frame.Width,
		frame.Height,
		frame.PixelTypeName(),
		len(frame.Data),
	)
}
```

## Commands

List connected cameras:

```powershell
go run ./cmd/mvs-list-devices
```

Grab one raw frame:

```powershell
go run ./cmd/mvs-grab-frame -out frame.raw
```

Grab one raw frame and write a JPEG image through the Hikrobot SDK:

```powershell
go run ./cmd/mvs-grab-frame -serial DB0612579 -out frame.raw -jpeg frame.jpg
```

Grab one raw frame and write another image format:

```powershell
go run ./cmd/mvs-grab-frame -serial DB0612579 -out frame.raw -image frame.png
```

## Verification

Regular checks:

```powershell
go test ./...
go vet ./...
```

Hardware integration check:

```powershell
go test -tags integration ./... -run TestCameraIntegration -count=1 -v
```

Release hardware integration check. This must be used before promoting the stable core because it fails when no camera is found or the selected camera is not accessible in exclusive mode:

```powershell
$env:MVS_TEST_REQUIRE_CAMERA = "1"
go test -tags integration ./... -count=1 -v
Remove-Item Env:\MVS_TEST_REQUIRE_CAMERA
```

Optional modules that may depend on device capability or write state are guarded by environment variables:

```powershell
$env:MVS_TEST_RECORD = "1"                 # records a temp AVI from a captured frame
$env:MVS_TEST_EVENT_NAME = "ExposureEnd"   # enables/disables a named event
$env:MVS_TEST_READ_DEVICE_FILE = "UserSet1"
$env:MVS_TEST_WRITE_DEVICE_FILE = "UserFile"
$env:MVS_TEST_WRITE_DEVICE_FILE_DATA = "payload"
$env:MVS_TEST_SERIAL = "1"
$env:MVS_TEST_SERIAL_WRITE = "payload"
$env:MVS_TEST_GENTL_CTI = "C:\path\to\producer.cti"
go test -tags integration ./... -run TestCameraIntegration -count=1 -v
```

Feature load is intentionally opt-in because it writes parameters back to the camera:

```powershell
$env:MVS_TEST_FEATURE_LOAD = "1"
go test -tags integration ./... -run TestCameraIntegration -count=1 -v
```

## Release Plan

See `docs/roadmap.md` for the release roadmap and `docs/release-checklist.md` for the validation gate.

Summary:

- `v0.1.0-beta.1`: freeze the common camera-management and acquisition API for integration use.
- `v0.1.0`: promote the stable core after hardware validation succeeds without skipped camera checks on a connected Hikrobot camera.
- Later minor releases: graduate recording, event, file-access, serial, FrameGrabber/GenTL, advanced ISP, and point-cloud helpers when each module has enough positive hardware validation.

## Project Layout

- `raw_windows.go`: DLL loading and low-level procedure calls
- `raw_types_windows.go`: raw Hikrobot SDK struct layouts and size checks
- `device_windows.go`: device transport names and raw device info conversion
- `interface_windows.go`: interface/frame-grabber and GenTL enumeration/open helpers
- `sdk_windows.go`: SDK lifecycle, device enumeration, and device open helpers
- `camera_windows.go`: camera lifecycle, frame acquisition, stream options, node APIs, SDK conversion, and SDK image export
- `event_windows.go`: event callback registration and event notification helpers
- `file_windows.go`: camera file access helpers
- `record_windows.go`: SDK recording helpers
- `serial_windows.go`: CameraLink and camera serial-port helpers
- `image_process_windows.go`: SDK image-processing, high-bandwidth decode, and reconstruction helpers
- `pixel.go`: pixel type utilities and preview image conversion helpers
- `cmd/mvs-list-devices`: enumerate connected cameras
- `cmd/mvs-grab-frame`: grab one frame and save raw bytes plus optional SDK image export
- `docs/recording-validation.md`: local recording validation notes against official Hikrobot samples
- `docs/roadmap.md`: stability boundaries and release plan
- `docs/release-checklist.md`: release validation checklist

## Limitations

- Windows amd64 only.
- `MV_CC_RegisterImageCallBackEx` callback mode is supported, but callback mode is mutually exclusive with active pull APIs on the same camera handle.
- FrameGrabber support is limited to the interface and GenTL functions present in this machine's `MvCameraControl.h`; no separate `MV_FG_*` API is wrapped because the installed SDK headers do not expose it.
- Recording is bound and parameter-validated, but `MV_CC_StartRecord` returned `MV_E_PARAMETER` on the local `MV-CS200-10GM` validation camera. The same result was reproduced through Hikrobot's official Python ctypes wrapper across ROI, pixel format, frame-rate, and bitrate probes; see `docs/recording-validation.md`.
- File access, event notification by name, serial-port read/write, GenTL CTI loading, ISP config files, high-bandwidth decode, and destructive device writes are opt-in integration checks because they depend on device capability or external files.
- Point-cloud support currently exposes MVS 3D pixel constants and MultiPart/SubImage metadata. Higher-level export formats such as PCD/PLY are not included yet.
- Pure Go image conversion is a preview helper only; production image export should use `SaveFrameToFile` or `ConvertFrame`.
- API compatibility is not guaranteed until v1.0.0.
