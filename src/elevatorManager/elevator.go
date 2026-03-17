package elevatormanager

import (
	"elevatorproject/src/elevio"
)

type Behaviour int

const (
	EB_Idle     Behaviour = 0
	EB_DoorOpen Behaviour = 1
	EB_Moving   Behaviour = 2
)

type Elevator struct {
	floor            int
	dirn             elevio.MotorDirection
	doorOpenDuration float32
	behaviour        Behaviour
	requests         [N_FLOORS][N_BUTTONS]bool
	lights_on        [N_FLOORS][N_BUTTONS]bool
	obstructed       bool
}

func elevator_init() {
	doortimer_init()
	elevio.SetDoorOpenLamp(false)
	fsm_setAllLights()
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
		return "down"
	case elevio.MD_Up:
		return "up"
	case elevio.MD_Stop:
		return "stop"
	default:
		return "undefined"
	}
}

func behaviourToString(behaviour Behaviour) string {
	switch behaviour {
	case EB_Idle:
		return "idle"
	case EB_DoorOpen:
		return "doorOpen"
	case EB_Moving:
		return "moving"
	default:
		return "undefined"
	}
}

func elevator_print() {
	/*
		fmt.Println("--------------------------------")
		fmt.Printf("floor: %d\n"+
			"dirn: %s\n"+
			"behaviour: %s\n", elevator.floor, dirnToString(elevator.dirn), behaviourToString(elevator.behaviour))
		fmt.Println("Requests:")
		fmt.Print(elevator.requests)
		fmt.Println("\n-------------------------------") */
}
