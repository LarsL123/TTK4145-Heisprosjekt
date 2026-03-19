package controllers

import (
	"Network-go/network/bcast"
	"elevatorproject/src/config"
	elevatormanager "elevatorproject/src/elevatorManager"
	"elevatorproject/src/types"
	"fmt"
	"strconv"
	"time"
)

type Slave struct {
	id                         string
	messageCount               int
	pendingOrders              map[int]types.OrderEnvelope
	pendingFinishedAssignments map[int]types.FinishedHallAssignments

	net      NetworkChannels
	elevator ElevatorChannels
}

type NetworkChannels struct {
	sendElevatorState         chan types.ElevatorState
	sendOrdersCh              chan types.OrderEnvelope
	sendFinishedAssignmentsCh chan types.FinishedHallAssignments

	receiveAssignmentsFromMaster chan types.Assignments
	hallOrderAck                 chan types.OrderAck
	finishedAssignmentsAck       chan types.FinishedHallAssignmentsAck
}

type ElevatorChannels struct {
	receiveElevatorState          chan types.ElevatorState
	receiveOrdersCh               chan types.Order
	receiveCompletedAssignmentsCh chan []types.Order

	sendAssignmentsCh chan [N_FLOORS][N_BUTTONS]bool
	sendLightsCh      chan [N_FLOORS][N_BUTTONS]bool
}

func NewSlave(id string) *Slave {
	elev := ElevatorChannels{
		receiveElevatorState:          make(chan types.ElevatorState),
		receiveOrdersCh:               make(chan types.Order, 10),
		receiveCompletedAssignmentsCh: make(chan []types.Order, 10),
		sendAssignmentsCh:             make(chan [N_FLOORS][N_BUTTONS]bool),
		sendLightsCh:                  make(chan [N_FLOORS][N_BUTTONS]bool),
	}

	net := NetworkChannels{
		sendElevatorState:            make(chan types.ElevatorState),
		sendOrdersCh:                 make(chan types.OrderEnvelope, 10),
		sendFinishedAssignmentsCh:    make(chan types.FinishedHallAssignments),
		receiveAssignmentsFromMaster: make(chan types.Assignments),
		hallOrderAck:                 make(chan types.OrderAck, 10),
		finishedAssignmentsAck:       make(chan types.FinishedHallAssignmentsAck),
	}

	return &Slave{
		id:            id,
		messageCount:  0,
		pendingOrders: make(map[int]types.OrderEnvelope),
		net:           net,
		elevator:      elev,
	}
}

func (s *Slave) Start(transferDeadMaster chan types.OrderEnvelope, aliveCh chan struct{}) {
	var slaveRequests [N_FLOORS][N_BUTTONS]bool

	//Resending logic:
	resendTicker := time.NewTicker(config.Cfg.AckRetryRate)
	s.pendingOrders = make(map[int]types.OrderEnvelope)
	s.pendingFinishedAssignments = make(map[int]types.FinishedHallAssignments)

	net := s.net
	elevator := s.elevator

	// Start ElevatorManager
	go elevatormanager.ElevatorManager(s.elevator.receiveElevatorState, s.elevator.receiveOrdersCh, s.elevator.receiveCompletedAssignmentsCh, s.elevator.sendAssignmentsCh, s.elevator.sendLightsCh)

	// Broadcast transmitter & receiver
	go bcast.Transmitter(config.Cfg.MasterListenPort, s.net.sendElevatorState, s.net.sendOrdersCh, s.net.sendFinishedAssignmentsCh)
	go bcast.Receiver(config.Cfg.SlaveListenPort, s.net.receiveAssignmentsFromMaster, s.net.hallOrderAck, s.net.finishedAssignmentsAck)

	for {
		select {

		case state := <-elevator.receiveElevatorState:
			state.ID = s.id
			net.sendElevatorState <- state

		case order := <-transferDeadMaster:
			s.transferToMaster(order)

		case order := <-elevator.receiveOrdersCh:
			s.sendOrder(s.id, order)

		case assignment := <-net.receiveAssignmentsFromMaster:
			s.handleNewAssignments(assignment, &slaveRequests)

		case finishedAssignments := <-elevator.receiveCompletedAssignmentsCh:
			s.returnCompletedAssignment(finishedAssignments, &slaveRequests)

		case ack := <-net.hallOrderAck:
			s.clearPendingOrders(ack)

		case ack := <-net.finishedAssignmentsAck:
			s.clearPendingCompletedAssignments(ack)

		case <-resendTicker.C:
			s.resendLostPackets()
			pingProcessPair(aliveCh)

		}
	}
}

