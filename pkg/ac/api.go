package ac

import (
	"errors"
	"fmt"
	"strconv"
)

func (a *AC) IsOn() (bool, error) {
	if !a.IsConnected() {
		if err := a.Connect(); err != nil {
			return false, err
		}
		defer a.Disconnect()
	}
	currentStatus := a.GetCurrentStatus()
	if onValue, ok := currentStatus[onDpsIndex]; ok {
		return onValue.(bool), nil
	} else {
		return false, errors.New(fmt.Sprintf("dps index %s not contained in: %s", onDpsIndex, currentStatus))
	}
}

func (a *AC) CurrentTemperature() (float64, error) {
	if !a.IsConnected() {
		if err := a.Connect(); err != nil {
			return 0, err
		}
		defer a.Disconnect()
	}
	currentStatus := a.GetCurrentStatus()
	if onValue, ok := currentStatus[temperatureDpsIndex]; ok {
		return onValue.(float64), nil
	} else {
		return 0, errors.New(fmt.Sprintf("dps index %s not contained in: %s", onDpsIndex, currentStatus))
	}
}

func (a *AC) Power(powerOn bool) error {
	if !a.IsConnected() {
		if err := a.Connect(); err != nil {
			return err
		}
		defer a.Disconnect()
	}
	return a.Set(map[string]interface{}{
		onDpsIndex: powerOn,
	})
}

func (a *AC) SetTemperature(temperature int) error {
	if !a.IsConnected() {
		if err := a.Connect(); err != nil {
			return err
		}
		defer a.Disconnect()
	}
	return a.Set(map[string]interface{}{
		onDpsIndex:          true,
		temperatureDpsIndex: temperature,
	})
}

// SetFanIntensity sets the fan intensity on a scale between 1 and 4.
// 1 is low,
// 4 is high
func (a *AC) SetFanIntensity(intensity int) error {
	if intensity < 1 || intensity > 4 {
		return errors.New("intensity must be between 1 and 4")
	}
	if !a.IsConnected() {
		if err := a.Connect(); err != nil {
			return err
		}
		defer a.Disconnect()
	}
	return a.Set(map[string]interface{}{
		onDpsIndex:           true,
		fanIntensityDpsIndex: intensity,
	})
}

// GetFanIntensity gets the current fan intensity on a scale between 1 and 4.
func (a *AC) GetFanIntensity() (int, error) {
	if !a.IsConnected() {
		if err := a.Connect(); err != nil {
			return 0, err
		}
		defer a.Disconnect()
	}
	currentStatus := a.GetCurrentStatus()
	if intensityValue, ok := currentStatus[fanIntensityDpsIndex]; ok {
		return strconv.Atoi(intensityValue.(string))
	} else {
		return 0, errors.New(fmt.Sprintf("dps index %s not contained in: %s", onDpsIndex, currentStatus))
	}
}

func (a *AC) SetFanSwing(swing bool) error {
	if !a.IsConnected() {
		if err := a.Connect(); err != nil {
			return err
		}
		defer a.Disconnect()
	}
	return a.Set(map[string]interface{}{
		onDpsIndex:       true,
		fanSwingDpsIndex: swing,
	})
}

func (a *AC) GetFanSwinging() (bool, error) {
	if !a.IsConnected() {
		if err := a.Connect(); err != nil {
			return false, err
		}
		defer a.Disconnect()
	}
	currentStatus := a.GetCurrentStatus()
	if intensityValue, ok := currentStatus[fanSwingDpsIndex]; ok {
		return intensityValue.(bool), nil
	} else {
		return false, errors.New(fmt.Sprintf("dps index %s not contained in: %s", onDpsIndex, currentStatus))
	}
}

func (a *AC) SetTurboMode(turbo bool) error {
	if !a.IsConnected() {
		if err := a.Connect(); err != nil {
			return err
		}
		defer a.Disconnect()
	}
	return a.Set(map[string]interface{}{
		onDpsIndex:        true,
		turboModeDpsIndex: turbo,
	})
}

func (a *AC) GetIsTurboEnabled() (bool, error) {
	if !a.IsConnected() {
		if err := a.Connect(); err != nil {
			return false, err
		}
		defer a.Disconnect()
	}
	currentStatus := a.GetCurrentStatus()
	if turboValue, ok := currentStatus[turboModeDpsIndex]; ok {
		return turboValue.(bool), nil
	} else {
		return false, errors.New(fmt.Sprintf("dps index %s not contained in: %s", onDpsIndex, currentStatus))
	}
}

func (a *AC) SetNightMode(nightMode bool) error {
	if !a.IsConnected() {
		if err := a.Connect(); err != nil {
			return err
		}
		defer a.Disconnect()
	}
	return a.Set(map[string]interface{}{
		onDpsIndex:        true,
		nightModeDpsIndex: nightMode,
	})
}

func (a *AC) GetIsNightModeEnabled() (bool, error) {
	if !a.IsConnected() {
		if err := a.Connect(); err != nil {
			return false, err
		}
		defer a.Disconnect()
	}
	currentStatus := a.GetCurrentStatus()
	if nightValue, ok := currentStatus[nightModeDpsIndex]; ok {
		return nightValue.(bool), nil
	} else {
		return false, errors.New(fmt.Sprintf("dps index %s not contained in: %s", onDpsIndex, currentStatus))
	}
}
