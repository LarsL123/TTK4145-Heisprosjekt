package donaldtrump

import (
	"Network-go/network/bcast"
	"elevatorproject/src/config"
	elevatormanager "elevatorproject/src/elevatorManager"
	"elevatorproject/src/elevio"
	"elevatorproject/src/network"
	"elevatorproject/src/types"
	"fmt"
	"strconv"
	"time"
)

func RunSlaveBrain(id string) {

	var readyToSendOrder bool = true
	var messageCount int = 0
	// Recive from elevatorManager, send to master.

	receiveOrdersFromElevCh := make(chan types.Order)
	receiveFinishedAssignmentsCh := make(chan []elevio.ButtonEvent)
	receiveElevatorState := make(chan types.ElevatorState)

	sendAssignmentsToElevCh := make(chan [N_FLOORS][N_BUTTONS]bool)
	sendElevatorState := make(chan types.ElevatorState)


	go elevatormanager.ElevatorManager(receiveElevatorState, receiveOrdersFromElevCh, receiveFinishedAssignmentsCh, sendAssignmentsToElevCh)



	sendOrdersToMasterCh := make(chan types.HallOrder)
	hallOrderAckCh := make(chan types.HallOrderAck)
	sendOrdersFromElevatorch := make(chan types.HallOrder)
	sendFinishedAssignmentsCh := make(chan types.FinishedHallAssignments)



	go network.SendOrdersWithAck(sendOrdersFromElevatorch, sendOrdersToMasterCh,hallOrderAckCh)

	
	// //-------------Init Order ack (generic sender) ------------------- 
	// sendOrdersCh := make(chan types.HallOrder)
	// hallOrderAck := make(chan types.HallOrderAck)

	// orderSender := &network.GenericSender[types.HallOrder, types.HallOrderAck]{
	// 	SendCh:     sendOrdersCh,
	// 	AckIn:      hallOrderAck,
	// 	AckResults: make(chan network.AckResult, 10), // buffered
	// }
	// orderSender.StartAckDispatcher()



	// // Finished Assignments setup
	// sendFinishedAssignmentsCh := make(chan types.FinishedHallAssignments)
	// finishedOrdersAckCh := make(chan types.FinishedHallAssignmentsAck)

	// completeAssignmentSender := &network.GenericSender[types.FinishedHallAssignments, types.FinishedHallAssignmentsAck]{
	// 	SendCh:     sendFinishedAssignmentsCh,
	// 	AckIn:      finishedOrdersAckCh,
	// 	AckResults: make(chan network.AckResult, 10), // buffered
	// }
	// completeAssignmentSender.StartAckDispatcher()





	go bcast.Transmitter(config.Cfg.MasterListenPort, sendElevatorState, sendOrdersToMasterCh, sendFinishedAssignmentsCh)

	receiveAssignmentsFromMasterCh := make(chan types.Assignements) //Denne skal vel egentlig bli passet som funksjonsparameter


	// TODO: Daniel fortsett her imorra med å legge inn ferdig ack
	go bcast.Receiver(config.Cfg.SlaveListenPort, receiveAssignmentsFromMasterCh, hallOrderAckCh, finishedOrdersAckCh)

	var slaveRequests [N_FLOORS][N_BUTTONS]bool

	for {
		select {

		case state := <-receiveElevatorState:
			state.ID = id
			sendElevatorState <- state

		case order := <-receiveOrdersFromElevCh:
			if order.Type == types.Cab {
				slaveRequests[order.Floor][order.Type] = true //TODO: have to save this somewhere if the elevator dies and is revived
				sendAssignmentsToElevCh <- slaveRequests
			} else if readyToSendOrder {
				messageCount += 1
				//readyToSendOrder = false
				idtoInt, _ := strconv.Atoi(id)
				ho := types.HallOrder{
					Floor:     order.Floor,
					Direction: int(order.Type),
					Timestamp: time.Now(),
					UpdateNr:  idtoInt*1000000 + messageCount,
				}
				orderSender.SendAsyncWithAck(ho)
			}

		case ackResult := <-orderSender.AckResults:
			if ackResult.Err != nil {
				fmt.Printf("ORDER WAS NOT ACKED BY MASTER: %s\n", ackResult.Err)
			}
			readyToSendOrder = true

			// TODO: Send order to master through network (Lars åssen bruker jeg network generic senderen din?)
			// Need to agree on format
			// Er også mulig å kjøre requests[order.Floor][order.Button] = true
			// For så å sende?? Dette blir nok buggy siden heisen kanskje tar requesten med en gang i så fall.
		case finishedOrders := <-receiveFinishedAssignmentsCh:
			fmt.Println("DEN GÅR GJENNOM FSM!!")
			idtoInt, _ := strconv.Atoi(id)
			messageCount += 1
			sendToMaster := types.FinishedHallAssignments{
				UpdateNr:  idtoInt*1000000 + messageCount,
				Timestamp: time.Now(),
				Orders:    make([]types.Order, len(finishedOrders)),
			}

			for i, request := range finishedOrders {
				slaveRequests[request.Floor][request.Button] = false
				sendToMaster.Orders[i] = types.Order{
					Floor: request.Floor,
					Type:  types.OrderType(request.Button),
				}
			}

			completeAssignmentSender.SendAsyncWithAck(sendToMaster)

			// TODO: enten sende hele requests eller bare sende endringen videre til master
			// Hvem vet hva som er best
		case assignments := <-receiveAssignmentsFromMasterCh:
			fmt.Println("Receved assignements. Doing the work")
			// TODO: kombinere hallrequests og cabrequest før man sender
			fmt.Println(assignments.Data)
			// slaveRequests = assignments
			for i := range slaveRequests {
				slaveRequests[i][0] = assignments.Data[id][i][0]
				slaveRequests[i][1] = assignments.Data[id][i][1]
			}
			// TODO: turn on lights of other
			sendAssignmentsToElevCh <- slaveRequests
		}
	}
}

// func RunSlaveBrain(id string){

// 			for {
// 				select {
// 				case <-time.After(1000 *time.Millisecond): //Simulates updating orders.

// 						count++
// 						msg := network.OrdersAndStateUpdate{
// 							SourceId: id,
// 							UpdateNr: count,
// 							OrdersAndState: "Moren din er mann",
// 						}

// 						orderSender.UpdateAsyncGeneric(msg)

// 				case res := <-orderSender.AckResults:

// 					//Check res.UpdateNr to check that the ACk is from the latest ctate change and is not old.
// 					if res.Err != nil {
// 						fmt.Println("Failed:", res.Err)
// 					} else {
// 						fmt.Println("ACK received for", res.UpdateNr)
// 					}
// 				}
// 			}
// }
