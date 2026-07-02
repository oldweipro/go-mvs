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

func interfaceLayerName(layer uint32) string {
	switch layer {
	case InterfaceTypeGigE:
		return "GigE Vision Interface"
	case InterfaceTypeCameraLink:
		return "CameraLink Interface"
	case InterfaceTypeCXP:
		return "CoaXPress Interface"
	case InterfaceTypeXOF:
		return "XoFLink Interface"
	case InterfaceTypeVirtual:
		return "Virtual Interface"
	case InterfaceTypeLightCtrl:
		return "Light Controller Interface"
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

func interfaceInfoFromRaw(index int, raw mvInterfaceInfo) InterfaceInfo {
	return InterfaceInfo{
		Index:              index,
		TransportLayer:     raw.TLayerType,
		TransportLayerName: interfaceLayerName(raw.TLayerType),
		PCIEInfo:           raw.PCIEInfo,
		InterfaceID:        byteString(raw.InterfaceID[:]),
		DisplayName:        byteString(raw.DisplayName[:]),
		SerialNumber:       byteString(raw.SerialNumber[:]),
		ModelName:          byteString(raw.ModelName[:]),
		Manufacturer:       byteString(raw.Manufacturer[:]),
		DeviceVersion:      byteString(raw.DeviceVersion[:]),
		UserDefinedName:    byteString(raw.UserDefinedName[:]),
	}
}

func genTLInterfaceInfoFromRaw(index int, raw mvGenTLIFInfo) GenTLInterfaceInfo {
	return GenTLInterfaceInfo{
		Index:       index,
		InterfaceID: byteString(raw.InterfaceID[:]),
		TLType:      byteString(raw.TLType[:]),
		DisplayName: byteString(raw.DisplayName[:]),
		CtiIndex:    raw.CtiIndex,
	}
}

func genTLDeviceInfoFromRaw(index int, raw mvGenTLDevInfo) GenTLDeviceInfo {
	return GenTLDeviceInfo{
		Index:           index,
		InterfaceID:     byteString(raw.InterfaceID[:]),
		DeviceID:        byteString(raw.DeviceID[:]),
		VendorName:      byteString(raw.VendorName[:]),
		ModelName:       byteString(raw.ModelName[:]),
		TLType:          byteString(raw.TLType[:]),
		DisplayName:     byteString(raw.DisplayName[:]),
		UserDefinedName: byteString(raw.UserDefinedName[:]),
		SerialNumber:    byteString(raw.SerialNumber[:]),
		DeviceVersion:   byteString(raw.DeviceVersion[:]),
		CtiIndex:        raw.CtiIndex,
	}
}

func rawInterfaceInfoFromPublic(info InterfaceInfo) mvInterfaceInfo {
	raw := mvInterfaceInfo{
		TLayerType: info.TransportLayer,
		PCIEInfo:   info.PCIEInfo,
	}
	copyFixedString(raw.InterfaceID[:], info.InterfaceID)
	copyFixedString(raw.DisplayName[:], info.DisplayName)
	copyFixedString(raw.SerialNumber[:], info.SerialNumber)
	copyFixedString(raw.ModelName[:], info.ModelName)
	copyFixedString(raw.Manufacturer[:], info.Manufacturer)
	copyFixedString(raw.DeviceVersion[:], info.DeviceVersion)
	copyFixedString(raw.UserDefinedName[:], info.UserDefinedName)
	return raw
}

func rawGenTLIFInfoFromPublic(info GenTLInterfaceInfo) mvGenTLIFInfo {
	raw := mvGenTLIFInfo{CtiIndex: info.CtiIndex}
	copyFixedString(raw.InterfaceID[:], info.InterfaceID)
	copyFixedString(raw.TLType[:], info.TLType)
	copyFixedString(raw.DisplayName[:], info.DisplayName)
	return raw
}

func rawGenTLDevInfoFromPublic(info GenTLDeviceInfo) mvGenTLDevInfo {
	raw := mvGenTLDevInfo{CtiIndex: info.CtiIndex}
	copyFixedString(raw.InterfaceID[:], info.InterfaceID)
	copyFixedString(raw.DeviceID[:], info.DeviceID)
	copyFixedString(raw.VendorName[:], info.VendorName)
	copyFixedString(raw.ModelName[:], info.ModelName)
	copyFixedString(raw.TLType[:], info.TLType)
	copyFixedString(raw.DisplayName[:], info.DisplayName)
	copyFixedString(raw.UserDefinedName[:], info.UserDefinedName)
	copyFixedString(raw.SerialNumber[:], info.SerialNumber)
	copyFixedString(raw.DeviceVersion[:], info.DeviceVersion)
	return raw
}

func copyFixedString(dst []byte, value string) {
	copy(dst, []byte(value))
	if len(value) < len(dst) {
		dst[len(value)] = 0
	}
}

func byteString(buf []byte) string {
	buf = bytes.TrimRight(buf, "\x00")
	return string(buf)
}

func ipv4String(ip uint32) string {
	return fmt.Sprintf("%d.%d.%d.%d", byte(ip>>24), byte(ip>>16), byte(ip>>8), byte(ip))
}
