# Release Checklist

This checklist is for `v0.1.0-alpha.2` and later pre-v1 releases.

## Required Local Checks

```powershell
go test ./...
go vet ./...
```

## Required Hardware Checks

Use a machine with Hikrobot MVS Runtime Package installed and at least one connected camera.

```powershell
go run ./cmd/mvs-list-devices
go test -tags integration ./... -count=1 -v
go run ./cmd/mvs-grab-frame -serial DB0612579 -out frame.release-check.raw -jpeg frame.release-check.jpg
```

Expected result:

- SDK version and DLL path are printed by `mvs-list-devices`.
- At least one camera is listed.
- Integration tests pass for enumeration, open, node access, feature save, pull acquisition, SDK conversion, SDK image export, callback acquisition, event callback registration, safe image-processing helpers, and image reconstruction where supported by the SDK/device.
- `frame.release-check.raw` and `frame.release-check.jpg` are created and non-empty.

Remove release-check output files after validation:

```powershell
Remove-Item frame.release-check.raw, frame.release-check.jpg -ErrorAction SilentlyContinue
```

## Optional Hardware Checks

Feature load writes parameters back to the camera. Run it only when that is acceptable for the connected device.

```powershell
$env:MVS_TEST_FEATURE_LOAD = "1"
go test -tags integration ./... -run TestCameraIntegration -count=1 -v
Remove-Item Env:\MVS_TEST_FEATURE_LOAD
```

Module-specific checks may require device support, CTI files, serial wiring, or device-file names. Enable only the checks that are valid for the connected equipment.

```powershell
$env:MVS_TEST_RECORD = "1"
go test -tags integration ./... -run TestCameraIntegration -count=1 -v
Remove-Item Env:\MVS_TEST_RECORD
```

```powershell
$env:MVS_TEST_GENTL_CTI = "C:\path\to\producer.cti"
go test -tags integration ./... -run TestCameraIntegration -count=1 -v
Remove-Item Env:\MVS_TEST_GENTL_CTI
```

```powershell
$env:MVS_TEST_EVENT_NAME = "ExposureEnd"
go test -tags integration ./... -run TestCameraIntegration -count=1 -v
Remove-Item Env:\MVS_TEST_EVENT_NAME
```

```powershell
$env:MVS_TEST_SERIAL = "1"
$env:MVS_TEST_SERIAL_WRITE = "payload"
go test -tags integration ./... -run TestCameraIntegration -count=1 -v
Remove-Item Env:\MVS_TEST_SERIAL, Env:\MVS_TEST_SERIAL_WRITE
```

```powershell
$env:MVS_TEST_READ_DEVICE_FILE = "UserSet1"
$env:MVS_TEST_WRITE_DEVICE_FILE = "UserFile"
$env:MVS_TEST_WRITE_DEVICE_FILE_DATA = "payload"
go test -tags integration ./... -run TestCameraIntegration -count=1 -v
Remove-Item Env:\MVS_TEST_READ_DEVICE_FILE, Env:\MVS_TEST_WRITE_DEVICE_FILE, Env:\MVS_TEST_WRITE_DEVICE_FILE_DATA
```

## Release Scope Boundary

The release can be tagged as an alpha when the checks above pass and validation notes are accurate. Do not describe it as a complete MVS SDK binding.

Current boundaries:

- FrameGrabber support is limited to `MV_CC_*` interface and GenTL APIs exposed by the installed headers; no `MV_FG_*` API is present in this SDK installation.
- Recording wrappers are included, but positive hardware validation is pending if `MV_CC_StartRecord` returns `MV_E_PARAMETER` on the connected camera. Record official-wrapper comparison results in `docs/recording-validation.md`.
- Camera file access, serial, named events, CTI loading, ISP config, high-bandwidth decode, and destructive writes require opt-in validation with known-safe device settings.
- Point-cloud support is limited to MVS 3D pixel constants plus MultiPart/SubImage metadata extraction; higher-level PCD/PLY export remains outside this scope.
- Linux or non-amd64 support

## Tagging

```powershell
git status --short
git tag v0.1.0-alpha.2
git push origin v0.1.0-alpha.2
```
