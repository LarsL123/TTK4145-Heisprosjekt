package elevatormanager



import (
	"elevatorproject/src/elevio"
)

type Behaviour int

const (
	EB_Idle Behaviour = 0 
	EB_DoorOpen Behaviour = 1
	EB_Moving Behaviour = 2
)

type elevator struct {
	floor int
	Dirn elevio.MotorDirection
	doorOpenDuration float32
	behaviour Behaviour
	requests [][]int
}

func elevator_uninitialized(addr string, numFloors int) elevator{
	
	thisElevator := elevator{floor: -1, Dirn: elevio.MD_Stop, behaviour: EB_Idle, doorOpenDuration: 3.0 }
	return thisElevator
}

func buttonToString(button elevio.ButtonType) string{
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