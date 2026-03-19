package elevatormanager

import (
	"elevatorproject/src/elevio"
	"elevatorproject/src/types"
)

//TODO: i guess denne skal ta imot knapperequests og sende de til main

func request_above() bool {
	for f := elevator.floor + 1; f < N_FLOORS; f++ {
		for btn := 0; btn < N_BUTTONS; btn++ {
			if elevator.requests[f][btn] {
				return true
			}
		}
	}
	return false
}

func request_below() bool {
	for f := 0; f < elevator.floor; f++ {
		for btn := 0; btn < N_BUTTONS; btn++ {
			if elevator.requests[f][btn] {
				return true
			}
		}
	}
	return false
}

func request_here() bool {
	for btn := 0; btn < N_BUTTONS; btn++ {
		if elevator.requests[elevator.floor][btn] {
			return true
		}
	}
	return false
}

// func requests_chooseDirection() (elevio.MotorDirection){

// }

func requests_chooseDirection() (elevio.MotorDirection, Behaviour) {
	switch elevator.dirn {
	case elevio.MD_Up:
		if request_above() {
			return elevio.MD_Up, EB_Moving
		} else if request_here() {
			return elevio.MD_Down, EB_DoorOpen
		} else if request_below() {
			return elevio.MD_Down, EB_Moving
		} else {
			return elevio.MD_Stop, EB_Idle
		}
	case elevio.MD_Down:
		if request_below() {
			return elevio.MD_Down, EB_Moving
		} else if request_here() {
			return elevio.MD_Up, EB_DoorOpen
		} else if request_above() {
			return elevio.MD_Up, EB_Moving
		} else {
			return elevio.MD_Stop, EB_Idle
		}
	case elevio.MD_Stop:
		if request_here() {
			return elevio.MD_Stop, EB_DoorOpen
		} else if request_above() {
			return elevio.MD_Up, EB_Moving
		} else if request_below() {
			return elevio.MD_Down, EB_Moving
		} else {
			return elevio.MD_Stop, EB_Idle
		}
	default:
		return elevio.MD_Stop, EB_Idle
	}
}

// int requests_shouldStop(Elevator e) __attribute__((pure));
func requestsShouldStop() bool {
	switch elevator.dirn {
	case elevio.MD_Down:
		return (elevator.requests[elevator.floor][elevio.BT_HallDown] ||
			elevator.requests[elevator.floor][elevio.BT_Cab] ||
			!request_below())
	case elevio.MD_Up:
		return (elevator.requests[elevator.floor][elevio.BT_HallUp] ||
			elevator.requests[elevator.floor][elevio.BT_Cab] ||
			!request_above())
	default:
		return true
	}
}

// NB Side effect: requests på samme heis til samme etasje som heisen er på med dør åpen blir ikke sendt til master
func requestShouldClearImmediately(buttonRequest elevio.ButtonEvent) bool {
	//TODO: fix chooseDirection
	return (elevator.floor == buttonRequest.Floor) && ((elevator.dirn == elevio.MD_Up && buttonRequest.Button == elevio.BT_HallUp) ||
		(elevator.dirn == elevio.MD_Down && buttonRequest.Button == elevio.BT_HallDown) ||
		elevator.dirn == elevio.MD_Stop ||
		buttonRequest.Button == elevio.BT_Cab)
}

//Lars: Wow for et monster ja. Håper ikke sverre ser den greia her...

// Denne sender per nå også til orderHandler, burde kanskje implementeres i annen kode, men nå er det sånn.
// Må også cleane opp spaghettikoden.
func requests_clearAtCurrentFloor(sendClearedRequests chan []types.Order) {
	// TODO: fix clearatcurrentfloor
	clearedRequestArray := make([]types.Order, 0)
	//var clearedRequest elevio.ButtonEvent
	elevator.requests[elevator.floor][elevio.BT_Cab] = false

	clearedRequestArray = append(clearedRequestArray, types.Order{
		Floor: elevator.floor,
		Type:  types.Cab})

	switch elevator.dirn {
	case elevio.MD_Up:
		if !request_above() && !elevator.requests[elevator.floor][elevio.BT_HallUp] {
			elevator.requests[elevator.floor][elevio.BT_HallDown] = false
			clearedRequestArray = append(clearedRequestArray, types.Order{
				Floor: elevator.floor,
				Type:  types.HallDown})
		}

		elevator.requests[elevator.floor][elevio.BT_HallUp] = false
		clearedRequestArray = append(clearedRequestArray, types.Order{
			Floor: elevator.floor,
			Type:  types.HallUp})
	case elevio.MD_Down:
		if !request_below() && !elevator.requests[elevator.floor][elevio.BT_HallDown] {
			elevator.requests[elevator.floor][elevio.BT_HallUp] = false
			clearedRequestArray = append(clearedRequestArray, types.Order{
				Floor: elevator.floor,
				Type:  types.HallUp})
		}
		elevator.requests[elevator.floor][elevio.BT_HallDown] = false

		clearedRequestArray = append(clearedRequestArray, types.Order{
			Floor: elevator.floor,
			Type:  types.HallDown})

	default:
		elevator.requests[elevator.floor][elevio.BT_HallUp] = false
		elevator.requests[elevator.floor][elevio.BT_HallDown] = false

		clearedRequestArray = append(clearedRequestArray, types.Order{
			Floor: elevator.floor,
			Type:  types.HallUp})
		clearedRequestArray = append(clearedRequestArray, types.Order{
			Floor: elevator.floor,
			Type:  types.HallDown})
	}

	sendClearedRequests <- clearedRequestArray
}
