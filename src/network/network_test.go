package network

import (
	"Network-go/network/bcast"
	"elevatorproject/src/config"
	elevatormanager "elevatorproject/src/elevatorManager"
	"elevatorproject/src/elevio"
	"elevatorproject/src/network"
	"elevatorproject/src/types"
	"fmt"
	"testing"
)

func TestSender(t *testing.T){
	config.Load()

	sendOrdersCh := make(chan OrdersAndStateUpdate)
	ackCh := make(chan OrdersAndStateAck)

	go bcast.Transmitter(config.Cfg.MasterListenPort, sendOrdersCh)
	go bcast.Receiver(config.Cfg.MasterListenPort, ackCh)

	orderSender := &GenericSender[OrdersAndStateUpdate, OrdersAndStateAck]{
		SendCh: sendOrdersCh,
		AckIn: ackCh,
		AckResults: make(chan AckResult, 10), // buffered
	}

	msg := OrdersAndStateUpdate{"1", 2, "Hei ehi"}

	fmt.Println(msg.OrdersAndState)

	orderSender.SendAsyncWithAck(msg)

	res:=  <- orderSender.AckResults


	fmt.Println(res.UpdateNr)
}

func masterSender(){
		var readyToSendOrder bool = true
	var messageCount int = 0
	// Recive from elevatorManager, send to master.

	receiveOrdersCh := make(chan types.Order)
	receiveFinishedAssignmentsCh := make(chan []elevio.ButtonEvent)
	receiveElevatorState := make(chan types.ElevatorState)

	sendAssignmentsCh := make(chan [N_FLOORS][N_BUTTONS]bool)
	sendElevatorState := make(chan types.ElevatorState)


	go elevatormanager.ElevatorManager(receiveElevatorState, receiveOrdersCh, receiveFinishedAssignmentsCh, sendAssignmentsCh)

	//Init Order ack
	sendOrdersCh := make(chan types.HallOrder)
	hallOrderAck := make(chan types.HallOrderAck)

	orderSender := &network.GenericSender[types.HallOrder, types.HallOrderAck]{
		SendCh:     sendOrdersCh,
		AckIn:      hallOrderAck,
		AckResults: make(chan network.AckResult, 10), // buffered
	}


	// Finished Assignments setup 
	sendFinishedAssignmentsCh := make(chan types.FinishedHallAssignments)
	finishedOrdersAckCh := make(chan types.FinishedHallAssignmentsAck)

	completeAssignmentSender := &network.GenericSender[types.FinishedHallAssignments, types.FinishedHallAssignmentsAck]{
		SendCh:     sendFinishedAssignmentsCh,
		AckIn:      finishedOrdersAckCh,
		AckResults: make(chan network.AckResult, 10), // buffered
	}



	
	go bcast.Transmitter(config.Cfg.MasterListenPort, sendElevatorState, sendOrdersCh, sendFinishedAssignmentsCh)

	receiveAssignmentsFromMasterCh := make(chan types.Assignements) //Denne skal vel egentlig bli passet som funksjonsparameter

	go bcast.Receiver(config.Cfg.SlaveListenPort, receiveAssignmentsFromMasterCh, hallOrderAck, finishedOrdersAckCh)

	var slaveRequests [N_FLOORS][N_BUTTONS]bool
}