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
	isMasterCh                          chan bool
	transferOrdersWhenMasterDowngradeCh chan types.OrderEnvelope
	forwardOrdersFromBackup             chan types.OrderEnvelope

	calculateAssignmentsCh chan ordermanager.HRAInput
	rawAssignmentsCh       chan map[string][N_FLOORS][2]bool

	updateStreamCh          chan types.ElevatorState
	receiveElevatorOrdersCh chan types.OrderEnvelope
	completedAssignmentCh   chan types.FinishedHallAssignments

	sendAssignmentsCh        chan types.Assignments
	ackOrderCh               chan types.OrderAck
	ackAssignmentCompletedCh chan types.FinishedHallAssignmentsAck

	sendBackupDataCh    chan types.BackupData
	receiveBackupAckCh  chan types.BackupDataAck
	pendingBackupOrders map[int]pendingAssignment
	backupUpdateNr      int
}

type masterData struct {
	hallRequests              [N_FLOORS][2]bool
	states                    map[string]types.ElevatorState
	cabRequests               map[string][N_FLOORS]bool
	timeSinceAssignmentUpdate [N_FLOORS][2]types.AssignedToAtTime
	suspendedElevators        map[string]types.SuspendedType
}

// Should this be moved somewhere else?
type pendingAssignment struct {
	assignments map[string][N_FLOORS][2]bool
	createdAt   time.Time
}

func NewMaster(id string, isMasterCh chan bool, transferMasterOrders chan types.OrderEnvelope, forwardOrdersFromBAckup chan types.OrderEnvelope) *Master {
	m := &Master{
		id:         id,
		isMasterCh: isMasterCh,
		isMaster:   false,

		data: masterData{
			//hallRequests: [N_FLOORS][2]bool{{false, false}, {false, false}, {false, false}, {false, false}},
			states:             make(map[string]types.ElevatorState),
			cabRequests:        make(map[string][N_FLOORS]bool),
			suspendedElevators: make(map[string]types.SuspendedType),

			//timeSinceAssignmentUpdate: [N_FLOORS][2]types.AssignedToAtTime, //TODO - Trenger vi denne? idk, den blir default assigned til "" og jesu fødsel, så i guess det går fint
		},

		transferOrdersWhenMasterDowngradeCh: transferMasterOrders,
		forwardOrdersFromBackup:             forwardOrdersFromBAckup,

		//Order calculation
		calculateAssignmentsCh: make(chan ordermanager.HRAInput),
		rawAssignmentsCh:       make(chan map[string][N_FLOORS][2]bool, 10),

		//receiving channel
		updateStreamCh:          make(chan types.ElevatorState),
		receiveElevatorOrdersCh: make(chan types.OrderEnvelope),
		completedAssignmentCh:   make(chan types.FinishedHallAssignments),
		receiveBackupAckCh:      make(chan types.BackupDataAck, 10),

		//Sending channel
		sendAssignmentsCh:        make(chan types.Assignments),
		ackAssignmentCompletedCh: make(chan types.FinishedHallAssignmentsAck),
		ackOrderCh:               make(chan types.OrderAck, 10),
		sendBackupDataCh:         make(chan types.BackupData, 10),

		// Backup
		pendingBackupOrders: make(map[int]pendingAssignment),
	}

	return m
}

func (m *Master) Start(masterAliveCh chan struct{}) {
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

	go bcast.Transmitter(config.Cfg.BackupSendPort, m.sendBackupDataCh)
	go bcast.Receiver(config.Cfg.BackupReceivePort, m.receiveBackupAckCh)

	m.runLoop(masterAliveCh)

}

