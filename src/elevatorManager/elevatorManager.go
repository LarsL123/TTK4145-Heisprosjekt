package elevatormanager

// Denne modulen skal være inspirert av utdelt single elevator control kode

import (
	"fmt"

	"elevatorproject/src/elevio"
)

func main (){

	//TODO: Lage en config_load funksjon
	address := "123.123.123.123"
	N_FLOORS := 4
	//pollRate_ms := 25


	elev := elevator_uninitialized(address, N_FLOORS)
	elevio.Init(address, N_FLOORS)
	
	chDriverButtonRequests := make(chan elevio.ButtonEvent)
	chDriverFloorSensor := make(chan int)
	chDriverObstruction := make(chan bool)
	chDriverTimeout := make(chan int)
	chDriverStopButton := make(chan bool)
	

	go elevio.PollFloorSensor(chDriverFloorSensor)
	go elevio.PollButtons(chDriverButtonRequests)
	go elevio.PollObstructionSwitch(chDriverObstruction)
	go elevio.PollStopButton(chDriverStopButton)

	for {
		select{
		case newButtonRequest := <- chDriverButtonRequests:
			fmt.Println("Button pressed")
			fsm_onNewButtonRequest(&elev, newButtonRequest)
			// TODO: fikse at requesten sendes til master

		case floorArrivedAt := <- chDriverFloorSensor:
			fmt.Println("Arrived at floor %d\n", floorArrivedAt)
			fsm_onFloorArrival(&elev, floorArrivedAt)

		case timedout := <- chDriverTimeout:
			fsm_onDoorTimeout(&elev,timedout)
			//timer

		case  <- chDriverObstruction :
			fsm_onObstruction(&elev)
		}

	}
}

