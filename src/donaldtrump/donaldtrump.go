package donaldtrump

import (
	"Network-go/network/bcast"
	"elevatorproject/src/config"
	"elevatorproject/src/ordermanager"
	"elevatorproject/src/types"
	"fmt"
	"time"
)

/*

 Im gonne create the greates eleavtor youll ever see. Its gonne go to 1000 floors, yeah it can go to the moon.
 It will be the greatest and fastest elevator in the history of elevators.
 And we gonne make china pay for it.

*/

const N_FLOORS = 4
const N_BUTTONS = 3

// Set test input
// input := HRAInput{
//     HallRequests: [4][2]bool{{false, false}, {true, false}, {false, false}, {false, true}},
//     States: map[string]HRAElevState{
//         "one": {
//             Behavior:       "moving",
//             Floor:          3,
//             Direction:      "down",
//             CabRequests:    [4]bool{false, false, false, true},
//         },
//         "two": {
//             Behavior:       "idle",
//             Floor:          0,
//             Direction:      "stop",
//             CabRequests:    [4]bool{false, false, false, false},
//         },
//     },
// }

type masterData struct {
	hallRequests    [4][2]bool
	states          map[string]types.ElevatorState
	timeSinceUpdate map[string]time.Time
}

func RunMasterBrain(id string) {
	masterData := masterData{
		hallRequests:    [4][2]bool{{false, false}, {false, false}, {false, false}, {false, false}},
		states:          make(map[string]types.ElevatorState),
		timeSinceUpdate: make(map[string]time.Time),
	}

	ordersCh := make(chan ordermanager.HRAInput)
	calculatedAssignementsCh := make(chan map[string][4][2]bool)
	go ordermanager.ManageOrders(ordersCh, calculatedAssignementsCh)

	stateUpdateCh := make(chan types.ElevatorState)

	receiveElevatorOrdersCh := make(chan types.HallOrder)
	sendOrderAckCh := make(chan types.HallOrderAck)

	reciveAssignmentComplete := make(chan types.FinishedHallAssignments)
	go bcast.Receiver(config.Cfg.MasterListenPort, stateUpdateCh, receiveElevatorOrdersCh, reciveAssignmentComplete)

	sendAssignemnetsCh := make(chan types.Assignements)
	ackAssignementCompleted := make(chan types.FinishedHallAssignmentsAck)

	go bcast.Transmitter(config.Cfg.SlaveListenPort, sendAssignemnetsCh, sendOrderAckCh, ackAssignementCompleted)

	for {
		select {
		case elevatorData := <-stateUpdateCh:
			masterData.states[elevatorData.ID] = elevatorData
			if elevatorData.Floor == -1 {
				continue
			}
			fmt.Println("Recived data from: ", elevatorData.ID)
			//ordersCh <- ordermanager.ToHRAInput(masterData.hallRequests, masterData.states)

		case orderReceived := <-receiveElevatorOrdersCh:
			fmt.Println("Reciving order")
			sendOrderAckCh <- types.HallOrderAck{UpdateNr: orderReceived.GetUpdateNr()}
			masterData.hallRequests[orderReceived.Floor][orderReceived.Direction] = true
			ordersCh <- ordermanager.ToHRAInput(masterData.hallRequests, masterData.states)

		case completedAssignments := <-reciveAssignmentComplete:
			for _, order := range completedAssignments.Orders {
				if order.Type == types.Cab {
					continue
				}
				masterData.hallRequests[order.Floor][order.Type] = false
				fmt.Printf("Assignment cleared, floor: %d, type: %d ", order.Floor, order.Type)
			}

			ackAssignementCompleted <- types.FinishedHallAssignmentsAck{
				UpdateNr: completedAssignments.GetUpdateNr(),
			}

		case assignment := <-calculatedAssignementsCh:
			fmt.Println(assignment)
			fmt.Println("Sending back")
			sendAssignemnetsCh <- types.Assignements{Data: assignment}

		}

	}
}

// func RunMasterBrain(id string){

// 	ordersCh := make(chan ordermanager.HRAInput)
// 	assignmentsCh := make(chan map[string][][2]bool)
// 	go ordermanager.ManageOrders(ordersCh, assignmentsCh)

// 	isMaster := make(chan bool)
// 	slaveUpdate := network.StartMaster(id, isMaster)

// 	receiveOrdersAndStateUpdateCh := make(chan network.OrdersAndStateUpdate)
// 	orderAndStateAckCh := make(chan network.OrdersAndStateAck)

// 	sendAssignmentCh := make(chan network.AssignmentsAndOrders)
// 	assignmentAckCh := make(chan network.AssignementsAndOrdersAck)

// 	go bcast.Receiver(config.Cfg.MasterListenPort,receiveOrdersAndStateUpdateCh, assignmentAckCh)
// 	go bcast.Transmitter(config.Cfg.SlaveListenPort, orderAndStateAckCh, sendAssignmentCh)

// 	assignmentSender := &network.GenericSender[network.AssignmentsAndOrders, network.AssignementsAndOrdersAck]{
// 		SendCh: sendAssignmentCh,
// 		AckIn: assignmentAckCh,
// 		AckResults: make(chan network.AckResult, 10), // buffered OBS-OBS!! DO i need this??
// 	}

// 	msg := network.AssignmentsAndOrders{
// 				SourceId: id,
// 				UpdateNr: 1,
// 				OrdersAndState: "Ice will come to your home",
// 			}
// 	assignmentSender.UpdateAsyncGeneric(msg)

// 	for {
// 		select{
// 		case p := <-slaveUpdate:
// 			fmt.Printf("Slave update:\n")
// 			fmt.Printf("  Slaves:    %q\n", p.Slaves)
// 			fmt.Printf("  New:      %q\n", p.New)
// 			fmt.Printf("  Lost:     %q\n", p.Lost)

// 		case data := <- receiveOrdersAndStateUpdateCh: //Constant ack
// 			fmt.Printf("Received from slave: %s \n", data.OrdersAndState)

// 		case <-time.After(6*time.Second): //New assignment to be distrebuted.
// 			fmt.Println("Sending new assignments. ")

// 			msg := network.AssignmentsAndOrders{
// 				SourceId: id,
// 				UpdateNr: 1,
// 				OrdersAndState: "Ice will come to your home",
// 			}

// 			assignmentSender.UpdateAsyncGeneric(msg)
// 		}
// 	}
// }
