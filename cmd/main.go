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
