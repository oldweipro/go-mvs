//go:build windows && amd64

package mvs

import (
	"fmt"
	"sync"
)

type Config struct {
	DLLPath string
}

type Version struct {
	Major    uint8
	Minor    uint8
	Revision uint8
	Build    uint8
}

func (v Version) String() string {
	return fmt.Sprintf("%d.%d.%d.%d", v.Major, v.Minor, v.Revision, v.Build)
}

type DeviceInfo struct {
	Index              int
	TransportLayer     uint32
	TransportLayerName string
	ModelName          string
	SerialNumber       string
	UserDefinedName    string
	ManufacturerName   string
	CurrentIP          string
	InterfaceID        string
	DeviceID           string

	raw mvCCDeviceInfo
}

type Frame struct {
	Width           uint32
	Height          uint32
	PixelType       uint32
	FrameNum        uint32
	DataLength      uint64
	DeviceTimestamp uint64
	HostTimestamp   int64
	ExposureTime    float32
	Gain            float32
	LostPacketCount uint32
	ExtraType       uint32
	SubImageNum     uint32
	Parts           []FramePart
	Data            []byte
}

type FloatValue struct {
	Current float32
	Max     float32
	Min     float32
}

type IntValue struct {
	Current   int64
	Max       int64
	Min       int64
	Increment int64
}

type EnumValue struct {
	Current   uint32
	Supported []uint32
}

type StringValue struct {
	Current   string
	MaxLength int64
}

type SDK struct {
	driver      *driver
	initialized bool
	mu          sync.Mutex
}

type Camera struct {
	sdk            *SDK
	handle         uintptr
	info           DeviceInfo
	open           bool
	grabbing       bool
	recording      bool
	serialOpen     bool
	callbackPtr    uintptr
	callback       FrameCallback
	eventCallbacks []EventCallback
	eventPtrs      []uintptr
	mu             sync.Mutex
}

type Interface struct {
	sdk    *SDK
	handle uintptr
	info   InterfaceInfo
	open   bool
	mu     sync.Mutex
}
