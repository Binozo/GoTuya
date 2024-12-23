# GoTuya 🔓🌩️

A Go Api to control your Tuya devices through your Local Network (LAN)

Features:
- 🌧🙅‍♂️ Works completely offline without any connection to the 💩 and unsecure chinese clouds
- 🔌 Easily extendable: Control any Tuya device with ease
- 🔬 Battle tested with a Tuya TCL A/C

> [!IMPORTANT]
> Currently only Tuya devices with Version 3.3 are supported because I don't own any devices that are using different protocol versions.
> Feel free to add them. If you need help you can take a look at the [tuyapi project](https://github.com/codetheweb/tuyapi/tree/master/lib) or create an Issue.

> [!NOTE]
> This Go api wouldn't be possible without the amazing [tuyapi project](https://github.com/codetheweb/tuyapi)

I use this package myself at my home with two A/Cs, so don't worry if there is no commit for a longer period of time.
If this package breaks I will fix it :)

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

### 🔌 Extending with your own Tuya device
Tuya devices work all the same way. They work with _dictionaries_.
Let's explain this with an example:

```go
state := map[string]int{
	"TEMPERATURE": 2,
	"SPEED": 5,
	"FEATURE_C": 0,
	"FEATURE_D": 1,
}
```

Now if we want to set the temperature to 10°C we would just set
```go
state["TEMPERATURE"] = 10
```
programmatically and upload that dictionary to the Tuya device. We don't need to send every value, we can just send the key value pair we just changed.

This is the way it works under the hood. Easy, isn't it? 😄

Now if you want to add your own device you just have to find out which keys exist in that map and which values the can have:

```go
package main

import (
	"fmt"
	"github.com/Binozo/GoTuya/pkg/tuya"
)

func main() {
	device := tuya.CreateDevice("IP", "DEVICEID", "KEY", tuya.Version_3_3)
	state := device.GetCurrentStatus() // TODO
	fmt.Println("Current state:", state)
}
```
Execute this code everytime you changed a parameter (e.g. Temeprature) of your device. This way you can find out which key stands for what feature.

The `ac` package is a prebuilt wrapper for A/Cs. Looking at the implementation could help you implement your own device.

PRs are always welcome!