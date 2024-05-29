# GoTuya 🔓🌩️

A Go Api to control your Tuya devices through your Local Network (LAN)

Features:
- 🌧🙅‍♂️ Works completely offline without any connection to the 💩 and unsecure chinese clouds
- 🔌 Easily extendable: Control any Tuya device with ease
- 🔬 Battle tested with a Tuya TCL A/C

> [!IMPORTANT]
> Currently only Tuya devices with Version 3.3 are supported

> [!NOTE]
> This Go api wouldn't be possible without the amazing [tuyapi project](https://github.com/codetheweb/tuyapi)

## Quickstart
### Install the Package:
```bash
$ go get -u github.com/Binozo/GoTuya
```

### Setup your Tuya Device
You will need the following data to control your Tuya device:
- The device's ip in your local network
- The `deviceId`
- The `localKey`

To retrieve those values you have to follow [those steps](https://github.com/codetheweb/tuyapi/blob/master/docs/SETUP.md#linking-a-tuya-device-with-smart-link).

### Implement the Api
Here is an example for an Air Conditioner:

```go
package main

import (
	"fmt"
	"github.com/Binozo/GoTuya/pkg/ac"
	"time"
)

func main() {
	myTclAc := ac.CreateAC("192.168.178.30", "15580880bcaac262j6eg", "A2In><,:-{Hy:[%K7")
	isOn, err := myTclAc.IsOn()
	if err != nil {
		panic(err)
	}
	fmt.Println("A/C is on:", isOn)

	if !isOn {
		targetTemp := 20
		fmt.Printf("The A/C is not on. Starting cooling with %d°C...\n", targetTemp)
		if err := myTclAc.SetTemperature(targetTemp); err != nil {
			fmt.Println("Couldn't turn on A/C:", err.Error())
			return
		}
		fmt.Println("A/C has been turn on. Waiting for 5 seconds to cool")
		time.Sleep(time.Second * 5)

		fmt.Println("Turning off")
		if err := myTclAc.Power(false); err != nil {
			fmt.Println("Couldn't turn off A/C:", err.Error())
		}
	}
}
```

### 🔌 Extending 
To use this Api to connect to your tuya device take a look at the tuya package.
The `ac` package is a wrapper for that.

Generally looking at the implementation of the `ac` device will help you implement your specific device.

Contributions are always welcome!