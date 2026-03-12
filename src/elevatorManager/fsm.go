package elevatormanager

import (
	"elevatorproject/src/elevio"
	"fmt"
)

// TODO: Denne fila skal ha kontroll på FSMen til heisen:)
var elevator = Elevator{floor: -1,
	behaviour:        EB_Idle,
	doorOpenDuration: 3.0,
	dirn:             elevio.MD_Down,
	obstructed:       false}

func fsm_setAllLights() {
	for floor := 0; floor < N_FLOORS; floor++ {
		for btn := 0; btn < N_BUTTONS; btn++ {
			elevio.SetButtonLamp(elevio.ButtonType(btn), floor, elevator.requests[floor][btn])
		}
	}
}

func fsm_onInitBetweenFloors() {
	elevio.SetMotorDirection(elevio.MD_Down)
	elevator.dirn = elevio.MD_Down
	elevator.behaviour = EB_Moving
}

// Skal motta requests i form av en array [N_FLOORS][N_BUTTONS]bool
func fsm_onNewAssignment(requests [N_FLOORS][N_BUTTONS]bool, sendClearedRequests chan elevio.ButtonEvent) {
	fmt.Printf("Requests updated")
	//TODO: request sendes videre til donaldtrump, deretter til master

	switch elevator.behaviour {
	case EB_DoorOpen:
		elevator.requests = requests
	case EB_Moving:
		elevator.requests = requests
	case EB_Idle:
		elevator.requests = requests
		direction, behaviour := requests_chooseDirection()
		elevator.dirn = direction
		elevator.behaviour = behaviour
		switch behaviour {
		case EB_DoorOpen:
			elevio.SetDoorOpenLamp(true)
			doortimer_start()
			requests_clearAtCurrentFloor(sendClearedRequests) // Denne sender per nå også til orderHandler, burde kanskje implementeres i annen kode, men nå er det sånn.

		case EB_Moving:
			elevio.SetMotorDirection(elevator.dirn)
		case EB_Idle:
		}
	}
	fsm_setAllLights()
	fmt.Println("\nNew state:")
	elevator_print()
}

func fsm_onFloorArrival(newFloor int, sendClearedRequests chan elevio.ButtonEvent) {
	fmt.Printf("Reached new floor: %d", newFloor)
	elevator.floor = newFloor
	elevator_print()


	elevio.SetFloorIndicator(elevator.floor)

	switch elevator.behaviour {
	case EB_Moving:
		if requestsShouldStop() {
			elevio.SetMotorDirection(elevio.MD_Stop)
			elevio.SetDoorOpenLamp(true)
			requests_clearAtCurrentFloor(sendClearedRequests) // Denne sender per nå også til orderHandler, burde kanskje implementeres i annen kode, men nå er det sånn.
			doortimer_start()
			fsm_setAllLights()
			elevator.behaviour = EB_DoorOpen
		}
	default:
		break
	}
	fmt.Printf("\nNew state:\n")
	elevator_print()
}

func fsm_onDoorTimeout(sendClearedRequests chan elevio.ButtonEvent) {
	fmt.Println("Door timed out")
	elevator_print()

	switch elevator.behaviour {
	case EB_DoorOpen:
		if elevator.obstructed {
			doortimer_start()
			return
		}

		dirn, behaviour := requests_chooseDirection()
		elevator.dirn = dirn
		elevator.behaviour = behaviour

		switch elevator.behaviour {
		case EB_DoorOpen:
			doortimer_start()
			requests_clearAtCurrentFloor(sendClearedRequests) // Denne sender per nå også til orderHandler, burde kanskje implementeres i annen kode, men nå er det sånn.
			fsm_setAllLights()
		case EB_Moving:
			elevio.SetMotorDirection(elevator.dirn)
			elevio.SetDoorOpenLamp(false)
		case EB_Idle:
			elevio.SetDoorOpenLamp(false)
			elevio.SetMotorDirection(elevator.dirn)
		}
	default:
		break
	}
}

func fsm_onObstruction(obstruction bool) {
	elevator.obstructed = obstruction
}

func fsm_onNewButtonRequest(buttonRequest elevio.ButtonEvent, sendOrderCh chan<- elevio.ButtonEvent) {
	fmt.Printf("New %s order on floor %d", buttonToString(buttonRequest.Button), buttonRequest.Floor)
	if elevator.behaviour == EB_DoorOpen && requestShouldClearImmediately(buttonRequest) {
		doortimer_start()
	} else {
		sendOrderCh <- buttonRequest
	}
}