func (m *Master) runLoop(aliveCh chan struct{}) {
	suspensionTicker := time.NewTicker(500 * time.Millisecond)
	defer suspensionTicker.Stop()

	resendToBackupTicker := time.NewTicker(15 * time.Millisecond) //TODO: Add to config
	defer resendToBackupTicker.Stop()
	for {
		if !m.isMaster {
			m.drainChannels()
			continue
		}

		select {
		case m.isMaster = <-m.isMasterCh:
			if m.isMaster == false { //May be redundant.
				m.pushOrdersToNewMaster()
			} else {
				m.data.suspendedElevators = make(map[string]types.SuspendedType)
				m.data.timeSinceAssignmentUpdate = [N_FLOORS][2]types.AssignedToAtTime{}
				m.runReassignment()
			}

		case <-suspensionTicker.C:
			m.removeDeadElevators()

			if m.suspendTimedOutElevators() {
				m.runReassignment()
			}
			if m.data.unsuspendElevators() {
				m.runReassignment()
			}

			//Random placed just have to called.
			select { // should probably have own ticker if we think they need other ticker rate
			case aliveCh <- struct{}{}:
			default:
			}
		case orderReceived := <-m.receiveElevatorOrdersCh:
			fmt.Println("Receiving order, sending ack")
			m.ackOrderCh <- types.OrderAck{UpdateNr: orderReceived.UpdateNr}

			hasChanged := m.data.storeOrder(orderReceived.Order, orderReceived.ElevatorID)
			if hasChanged {
				m.runReassignment() //Loops back to case
			}

		case completedAssignments := <-m.completedAssignmentCh:
			m.ackAssignmentCompletedCh <- types.FinishedHallAssignmentsAck{UpdateNr: completedAssignments.UpdateNr}

			if m.data.removeOrders(completedAssignments.Orders, completedAssignments.ElevatorID) {
				m.data.clearAssignmentTimestamps(completedAssignments.Orders)
				m.runReassignment()
			}

		case assignments := <-m.rawAssignmentsCh:
			m.backupUpdateNr++

			m.sendBackupDataCh <- types.BackupData{
				UpdateNr:     m.backupUpdateNr,
				HallRequests: m.data.hallRequests,
				CabRequests:  m.data.cabRequests,
			}

			m.pendingBackupOrders[m.backupUpdateNr] = pendingAssignment{
				assignments: assignments,
				createdAt:   time.Now(),
			}

		case ack := <-m.receiveBackupAckCh:
			pending, ok := m.pendingBackupOrders[ack.UpdateNr]
			if !ok {
				continue
			}
			delete(m.pendingBackupOrders, ack.UpdateNr)

			m.updateAssignmentTimestamps(pending.assignments)
			m.sendAssignmentsCh <- types.Assignments{Assignments: m.mergeAssignmentsWithCabRequests(pending.assignments)} // TODO: Rename this channel? Might be inaccurate

		case order := <-m.transferOrdersWhenMasterDowngradeCh:
			hasChanged := m.data.storeOrder(order.Order, order.ElevatorID)
			if hasChanged {
				m.runReassignment()
			}

		case envelope := <-m.forwardOrdersFromBackup:
			hasChanged := m.data.storeOrder(envelope.Order, envelope.ElevatorID)
			if hasChanged {
				m.runReassignment()
			}
		case <-resendToBackupTicker.C:
			for updateNr, pending := range m.pendingBackupOrders {
				if time.Since(pending.createdAt) > config.Cfg.AckTimeout {
					// Backup is dead // DANIEL CONTINUE HERE YOU MAFAKKA!!!
					if len(m.data.states) <= 1 {
						// No other elevators alive, no backup -> send without backup ack
						m.sendAssignmentsCh <- types.Assignments{
							Assignments: m.mergeAssignmentsWithCabRequests(pending.assignments),
						}
						delete(m.pendingBackupOrders, updateNr)
					} else {
						// Backup should exist -> resend if no ack was received
						fmt.Println("Backup not acking, resending backup data...")
						m.sendBackupDataCh <- types.BackupData{
							UpdateNr:     updateNr,
							HallRequests: m.data.hallRequests,
							CabRequests:  m.data.cabRequests,
						}
					}
				}
			}
		case elevatorData := <-m.updateStreamCh:
			// TODO: if message from suspended elevator, unsuspend if elevatordata != from m.data.states[elevatorid]
			if elevatorData.Floor == -1 { //TODO: Should maybe change this, might be an idea to let the elevatordata be saved but not used for reassignment if between floors
				continue
			}

			_, wasKnown := m.data.states[elevatorData.ID]
			m.data.states[elevatorData.ID] = elevatorData

			if !wasKnown {
				// New elevator connected -> reassing so it gets its orders
				fmt.Printf("New elevator %s connected -> reassigning", elevatorData.ID)
				m.runReassignment()
			}
			fmt.Println("Received data from: ", elevatorData.ID)
		}
	}
}

func (m *Master) removeDeadElevators() {
	for id, state := range m.data.states {
		if time.Since(state.CreatedAt) > config.Cfg.ElevatorDeadTimeout {
			fmt.Printf("No state received from elevator %s - removing\n", id)
			delete(m.data.states, id)
			//delete(m.data.cabRequests,id) Could probably delete this, per now they are saved two places, which is probably not good, but what can we do...
			delete(m.data.suspendedElevators, id)
			m.runReassignment()
		}
	}
}

func (d *masterData) unsuspendElevators() bool {
	changed := false
	for id, suspend := range d.suspendedElevators {
		if suspend.IsSuspended && time.Since(suspend.TimeStamp) > config.Cfg.MaxElevatorSuspendTime {
			d.suspendedElevators[id] = types.SuspendedType{
				IsSuspended: false,
				TimeStamp:   time.Now(),
			}
			changed = true
		}
	}
	return changed
}

