# Stable Core API Review

Review target: `v0.1.0-rc.1`

This review covers the API surface that is intended to remain stable through the `v0.1.x` common camera-management release line. It does not create a v1 compatibility promise.

## Decision

No stable-core public API rename is required before `v0.1.0-rc.1`.

The current naming is acceptable for application integration because it follows the Hikrobot SDK concepts closely while keeping the main Go entry points direct:

- `New`, `SDK.Initialize`, `SDK.Finalize`, and `SDK.Close` for lifecycle management
- `EnumerateDefaultDevices`, `EnumerateDevices`, `FindDeviceBySerial`, `FindDeviceByIP`, and `FindDeviceByUserDefinedName` for discovery
- `OpenDevice`, `OpenDeviceBySerial`, `OpenDeviceByIP`, and `OpenDeviceByUserDefinedName` for connection
- `Camera.StartGrabbing`, `Camera.StopGrabbing`, `Camera.GetFrame`, `Camera.GetOneFrameTimeout`, and `Camera.GetOneFrameInto` for acquisition
- `Camera.RegisterFrameCallback` for callback acquisition
- `GetInt`, `SetInt`, `GetFloat`, `SetFloat`, `GetEnum`, `SetEnum`, `SetEnumByString`, `GetBool`, `SetBool`, `GetString`, `SetString`, and `Command` for GenICam node access
- `SaveFrameToFile`, `ConvertFrame`, `ConvertPixelType`, `ImageFromBuffer`, and `Frame.Image` for image export and conversion

## Stable Core Surface

### SDK lifecycle

- `Config`
- `New`
- `SDK.Initialize`
- `SDK.Finalize`
- `SDK.Close`
- `SDK.Version`
- `SDK.DLLPath`

### Device discovery and opening

- `DeviceInfo`
- `SDK.EnumerateDefaultDevices`
- `SDK.EnumerateDevices`
- `SDK.IsDeviceAccessible`
- `SDK.FindDeviceBySerial`
- `SDK.FindDeviceByIP`
- `SDK.FindDeviceByUserDefinedName`
- `SDK.OpenDevice`
- `SDK.OpenDeviceBySerial`
- `SDK.OpenDeviceByIP`
- `SDK.OpenDeviceByUserDefinedName`

### Camera lifecycle and stream control

- `Camera.Info`
- `Camera.IsConnected`
- `Camera.Close`
- `Camera.StartGrabbing`
- `Camera.StopGrabbing`
- `Camera.ConfigureOptimalPacketSize`
- `Camera.SetImageNodeNum`
- `Camera.SetGrabStrategy`
- `Camera.SetOutputQueueSize`
- `Camera.ClearImageBuffer`

### Acquisition

- `Frame`
- `FrameCallback`
- `Camera.GetFrame`
- `Camera.GetOneFrameTimeout`
- `Camera.GetOneFrameInto`
- `Camera.CurrentFrameBufferSize`
- `Camera.RegisterFrameCallback`

`GetOneFrameTimeout` is kept for direct SDK parity. `GetOneFrameInto` is kept for callers that need buffer reuse.

### Common GenICam access

- `FloatValue`
- `IntValue`
- `EnumValue`
- `EnumEntry`
- `StringValue`
- `Camera.GetFloat`
- `Camera.SetFloat`
- `Camera.GetInt`
- `Camera.SetInt`
- `Camera.GetEnum`
- `Camera.SetEnum`
- `Camera.SetEnumByString`
- `Camera.GetEnumEntrySymbolic`
- `Camera.GetEnumEntries`
- `Camera.GetBool`
- `Camera.SetBool`
- `Camera.GetString`
- `Camera.SetString`
- `Camera.Command`
- `Camera.SetTriggerMode`
- `Camera.TriggerOnce`

GenICam node names are intentionally passed as strings so the package can work with camera-specific features without a release for every node.

### Image export and conversion

- `ImageSaveOptions`
- `PixelConvertOptions`
- `ImageType`
- `ImageTypeFromExtension`
- `Camera.SaveFrameToFile`
- `Camera.ConvertFrame`
- `Camera.ConvertPixelType`
- `Frame.Image`
- `Frame.PixelTypeName`
- `PixelTypeName`
- `PixelBitCount`
- `ExpectedFrameDataLength`
- `ImageFromBuffer`
- `IsMonoPixelType`
- `IsColorPixelType`
- `IsBayerPixelType`
- `IsCompressedPixelType`

Pure Go image conversion remains a preview helper. Production export should prefer the Hikrobot SDK path through `SaveFrameToFile` or `ConvertFrame`.

### Feature files

- `Camera.FeatureSave`
- `Camera.FeatureLoad`

`FeatureLoad` writes camera parameters and remains opt-in in integration tests.

### Errors

- Sentinel errors such as `ErrSDKNotInitialized`, `ErrCameraClosed`, `ErrCameraNotGrabbing`, `ErrDeviceNotFound`, `ErrInvalidArgument`, and `ErrAcquisitionModeConflict`
- `SDKError`
- `IsSDKErrorCode`

SDK numeric error codes are preserved so application code can handle Hikrobot-specific failures such as `ErrAccessDenied`, `ErrGCAccess`, and `ErrParameter`.

## Experimental Surface

The following APIs are available but not part of the `v0.1.x` stable-core promise:

- Interface/frame-grabber and GenTL helpers
- Event callbacks and event notification helpers
- Camera file access helpers
- CameraLink and camera serial helpers
- SDK recording helpers
- Rotate, flip, Bayer, Gamma, CCM, contrast, purple-fringing, ISP helpers
- High-bandwidth decode and image reconstruction
- MultiPart/SubImage and 3D metadata helpers

These modules can graduate in later minor releases after positive validation on suitable hardware.

## Compatibility Policy Before v1

For `v0.1.x`, avoid breaking the stable core unless a change is required to fix data corruption, unsafe lifecycle behavior, or a severe mismatch with the Hikrobot SDK.

Experimental APIs may still change before v1 when validation exposes better naming, safer arguments, or missing constraints.
