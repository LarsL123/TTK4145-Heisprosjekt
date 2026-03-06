package elevatormanager

import (
	"elevatorproject/src/elevio"
	"fmt"
)

// TODO: Denne fila skal ha kontroll på FSMen til heisen:)

func fsm_onFloorArrival(elev *elevator, floor int){

}

func fsm_onDoorTimeout(elev *elevator, timeout bool){

}

func fsm_onNewButtonRequest(elev *elevator, buttonRequest elevio.ButtonEvent){
	fmt.Printf("New %s request on floor %d", buttonToString(buttonRequest.Button), buttonRequest.Floor)
	//TODO: request sendes videre til donaldtrump, deretter til master
	switch elev.behaviour{
		case EB_DoorOpen:
			if requestShouldClearImmediately(elev, buttonRequest){
				//TODO: Start timer
			}else {
				elev.requests[buttonRequest.Floor][buttonRequest.Button] = 1
			}
			break;

		case EB_Moving:
			elev.requests[buttonRequest.Floor][buttonRequest.Button] = 1 
			break

		case EB_Idle:
			elev.requests[buttonRequest.Floor][buttonRequest.Button] = 1
			direction, behaviour := requests_chooseDirection(*elev)
			elev.Dirn = direction
			elev.behaviour = behaviour
			switch behaviour{
				case EB_DoorOpen:
					elevio.SetDoorOpenLamp(true)
					// TODO: Start timer
					elev.requests = requests_clearAtCurrentFloor(elev)
					break
				case EB_Moving
					elevator_motorDirection(elev.Dirn)
					break
				
				case EB_Idle:
					break
				}
			break
	}
}

func fsm_onObstruction(elev *elevator){

}

func fsm_onInitBetweenFloors(elev *elevator){

}
