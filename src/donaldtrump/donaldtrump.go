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
	hallRequests              [N_FLOORS][2]bool
	states                    map[string]types.ElevatorState
	cabRequests               map[string][N_FLOORS]bool
	timeSinceAssignmentUpdate [N_FLOORS][2]types.AssignedToAtTime
}

func RunMasterBrain(id string, isMasterCh chan bool) {
	masterData := masterData{
		hallRequests: [N_FLOORS][2]bool{{false, false}, {false, false}, {false, false}, {false, false}},
		states:       make(map[string]types.ElevatorState),
		cabRequests:  make(map[string][N_FLOORS]bool),
		//timeSinceAssignmentUpdate: [N_FLOORS][2]types.AssignedToAtTime, //TODO - Trenger vi denne? idk, den blir default assigned til "" og jesu fødsel, så i guess det går fint
	}
	var isMaster = false

	//Order calculation
	ordersCh := make(chan ordermanager.HRAInput)
	calculatedAssignmentsCh := make(chan map[string][N_FLOORS][2]bool, 10)
	go ordermanager.ManageOrders(ordersCh, calculatedAssignmentsCh)

	//reciving channel
	updateStreamCh := make(chan types.ElevatorState)
	receiveElevatorOrdersCh := make(chan types.OrderEnvelope)
	receiveAssignmentComplete := make(chan types.FinishedHallAssignments)
	go bcast.Receiver(config.Cfg.MasterListenPort, updateStreamCh, receiveElevatorOrdersCh, receiveAssignmentComplete)

	//Sending channel
	sendAssignmentsCh := make(chan types.Assignments)
	ackAssignmentCompleted := make(chan types.FinishedHallAssignmentsAck)
	sendOrderAckCh := make(chan types.OrderAck, 10)
	go bcast.Transmitter(config.Cfg.SlaveListenPort, sendAssignmentsCh, sendOrderAckCh, ackAssignmentCompleted)

	for {
		if !isMaster { //Bro sverre ikke se her
			select {
			case isMaster = <-isMasterCh:
			case <-receiveElevatorOrdersCh:
			case <-receiveAssignmentComplete:
			case <-calculatedAssignmentsCh:
			case <-updateStreamCh:
			default:
			}
			continue
		}

		select {
		case isMaster = <-isMasterCh:

		case orderReceived := <-receiveElevatorOrdersCh:
			fmt.Println("Receiving order, sending ack")
			sendOrderAckCh <- types.OrderAck{UpdateNr: orderReceived.UpdateNr}

			if orderReceived.Order.Type == types.Cab {
				arr := masterData.cabRequests[orderReceived.ElevatorID]
				arr[orderReceived.Order.Floor] = true
				masterData.cabRequests[orderReceived.ElevatorID] = arr
			} else {
				masterData.hallRequests[orderReceived.Order.Floor][orderReceived.Order.Type] = true
			}

			ordersCh <- ordermanager.ToHRAInput(masterData.hallRequests, masterData.cabRequests, masterData.states) //Loops back to case

		case completedAssignments := <-receiveAssignmentComplete:
			dataChanged := false

			for _, order := range completedAssignments.Orders {
				if order.Type == types.Cab {
					arr := masterData.cabRequests[completedAssignments.ElevatorID]
					if arr[order.Floor] {
						dataChanged = true
					}
					arr[order.Floor] = false
					masterData.cabRequests[completedAssignments.ElevatorID] = arr
				} else {
					if masterData.hallRequests[order.Floor][order.Type] {
						dataChanged = true
					}
					fmt.Printf("Clearing order: Floor: %d, Dirn: %s Assigned to: %s \n", order.Floor, types.OrderTypeToString(order.Type), masterData.timeSinceAssignmentUpdate[order.Floor][order.Type].ElevatorId)
					masterData.timeSinceAssignmentUpdate[order.Floor][order.Type] = types.AssignedToAtTime{
						ElevatorId: "",
						TimeStamp:  time.Now(),
					}

					masterData.hallRequests[order.Floor][order.Type] = false
					fmt.Printf("Assignment cleared, floor: %d, type: %d \n", order.Floor, order.Type)

				}
			}

			ackAssignmentCompleted <- types.FinishedHallAssignmentsAck{
				UpdateNr: completedAssignments.UpdateNr,
			}
			if dataChanged {
				ordersCh <- ordermanager.ToHRAInput(masterData.hallRequests, masterData.cabRequests, masterData.states)
			}

		case assignments := <-calculatedAssignmentsCh:
			fmt.Println("Sending back assignment: ")

			assignmentOut := make(map[string][N_FLOORS][N_BUTTONS]bool)

			// Looping and creating map of assignments to send out to elevators
			for id, _ := range assignments {
				arr := [N_FLOORS][N_BUTTONS]bool{}
				for i := 0; i < N_FLOORS; i++ {
					arr[i] = [N_BUTTONS]bool{assignments[id][i][types.HallUp], assignments[id][i][types.HallDown], masterData.cabRequests[id][i]}

					// Checking if last assigned elevator to assignment is not the same and updating timeStamp if there is a new elevator assigned to it
					if masterData.timeSinceAssignmentUpdate[i][types.HallUp].ElevatorId != id && assignments[id][i][types.HallUp] {
						masterData.timeSinceAssignmentUpdate[i][types.HallUp] = types.AssignedToAtTime{
							ElevatorId: id,
							TimeStamp:  time.Now(),
						}
					}
					if masterData.timeSinceAssignmentUpdate[i][types.HallDown].ElevatorId != id && assignments[id][i][types.HallDown] {
						masterData.timeSinceAssignmentUpdate[i][types.HallDown] = types.AssignedToAtTime{
							ElevatorId: id,
							TimeStamp:  time.Now(),
						}
					}
				}
				assignmentOut[id] = arr
			}

			sendAssignmentsCh <- types.Assignments{Assignments: assignmentOut}

		case elevatorData := <-updateStreamCh:
			masterData.states[elevatorData.ID] = elevatorData
			if elevatorData.Floor == -1 {
				continue
			}
			fmt.Println("Received data from: ", elevatorData.ID)

			// Suspending orders over maxOrderSuspendTime
			for floor := range N_FLOORS {
				for orderType := range 2 {
					currentorder := masterData.timeSinceAssignmentUpdate[floor][orderType]
					if time.Since(currentorder.TimeStamp) > config.Cfg.MaxOrderSuspendTime && currentorder.ElevatorId != "" {

						// Suspend the correct elevator
						//suspendElevator(masterData.timeSinceAssignmentUpdate[floor][orderType].ElevatorId)
						currentElevId := masterData.timeSinceAssignmentUpdate[floor][orderType].ElevatorId
						tempState := masterData.states[currentElevId]
						tempState.Suspended = types.SuspendedType{
							IsSuspended: true,
							TimeStamp:   time.Now(),
						}
						masterData.states[currentElevId] = tempState
						fmt.Printf("Suspended elevator: %s\n", currentElevId)
					}
				}
			}
		}
	}
}

func suspendElevator(elevatorId string, masterData masterData) {

}
