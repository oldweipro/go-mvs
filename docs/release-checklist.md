# Release Checklist

This checklist is for `v0.1.0-alpha.1` and later pre-v1 releases.

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
- Integration tests pass for enumeration, open, node access, feature save, pull acquisition, SDK conversion, SDK image export, and callback acquisition.
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

## Release Scope Boundary

The release can be tagged as an alpha when the checks above pass. Do not describe it as a complete MVS SDK binding.

Not included in this release scope:

- Interface and frame-grabber management APIs
- Recording APIs
- Event callbacks
- Camera file access APIs
- Serial port APIs
- Advanced ISP tuning beyond SDK image conversion/save defaults
- Rotate, flip, and image reconstruction helpers
- Point-cloud specialized helpers
- Linux or non-amd64 support

## Tagging

```powershell
git status --short
git tag v0.1.0-alpha.1
git push origin v0.1.0-alpha.1
```
