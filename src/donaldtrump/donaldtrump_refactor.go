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

type Master struct {
	id       string
	isMaster bool

	// State
	data masterData

	// Channels
	isMasterCh chan bool

	calculateAssignmentsCh     chan ordermanager.HRAInput
	rawAssignmentsCh           chan map[string][N_FLOORS][2]bool

	updateStreamCh             chan types.ElevatorState
	receiveElevatorOrdersCh    chan types.OrderEnvelope
	completedAssignmentCh      chan types.FinishedHallAssignments

	sendAssignmentsCh          chan types.Assignments
	ackOrderCh                 chan types.OrderAck
	ackAssignmentCompletedCh   chan types.FinishedHallAssignmentsAck
}

type masterData struct {
	hallRequests              [N_FLOORS][2]bool
	states                    map[string]types.ElevatorState
	cabRequests               map[string][N_FLOORS]bool
	timeSinceAssignmentUpdate [N_FLOORS][2]types.AssignedToAtTime
}

func NewMaster(id string, isMasterCh chan bool) *Master {
	m := &Master{
		id:         id,
		isMasterCh: isMasterCh,
		isMaster: false,

		data: masterData{
			//hallRequests: [N_FLOORS][2]bool{{false, false}, {false, false}, {false, false}, {false, false}},
			states:       make(map[string]types.ElevatorState),
			cabRequests:  make(map[string][N_FLOORS]bool),
			//timeSinceAssignmentUpdate: [N_FLOORS][2]types.AssignedToAtTime, //TODO - Trenger vi denne? idk, den blir default assigned til "" og jesu fødsel, så i guess det går fint
		},

	//Order calculation
	calculateAssignmentsCh: make(chan ordermanager.HRAInput),
	rawAssignmentsCh: make(chan map[string][N_FLOORS][2]bool, 10),

	//reciving channel
	updateStreamCh: make(chan types.ElevatorState),
	receiveElevatorOrdersCh: make(chan types.OrderEnvelope),
	completedAssignmentCh: make(chan types.FinishedHallAssignments),

	//Sending channel
	sendAssignmentsCh: make(chan types.Assignments),
	ackAssignmentCompletedCh: make(chan types.FinishedHallAssignmentsAck),
	ackOrderCh: make(chan types.OrderAck, 10),
	}

	return m
}

func (m *Master) Start() {
	go ordermanager.ManageOrders(m.calculateAssignmentsCh, m.rawAssignmentsCh)

	go bcast.Receiver(
		config.Cfg.MasterListenPort,
		m.updateStreamCh,
		m.receiveElevatorOrdersCh,
		m.completedAssignmentCh,
	)

	go bcast.Transmitter(
		config.Cfg.SlaveListenPort,
		m.sendAssignmentsCh,
		m.ackAssignmentCompletedCh,
		m.ackOrderCh,
	)

	m.runLoop()

}

func (data *masterData) storeOrder(order types.Order, elevatorId string) bool{
	hasChanged := false

	if order.Type == types.Cab {
		arr := data.cabRequests[elevatorId]

		hasChanged = !arr[order.Floor]
		arr[order.Floor] = true

		data.cabRequests[elevatorId] = arr
	} else {
		hasChanged = !data.hallRequests[order.Floor][order.Type]
		data.hallRequests[order.Floor][order.Type] = true
	}

	return hasChanged
}

func (data *masterData) removeOrders(orders []types.Order, elevatorID string) bool{
	hasChanged := false

	for _, order := range orders {
		if order.Type == types.Cab {
			arr := data.cabRequests[elevatorID]
			if arr[order.Floor] { //TODO: Spøre gutta om de liker denne eller den uten if
				hasChanged = true
			}
			arr[order.Floor] = false
			data.cabRequests[elevatorID] = arr
		} else {
			if data.hallRequests[order.Floor][order.Type] {
				hasChanged = true
			}
			data.hallRequests[order.Floor][order.Type] = false
			fmt.Printf("Assignment cleared, floor: %d, type: %d \n", order.Floor, order.Type)
		}
	}
	return hasChanged
}

