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
	requests         [][]bool
	obstructed       bool
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

// Elevator_floorsensor()
// int elevator_floorSensor(void);
// int elevator_requestButton(int f, Button b);
// int elevator_stopButton(void);
// int elevator_obstruction(void);

// void elevator_floorIndicator(int f);
// void elevator_requestButtonLight(int f, Button b, int v);
// void elevator_doorLight(int v);
// void elevator_stopButtonLight(int v);
// void elevator_motorDirection(Dirn d);
