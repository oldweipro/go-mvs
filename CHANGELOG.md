# Changelog

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
- Interface/frame-grabber management, recording, event callbacks, camera file access, serial port APIs, advanced ISP tuning, rotation/flip, and point-cloud helpers are not wrapped yet.
- Pure Go image conversion is intended for preview only; SDK image export/conversion should be used for production export paths.
- The public API is still allowed to change before v1.0.0.
