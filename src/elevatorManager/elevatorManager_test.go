package elevatormanager

import (
	"elevatorproject/src/elevio"
	"testing"
	"fmt"
)

func TestButtonToString(t *testing.T){
	var button elevio.ButtonType
	button = elevio.BT_HallDown
	fmt.Println(buttonToString(button))
	button = elevio.BT_HallUp
	fmt.Println(buttonToString(button))
}