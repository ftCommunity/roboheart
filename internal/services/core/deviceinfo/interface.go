package deviceinfo

type DeviceInfo interface {
	GetPlatform() string
	GetDevice() string
}
