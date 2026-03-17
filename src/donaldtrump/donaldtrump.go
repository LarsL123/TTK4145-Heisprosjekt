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
const MAX_OBSTRUCTED_TIME = 20 // [seconds]

type masterData struct {
	hallRequests    [N_FLOORS][2]bool
	states          map[string]types.ElevatorState
	timeSinceUpdate map[string]time.Time
	obstructedTime map[string]time.Time
}

func RunMasterBrain(id string) {
	masterData := masterData{
		hallRequests:    [N_FLOORS][2]bool{{false, false}, {false, false}, {false, false}, {false, false}},
		states:          make(map[string]types.ElevatorState),
		timeSinceUpdate: make(map[string]time.Time),
		obstructedTime: make(map[string]time.Time),
	}

	//Order calculation
	ordersCh := make(chan ordermanager.HRAInput)
	calculatedAssignementsCh := make(chan map[string][N_FLOORS][2]bool)
	go ordermanager.ManageOrders(ordersCh, calculatedAssignementsCh)

	//reciving channel
	updateStreamCh := make(chan types.ElevatorState)
	receiveElevatorOrdersCh := make(chan types.HallOrder)
	reciveAssignmentComplete := make(chan types.FinishedHallAssignments)
	go bcast.Receiver(config.Cfg.MasterListenPort, updateStreamCh, receiveElevatorOrdersCh, reciveAssignmentComplete)

	//Sending channel
	sendAssignemnetsCh := make(chan types.Assignments)
	ackAssignementCompleted := make(chan types.FinishedHallAssignmentsAck)
	sendOrderAckCh := make(chan types.HallOrderAck, 10)
	go bcast.Transmitter(config.Cfg.SlaveListenPort, sendAssignemnetsCh, sendOrderAckCh, ackAssignementCompleted)

	for {
		select {

		case orderReceived := <-receiveElevatorOrdersCh:
			fmt.Println("Reciving order ack back")
			sendOrderAckCh <- types.HallOrderAck{UpdateNr: orderReceived.UpdateNr}

			masterData.hallRequests[orderReceived.Floor][orderReceived.Direction] = true
			ordersCh <- ordermanager.ToHRAInput(masterData.hallRequests, masterData.states) //Loopes back to case

		case completedAssignments := <-reciveAssignmentComplete:
			for _, order := range completedAssignments.Orders {
				if order.Type == types.Cab {
					continue
				}
				masterData.hallRequests[order.Floor][order.Type] = false
				fmt.Printf("Assignment cleared, floor: %d, type: %d \n", order.Floor, order.Type)
			}

			ackAssignementCompleted <- types.FinishedHallAssignmentsAck{
				UpdateNr: completedAssignments.UpdateNr,
			}

		case assignment := <-calculatedAssignementsCh:
			fmt.Println("Sending back assignment: ")
			sendAssignemnetsCh <- types.Assignments{Data: assignment}

		case elevatorData := <-updateStreamCh:
			masterData.states[elevatorData.ID] = elevatorData
			if elevatorData.Floor == -1 {
				continue
			}
			if elevatorData.Obstructed{
				masterData.obstructedTime[elevatorData.ID] = time.Now()
			}
			obstructedTime, wasObstructed := masterData.obstructedTime[elevatorData.ID]
			if wasObstructed{
				if !elevatorData.Obstructed{
					delete(masterData.obstructedTime,elevatorData.ID)
				}
				if time.Since(obstructedTime) > MAX_OBSTRUCTED_TIME*time.Second{
					//REMOVE ELLERNOSÅNT
				}
			}
			fmt.Println("Recived data from: ", elevatorData.ID)

		}
	}
}
