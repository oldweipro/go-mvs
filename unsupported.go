//go:build !windows || !amd64

package mvs

import (
	"errors"
	"image"
	"time"
)

var ErrUnsupportedPlatform = errors.New("mvs: Hikrobot MVS is supported on windows/amd64 only")

type Config struct {
	DLLPath string
}

type Version struct {
	Major    uint8
	Minor    uint8
	Revision uint8
	Build    uint8
}

type DeviceInfo struct{}
type Frame struct{}
type FloatValue struct{}
type IntValue struct{}
type EnumValue struct{}
type StringValue struct{}
type SDK struct{}
type Camera struct{}

func New(Config) (*SDK, error) {
	return nil, ErrUnsupportedPlatform
}

func (v Version) String() string {
	return ""
}

func (s *SDK) DLLPath() string {
	return ""
}

func (s *SDK) Initialize() error {
	return ErrUnsupportedPlatform
}

func (s *SDK) Finalize() error {
	return nil
}

func (s *SDK) Close() error {
	return nil
}

func (s *SDK) Version() Version {
	return Version{}
}

func (s *SDK) EnumerateDevices(uint32) ([]DeviceInfo, error) {
	return nil, ErrUnsupportedPlatform
}

func (s *SDK) EnumerateDefaultDevices() ([]DeviceInfo, error) {
	return nil, ErrUnsupportedPlatform
}

func (s *SDK) IsDeviceAccessible(DeviceInfo, uint32) bool {
	return false
}

func (s *SDK) OpenDevice(DeviceInfo, uint32) (*Camera, error) {
	return nil, ErrUnsupportedPlatform
}

func (s *SDK) FindDeviceBySerial(string) (DeviceInfo, bool, error) {
	return DeviceInfo{}, false, ErrUnsupportedPlatform
}

func (s *SDK) FindDeviceByIP(string) (DeviceInfo, bool, error) {
	return DeviceInfo{}, false, ErrUnsupportedPlatform
}

func (s *SDK) FindDeviceByUserDefinedName(string) (DeviceInfo, bool, error) {
	return DeviceInfo{}, false, ErrUnsupportedPlatform
}

func (s *SDK) OpenDeviceBySerial(string, uint32) (*Camera, error) {
	return nil, ErrUnsupportedPlatform
}

func (s *SDK) OpenDeviceByIP(string, uint32) (*Camera, error) {
	return nil, ErrUnsupportedPlatform
}

func (s *SDK) OpenDeviceByUserDefinedName(string, uint32) (*Camera, error) {
	return nil, ErrUnsupportedPlatform
}

func (c *Camera) Info() DeviceInfo {
	return DeviceInfo{}
}

func (c *Camera) IsConnected() (bool, error) {
	return false, ErrUnsupportedPlatform
}

func (c *Camera) SetImageNodeNum(uint32) error {
	return ErrUnsupportedPlatform
}

func (c *Camera) SetGrabStrategy(GrabStrategy) error {
	return ErrUnsupportedPlatform
}

func (c *Camera) SetOutputQueueSize(uint32) error {
	return ErrUnsupportedPlatform
}

func (c *Camera) ClearImageBuffer() error {
	return ErrUnsupportedPlatform
}

func (c *Camera) ConfigureOptimalPacketSize() error {
	return ErrUnsupportedPlatform
}

func (c *Camera) StartGrabbing() error {
	return ErrUnsupportedPlatform
}

func (c *Camera) StopGrabbing() error {
	return nil
}

func (c *Camera) Close() error {
	return nil
}

func (c *Camera) GetFrame(time.Duration) (*Frame, error) {
	return nil, ErrUnsupportedPlatform
}

func (c *Camera) GetOneFrameTimeout(time.Duration) (*Frame, error) {
	return nil, ErrUnsupportedPlatform
}

func (c *Camera) GetOneFrameInto([]byte, time.Duration) (*Frame, error) {
	return nil, ErrUnsupportedPlatform
}

func (c *Camera) CurrentFrameBufferSize() (int, error) {
	return 0, ErrUnsupportedPlatform
}

func (c *Camera) GetFloat(string) (FloatValue, error) {
	return FloatValue{}, ErrUnsupportedPlatform
}

func (c *Camera) SetFloat(string, float32) error {
	return ErrUnsupportedPlatform
}

func (c *Camera) GetInt(string) (IntValue, error) {
	return IntValue{}, ErrUnsupportedPlatform
}

func (c *Camera) SetInt(string, int64) error {
	return ErrUnsupportedPlatform
}

func (c *Camera) GetEnum(string) (EnumValue, error) {
	return EnumValue{}, ErrUnsupportedPlatform
}

func (c *Camera) SetEnum(string, uint32) error {
	return ErrUnsupportedPlatform
}

func (c *Camera) SetEnumByString(string, string) error {
	return ErrUnsupportedPlatform
}

func (c *Camera) GetEnumEntrySymbolic(string, uint32) (string, error) {
	return "", ErrUnsupportedPlatform
}

func (c *Camera) GetEnumEntries(string) ([]EnumEntry, error) {
	return nil, ErrUnsupportedPlatform
}

func (c *Camera) GetBool(string) (bool, error) {
	return false, ErrUnsupportedPlatform
}

func (c *Camera) SetBool(string, bool) error {
	return ErrUnsupportedPlatform
}

func (c *Camera) GetString(string) (StringValue, error) {
	return StringValue{}, ErrUnsupportedPlatform
}

func (c *Camera) SetString(string, string) error {
	return ErrUnsupportedPlatform
}

func (c *Camera) Command(string) error {
	return ErrUnsupportedPlatform
}

func (c *Camera) SetTriggerMode(bool) error {
	return ErrUnsupportedPlatform
}

func (c *Camera) TriggerOnce() error {
	return ErrUnsupportedPlatform
}

func (c *Camera) FeatureSave(string) error {
	return ErrUnsupportedPlatform
}

func (c *Camera) FeatureLoad(string) error {
	return ErrUnsupportedPlatform
}

func (c *Camera) SaveFrameToFile(*Frame, string, ImageSaveOptions) error {
	return ErrUnsupportedPlatform
}

func (c *Camera) ConvertPixelType(*Frame, uint32) (*Frame, error) {
	return nil, ErrUnsupportedPlatform
}

func (c *Camera) ConvertFrame(*Frame, PixelConvertOptions) (*Frame, error) {
	return nil, ErrUnsupportedPlatform
}

func (c *Camera) RegisterFrameCallback(FrameCallback) error {
	return ErrUnsupportedPlatform
}

func (f *Frame) Image() (image.Image, error) {
	return nil, ErrUnsupportedPlatform
}

func (f *Frame) PixelTypeName() string {
	return ""
}
