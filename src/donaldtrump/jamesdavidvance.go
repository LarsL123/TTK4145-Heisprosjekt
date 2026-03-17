package donaldtrump

import (
	"Network-go/network/bcast"
	"elevatorproject/src/config"
	elevatormanager "elevatorproject/src/elevatorManager"
	"elevatorproject/src/types"
	"fmt"
	"strconv"
	"time"
)

func RunSlaveBrain(id string) {
	var messageCount int = 0
	var slaveRequests [N_FLOORS][N_BUTTONS]bool
	var lightsOn [N_FLOORS][N_BUTTONS]bool

	// Channels
	receiveOrdersCh := make(chan types.Order)
	receiveFinishedAssignmentsCh := make(chan []types.Order)
	receiveElevatorState := make(chan types.ElevatorState)

	sendAssignmentsCh := make(chan [N_FLOORS][N_BUTTONS]bool)
	sendElevatorState := make(chan types.ElevatorState)
	sendLightsCh := make(chan [N_FLOORS][N_BUTTONS]bool)

	// Send orders to masters
	sendOrdersCh := make(chan types.HallOrder, 10)
	hallOrderAck := make(chan types.HallOrderAck, 10)

	// Finished Assignments
	sendFinishedAssignmentsCh := make(chan types.FinishedHallAssignments)
	finishedAssignmentsAckCh := make(chan types.FinishedHallAssignmentsAck)

	// Newly calculated assignments
	receiveAssignmentsFromMasterCh := make(chan types.Assignments)

	// Start ElevatorManager
	go elevatormanager.ElevatorManager(receiveElevatorState, receiveOrdersCh, receiveFinishedAssignmentsCh, sendAssignmentsCh, sendLightsCh)

	// Broadcast transmitter & receiver
	go bcast.Transmitter(config.Cfg.MasterListenPort, sendElevatorState, sendOrdersCh, sendFinishedAssignmentsCh)
	go bcast.Receiver(config.Cfg.SlaveListenPort, receiveAssignmentsFromMasterCh, hallOrderAck, finishedAssignmentsAckCh)

	//Resending logic:
	pendingOrders := make(map[int]types.HallOrder)
	pendingFinishedAssignments := make(map[int]types.FinishedHallAssignments)
	resendTicker := time.NewTicker(config.Cfg.AckRetryRate)

	for {
		select {

		case state := <-receiveElevatorState:
			state.ID = id
			sendElevatorState <- state

		case order := <-receiveOrdersCh:
			switch order.Type {
			case types.Cab:
				slaveRequests[order.Floor][order.Type] = true // TODO: persist if elevator dies
				sendAssignmentsCh <- slaveRequests
			default:
				messageCount++
				ho := createHallOrder(id, order, messageCount)

				sendOrdersCh <- ho
				pendingOrders[ho.UpdateNr] = ho
			}
		case <-resendTicker.C:
			removeTimeouts(pendingOrders)
			removeTimeouts(pendingFinishedAssignments)

			for updateNr, ho := range pendingOrders {
				fmt.Println("Resending Order: ", updateNr)
				sendOrdersCh <- ho
			}

			for updateNr, ho := range pendingFinishedAssignments {
				fmt.Println("Resending Assignment: ", updateNr)
				sendFinishedAssignmentsCh <- ho
			}

		case ack := <-hallOrderAck:
			fmt.Println("Recived ACK for order", ack.UpdateNr)
			delete(pendingOrders, ack.UpdateNr)

		case finishedOrders := <-receiveFinishedAssignmentsCh:
			for _, request := range finishedOrders {
				slaveRequests[request.Floor][request.Type] = false
				lightsOn[request.Floor][request.Type] = false
			}
			sendLightsCh <- lightsOn

			messageCount++
			finishedAssigment := createFinishedAssignments(id, finishedOrders, messageCount)

			fmt.Println("Clearing assignment")
			sendFinishedAssignmentsCh <- finishedAssigment
			pendingFinishedAssignments[finishedAssigment.UpdateNr] = finishedAssigment

			//TODO: choose whether to send full requests or only changes

		case ack := <-finishedAssignmentsAckCh:
			fmt.Println("Recived ACK for assignemnt", ack.UpdateNr)
			delete(pendingFinishedAssignments, ack.UpdateNr)

		case assignments := <-receiveAssignmentsFromMasterCh:
			fmt.Println("Received assignments. Doing the work")

			// Combine hall requests and cab requests before sending
			for i := range slaveRequests {
				slaveRequests[i][0] = assignments.Data[id][i][0]
				slaveRequests[i][1] = assignments.Data[id][i][1]
			}

			// Send assignments to elevator
			sendAssignmentsCh <- slaveRequests

			// Prepare lights on
			lightsOn = lightsFromAssignments(assignments.Data, slaveRequests)
			sendLightsCh <- lightsOn
		}
	}
}

