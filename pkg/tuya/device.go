package tuya

import "net"

// Device generic tuya api interface
type Device struct {
	// IP of the Device
	IP string
	// DeviceID used for encryption
	DeviceID string
	// Key for the encrypted traffic
	Key []byte
	// Version for the different tuya api specifications
	Version Version
	// currentSequenceNr used for communication
	currentSequenceNr int
	conn              *net.Conn
	// currentStatus cached responses and device status
	currentStatus response
}

// CreateDevice for generic tuya devices
func CreateDevice(ip string, deviceId string, key string, version Version) *Device {
	return &Device{
		IP:       ip,
		DeviceID: deviceId,
		Key:      []byte(key),
		Version:  version,
	}
}
