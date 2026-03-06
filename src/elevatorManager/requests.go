package elevatormanager

import "elevatorproject/src/elevio"

//TODO: i guess denne skal ta imot knapperequests og sende de til main

func request_above(elevator Elevator) bool {
	for f := elevator.floor +1 ; f < N_FLOORS ; f++ {
		for btn := 0; btn < N_BUTTONS; btn++{
			if elevator.requests[f][btn]{
				return true
			}
		}
	}
	return false
}

func request_below(elevator Elevator) bool {
	for f := 0 ; f < elevator.floor ; f++ {
		for btn := 0; btn < N_BUTTONS; btn++{
			if elevator.requests[f][btn]{
				return true
			}
		}
	}
	return false
}

func requests_here(elevator Elevator) bool{
		for btn := 0; btn < N_BUTTONS; btn++{
		if elevator.requests[elevator.floor][btn]{
			return true
		}
	}
	return false
}


// int requests_shouldStop(Elevator e) __attribute__((pure));
func requestsShouldStop(elevator Elevator) bool {
	switch elevator.Dirn {
	case elevio.MD_Down:
		return (elevator.requests[elevator.floor][elevio.BT_HallDown] ||
			elevator.requests[elevator.floor][elevio.BT_Cab] ||
			!request_below(elevator))
	case elevio.MD_Up:
		return (elevator.requests[elevator.floor][elevio.BT_HallUp] ||
			elevator.requests[elevator.floor][elevio.BT_Cab] ||
			!request_above(elevator))
	default:
		return true
	}
}

func requestShouldClearImmediately(elev Elevator, buttonRequest elevio.ButtonEvent) bool {
	//TODO: fix chooseDirection

}

func requests_chooseDirection(elev Elevator) (elevio.MotorDirection, Behaviour) {
	// TODO: fix chooseDirection

}

// Elevator requests_clearAtCurrentFloor(Elevator e) __attribute__((pure));
func requests_clearAtCurrentFloor(elevator Elevator) bool{
	// TODO: fix clearatcurrentfloor

}