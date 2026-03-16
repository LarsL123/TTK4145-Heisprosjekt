package donaldtrump

import (
	"Network-go/network/bcast"
	"elevatorproject/src/config"
	elevatormanager "elevatorproject/src/elevatorManager"
	"elevatorproject/src/elevio"
	"elevatorproject/src/types"
	"fmt"
)

func RunSlaveBrain(id string) {
	// Recive from elevatorManager, send to master.

	receiveOrdersCh := make(chan elevio.ButtonEvent)
	receiveFinishedOrderCh := make(chan []elevio.ButtonEvent)

	sendAssignmentsCh := make(chan [N_FLOORS][N_BUTTONS]bool)
	receiveElevatorState := make(chan types.ElevatorState)

	go elevatormanager.ElevatorManager(receiveElevatorState, receiveOrdersCh, receiveFinishedOrderCh, sendAssignmentsCh)

	sendElevatorState := make(chan types.ElevatorState)
	go bcast.Transmitter(config.Cfg.MasterListenPort, sendElevatorState)

	receiveAssignmentsFromMasterCh := make(chan types.Assignments) //Denne skal vel egentlig bli passet som funksjonsparameter

	go bcast.Receiver(config.Cfg.SlaveListenPort, receiveAssignmentsFromMasterCh )

	var slaveRequests [N_FLOORS][N_BUTTONS]bool

	for {
		select {
		case state := <-receiveElevatorState:
			state.ID = id
			sendElevatorState <- state
		case /*order :=*/ <-receiveOrdersCh:
			// TODO: Send order to master through network (Lars åssen bruker jeg network generic senderen din?)
			// Need to agree on format
			// Er også mulig å kjøre requests[order.Floor][order.Button] = true
			// For så å sende?? Dette blir nok buggy siden heisen kanskje tar requesten med en gang i så fall.
		case finishedOrders := <-receiveFinishedOrderCh:
			for _, request := range finishedOrders {
				slaveRequests[request.Floor][request.Button] = false
			}
			// TODO: enten sende hele requests eller bare sende endringen videre til master
			// Hvem vet hva som er best
		case assignments := <-receiveAssignmentsFromMasterCh:
			fmt.Println("Receved assignements. Doing the work")
			// TODO: kombinere hallrequests og cabrequest før man sender
			fmt.Println(assignments.Data)
			// slaveRequests = assignments


			for _, 

			sendAssignmentsCh <- assignments.Data[id]
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
