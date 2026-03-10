package elevatormanager

import (
	//"elevatorproject/src/elevio"
	"fmt"
	"testing"
	"time"
)

/*func TestButtonToString(t *testing.T){
	var button elevio.ButtonType
	button = elevio.BT_HallDown
	fmt.Println(buttonToString(button))
	button = elevio.BT_HallUp
	fmt.Println(buttonToString(button))
}*/

// func TestTimer(t *testing.T) {
// 	doortimer_start()
// 	<-doortimer.C
// 	fmt.Println("timer stopped")
// 	go test2()
// 	doortimer_start()
// 	<-doortimer.C
// 	fmt.Println("timer stopped")
// }

// func test2() {
// 	fmt.Println("Started test2")
// 	time.Sleep(time.Second)
// 	fmt.Println("Waited 1 sec")
// 	doortimer_start()
// 	<-doortimer.C
// 	fmt.Println("timer stopped from 2")
// }

// func TestTimer(t *testing.T){
// 	doortimer_init()
// 	doortimer_start()
// 	fmt.Println("timer started")
// 	go resetter_timer()
// 	<- doortimer.C
// 	fmt.Println("timer stopped")
// }

func resetter_timer(){
	for i := 0; i <5; i++{
		time.Sleep(time.Second)
		doortimer_start()
		fmt.Println("timer reset")
	}
}

func TestSingleElevator(t *testing.T){
	go main()
	select{}
}