package elevatormanager

// Denne modulen skal være inspirert av utdelt single elevator control kode

import (
	"fmt"

	"elevatorproject/src/elevio"
)

func main (){
	//TODO: finne ut hvor adressen skal komme fra
	//TODO: finne ut hvor N_FLOORS skal komme fra
	address := "123.123.123.123"
	N_FLOORS := 4


	elevator := elevator_uninitialized(address, N_FLOORS)

	
	driverButtonRequests := make(chan elevio.ButtonEvent)
	driverFloorSensor := make(chan int)
	driverObstruction := make(chan bool)
	driverTimeout := make(chan int)
	driverStopButton := make(chan bool)
	

	go elevio.PollFloorSensor(driverFloorSensor)
	go elevio.PollButtons(driverButtonRequests)
	go elevio.PollObstructionSwitch(driverObstruction)
	go elevio.PollStopButton(driverStopButton)

	for {
		select{
		case newButtonRequest := <- driverButtonRequests:
			fmt.Println("Button pressed")
			fsm_onNewButtonRequest(newButtonRequest)
			// TODO: fikse at requesten sendes til master

		case floorArrivedAt := <- driverFloorSensor:
			fmt.Println("Arrived at floor %d\n", floorArrivedAt)
			fsm_onFloorArrival(floorArrivedAt)
		case timedout := <- driverTimeout:
			fsm_onDoorTimeout(timedout)
			//timer
		case  <- driverObstruction :
			fsm_onObstruction()
		}

	}
}

