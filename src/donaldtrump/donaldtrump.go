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

type masterData struct {
	hallRequests    [N_FLOORS][2]bool
	states          map[string]types.ElevatorState
	cabRequests     map[string][N_FLOORS]bool
	timeSinceUpdate map[string]time.Time
}

func RunMasterBrain(id string, isMasterCh chan bool) {
	masterData := masterData{
		hallRequests:    [N_FLOORS][2]bool{{false, false}, {false, false}, {false, false}, {false, false}},
		states:          make(map[string]types.ElevatorState),
		cabRequests:     make(map[string][N_FLOORS]bool),
		timeSinceUpdate: make(map[string]time.Time), //TODO - Trenger vi denne?
	}
	var isMaster = false

	//Order calculation
	ordersCh := make(chan ordermanager.HRAInput)
	calculatedAssignementsCh := make(chan map[string][N_FLOORS][2]bool, 10)
	go ordermanager.ManageOrders(ordersCh, calculatedAssignementsCh)

	//reciving channel
	updateStreamCh := make(chan types.ElevatorState)
	receiveElevatorOrdersCh := make(chan types.OrderEnvelope)
	reciveAssignmentComplete := make(chan types.FinishedHallAssignments)
	go bcast.Receiver(config.Cfg.MasterListenPort, updateStreamCh, receiveElevatorOrdersCh, reciveAssignmentComplete)

	//Sending channel
	sendAssignemnetsCh := make(chan types.Assignments)
	ackAssignementCompleted := make(chan types.FinishedHallAssignmentsAck)
	sendOrderAckCh := make(chan types.OrderAck, 10)
	go bcast.Transmitter(config.Cfg.SlaveListenPort, sendAssignemnetsCh, sendOrderAckCh, ackAssignementCompleted)

	for {
		if !isMaster { //Bro sverre ikke se her
			select {
			case isMaster = <-isMasterCh:
			case <-receiveElevatorOrdersCh:
			case <-reciveAssignmentComplete:
			case <-calculatedAssignementsCh:
			case <-updateStreamCh:
			default:
			}
			continue
		}

		select {
		case isMaster = <-isMasterCh:

		case orderReceived := <-receiveElevatorOrdersCh:
			fmt.Println("Reciving order, sending ack")
			sendOrderAckCh <- types.OrderAck{UpdateNr: orderReceived.UpdateNr}

			if orderReceived.Order.Type == types.Cab {
				arr := masterData.cabRequests[orderReceived.ElevatorID]
				arr[orderReceived.Order.Floor] = true
				masterData.cabRequests[orderReceived.ElevatorID] = arr
			} else {
				masterData.hallRequests[orderReceived.Order.Floor][orderReceived.Order.Type] = true
			}

			ordersCh <- ordermanager.ToHRAInput(masterData.hallRequests, masterData.cabRequests, masterData.states) //Loopes back to case

		case completedAssignments := <-reciveAssignmentComplete:
			for _, order := range completedAssignments.Orders {
				if order.Type == types.Cab {
					arr := masterData.cabRequests[completedAssignments.ElevatorID]
					arr[order.Floor] = false
					masterData.cabRequests[completedAssignments.ElevatorID] = arr
				} else {
					masterData.hallRequests[order.Floor][order.Type] = false
					fmt.Printf("Assignment cleared, floor: %d, type: %d \n", order.Floor, order.Type)
				}
			}

			ackAssignementCompleted <- types.FinishedHallAssignmentsAck{
				UpdateNr: completedAssignments.UpdateNr,
			}

			ordersCh <- ordermanager.ToHRAInput(masterData.hallRequests, masterData.cabRequests, masterData.states) //TODO: ER dette lurt? Bør sette en if(hasChanged) greie her for å ikke spamme.

		case assignment := <-calculatedAssignementsCh:
			fmt.Println("Sending back assignment: ")

			assignmenetOut := make(map[string][N_FLOORS][N_BUTTONS]bool)

			for id, _ := range assignment {
				arr := [N_FLOORS][N_BUTTONS]bool{}
				for i := 0; i < N_FLOORS; i++ {
					arr[i] = [N_BUTTONS]bool{assignment[id][i][types.HallUp], assignment[id][i][types.HallDown], masterData.cabRequests[id][i]}
				}
				assignmenetOut[id] = arr

			}

			sendAssignemnetsCh <- types.Assignments{Assignments: assignmenetOut}

		case elevatorData := <-updateStreamCh:
			masterData.states[elevatorData.ID] = elevatorData
			if elevatorData.Floor == -1 {
				continue
			}
			fmt.Println("Recived data from: ", elevatorData.ID)
		}

	}
}
