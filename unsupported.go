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
type Interface struct{}

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

func (s *SDK) EnumerateInterfaces(uint32) ([]InterfaceInfo, error) {
	return nil, ErrUnsupportedPlatform
}

func (s *SDK) OpenInterface(InterfaceInfo) (*Interface, error) {
	return nil, ErrUnsupportedPlatform
}

func (s *SDK) OpenInterfaceByID(string) (*Interface, error) {
	return nil, ErrUnsupportedPlatform
}

func (s *SDK) EnumerateGenTLInterfaces(string) ([]GenTLInterfaceInfo, error) {
	return nil, ErrUnsupportedPlatform
}

func (s *SDK) UnloadGenTLLibrary(string) error {
	return ErrUnsupportedPlatform
}

func (s *SDK) EnumerateGenTLDevices(GenTLInterfaceInfo) ([]GenTLDeviceInfo, error) {
	return nil, ErrUnsupportedPlatform
}

func (s *SDK) OpenGenTLDevice(GenTLDeviceInfo, uint32) (*Camera, error) {
	return nil, ErrUnsupportedPlatform
}

func (s *SDK) GetCameraLinkSerialPortList() ([]CameraLinkSerialPort, error) {
	return nil, ErrUnsupportedPlatform
}

func (s *SDK) SetCameraLinkEnumSerialPorts([]CameraLinkSerialPort) error {
	return ErrUnsupportedPlatform
}

func (i *Interface) Info() InterfaceInfo {
	return InterfaceInfo{}
}

func (i *Interface) EnumerateDevices() ([]DeviceInfo, error) {
	return nil, ErrUnsupportedPlatform
}

func (i *Interface) Close() error {
	return nil
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

func (c *Camera) RegisterAllEventCallback(EventCallback) error {
	return ErrUnsupportedPlatform
}

func (c *Camera) RegisterEventCallback(string, EventCallback) error {
	return ErrUnsupportedPlatform
}

func (c *Camera) EventNotificationOn(string) error {
	return ErrUnsupportedPlatform
}

func (c *Camera) EventNotificationOff(string) error {
	return ErrUnsupportedPlatform
}

func (c *Camera) ReadDeviceFileToFile(string, string) error {
	return ErrUnsupportedPlatform
}

func (c *Camera) WriteDeviceFileFromFile(string, string) error {
	return ErrUnsupportedPlatform
}

func (c *Camera) ReadDeviceFile(string, int) ([]byte, FileAccessProgress, error) {
	return nil, FileAccessProgress{}, ErrUnsupportedPlatform
}

func (c *Camera) WriteDeviceFile(string, []byte) (FileAccessProgress, error) {
	return FileAccessProgress{}, ErrUnsupportedPlatform
}

func (c *Camera) GetFileAccessProgress() (FileAccessProgress, error) {
	return FileAccessProgress{}, ErrUnsupportedPlatform
}

func (c *Camera) StartRecord(RecordOptions) error {
	return ErrUnsupportedPlatform
}

func (c *Camera) InputRecordFrame(*Frame) error {
	return ErrUnsupportedPlatform
}

func (c *Camera) InputRecordData([]byte) error {
	return ErrUnsupportedPlatform
}

func (c *Camera) StopRecord() error {
	return nil
}

func (c *Camera) SetCameraLinkBaudrate(uint32) error {
	return ErrUnsupportedPlatform
}

func (c *Camera) GetCameraLinkBaudrate() (uint32, error) {
	return 0, ErrUnsupportedPlatform
}

func (c *Camera) GetSupportedCameraLinkBaudrates() (uint32, error) {
	return 0, ErrUnsupportedPlatform
}

func (c *Camera) SetCameraLinkGenCPTimeout(time.Duration) error {
	return ErrUnsupportedPlatform
}

func (c *Camera) OpenSerialPort() error {
	return ErrUnsupportedPlatform
}

func (c *Camera) WriteSerialPort([]byte) (int, error) {
	return 0, ErrUnsupportedPlatform
}

func (c *Camera) ReadSerialPort(int, time.Duration) ([]byte, error) {
	return nil, ErrUnsupportedPlatform
}

func (c *Camera) ReadSerialPortInto([]byte, time.Duration) (int, error) {
	return 0, ErrUnsupportedPlatform
}

func (c *Camera) ClearSerialPort() error {
	return ErrUnsupportedPlatform
}

func (c *Camera) CloseSerialPort() error {
	return nil
}

func (c *Camera) RotateFrame(*Frame, RotationAngle) (*Frame, error) {
	return nil, ErrUnsupportedPlatform
}

func (c *Camera) FlipFrame(*Frame, FlipType) (*Frame, error) {
	return nil, ErrUnsupportedPlatform
}

func (c *Camera) SetBayerCvtQuality(InterpolationMethod) error {
	return ErrUnsupportedPlatform
}

func (c *Camera) SetBayerFilterEnable(bool) error {
	return ErrUnsupportedPlatform
}

func (c *Camera) SetBayerGammaValue(float32) error {
	return ErrUnsupportedPlatform
}

func (c *Camera) SetGammaValue(uint32, float32) error {
	return ErrUnsupportedPlatform
}

func (c *Camera) SetBayerGammaParam(GammaOptions) error {
	return ErrUnsupportedPlatform
}

func (c *Camera) SetBayerCCMParam(CCMOptions) error {
	return ErrUnsupportedPlatform
}

func (c *Camera) SetBayerCCMParamEx(CCMOptionsEx) error {
	return ErrUnsupportedPlatform
}

func (c *Camera) AdjustContrast(*Frame, ContrastOptions) (*Frame, error) {
	return nil, ErrUnsupportedPlatform
}

func (c *Camera) CorrectPurpleFringing(*Frame, PurpleFringingOptions) (*Frame, error) {
	return nil, ErrUnsupportedPlatform
}

func (c *Camera) SetISPConfig(string) error {
	return ErrUnsupportedPlatform
}

func (c *Camera) ISPProcess(*Frame, PixelConvertOptions) (*Frame, error) {
	return nil, ErrUnsupportedPlatform
}

func (c *Camera) DecodeHighBandwidthData([]byte, HBDecodeOptions) (*Frame, FrameSpecInfo, error) {
	return nil, FrameSpecInfo{}, ErrUnsupportedPlatform
}

func (c *Camera) HBDecode([]byte, HBDecodeOptions) (*Frame, FrameSpecInfo, error) {
	return nil, FrameSpecInfo{}, ErrUnsupportedPlatform
}

func (c *Camera) DecodeHighBandwidthFrame(*Frame, HBDecodeOptions) (*Frame, FrameSpecInfo, error) {
	return nil, FrameSpecInfo{}, ErrUnsupportedPlatform
}

func (c *Camera) HBDecodeFrame(*Frame, HBDecodeOptions) (*Frame, FrameSpecInfo, error) {
	return nil, FrameSpecInfo{}, ErrUnsupportedPlatform
}

func (c *Camera) ReconstructImage(*Frame, ReconstructImageOptions) ([]*Frame, error) {
	return nil, ErrUnsupportedPlatform
}

func (f *Frame) Image() (image.Image, error) {
	return nil, ErrUnsupportedPlatform
}

func (f *Frame) PixelTypeName() string {
	return ""
}
