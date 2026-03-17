package elevatormanager

import (
	"elevatorproject/src/config"
	//"elevatorproject/src/elevio"
	"elevatorproject/src/types"
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

func resetter_timer() {
	for i := 0; i < 5; i++ {
		time.Sleep(time.Second)
		doortimer_start()
		fmt.Println("timer reset")
	}
}

func TestSingleElevator(t *testing.T) {
	config.Load()

	receiveOrdersCh := make(chan types.Order)
	receiveFinishedOrderCh := make(chan []types.Order)
	sendAssignmentsCh := make(chan [N_FLOORS][N_BUTTONS]bool)
	receiveElevatorState := make(chan types.ElevatorState)

	go ElevatorManager(receiveElevatorState, receiveOrdersCh, receiveFinishedOrderCh, sendAssignmentsCh)

	var requests [N_FLOORS][N_BUTTONS]bool

	for {
		select {
		case order := <-receiveOrdersCh:
			requests[order.Floor][order.Type] = true
			sendAssignmentsCh <- requests
		case clearedOrders := <-receiveFinishedOrderCh:
			for _, request := range clearedOrders {
				requests[request.Floor][request.Type] = false
			}
		case data := <-receiveElevatorState:
			fmt.Println("The elevator is currently: ", data.Behaviour)
		}
	}
}
