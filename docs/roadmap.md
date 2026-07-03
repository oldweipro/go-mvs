# Roadmap

This roadmap keeps the release promise focused on the common camera-management path first. Optional MVS modules should graduate only after positive validation on suitable hardware.

## Stability Labels

### Stable Core Target

The stable core is the minimum API surface required by a production vision application to discover a Hikrobot camera, configure common nodes, acquire frames, and export images.

- SDK lifecycle, DLL resolution, SDK version reporting
- Device enumeration, accessibility checks, and open helpers
- Camera open, close, connection check, start/stop grabbing
- Pull acquisition through `MV_CC_GetImageBuffer`
- Caller-provided buffer acquisition through `MV_CC_GetOneFrameTimeout`
- Image callback acquisition through `MV_CC_RegisterImageCallBackEx`
- Common GenICam node access: int, float, enum, bool, string, command
- Enum symbolic lookup and enum-by-string writes
- Trigger mode helpers
- GigE optimal packet size and common stream-buffer helpers
- Feature file save and opt-in feature file load
- SDK image export and SDK pixel conversion
- Pixel metadata and preview conversion helpers

### Experimental

These modules are wrapped, but they are not part of the stable release promise yet.

- SDK recording
- Event notification by name
- Camera file access
- CameraLink local serial-port helpers and camera serial-port read/write
- Interface/frame-grabber and GenTL CTI helpers exposed through `MV_CC_*`
- Rotate, flip, Bayer interpolation/Gamma/CCM, contrast, purple-fringing, ISP process
- High-bandwidth decode and image reconstruction
- MultiPart/SubImage metadata, including 3D image part metadata

### Future

- Higher-level recording fallback when Hikrobot SDK recording is unavailable
- PCD/PLY or other point-cloud export helpers
- Broader hardware validation matrix
- Linux or non-amd64 support, if a supported deployment target requires it
- v1 API compatibility contract

## Release Targets

### v0.1.0-beta.1

Goal: make the common camera-management and acquisition API ready for application integration.

Required before tagging:

- `go test ./...` passes.
- `go vet ./...` passes.
- Hardware integration tests pass with `MVS_TEST_REQUIRE_CAMERA=1` and a connected Hikrobot camera.
- `mvs-list-devices` lists at least one camera.
- `mvs-grab-frame` saves a non-empty raw frame and a non-empty SDK-exported image.
- README, release checklist, and validation notes accurately separate stable core and experimental modules.

Non-goals:

- Do not claim complete Hikrobot MVS SDK coverage.
- Do not claim SDK recording as production-ready until `MV_CC_StartRecord` has positive validation.
- Do not promise FrameGrabber, serial, file access, event, ISP, or point-cloud behavior as stable.

### v0.1.0

Goal: publish the first stable pre-v1 release for the common camera-management and acquisition path.

Required before tagging:

- All `v0.1.0-beta.1` gates pass on the release machine.
- Public API names for the stable core are reviewed and intentionally accepted.
- Integration behavior is documented for the validated SDK version, DLL path, camera model, and camera serial number.
- Known limitations are documented without treating experimental modules as stable.

### v0.2.x

Goal: graduate device-dependent modules one by one after positive validation.

Candidates:

- Recording or a supported recording fallback
- Event notification by name
- Camera file access
- CameraLink and camera serial helpers

### v0.3.x

Goal: expand beyond common acquisition.

Candidates:

- Interface/frame-grabber and GenTL workflows
- Advanced ISP helpers
- High-bandwidth decode and image reconstruction
- Higher-level point-cloud utilities

### v1.0.0

Goal: freeze the public API compatibility contract.

Required before tagging:

- Stable package-level API policy.
- Backward-compatibility tests for stable APIs.
- Multi-device validation record.
- Clear support matrix for OS, architecture, SDK version, transport type, and optional modules.
