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

	//Order calculation
	ordersCh := make(chan ordermanager.HRAInput)
	calculatedAssignementsCh := make(chan map[string][4][2]bool)
	go ordermanager.ManageOrders(ordersCh, calculatedAssignementsCh)

	//reciving channel
	updateStreamCh := make(chan types.ElevatorState)
	receiveElevatorOrdersCh := make(chan types.HallOrder)
	reciveAssignmentComplete := make(chan types.FinishedHallAssignments)
	go bcast.Receiver(config.Cfg.MasterListenPort, updateStreamCh, receiveElevatorOrdersCh, reciveAssignmentComplete)

	//Sending channel
	sendAssignemnetsCh := make(chan types.Assignements)
	ackAssignementCompleted := make(chan types.FinishedHallAssignmentsAck)
	sendOrderAckCh := make(chan types.HallOrderAck)
	go bcast.Transmitter(config.Cfg.SlaveListenPort, sendAssignemnetsCh, sendOrderAckCh, ackAssignementCompleted)

	assignmentTicker := time.NewTicker(1 * time.Second)
	defer assignmentTicker.Stop()

	for {
		select {

		case <-assignmentTicker.C:
			//ordersCh <- ordermanager.ToHRAInput(masterData.hallRequests, masterData.states) //Loopes back to case

		case orderReceived := <-receiveElevatorOrdersCh:
			fmt.Println("Reciving order")
			sendOrderAckCh <- types.HallOrderAck{UpdateNr: orderReceived.GetUpdateNr()}

			masterData.hallRequests[orderReceived.Floor][orderReceived.Direction] = true

		case /*completedAssignments :=*/ <-reciveAssignmentComplete:
			// for _, order := range completedAssignments.Orders {
			// 	if order.Type == types.Cab {
			// 		continue
			// 	}
			// 	masterData.hallRequests[order.Floor][order.Type] = false
			// 	fmt.Printf("Assignment cleared, floor: %d, type: %d \n", order.Floor, order.Type)
			// }

			// ackAssignementCompleted <- types.FinishedHallAssignmentsAck{
			// 	UpdateNr: completedAssignments.GetUpdateNr(),
			// }

		case /*assignment := */<-calculatedAssignementsCh:
			// fmt.Println(assignment)
			// fmt.Println("Sending back")
			// sendAssignemnetsCh <- types.Assignements{Data: assignment}

		case elevatorData := <-updateStreamCh:
			masterData.states[elevatorData.ID] = elevatorData
			if elevatorData.Floor == -1 {
				continue
			}
			fmt.Println("Recived data from: ", elevatorData.ID)

		}

	}
}