func (m *Master) runLoop(){

	for {
		if !m.isMaster { 	
			m.drainChannels()
			continue
		}

		select {	
			case m.isMaster = <-m.isMasterCh:

			case orderReceived := <-m.receiveElevatorOrdersCh:
				fmt.Println("Receiving order, sending ack")
				m.ackOrderCh <- types.OrderAck{UpdateNr: orderReceived.UpdateNr}

				hasChanged := m.data.storeOrder(orderReceived.Order, orderReceived.ElevatorID)
				if hasChanged{
					m.calculateAssignmentsCh <- ordermanager.ToHRAInput(m.data.hallRequests, m.data.cabRequests, m.data.states) //Loops back to case
				}

			case completedAssignments := <-m.completedAssignmentCh:

				m.ackAssignmentCompletedCh <- types.FinishedHallAssignmentsAck{
					UpdateNr: completedAssignments.UpdateNr,
				}
			
				orders := completedAssignments.Orders
				dataChanged := m.data.removeOrders(orders, completedAssignments.ElevatorID)

				for _, order := range orders { //TODO: Skrive om til funksjone. Høre Brage og Daniel
					if order.Type != types.Cab {
						fmt.Printf("Clearing order: Floor: %d, Dirn: %s Assigned to: %s \n", order.Floor, order.Type.ToString(), m.data.timeSinceAssignmentUpdate[order.Floor][order.Type].ElevatorId)
						m.data.timeSinceAssignmentUpdate[order.Floor][order.Type] = types.AssignedToAtTime{
							ElevatorId: "",
							TimeStamp:  time.Now(),
						}
					}
				}


				if dataChanged {
					m.calculateAssignmentsCh <- ordermanager.ToHRAInput(m.data.hallRequests, m.data.cabRequests, m.data.states)
				}

			case assignments := <-m.rawAssignmentsCh:
				fmt.Println("Sending back assignment: ")

				assignmentOut := make(map[string][N_FLOORS][N_BUTTONS]bool)

				// Looping and creating map of assignments to send out to elevators
				for id, _ := range assignments {
					arr := [N_FLOORS][N_BUTTONS]bool{}
					for i := 0; i < N_FLOORS; i++ {
						arr[i] = [N_BUTTONS]bool{assignments[id][i][types.HallUp], assignments[id][i][types.HallDown], m.data.cabRequests[id][i]}

						// Checking if last assigned elevator to assignment is not the same and updating timeStamp if there is a new elevator assigned to it
						if m.data.timeSinceAssignmentUpdate[i][types.HallUp].ElevatorId != id && assignments[id][i][types.HallUp] {
							m.data.timeSinceAssignmentUpdate[i][types.HallUp] = types.AssignedToAtTime{
								ElevatorId: id,
								TimeStamp:  time.Now(),
							}
						}
						if m.data.timeSinceAssignmentUpdate[i][types.HallDown].ElevatorId != id && assignments[id][i][types.HallDown] {
							m.data.timeSinceAssignmentUpdate[i][types.HallDown] = types.AssignedToAtTime{
								ElevatorId: id,
								TimeStamp:  time.Now(),
							}
						}
					}
					assignmentOut[id] = arr
				}

				m.sendAssignmentsCh <- types.Assignments{Assignments: assignmentOut}

			case elevatorData := <-m.updateStreamCh:
				m.data.states[elevatorData.ID] = elevatorData
				if elevatorData.Floor == -1 {
					continue
				}
				fmt.Println("Received data from: ", elevatorData.ID)

				// Suspending orders over maxOrderSuspendTime
				for floor := range N_FLOORS {
					for orderType := range 2 {
						currentorder := m.data.timeSinceAssignmentUpdate[floor][orderType]
						if time.Since(currentorder.TimeStamp) > config.Cfg.MaxOrderSuspendTime && currentorder.ElevatorId != "" {

							// Suspend the correct elevator
							//suspendElevator(masterData.timeSinceAssignmentUpdate[floor][orderType].ElevatorId)
							currentElevId := m.data.timeSinceAssignmentUpdate[floor][orderType].ElevatorId
							tempState := m.data.states[currentElevId]
							tempState.Suspended = types.SuspendedType{
								IsSuspended: true,
								TimeStamp:   time.Now(),
							}
							m.data.states[currentElevId] = tempState
							fmt.Printf("Suspended elevator: %s\n", currentElevId)
						}
					}
				}
			}
	}
}

func (m *Master) drainChannels(){
	select {
		case m.isMaster = <-m.isMasterCh:
		case <-m.receiveElevatorOrdersCh:
		case <-m.completedAssignmentCh:
		case <-m.rawAssignmentsCh:
		case <-m.updateStreamCh:
		default:
	}
}


func suspendElevator(elevatorId string, masterData masterData) {
 //TODO: Hva soren er dette daniel/brage?
}
