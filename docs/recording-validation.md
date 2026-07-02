# Recording Validation

This note records the local validation of `MV_CC_StartRecord` against Hikrobot's official MVS samples.

## Environment

- MVS SDK: `4.8.0.3`
- Runtime DLL: `C:\Program Files (x86)\Common Files\MVS\Runtime\Win64_x64\MvCameraControl.dll`
- Camera: `MV-CS200-10GM`
- Camera pixel format during validation: `Mono8` (`0x01080001`)
- Original image size: `5472 x 3648`

## Official Sample Behavior

The official C++/C#/Python recording samples follow the same low-level pattern:

- Set `TriggerMode` to off.
- Try to set `ImageCompressionMode` to off when the node is available.
- Read `Width`, `Height`, `PixelFormat`, and `ResultingFrameRate` from the camera.
- Set `nBitRate` to `1000`.
- Set `enRecordFmtType` to `MV_FormatType_AVI`.
- Call `MV_CC_StartRecord`.
- Feed raw SDK frame buffers to `MV_CC_InputOneFrame`.

The Go integration path now mirrors this sequence.

## Local Results

`ImageCompressionMode` returned `MV_E_GC_ACCESS` (`0x80000106`), so the node is not writable/accessible on this camera state.

Using Hikrobot's official Python ctypes wrapper directly, `MV_CC_StartRecord` returned `MV_E_PARAMETER` (`0x80000004`) with:

- Full frame: `5472 x 3648`
- Reduced ROI: `4096 x 2160`, `3840 x 2160`, `2560 x 1440`, `1920 x 1080`, `1280 x 720`, `640 x 480`, `320 x 240`, `96 x 96`
- Pixel formats passed to the recorder: `Mono8`, `Mono10`, `Mono12`, `Mono16`, `RGB8Packed`, `BGR8Packed`, `YUV422Packed`, `YUV422YUYVPacked`
- Frame rates: `1/16`, `1`, `4.6637`, `5`, `10`, `25`, `30`, `100`
- Bitrates: `128`, `512`, `1000`, `4096`, `8192`, `16384`

## Conclusion

The failure is not caused by the Go struct layout or syscall wrapper, because the official Python wrapper returns the same SDK error under the same conditions.

No hidden resolution, frame-rate, bitrate, or common pixel-format constraint was confirmed on this machine. The remaining likely causes are SDK/runtime recorder capability for this installation, camera-specific recorder support, or an undocumented recorder dependency/condition outside the public header comments.

Keep the Go recording API as an alpha binding and require positive validation on another camera or a Hikrobot-supported recording environment before marking it production-ready.