func removeTimeouts[T types.LivingMessage](pending map[int]T) {
	for updateNr, ho := range pending {
		if time.Since(ho.GetCreationTime()) > config.Cfg.AckTimeout {
			fmt.Println("Dropping order:", updateNr)
			delete(pending, updateNr)
		}
	}
}

func lightsFromAssignments(assignments map[string][N_FLOORS][2]bool, slaveRequests [N_FLOORS][2]bool) [N_FLOORS][N_BUTTONS]bool {
	var lightsOn [N_FLOORS][N_BUTTONS]bool
	for _, assignment := range assignments {
		for i := range N_FLOORS {
			for j := range 2 {
				if assignment[i][j] {
					lightsOn[i][j] = true
				}
			}
		}
	}
	return lightsOn
}

func createHallOrder(id string, order types.Order, messageCount int) types.HallOrder {
	idInt, _ := strconv.Atoi(id)

	return types.HallOrder{
		Floor:     order.Floor,
		Direction: int(order.Type),
		CreatedAt: time.Now(),
		UpdateNr:  idInt*1000000 + messageCount,
	}
}

func createFinishedAssignments(id string, orders []types.Order, messageCount int) types.FinishedHallAssignments {
	idInt, _ := strconv.Atoi(id)
	sendToMaster := types.FinishedHallAssignments{
		UpdateNr:  idInt*1000000 + messageCount,
		CreatedAt: time.Now(),
		Orders:    make([]types.Order, len(orders)),
	}

	for i, request := range orders {
		sendToMaster.Orders[i] = types.Order{
			Floor: request.Floor,
			Type:  types.OrderType(request.Type),
		}
	}

	return sendToMaster
}

// func RunSlaveBrain(id string) {

// 	var readyToSendOrder bool = true
// 	var messageCount int = 0
// 	// Recive from elevatorManager, send to master.

// 	receiveOrdersCh := make(chan types.Order)
// 	receiveFinishedAssignmentsCh := make(chan []elevio.ButtonEvent)
// 	receiveElevatorState := make(chan types.ElevatorState)

// 	sendAssignmentsCh := make(chan [N_FLOORS][N_BUTTONS]bool)
// 	sendElevatorState := make(chan types.ElevatorState)

// 	go elevatormanager.ElevatorManager(receiveElevatorState, receiveOrdersCh, receiveFinishedAssignmentsCh, sendAssignmentsCh)

// 	//Init Order ack
// 	sendOrdersCh := make(chan types.HallOrder)
// 	hallOrderAck := make(chan types.HallOrderAck)

// 	orderSender := &network.GenericSender[types.HallOrder, types.HallOrderAck]{
// 		SendCh:     sendOrdersCh,
// 		AckIn:      hallOrderAck,
// 		AckResults: make(chan network.AckResult, 10), // buffered
// 	}

// 	// Finished Assignments setup
// 	sendFinishedAssignmentsCh := make(chan types.FinishedHallAssignments)
// 	finishedOrdersAckCh := make(chan types.FinishedHallAssignmentsAck)

