package elevatormanager

import (
	//"elevatorproject/src/elevio"
	"testing"
	"fmt"
)

/*func TestButtonToString(t *testing.T){
	var button elevio.ButtonType
	button = elevio.BT_HallDown
	fmt.Println(buttonToString(button))
	button = elevio.BT_HallUp
	fmt.Println(buttonToString(button))
}*/

func TestTimer(t *testing.T){
	fmt.Println("started timertest")
	doorCh := make(chan bool)
	go doortimer(doorCh)
	doorCh <- true
	<- doorCh
	fmt.Println("timer finished")
	doorCh <- true
	doorCh <- false
}