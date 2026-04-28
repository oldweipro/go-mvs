//go:build windows && amd64

package mvsdk

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
