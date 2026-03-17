package elevatormanager

// Denne modulen skal være inspirert av utdelt single elevator control kode

import (
	"fmt"
	"time"

	"elevatorproject/src/config"
	"elevatorproject/src/elevio"
	"elevatorproject/src/types"
)

const N_FLOORS = 4

// TODO: Lage en config_load funksjon
var address = "0.0.0.0:" + config.Cfg.ElevatorPort
const N_BUTTONS = 3
const DOOR_OPEN_DURATION = 3 // [seconds]

func ElevatorManager(elevStateCh chan<- types.ElevatorState, sendOrderCh chan types.Order, sendFinishedOrderch chan []elevio.ButtonEvent, receiveAssignmentsCh chan [N_FLOORS][N_BUTTONS]bool) {
	// Spørsmål til studass: er det greit å heller definere elevator på package level, slipper dermed å passe elevator pointer til alle funksjonene som skal endre på den??

	elevio.Init(address, N_FLOORS)
	elevator_init()

	driverButtonRequestsCh := make(chan elevio.ButtonEvent)
	driverFloorSensorCh := make(chan int)
	driverObstructionCh := make(chan bool)
	driverStopButtonCh := make(chan bool)
	//doorTimeoutCh 				:= make(chan int)

	go elevio.PollFloorSensor(driverFloorSensorCh)
	go elevio.PollButtons(driverButtonRequestsCh)
	go elevio.PollObstructionSwitch(driverObstructionCh)
	go elevio.PollStopButton(driverStopButtonCh)

	if elevio.GetFloor() == -1 {
		fsm_onInitBetweenFloors()
	}

	sendStateTicker := time.NewTicker(config.Cfg.ElevatorUpdateRate)
	defer sendStateTicker.Stop()

	for {
		select {
		case newAssignment := <-receiveAssignmentsCh:
			fmt.Println("Received Assignment")
			fsm_onNewAssignment(newAssignment, sendFinishedOrderch)

		case newButtonRequest := <-driverButtonRequestsCh:
			fmt.Println("Button pressed")
			fsm_onNewButtonRequest(newButtonRequest, sendOrderCh)
			// TODO: fikse at requesten sendes til master

		case floorArrivedAt := <-driverFloorSensorCh:
			fmt.Printf("Arrived at floor %d\n", floorArrivedAt)
			fsm_onFloorArrival(floorArrivedAt, sendFinishedOrderch)

		case <-doortimer.C:
			fmt.Println("Door timeout")
			fsm_onDoorTimeout(sendFinishedOrderch)

		case obstruction := <-driverObstructionCh:
			fsm_onObstruction(obstruction)

		case <-sendStateTicker.C:
			elevStateCh <- types.ElevatorState{ //Export interface so we dont have to expose Elevator internals. Daniel remove comment when read.
				Floor:       elevator.floor,
				Direction:   dirnToString(elevator.dirn),
				Behaviour:   behaviourToString(elevator.behaviour),
				CabRequests: fsm_getCabRequests(),
				CreatedAt:   time.Now(),
				Obstructed:  elevator.obstructed,
			}
		}

	}
}
