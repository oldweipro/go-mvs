//go:build windows && amd64

package mvs

import (
	"fmt"
	"time"
)

func (s *SDK) GetCameraLinkSerialPortList() ([]CameraLinkSerialPort, error) {
	var raw mvCamlSerialPortList
	if err := s.driver.camlGetSerialPortList(&raw); err != nil {
		return nil, err
	}

	count := int(raw.SerialPortNum)
	if count > len(raw.SerialPort) {
		count = len(raw.SerialPort)
	}
	ports := make([]CameraLinkSerialPort, 0, count)
	for i := 0; i < count; i++ {
		name := byteString(raw.SerialPort[i].SerialPort[:])
		if name == "" {
			continue
		}
		ports = append(ports, CameraLinkSerialPort{Name: name})
	}
	return ports, nil
}

func (s *SDK) SetCameraLinkEnumSerialPorts(ports []CameraLinkSerialPort) error {
	if len(ports) > mvMaxSerialPortNum {
		return fmt.Errorf("%w: too many serial ports: %d", ErrInvalidArgument, len(ports))
	}

	var raw mvCamlSerialPortList
	raw.SerialPortNum = uint32(len(ports))
	for i, port := range ports {
		if port.Name == "" {
			return fmt.Errorf("%w: serial port name is empty", ErrInvalidArgument)
		}
		copyFixedString(raw.SerialPort[i].SerialPort[:], port.Name)
	}
	return s.driver.camlSetEnumSerialPorts(&raw)
}

func (c *Camera) SetCameraLinkBaudrate(baudrate uint32) error {
	if baudrate == 0 {
		return fmt.Errorf("%w: baudrate is empty", ErrInvalidArgument)
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.open || c.handle == 0 {
		return ErrCameraClosed
	}
	return c.sdk.driver.camlSetDeviceBaudrate(c.handle, baudrate)
}

func (c *Camera) GetCameraLinkBaudrate() (uint32, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.open || c.handle == 0 {
		return 0, ErrCameraClosed
	}
	return c.sdk.driver.camlGetDeviceBaudrate(c.handle)
}

func (c *Camera) GetSupportedCameraLinkBaudrates() (uint32, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.open || c.handle == 0 {
		return 0, ErrCameraClosed
	}
	return c.sdk.driver.camlGetSupportBaudrates(c.handle)
}

func (c *Camera) SetCameraLinkGenCPTimeout(timeout time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.open || c.handle == 0 {
		return ErrCameraClosed
	}
	return c.sdk.driver.camlSetGenCPTimeout(c.handle, timeoutMilliseconds(timeout))
}

func (c *Camera) OpenSerialPort() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.open || c.handle == 0 {
		return ErrCameraClosed
	}
	if c.serialOpen {
		return nil
	}
	if err := c.sdk.driver.serialPortOpen(c.handle); err != nil {
		return err
	}
	c.serialOpen = true
	return nil
}

func (c *Camera) WriteSerialPort(data []byte) (int, error) {
	if len(data) == 0 {
		return 0, fmt.Errorf("%w: serial data is empty", ErrInvalidArgument)
	}
	if uint64(len(data)) > maxUint32Value {
		return 0, fmt.Errorf("%w: serial data is larger than UINT_MAX", ErrInvalidArgument)
	}

	c.mu.Lock()
	if !c.open || c.handle == 0 {
		c.mu.Unlock()
		return 0, ErrCameraClosed
	}
	handle := c.handle
	c.mu.Unlock()

	total := 0
	for total < len(data) {
		end := total + SerialPortMaxWriteSize
		if end > len(data) {
			end = len(data)
		}
		written, err := c.sdk.driver.serialPortWrite(handle, data[total:end])
		if err != nil {
			return total, err
		}
		if written == 0 {
			return total, fmt.Errorf("%w: serial port wrote 0 bytes", ErrInvalidFrameData)
		}
		if int(written) > end-total {
			return total, fmt.Errorf("%w: SDK reported %d written bytes for a %d-byte serial chunk", ErrInvalidFrameData, written, end-total)
		}
		total += int(written)
	}
	return total, nil
}

func (c *Camera) ReadSerialPort(bufferSize int, timeout time.Duration) ([]byte, error) {
	if bufferSize <= 0 || uint64(bufferSize) > maxUint32Value {
		return nil, fmt.Errorf("%w: invalid serial buffer size %d", ErrInvalidArgument, bufferSize)
	}
	buffer := make([]byte, bufferSize)
	n, err := c.ReadSerialPortInto(buffer, timeout)
	if err != nil {
		return nil, err
	}
	return buffer[:n], nil
}

func (c *Camera) ReadSerialPortInto(buffer []byte, timeout time.Duration) (int, error) {
	if len(buffer) == 0 {
		return 0, fmt.Errorf("%w: serial buffer is empty", ErrInvalidArgument)
	}
	if uint64(len(buffer)) > maxUint32Value {
		return 0, fmt.Errorf("%w: serial buffer is larger than UINT_MAX", ErrInvalidArgument)
	}

	c.mu.Lock()
	if !c.open || c.handle == 0 {
		c.mu.Unlock()
		return 0, ErrCameraClosed
	}
	handle := c.handle
	c.mu.Unlock()

	n, err := c.sdk.driver.serialPortRead(handle, buffer, timeoutMilliseconds(timeout))
	if int(n) > len(buffer) {
		return 0, fmt.Errorf("%w: SDK reported %d serial bytes into %d-byte buffer", ErrInvalidFrameData, n, len(buffer))
	}
	return int(n), err
}

func (c *Camera) ClearSerialPort() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.open || c.handle == 0 {
		return ErrCameraClosed
	}
	return c.sdk.driver.serialPortClearBuffer(c.handle)
}

func (c *Camera) CloseSerialPort() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.open || c.handle == 0 {
		return nil
	}
	if !c.serialOpen {
		return nil
	}
	if err := c.sdk.driver.serialPortClose(c.handle); err != nil {
		return err
	}
	c.serialOpen = false
	return nil
}
