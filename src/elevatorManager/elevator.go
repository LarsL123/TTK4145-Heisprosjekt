package elevatormanager

import (
	"elevatorproject/src/elevio"
	"fmt"
)

type Behaviour int

const (
	EB_Idle     Behaviour = 0
	EB_DoorOpen Behaviour = 1
	EB_Moving   Behaviour = 2
)

type Elevator struct {
	floor            	int
	dirn             	elevio.MotorDirection
	doorOpenDuration 	float32
	behaviour        	Behaviour
	requests         	[N_FLOORS][N_BUTTONS]bool
	//assignments			[N_FLOORS][N_BUTTONS]bool
	obstructed       	bool
}

func buttonToString(button elevio.ButtonType) string {
	switch button {
	case elevio.BT_HallUp:
		return "HallUp"
	case elevio.BT_HallDown:
		return "HallDown"
	case elevio.BT_Cab:
		return "Cab"
	}
	return "Invalid Button"
}

func dirnToString(dirn elevio.MotorDirection) string {
	switch dirn {
	case elevio.MD_Down:
		return "MD_Down"
	case elevio.MD_Up:
		return "MD_up"
	case elevio.MD_Stop:
		return "MD_Stop"
	default:
		return "MD_Undefined"
	}
}

func behaviourToString(behaviour Behaviour) string {
	switch behaviour {
	case EB_Idle:
		return "EB_Idle"
	case EB_DoorOpen:
		return "EB_DoorOpen"
	case EB_Moving:
		return "EB_Moving"
	default:
		return "EB_Undefined"
	}
}

func elevator_print() {
	fmt.Println("--------------------------------")
	fmt.Printf("floor: %d\n"+
		"dirn: %s\n"+
		"behaviour: %s\n", elevator.floor, dirnToString(elevator.dirn), behaviourToString(elevator.behaviour))
	fmt.Println("Requests:")
	fmt.Print(elevator.requests)
}
