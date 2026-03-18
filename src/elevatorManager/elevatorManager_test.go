package elevatormanager

import (
	"elevatorproject/src/config"
	//"elevatorproject/src/elevio"
	"elevatorproject/src/types"
	"fmt"
	"testing"
	"time"
)


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

	go ElevatorManager(receiveElevatorState, receiveOrdersCh, receiveFinishedOrderCh, sendAssignmentsCh,sendLightsCh)

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
