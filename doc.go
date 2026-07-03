// Package mvs provides a Go wrapper around Hikrobot MVS.
//
// Its stable-core target is camera management and acquisition on Windows amd64:
// SDK lifecycle management, device enumeration, camera opening, common GenICam
// node access, pull acquisition, caller-buffer acquisition, callback
// acquisition, stream tuning, SDK pixel conversion, SDK image export, and
// pixel metadata helpers.
//
// Optional modules such as SDK recording, events, camera file access, serial
// helpers, interface/frame-grabber helpers, GenTL CTI helpers, advanced image
// processing, and 3D metadata are available for integration work but should be
// treated as experimental until validated on suitable hardware.
package mvs
