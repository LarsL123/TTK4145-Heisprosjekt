package elevatormanager

import (
	"elevatorproject/src/elevio"
	"fmt"
)

// TODO: Denne fila skal ha kontroll på FSMen til heisen:)

func fsm_onFloorArrival(elev *Elevator, newFloor int) {
	fmt.Printf("Reached new floor: %d", newFloor)
	// Print elevator?
	elev.floor = newFloor

	elevio.SetFloorIndicator(elev.floor)

	switch elev.behaviour {
	case EB_Moving:
		if true { //TODO should be: requestsShouldStop(elev)
			elevio.SetMotorDirection(elevio.MD_Stop)
			elevio.SetDoorOpenLamp(true)
			//Start timer
			//Setall lights??
			elev.behaviour = EB_DoorOpen
		}
		break
	default:
		break
	}
	fmt.Printf("\nNew state:\n")
	// TODO: implement elevatorPrint(elev)
}

func fsm_onDoorTimeout(elev *Elevator, timeout bool) {
	fmt.Println("Door timed out")
	// TODO: implement elevatorPrint(elev)

	switch elev.behaviour {
	case EB_DoorOpen:
		dirn, behaviour := requests_chooseDirection(*elev)
		elev.Dirn = dirn
		elev.behaviour = behaviour

		switch elev.behaviour {
		case EB_DoorOpen:
			//Start timer

		}
	default:
		break
	}
}

func fsm_onNewButtonRequest(elev *Elevator, buttonRequest elevio.ButtonEvent) {
	fmt.Printf("New %s request on floor %d", buttonToString(buttonRequest.Button), buttonRequest.Floor)
	//TODO: request sendes videre til donaldtrump, deretter til master
	switch elev.behaviour {
	case EB_DoorOpen:
		if requestShouldClearImmediately(*elev, buttonRequest) {
			//TODO: Start timer
		} else {
			elev.requests[buttonRequest.Floor][buttonRequest.Button] = true
		}
		break

	case EB_Moving:
		elev.requests[buttonRequest.Floor][buttonRequest.Button] = true
		break

	case EB_Idle:
		elev.requests[buttonRequest.Floor][buttonRequest.Button] = true
		direction, behaviour := requests_chooseDirection(*elev)
		elev.Dirn = direction
		elev.behaviour = behaviour
		switch behaviour {
		case EB_DoorOpen:
			elevio.SetDoorOpenLamp(true)
			// TODO: Start timer
			elev.requests = requests_clearAtCurrentFloor(elev)
			break
		case EB_Moving:
			elevio.SetMotorDirection(elev.Dirn)
			break

		case EB_Idle:
			break
		}
		break
	}
}

func fsm_onObstruction(elev *Elevator) {

}

func fsm_onInitBetweenFloors(elev *Elevator) {

}
