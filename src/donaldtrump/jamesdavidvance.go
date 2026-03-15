package donaldtrump

import (
	"Network-go/network/bcast"
	"elevatorproject/src/config"
	elevatormanager "elevatorproject/src/elevatorManager"
	"elevatorproject/src/elevio"
	"elevatorproject/src/types"
)

func RunSlaveBrain(id string){
	// Recive from elevatorManager, send to master. 

	receiveOrdersCh := make(chan elevio.ButtonEvent)
	receiveFinishedOrderCh := make(chan []elevio.ButtonEvent)

	sendAssignmentsCh := make(chan [N_FLOORS][N_BUTTONS]bool)
	receiveElevatorState := make(chan types.ElevatorState)

	go elevatormanager.ElevatorManager(receiveElevatorState,receiveOrdersCh, receiveFinishedOrderCh, sendAssignmentsCh)

	sendElevatorState := make(chan types.ElevatorState)
	go bcast.Transmitter(config.Cfg.MasterListenPort, sendElevatorState)

	for {
		select{
		case state :=  <- receiveElevatorState:
			state.ID = id
			sendElevatorState <- state

		case <-receiveOrdersCh:
		case <- receiveFinishedOrderCh:
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