package ac

import "github.com/Binozo/GoTuya/pkg/tuya"

type AC struct {
	*tuya.Device
}

// CreateAC creates an A/C instance to control it easily with the included api
func CreateAC(ip, deviceId string, key string) *AC {
	return &AC{
		tuya.CreateDevice(ip, deviceId, key, tuya.Version_3_3),
	}
}
