package ac

import (
	"errors"
	"fmt"
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