func (m *Master) mergeAssignmentsWithCabRequests(rawAssignments map[string][N_FLOORS][2]bool) map[string][N_FLOORS][N_BUTTONS]bool {
	assignmentOut := make(map[string][N_FLOORS][N_BUTTONS]bool)

	// Looping and creating map of assignments to send out to elevators
	for id := range rawAssignments {
		arr := [N_FLOORS][N_BUTTONS]bool{}
		for i := 0; i < N_FLOORS; i++ {
			arr[i] = [N_BUTTONS]bool{rawAssignments[id][i][types.HallUp], rawAssignments[id][i][types.HallDown], m.data.cabRequests[id][i]}
		}
		assignmentOut[id] = arr
	}
	return assignmentOut
}

func (m *Master) updateAssignmentTimestamps(assignments map[string][N_FLOORS][2]bool) {
	for id := range assignments {
		for floor := range N_FLOORS {
			if m.data.timeSinceAssignmentUpdate[floor][types.HallUp].ElevatorId != id && assignments[id][floor][types.HallUp] {
				m.data.timeSinceAssignmentUpdate[floor][types.HallUp] = types.AssignedToAtTime{
					ElevatorId: id,
					TimeStamp:  time.Now(),
				}
			}

			if m.data.timeSinceAssignmentUpdate[floor][types.HallDown].ElevatorId != id && assignments[id][floor][types.HallDown] {
				m.data.timeSinceAssignmentUpdate[floor][types.HallDown] = types.AssignedToAtTime{
					ElevatorId: id,
					TimeStamp:  time.Now(),
				}
			}
		}
	}
}

func (m *Master) runReassignment() {
	input := ordermanager.ToHRAInput(m.data.hallRequests, m.data.cabRequests, m.data.states, m.data.suspendedElevators)
	select {
	case m.calculateAssignmentsCh <- input:
	default:
		//HRA is busy - resend ticker will fix
	}

}

func (d *masterData) clearAssignmentTimestamps(orders []types.Order) {
	for _, order := range orders {
		if order.Type != types.Cab {
			fmt.Printf("Clearing order: Floor: %d, Dirn: %s Assigned to: %s \n", order.Floor, order.Type.ToString(), d.timeSinceAssignmentUpdate[order.Floor][order.Type].ElevatorId)
			d.timeSinceAssignmentUpdate[order.Floor][order.Type] = types.AssignedToAtTime{
				ElevatorId: "",
				TimeStamp:  time.Now(),
			}
		}
	}
}

func (m *Master) suspendTimedOutElevators() bool {
	// Suspending elevators that have assignments over maxOrderSuspendTime (Could be moved to some other case)
	if len(m.data.states) <= 1 { // TODO: Might need to change this if m.data.states doesn't get updated when an elevator dies
		return false
	}

	var elevWasSuspended = false

	for floor := range N_FLOORS {
		for orderType := range 2 {
			cur := m.data.timeSinceAssignmentUpdate[floor][orderType]

			if time.Since(cur.TimeStamp) > config.Cfg.MaxOrderSuspendTime && cur.ElevatorId != "" {
				// Suspend the correct elevator
				id := m.data.timeSinceAssignmentUpdate[floor][orderType].ElevatorId

				tempState := m.data.states[id]
				tempState.Suspended = types.SuspendedType{
					IsSuspended: true,
					TimeStamp:   time.Now(),
				}

				m.data.suspendedElevators[id] = tempState.Suspended
				fmt.Printf("Suspended elevator: %s\n", id)

				m.data.timeSinceAssignmentUpdate[floor][orderType].TimeStamp = time.Now()

				elevWasSuspended = true
			}
		}
	}
	return elevWasSuspended
}

func (m *Master) pushOrdersToNewMaster() {
	for floor := range N_FLOORS {
		for button := range 2 {
			if m.data.hallRequests[floor][button] {
				m.transferOrdersWhenMasterDowngradeCh <- types.OrderEnvelope{
					Order: types.Order{Floor: floor, Type: types.OrderType(button)},
				}
			}
		}
		for elevID, cabReq := range m.data.cabRequests {
			if cabReq[floor] {
				m.transferOrdersWhenMasterDowngradeCh <- types.OrderEnvelope{
					ElevatorID: elevID,
					Order:      types.Order{Floor: floor, Type: types.Cab},
				}
			}
		}
	}
}

func (data *masterData) storeOrder(order types.Order, elevatorId string) bool {
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

func (data *masterData) removeOrders(orders []types.Order, elevatorID string) bool {
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

func (m *Master) drainChannels() {
	select {
	case m.isMaster = <-m.isMasterCh:
	case <-m.receiveElevatorOrdersCh:
	case <-m.completedAssignmentCh:
	case <-m.rawAssignmentsCh:
	case <-m.updateStreamCh:
	case <-m.transferOrdersWhenMasterDowngradeCh:
	case <-m.forwardOrdersFromBackup:
	case <-m.receiveBackupAckCh:
		// default: no default, this might be a bug
	}
}
