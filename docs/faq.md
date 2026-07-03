# FAQ

## `MvCameraControl.dll` cannot be found

Install the Hikrobot Machine Vision Industrial Camera SDK runtime package on the host.

DLL lookup order:

1. `MVS_SDK_DLL`
2. `C:\Program Files (x86)\Common Files\MVS\Runtime\Win64_x64\MvCameraControl.dll`
3. `C:\Program Files\Common Files\MVS\Runtime\Win64_x64\MvCameraControl.dll`
4. `MvCameraControl.dll` from `PATH`

Use `go run ./cmd/mvs-list-devices` to print the SDK version and resolved DLL path.

## No devices are found

Check the camera first with Hikrobot's official MVS client. If the official client cannot see the camera, fix the device, driver, network, power, or permissions before debugging this Go package.

For GigE cameras, verify that the camera and NIC are on a reachable subnet and that firewall or network isolation is not blocking discovery. Link-local addresses such as `169.254.x.x` can work, but the NIC must be connected to the same camera network.

For release validation, run:

```powershell
$env:MVS_TEST_REQUIRE_CAMERA = "1"
go test -tags integration ./... -count=1 -v
Remove-Item Env:\MVS_TEST_REQUIRE_CAMERA
```

This fails when no camera is found instead of skipping the hardware checks.

## `MV_E_ACCESS_DENIED`

`MV_E_ACCESS_DENIED` usually means another process already opened the camera or the requested access mode is not available.

Close Hikrobot MVS client, other vision software, and any other process using the same camera. Then retry with `AccessExclusive`.

## `MV_E_GC_ACCESS`

`MV_E_GC_ACCESS` means a GenICam node is not accessible in the current camera state. Some nodes are read-only, require a different acquisition state, or depend on the current device configuration.

For parameter writes, stop acquisition when the camera requires it and confirm the node is writable in the official MVS client.

## Frame acquisition times out

Check trigger mode first. In continuous acquisition, use:

```go
if err := cam.SetTriggerMode(false); err != nil {
	return err
}
```

For software trigger workflows, enable trigger mode, start grabbing, call `TriggerOnce`, then read a frame.

Also check exposure time, network packet settings, camera power, and cable quality.

## GigE packet size setup fails

`ConfigureOptimalPacketSize` calls Hikrobot's optimal packet-size helper and writes `GevSCPSPacketSize`. If this fails, verify NIC jumbo-frame settings and camera network configuration.

The helper returns nil for non-GigE transport layers.

## Pull acquisition and callback acquisition conflict

Do not mix active pull APIs and callback acquisition on the same camera handle.

Register `RegisterFrameCallback` before `StartGrabbing` when using callback mode. Use `GetFrame` or `GetOneFrameInto` only when no frame callback is registered.

## Pure Go preview conversion is limited

`Frame.Image` and `ImageFromBuffer` are preview helpers for common pixel formats. Production image export should prefer:

- `SaveFrameToFile`
- `ConvertFrame`
- `ConvertPixelType`

## Recording returns `MV_E_PARAMETER`

SDK recording is experimental in this package.

On the local validation camera `MV-CS200-10GM`, `MV_CC_StartRecord` returned `MV_E_PARAMETER`. The same result was reproduced with Hikrobot's official Python ctypes wrapper across common ROI, pixel format, frame-rate, and bitrate combinations.

See `docs/recording-validation.md` for details. Do not treat SDK recording as production-ready until it has positive validation on supported equipment.

## Unsupported platform

The active implementation supports Windows amd64. Other platforms compile with unsupported stubs and return `ErrUnsupportedPlatform`.
