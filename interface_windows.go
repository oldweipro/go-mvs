//go:build windows && amd64

package mvs

import "fmt"

func (s *SDK) EnumerateInterfaces(layerType uint32) ([]InterfaceInfo, error) {
	s.mu.Lock()
	initialized := s.initialized
	s.mu.Unlock()

	if !initialized {
		return nil, ErrSDKNotInitialized
	}
	if layerType == 0 {
		layerType = DefaultInterfaceTypeMask
	}

	var rawList mvInterfaceInfoList
	if err := s.driver.enumInterfaces(layerType, &rawList); err != nil {
		return nil, err
	}

	count := int(rawList.InterfaceNum)
	if count > len(rawList.Interface) {
		count = len(rawList.Interface)
	}
	interfaces := make([]InterfaceInfo, 0, count)
	for i := 0; i < count; i++ {
		if rawList.Interface[i] == nil {
			continue
		}
		interfaces = append(interfaces, interfaceInfoFromRaw(i, *rawList.Interface[i]))
	}
	return interfaces, nil
}

func (s *SDK) OpenInterface(info InterfaceInfo) (*Interface, error) {
	s.mu.Lock()
	initialized := s.initialized
	s.mu.Unlock()

	if !initialized {
		return nil, ErrSDKNotInitialized
	}

	var (
		handle uintptr
		err    error
	)
	if info.InterfaceID != "" {
		handle, err = s.driver.createInterfaceByID(info.InterfaceID)
	} else {
		raw := rawInterfaceInfoFromPublic(info)
		handle, err = s.driver.createInterface(&raw)
	}
	if err != nil {
		return nil, err
	}
	if err := s.driver.openInterface(handle); err != nil {
		_ = s.driver.destroyInterface(handle)
		return nil, err
	}

	return &Interface{
		sdk:    s,
		handle: handle,
		info:   info,
		open:   true,
	}, nil
}

func (s *SDK) OpenInterfaceByID(interfaceID string) (*Interface, error) {
	if interfaceID == "" {
		return nil, fmt.Errorf("%w: interface id is empty", ErrInvalidArgument)
	}
	return s.OpenInterface(InterfaceInfo{InterfaceID: interfaceID})
}

func (s *SDK) EnumerateGenTLInterfaces(ctiPath string) ([]GenTLInterfaceInfo, error) {
	s.mu.Lock()
	initialized := s.initialized
	s.mu.Unlock()

	if !initialized {
		return nil, ErrSDKNotInitialized
	}
	if ctiPath == "" {
		return nil, fmt.Errorf("%w: gentl cti path is empty", ErrInvalidArgument)
	}

	var rawList mvGenTLIFInfoList
	if err := s.driver.enumInterfacesByGenTL(ctiPath, &rawList); err != nil {
		return nil, err
	}

	count := int(rawList.InterfaceNum)
	if count > len(rawList.Interface) {
		count = len(rawList.Interface)
	}
	interfaces := make([]GenTLInterfaceInfo, 0, count)
	for i := 0; i < count; i++ {
		if rawList.Interface[i] == nil {
			continue
		}
		interfaces = append(interfaces, genTLInterfaceInfoFromRaw(i, *rawList.Interface[i]))
	}
	return interfaces, nil
}

func (s *SDK) UnloadGenTLLibrary(ctiPath string) error {
	if ctiPath == "" {
		return fmt.Errorf("%w: gentl cti path is empty", ErrInvalidArgument)
	}
	return s.driver.unloadGenTLLibrary(ctiPath)
}

func (s *SDK) EnumerateGenTLDevices(info GenTLInterfaceInfo) ([]GenTLDeviceInfo, error) {
	s.mu.Lock()
	initialized := s.initialized
	s.mu.Unlock()

	if !initialized {
		return nil, ErrSDKNotInitialized
	}

	rawIF := rawGenTLIFInfoFromPublic(info)
	var rawList mvGenTLDevInfoList
	if err := s.driver.enumDevicesByGenTL(&rawIF, &rawList); err != nil {
		return nil, err
	}

	count := int(rawList.DeviceNum)
	if count > len(rawList.Device) {
		count = len(rawList.Device)
	}
	devices := make([]GenTLDeviceInfo, 0, count)
	for i := 0; i < count; i++ {
		if rawList.Device[i] == nil {
			continue
		}
		devices = append(devices, genTLDeviceInfoFromRaw(i, *rawList.Device[i]))
	}
	return devices, nil
}

func (s *SDK) OpenGenTLDevice(info GenTLDeviceInfo, accessMode uint32) (*Camera, error) {
	s.mu.Lock()
	initialized := s.initialized
	s.mu.Unlock()

	if !initialized {
		return nil, ErrSDKNotInitialized
	}
	if accessMode == 0 {
		accessMode = AccessExclusive
	}

	raw := rawGenTLDevInfoFromPublic(info)
	handle, err := s.driver.createHandleByGenTL(&raw)
	if err != nil {
		return nil, err
	}
	if err := s.driver.openDevice(handle, accessMode, 0); err != nil {
		_ = s.driver.destroyHandle(handle)
		return nil, err
	}

	return &Camera{
		sdk:    s,
		handle: handle,
		info: DeviceInfo{
			Index:              info.Index,
			TransportLayerName: info.TLType,
			ModelName:          info.ModelName,
			SerialNumber:       info.SerialNumber,
			UserDefinedName:    info.UserDefinedName,
			ManufacturerName:   info.VendorName,
			InterfaceID:        info.InterfaceID,
			DeviceID:           info.DeviceID,
		},
		open: true,
	}, nil
}

func (i *Interface) Info() InterfaceInfo {
	i.mu.Lock()
	defer i.mu.Unlock()
	return i.info
}

func (i *Interface) EnumerateDevices() ([]DeviceInfo, error) {
	i.mu.Lock()
	if !i.open || i.handle == 0 {
		i.mu.Unlock()
		return nil, ErrCameraClosed
	}
	handle := i.handle
	i.mu.Unlock()

	var rawList mvCCDeviceInfoList
	if err := i.sdk.driver.enumDevicesByInterface(handle, &rawList); err != nil {
		return nil, err
	}

	count := int(rawList.DeviceNum)
	if count > len(rawList.DeviceInfo) {
		count = len(rawList.DeviceInfo)
	}
	devices := make([]DeviceInfo, 0, count)
	for idx := 0; idx < count; idx++ {
		if rawList.DeviceInfo[idx] == nil {
			continue
		}
		devices = append(devices, deviceInfoFromRaw(idx, *rawList.DeviceInfo[idx]))
	}
	return devices, nil
}

func (i *Interface) Close() error {
	i.mu.Lock()
	defer i.mu.Unlock()

	if !i.open && i.handle == 0 {
		return nil
	}

	var firstErr error
	if i.open {
		if err := i.sdk.driver.closeInterface(i.handle); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	if i.handle != 0 {
		if err := i.sdk.driver.destroyInterface(i.handle); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	i.handle = 0
	i.open = false
	return firstErr
}
