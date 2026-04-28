# MVS Go SDK

`mvs-go-sdk` is a Windows-first Go wrapper around Hikrobot MVS `MvCameraControl.dll`.

This scaffold focuses on the minimum acquisition chain:

1. Initialize SDK
2. Enumerate devices
3. Create device handle
4. Open device
5. Start grabbing
6. Pull frames in a loop
7. Stop grabbing
8. Close device and destroy handle

## Current scope

- Windows `amd64` only
- Direct DLL binding through Go's Windows syscall layer
- No cgo dependency
- No UI binding
- No callback mode yet

## Project layout

- `raw_windows.go`: DLL loading, raw structs, low-level procedure calls
- `sdk_windows.go`: SDK lifecycle and device enumeration
- `camera_windows.go`: camera lifecycle, frame acquisition, common node APIs

## DLL lookup

The SDK resolves the runtime DLL in this order:

1. `MVS_SDK_DLL` environment variable
2. `C:\Program Files (x86)\Common Files\MVS\Runtime\Win64_x64\MvCameraControl.dll`
3. `C:\Program Files\Common Files\MVS\Runtime\Win64_x64\MvCameraControl.dll`
4. `MvCameraControl.dll` from `PATH`

## Example

```go
package main

import (
	"fmt"
	"log"
	"time"

	mvsdk "github.com/oldweipro/mvs-go-sdk"
)

func main() {
	sdk, err := mvsdk.New(mvsdk.Config{})
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

	cam, err := sdk.OpenDevice(devices[0], mvsdk.AccessExclusive)
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

	for i := 0; i < 10; i++ {
		frame, err := cam.GetFrame(time.Second)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("frame=%d size=%dx%d bytes=%d\n", frame.FrameNum, frame.Width, frame.Height, len(frame.Data))
	}
}
```

## Notes

- The module path is intentionally local-friendly: `github.com/oldweipro/mvs-go-sdk`. If you move this into a dedicated repository, update `go.mod` to your final import path.
- This environment does not currently expose a Go toolchain in `PATH`, so the code was scaffolded against the installed MVS headers and Python ctypes definitions, but could not be compiled here yet.