// 	// completeAssignmentSender := &network.GenericSender[types.FinishedHallAssignments, types.FinishedHallAssignmentsAck]{
// 	// 	SendCh:     sendFinishedAssignmentsCh,
// 	// 	AckIn:      finishedOrdersAckCh,
// 	// 	AckResults: make(chan network.AckResult, 10), // buffered
// 	// }

// 	go bcast.Transmitter(config.Cfg.MasterListenPort, sendElevatorState, sendOrdersCh, sendFinishedAssignmentsCh)

// 	receiveAssignmentsFromMasterCh := make(chan types.Assignements) //Denne skal vel egentlig bli passet som funksjonsparameter

// 	go bcast.Receiver(config.Cfg.SlaveListenPort, receiveAssignmentsFromMasterCh, hallOrderAck, finishedOrdersAckCh)

// 	var slaveRequests [N_FLOORS][N_BUTTONS]bool

// 	for {
// 		select {

// 		case state := <-receiveElevatorState:
// 			state.ID = id
// 			sendElevatorState <- state

// 		case order := <-receiveOrdersCh:
// 			if order.Type == types.Cab {
// 				slaveRequests[order.Floor][order.Type] = true //TODO: have to save this somewhere if the elevator dies and is revived
// 				sendAssignmentsCh <- slaveRequests
// 			} else if readyToSendOrder {
// 				messageCount += 1
// 				//readyToSendOrder = false
// 				idtoInt, _ := strconv.Atoi(id)
// 				ho := types.HallOrder{
// 					Floor:     order.Floor,
// 					Direction: int(order.Type),
// 					Timestamp: time.Now(),
// 					UpdateNr:  idtoInt*1000000 + messageCount,
// 				}
// 				orderSender.SendAsyncWithAck(ho)
// 			}

// 		case ackResult := <-orderSender.AckResults:
// 			if ackResult.Err != nil {
// 				fmt.Printf("ORDER WAS NOT ACKED BY MASTER: %s\n", ackResult.Err)
// 			}
// 			// //readyToSendOrder = true

// 			// TODO: Send order to master through network (Lars åssen bruker jeg network generic senderen din?)
// 			// Need to agree on format
// 			// Er også mulig å kjøre requests[order.Floor][order.Button] = true
// 			// For så å sende?? Dette blir nok buggy siden heisen kanskje tar requesten med en gang i så fall.
// 		case /*finishedOrders := */<-receiveFinishedAssignmentsCh:
// 			// fmt.Println("DEN GÅR GJENNOM FSM!!")
// 			// idtoInt, _ := strconv.Atoi(id)
// 			// messageCount += 1
// 			// sendToMaster := types.FinishedHallAssignments{
// 			// 	UpdateNr:  idtoInt*1000000 + messageCount,
// 			// 	Timestamp: time.Now(),
// 			// 	Orders:    make([]types.Order, len(finishedOrders)),
// 			// }

// 			// for i, request := range finishedOrders {
// 			// 	slaveRequests[request.Floor][request.Button] = false
// 			// 	sendToMaster.Orders[i] = types.Order{
// 			// 		Floor: request.Floor,
// 			// 		Type:  types.OrderType(request.Button),
// 			// 	}
// 			// }

// 			// completeAssignmentSender.SendAsyncWithAck(sendToMaster)

// 			// TODO: enten sende hele requests eller bare sende endringen videre til master
// 			// Hvem vet hva som er best
// 		case /*assignments :=*/ <-receiveAssignmentsFromMasterCh:
// 			// fmt.Println("Receved assignements. Doing the work")
// 			// // TODO: kombinere hallrequests og cabrequest før man sender
// 			// fmt.Println(assignments.Data)
// 			// // slaveRequests = assignments
// 			// for i := range slaveRequests {
// 			// 	slaveRequests[i][0] = assignments.Data[id][i][0]
// 			// 	slaveRequests[i][1] = assignments.Data[id][i][1]
// 			// }
// 			// // TODO: turn on lights of other
// 			// sendAssignmentsCh <- slaveRequests
// 		}
// 	}
// }
