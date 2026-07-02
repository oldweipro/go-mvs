//go:build windows && amd64

package mvs

import "fmt"

func New(config Config) (*SDK, error) {
	driver, err := newDriver(config.DLLPath)
	if err != nil {
		return nil, err
	}
	return &SDK{driver: driver}, nil
}

func (s *SDK) DLLPath() string {
	return s.driver.dllPath
}

func (s *SDK) Initialize() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.initialized {
		return nil
	}
	if err := s.driver.initialize(); err != nil {
		return err
	}
	s.initialized = true
	return nil
}

func (s *SDK) Finalize() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.initialized {
		return nil
	}
	if err := s.driver.finalize(); err != nil {
		return err
	}
	s.initialized = false
	return nil
}

func (s *SDK) Close() error {
	return s.Finalize()
}

func (s *SDK) Version() Version {
	raw := s.driver.getSDKVersion()
	return Version{
		Major:    uint8(raw >> 24),
		Minor:    uint8(raw >> 16),
		Revision: uint8(raw >> 8),
		Build:    uint8(raw),
	}
}

func (s *SDK) EnumerateDevices(layerType uint32) ([]DeviceInfo, error) {
	s.mu.Lock()
	initialized := s.initialized
	s.mu.Unlock()

	if !initialized {
		return nil, ErrSDKNotInitialized
	}

	var rawList mvCCDeviceInfoList
	if err := s.driver.enumDevices(layerType, &rawList); err != nil {
		return nil, err
	}

	count := int(rawList.DeviceNum)
	if count > len(rawList.DeviceInfo) {
		count = len(rawList.DeviceInfo)
	}

	devices := make([]DeviceInfo, 0, count)
	for i := 0; i < count; i++ {
		if rawList.DeviceInfo[i] == nil {
			continue
		}
		devices = append(devices, deviceInfoFromRaw(i, *rawList.DeviceInfo[i]))
	}
	return devices, nil
}

func (s *SDK) EnumerateDefaultDevices() ([]DeviceInfo, error) {
	return s.EnumerateDevices(DefaultDeviceTransportLayer)
}

func (s *SDK) IsDeviceAccessible(info DeviceInfo, accessMode uint32) bool {
	if accessMode == 0 {
		accessMode = AccessExclusive
	}
	return s.driver.isDeviceAccessible(&info.raw, accessMode)
}

func (s *SDK) OpenDevice(info DeviceInfo, accessMode uint32) (*Camera, error) {
	s.mu.Lock()
	initialized := s.initialized
	s.mu.Unlock()

	if !initialized {
		return nil, ErrSDKNotInitialized
	}
	if accessMode == 0 {
		accessMode = AccessExclusive
	}

	handle, err := s.driver.createHandle(&info.raw)
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
		info:   info,
		open:   true,
	}, nil
}

func (s *SDK) FindDeviceBySerial(serial string) (DeviceInfo, bool, error) {
	return s.findDefaultDevice(func(device DeviceInfo) bool {
		return serial != "" && device.SerialNumber == serial
	})
}

func (s *SDK) FindDeviceByIP(ip string) (DeviceInfo, bool, error) {
	return s.findDefaultDevice(func(device DeviceInfo) bool {
		return ip != "" && device.CurrentIP == ip
	})
}

func (s *SDK) FindDeviceByUserDefinedName(name string) (DeviceInfo, bool, error) {
	return s.findDefaultDevice(func(device DeviceInfo) bool {
		return name != "" && device.UserDefinedName == name
	})
}

func (s *SDK) OpenDeviceBySerial(serial string, accessMode uint32) (*Camera, error) {
	device, ok, err := s.FindDeviceBySerial(serial)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("%w: serial %q", ErrDeviceNotFound, serial)
	}
	return s.OpenDevice(device, accessMode)
}

func (s *SDK) OpenDeviceByIP(ip string, accessMode uint32) (*Camera, error) {
	device, ok, err := s.FindDeviceByIP(ip)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("%w: ip %q", ErrDeviceNotFound, ip)
	}
	return s.OpenDevice(device, accessMode)
}

func (s *SDK) OpenDeviceByUserDefinedName(name string, accessMode uint32) (*Camera, error) {
	device, ok, err := s.FindDeviceByUserDefinedName(name)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("%w: user-defined name %q", ErrDeviceNotFound, name)
	}
	return s.OpenDevice(device, accessMode)
}

func (s *SDK) findDefaultDevice(match func(DeviceInfo) bool) (DeviceInfo, bool, error) {
	devices, err := s.EnumerateDefaultDevices()
	if err != nil {
		return DeviceInfo{}, false, err
	}
	for _, device := range devices {
		if match(device) {
			return device, true, nil
		}
	}
	return DeviceInfo{}, false, nil
}
