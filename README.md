# go-mvs

`go-mvs` is a Go wrapper around Hikrobot MVS `MvCameraControl.dll`.

Current release status: `v0.1.0-alpha.1`. The package is suitable for early integration and internal validation on Windows amd64. It is not a complete Hikrobot MVS SDK binding and should not be treated as a stable v1 API yet.

## Scope

- Windows amd64
- No cgo dependency
- Runtime DLL binding through Go's Windows syscall layer
- Device enumeration, accessibility checks, and open helpers
- Frame acquisition through SDK buffer pull mode, caller-provided buffer mode, and image callback mode
- Common GenICam node access: int, float, enum, bool, string, command
- Enum symbolic lookup and enum-by-string writes
- GigE optimal packet size configuration, stream buffer node count, grab strategy, output queue size, and image buffer clearing
- Feature file save/load through Hikrobot SDK
- SDK image export through `MV_CC_SaveImageToFileEx`
- SDK pixel conversion through `MV_CC_ConvertPixelTypeEx`
- Pixel helpers for Mono8, RGB8, BGR8, RGBA8, BGRA8, Bayer8 preview conversion, and MVS pixel type metadata

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
go get github.com/oldweipro/go-mvs@v0.1.0-alpha.1
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

Full hardware integration check, including callback acquisition:

```powershell
go test -tags integration ./... -count=1 -v
```

Feature load is intentionally opt-in because it writes parameters back to the camera:

```powershell
$env:MVS_TEST_FEATURE_LOAD = "1"
go test -tags integration ./... -run TestCameraIntegration -count=1 -v
```

## Project Layout

- `raw_windows.go`: DLL loading and low-level procedure calls
- `raw_types_windows.go`: raw Hikrobot SDK struct layouts and size checks
- `device_windows.go`: device transport names and raw device info conversion
- `sdk_windows.go`: SDK lifecycle, device enumeration, and device open helpers
- `camera_windows.go`: camera lifecycle, frame acquisition, stream options, node APIs, SDK conversion, and SDK image export
- `pixel.go`: pixel type utilities and preview image conversion helpers
- `cmd/mvs-list-devices`: enumerate connected cameras
- `cmd/mvs-grab-frame`: grab one frame and save raw bytes plus optional SDK image export
- `docs/release-checklist.md`: release validation checklist

## Limitations

- Windows amd64 only.
- `MV_CC_RegisterImageCallBackEx` callback mode is supported, but callback mode is mutually exclusive with active pull APIs on the same camera handle.
- Interface/frame-grabber management, recording, event callbacks, camera file access, serial port APIs, advanced ISP tuning, rotation/flip, and point-cloud helpers are not wrapped yet.
- Pure Go image conversion is a preview helper only; production image export should use `SaveFrameToFile` or `ConvertFrame`.
- API compatibility is not guaranteed until v1.0.0.
