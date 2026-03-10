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

func fsm_onNewButtonRequest(buttonRequest elevio.ButtonEvent) {
	fmt.Printf("New %s request on floor %d", buttonToString(buttonRequest.Button), buttonRequest.Floor)
	//TODO: request sendes videre til donaldtrump, deretter til master
	switch elevator.behaviour {
	case EB_DoorOpen:
		if requestShouldClearImmediately(buttonRequest) {
			doortimer_start()
		} else {
			elevator.requests[buttonRequest.Floor][buttonRequest.Button] = true
		}

	case EB_Moving:
		elevator.requests[buttonRequest.Floor][buttonRequest.Button] = true
	case EB_Idle:
		elevator.requests[buttonRequest.Floor][buttonRequest.Button] = true
		direction, behaviour := requests_chooseDirection()
		elevator.dirn = direction
		elevator.behaviour = behaviour
		switch behaviour {
		case EB_DoorOpen:
			elevio.SetDoorOpenLamp(true)
			doortimer_start()
			requests_clearAtCurrentFloor()

		case EB_Moving:
			elevio.SetMotorDirection(elevator.dirn)

		case EB_Idle:

		}
	}
	//setAllLights(elev)//?
	fmt.Println("\nNew state:")
	//print elevator()
}

func fsm_onFloorArrival(newFloor int) {
	fmt.Printf("Reached new floor: %d", newFloor)
	// Print elevator?
	elevator.floor = newFloor

	elevio.SetFloorIndicator(elevator.floor)

	switch elevator.behaviour {
	case EB_Moving:
		if true { //TODO should be: requestsShouldStop(elev)
			elevio.SetMotorDirection(elevio.MD_Stop)
			elevio.SetDoorOpenLamp(true)
			doortimer_start()
			//Setall lights??
			elevator.behaviour = EB_DoorOpen
		}
	default:
		break
	}
	fmt.Printf("\nNew state:\n")
	// TODO: implement elevatorPrint(elev)
}

func fsm_onDoorTimeout() {
	fmt.Println("Door timed out")
	// TODO: implement elevatorPrint(elev)

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

			requests_clearAtCurrentFloor()
			//TODO: setAllLights(elev)
		case EB_Moving:
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
