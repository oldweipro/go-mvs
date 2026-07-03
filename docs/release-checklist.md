# Release Checklist

This checklist is for pre-v1 releases. It separates the stable core release gate from experimental module validation.

## Version Decision

- Tag an alpha release when the API or validation boundary is still moving.
- Tag `v0.1.0-beta.1` when the common camera-management and acquisition API is ready for application integration.
- Tag `v0.1.0-rc.1` when the stable-core API has been reviewed, release documentation is complete, and connected-camera validation passes.
- Tag `v0.1.0` only after the stable core hardware checks pass on a connected Hikrobot camera without "no devices found" skips.
- Keep recording, events, file access, serial, FrameGrabber/GenTL, advanced ISP, and point-cloud helpers outside the stable promise until their optional checks pass on suitable hardware.

## Required Local Checks

```powershell
git status --short
go test ./...
go vet ./...
```

Release candidates also require:

- `docs/api-review.md` is current.
- `docs/faq.md` is current.
- README, roadmap, changelog, and this checklist agree on the release scope.

## Required Hardware Checks

Use a machine with Hikrobot MVS Runtime Package installed and at least one connected camera.

```powershell
go run ./cmd/mvs-list-devices
$env:MVS_TEST_REQUIRE_CAMERA = "1"
go test -tags integration ./... -count=1 -v
Remove-Item Env:\MVS_TEST_REQUIRE_CAMERA
go run ./cmd/mvs-grab-frame -serial DB0612579 -out frame.release-check.raw -jpeg frame.release-check.jpg
```

Expected result:

- SDK version and DLL path are printed by `mvs-list-devices`.
- At least one camera is listed.
- Integration tests pass for enumeration, open, node access, feature save, pull acquisition, SDK conversion, SDK image export, callback acquisition, event callback registration, safe image-processing helpers, and image reconstruction where supported by the SDK/device.
- Stable core release candidates must not skip hardware integration because no camera was found or the selected camera is not accessible in exclusive mode.
- `frame.release-check.raw` and `frame.release-check.jpg` are created and non-empty.

Remove release-check output files after validation:

```powershell
Remove-Item frame.release-check.raw, frame.release-check.jpg -ErrorAction SilentlyContinue
```

## Optional Hardware Checks

Optional checks are useful for roadmap modules, but they do not promote those modules to stable by themselves. Record the SDK version, DLL path, camera model, serial number, node settings, and result when any optional check is used as release evidence.

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

Do not describe any pre-v1 release as a complete MVS SDK binding.

Stable core boundary:

- SDK lifecycle, DLL resolution, SDK version reporting
- Device enumeration, accessibility checks, and open helpers
- Camera open/close, connection check, start/stop grabbing
- Pull, caller-buffer, and callback acquisition
- Common GenICam node access and enum symbolic helpers
- Trigger mode, GigE packet size, and common stream-buffer helpers
- Feature file save and opt-in feature file load
- SDK image export and SDK pixel conversion
- Pixel metadata and preview conversion helpers

Current boundaries:

- FrameGrabber support is limited to `MV_CC_*` interface and GenTL APIs exposed by the installed headers; no `MV_FG_*` API is present in this SDK installation.
- Recording wrappers are included, but positive hardware validation is pending if `MV_CC_StartRecord` returns `MV_E_PARAMETER` on the connected camera. Record official-wrapper comparison results in `docs/recording-validation.md`.
- Camera file access, serial, named events, CTI loading, ISP config, high-bandwidth decode, and destructive writes require opt-in validation with known-safe device settings.
- Point-cloud support is limited to MVS 3D pixel constants plus MultiPart/SubImage metadata extraction; higher-level PCD/PLY export remains outside this scope.
- Linux or non-amd64 support

## Tagging

```powershell
git status --short
$version = "v0.1.0-rc.1"
git tag $version
git push origin $version
```