func (s *Slave) clearPendingOrders(ack types.OrderAck) {
	fmt.Println("Received ACK for order", ack.UpdateNr)
	delete(s.pendingOrders, ack.UpdateNr)
}

func (s *Slave) clearPendingCompletedAssignments(ack types.FinishedHallAssignmentsAck) {
	fmt.Println("Received ACK for assignment", ack.UpdateNr)
	delete(s.pendingFinishedAssignments, ack.UpdateNr)

}

func (s *Slave) returnCompletedAssignment(orders []types.Order, slaveRequests *[N_FLOORS][N_BUTTONS]bool) {
	for _, request := range orders {
		slaveRequests[request.Floor][request.Type] = false
	}

	s.messageCount++
	finishedAssigment := createFinishedAssignments(s.id, orders, s.messageCount)

	fmt.Println("Clearing assignment")
	s.net.sendFinishedAssignmentsCh <- finishedAssigment
	s.pendingFinishedAssignments[finishedAssigment.UpdateNr] = finishedAssigment
}

func (s *Slave) transferToMaster(order types.OrderEnvelope) {
	fmt.Println("Transfering order from dead master", order)

	if order.Order.Type == types.Cab {
		s.sendOrder(order.ElevatorID, order.Order)
	} else {
		s.sendOrder(order.ElevatorID, order.Order)
		s.elevator.receiveOrdersCh <- order.Order
	}
}

func (s *Slave) resendLostPackets() {
	cleanExpiredMessages(s.pendingOrders)
	cleanExpiredMessages(s.pendingFinishedAssignments)

	for updateNr, ho := range s.pendingOrders {
		fmt.Println("Resending Order: ", updateNr)
		s.net.sendOrdersCh <- ho
	}

	for updateNr, ass := range s.pendingFinishedAssignments {
		fmt.Println("Resending Assignment: ", updateNr)
		s.net.sendFinishedAssignmentsCh <- ass
	}
}

func (s *Slave) handleNewAssignments(as types.Assignments, requests *[N_FLOORS][N_BUTTONS]bool) {
	// Update local request state
	for f := 0; f < N_FLOORS; f++ {
		
		requests[f][types.HallDown] = as.Assignments[s.id][f][types.HallDown]
		requests[f][types.HallUp] = as.Assignments[s.id][f][types.HallUp]
		if as.Assignments[s.id][f][types.Cab] {
			requests[f][types.Cab] = true
		}
	}

	s.elevator.sendAssignmentsCh <- *requests
	s.elevator.sendLightsCh <- calculateLights(as.Assignments, *requests)
}

func (s *Slave) sendOrder(id string, order types.Order) {
	s.messageCount++
	ho := createHallOrder(id, order, s.messageCount)

	s.net.sendOrdersCh <- ho
	s.pendingOrders[ho.UpdateNr] = ho
}

func cleanExpiredMessages[T types.LivingMessage](pending map[int]T) {
	for updateNr, msg := range pending {
		if time.Since(msg.GetCreationTime()) > config.Cfg.AckTimeout {
			delete(pending, updateNr)
			fmt.Println("Dropping order:", updateNr)
		}
	}
}

func calculateLights(assignments map[string][N_FLOORS][3]bool, slaveRequests [N_FLOORS][N_BUTTONS]bool) [N_FLOORS][N_BUTTONS]bool {
	var lightsOn [N_FLOORS][N_BUTTONS]bool
	for _, assignment := range assignments {
		for floor := range N_FLOORS {
			for btn := range 2 {
				if assignment[floor][btn] {
					lightsOn[floor][btn] = true
				}
			}
		}
	}
	for floor := range N_FLOORS {
		lightsOn[floor][2] = slaveRequests[floor][2]
	}
	return lightsOn
}

func createHallOrder(id string, order types.Order, messageCount int) types.OrderEnvelope {
	idInt, _ := strconv.Atoi(id)

	return types.OrderEnvelope{
		ElevatorID: id,
		Order:      order,
		CreatedAt:  time.Now(),
		UpdateNr:   idInt*1000000 + messageCount,
	}
}

func createFinishedAssignments(id string, orders []types.Order, messageCount int) types.FinishedHallAssignments {
	idInt, _ := strconv.Atoi(id)
	sendToMaster := types.FinishedHallAssignments{
		ElevatorID: id,
		UpdateNr:   idInt*1000000 + messageCount,
		CreatedAt:  time.Now(),
		Orders:     make([]types.Order, len(orders)),
	}

	for i, request := range orders {
		sendToMaster.Orders[i] = types.Order{
			Floor: request.Floor,
			Type:  types.OrderType(request.Type),
		}
	}

	return sendToMaster
}
