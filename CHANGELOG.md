# Changelog

## v0.1.0-beta.1 - 2026-07-03

### Changed

- Split release documentation into a stable core target and experimental modules.
- Added `docs/roadmap.md` to define the `v0.1.0-beta.1`, `v0.1.0`, later pre-v1, and v1 release gates.
- Added `MVS_TEST_REQUIRE_CAMERA=1` so stable core release candidates fail integration checks instead of skipping when no camera is found or the selected camera is not accessible in exclusive mode.

### Validation Notes

- Stable core hardware validation passed with MVS SDK `4.8.0.3`, runtime DLL `C:\Program Files (x86)\Common Files\MVS\Runtime\Win64_x64\MvCameraControl.dll`, and camera `MV-CS200-10GM` serial `DB0612579`.
- Release CLI validation saved a non-empty raw frame and SDK-exported JPEG from `5472 x 3648` `Mono8` input.
- Experimental modules remain outside the stable promise until their module-specific hardware checks pass on suitable equipment.

## v0.1.0-alpha.2 - 2026-07-02

### Added

- Interface/frame-grabber enumeration and open helpers using the `MV_CC_EnumInterfaces`, `MV_CC_CreateInterface*`, `MV_CC_OpenInterface`, and `MV_CC_EnumDevicesByInterface` APIs available from `MvCameraControl.h`.
- GenTL CTI enumeration and open helpers using `MV_CC_EnumInterfacesByGenTL`, `MV_CC_EnumDevicesByGenTL`, and `MV_CC_CreateHandleByGenTL`.
- Event callback registration and event notification helpers.
- Camera file access read/write helpers for file path and in-memory buffer modes.
- CameraLink serial-port enumeration, CameraLink baudrate helpers, and camera serial-port open/read/write/clear/close helpers.
- SDK recording wrappers with Go-side validation for path, size, frame rate, bitrate, and format.
- SDK image processing helpers for rotate, flip, Bayer interpolation/Gamma/CCM, contrast, purple-fringing, ISP process, high-bandwidth decode, and image reconstruction.
- Frame MultiPart/SubImage extraction into `Frame.Parts`, including 3D image part metadata.
- Unit tests for record option validation and MultiPart/SubImage parsing.
- Integration coverage for interface enumeration, event callback registration, rotate/flip, contrast, reconstruction, and gated record/file/serial/GenTL checks.

### Validation Notes

- Default hardware integration passed with MVS SDK `4.8.0.3` and camera `MV-CS200-10GM`.
- No frame-grabber/interface entries were present on the local validation machine.
- Recording is bound but did not pass positive hardware validation on the local camera: `MV_CC_StartRecord` returned `MV_E_PARAMETER`. The same result was reproduced with Hikrobot's official Python ctypes wrapper for full frame, reduced ROI, multiple pixel types, and frame-rate/bitrate combinations.
- File access, serial write/read, named event notification, GenTL CTI loading, ISP config processing, and high-bandwidth decode remain opt-in validation paths because they require device-specific configuration or external inputs.

## v0.1.0-alpha.1 - 2026-07-02

Initial alpha release for Hikrobot MVS acquisition on Windows amd64.

### Added

- MVS runtime loading from `MVS_SDK_DLL`, the default runtime install path, or `PATH`.
- SDK lifecycle APIs: initialize, finalize, version, DLL path.
- Device enumeration for GigE, USB3, CameraLink, and GenTL transport layers.
- Device accessibility checks before opening a camera.
- Camera lifecycle APIs: open, close, connection check, start/stop grabbing.
- Frame acquisition through `MV_CC_GetImageBuffer` / `MV_CC_FreeImageBuffer`.
- Caller-provided buffer acquisition through `MV_CC_GetOneFrameTimeout`.
- Callback acquisition through `MV_CC_RegisterImageCallBackEx`.
- Common node APIs: int, float, enum, bool, string, command.
- Enum symbolic lookup and enum-by-string writes.
- Device open helpers by serial number, current IP address, and user-defined name.
- Trigger helpers and GigE optimal packet size configuration.
- Stream buffer helpers for image node count, grab strategy, output queue size, and image buffer clearing.
- Feature file save/load through `MV_CC_FeatureSave` and `MV_CC_FeatureLoad`.
- SDK image export through `MV_CC_SaveImageToFileEx`.
- SDK pixel conversion through `MV_CC_ConvertPixelTypeEx`.
- Expanded MVS pixel type constants and image conversion helpers for Mono8, RGB8, BGR8, RGBA8, BGRA8, and Bayer8 preview conversion.
- CLI examples: `mvs-list-devices` and `mvs-grab-frame`.
- Hardware integration test behind the `integration` build tag.

### Known Limitations

- Windows amd64 only.
- Callback acquisition and active pull APIs are mutually exclusive on the same camera handle.
- Interface/frame-grabber management, recording, event callbacks, camera file access, serial port APIs, advanced ISP tuning, rotation/flip, and point-cloud helpers were outside the initial alpha scope.
- Pure Go image conversion is intended for preview only; SDK image export/conversion should be used for production export paths.
- The public API is still allowed to change before v1.0.0.
