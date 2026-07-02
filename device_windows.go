//go:build windows && amd64

package mvs

import (
	"bytes"
	"fmt"
	"unsafe"
)

func transportLayerName(layer uint32) string {
	switch layer {
	case DeviceTypeGigE:
		return "GigE"
	case DeviceTypeUSB:
		return "USB3"
	case DeviceTypeCameraLink:
		return "CameraLink"
	case DeviceTypeGentlGigE:
		return "GenTL GigE"
	case DeviceTypeGentlCameraLink:
		return "GenTL CameraLink"
	case DeviceTypeGentlCXP:
		return "GenTL CoaXPress"
	case DeviceTypeGentlXOF:
		return "GenTL XoF"
	case DeviceTypeGentlVirtual:
		return "GenTL Virtual"
	case DeviceTypeVirtualGigE:
		return "Virtual GigE"
	case DeviceTypeVirtualUSB:
		return "Virtual USB"
	default:
		return fmt.Sprintf("Unknown(0x%08X)", layer)
	}
}

func deviceInfoFromRaw(index int, raw mvCCDeviceInfo) DeviceInfo {
	info := DeviceInfo{
		Index:              index,
		TransportLayer:     raw.TLayerType,
		TransportLayerName: transportLayerName(raw.TLayerType),
		raw:                raw,
	}

	switch raw.TLayerType {
	case DeviceTypeGigE, DeviceTypeGentlGigE:
		spec := (*mvGigeDeviceInfo)(unsafe.Pointer(&raw.SpecialInfo[0]))
		info.ModelName = byteString(spec.ModelName[:])
		info.SerialNumber = byteString(spec.SerialNumber[:])
		info.UserDefinedName = byteString(spec.UserDefinedName[:])
		info.ManufacturerName = byteString(spec.ManufacturerName[:])
		info.CurrentIP = ipv4String(spec.CurrentIP)
	case DeviceTypeUSB:
		spec := (*mvUSB3DeviceInfo)(unsafe.Pointer(&raw.SpecialInfo[0]))
		info.ModelName = byteString(spec.ModelName[:])
		info.SerialNumber = byteString(spec.SerialNumber[:])
		info.UserDefinedName = byteString(spec.UserDefinedName[:])
		info.ManufacturerName = byteString(spec.ManufacturerName[:])
	case DeviceTypeCameraLink:
		spec := (*mvCamLDevInfo)(unsafe.Pointer(&raw.SpecialInfo[0]))
		info.ModelName = byteString(spec.ModelName[:])
		info.SerialNumber = byteString(spec.SerialNumber[:])
		info.ManufacturerName = byteString(spec.ManufacturerName[:])
	case DeviceTypeGentlCameraLink:
		spec := (*mvCmlDeviceInfo)(unsafe.Pointer(&raw.SpecialInfo[0]))
		info.InterfaceID = byteString(spec.InterfaceID[:])
		info.ModelName = byteString(spec.ModelName[:])
		info.SerialNumber = byteString(spec.SerialNumber[:])
		info.UserDefinedName = byteString(spec.UserDefinedName[:])
		info.DeviceID = byteString(spec.DeviceID[:])
		info.ManufacturerName = byteString(spec.VendorName[:])
	case DeviceTypeGentlCXP:
		spec := (*mvCxpDeviceInfo)(unsafe.Pointer(&raw.SpecialInfo[0]))
		info.InterfaceID = byteString(spec.InterfaceID[:])
		info.ModelName = byteString(spec.ModelName[:])
		info.SerialNumber = byteString(spec.SerialNumber[:])
		info.UserDefinedName = byteString(spec.UserDefinedName[:])
		info.DeviceID = byteString(spec.DeviceID[:])
		info.ManufacturerName = byteString(spec.VendorName[:])
	case DeviceTypeGentlXOF:
		spec := (*mvXofDeviceInfo)(unsafe.Pointer(&raw.SpecialInfo[0]))
		info.InterfaceID = byteString(spec.InterfaceID[:])
		info.ModelName = byteString(spec.ModelName[:])
		info.SerialNumber = byteString(spec.SerialNumber[:])
		info.UserDefinedName = byteString(spec.UserDefinedName[:])
		info.DeviceID = byteString(spec.DeviceID[:])
		info.ManufacturerName = byteString(spec.VendorName[:])
	case DeviceTypeGentlVirtual:
		spec := (*mvGentlVirDeviceInfo)(unsafe.Pointer(&raw.SpecialInfo[0]))
		info.InterfaceID = byteString(spec.InterfaceID[:])
		info.ModelName = byteString(spec.ModelName[:])
		info.SerialNumber = byteString(spec.SerialNumber[:])
		info.UserDefinedName = byteString(spec.UserDefinedName[:])
		info.DeviceID = byteString(spec.DeviceID[:])
		info.ManufacturerName = byteString(spec.VendorName[:])
	}

	return info
}

func byteString(buf []byte) string {
	buf = bytes.TrimRight(buf, "\x00")
	return string(buf)
}

func ipv4String(ip uint32) string {
	return fmt.Sprintf("%d.%d.%d.%d", byte(ip>>24), byte(ip>>16), byte(ip>>8), byte(ip))
}
