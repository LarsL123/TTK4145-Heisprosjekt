package elevatormanager
//TODO: lage elevator class
import (
	"fmt"
	"Driver-go/elevio"

)



type Behaviour int

const (
	EB_Idle Behaviour = 0 
	EB_DoorOpen Behaviour = 1
	EB_Moving Behaviour = 2
)



type Elevator struct {
	floor int
	Dirn elevio.MotorDirection
	Btn elevio.ButtonType	
}