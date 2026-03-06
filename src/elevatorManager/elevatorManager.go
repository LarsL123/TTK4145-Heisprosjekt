package elevatormanager

// Denne modulen skal være inspirert av utdelt single elevator control kode

import (
	"fmt"

	"elevatorproject/src/elevio"
)
const N_FLOORS = 4
//TODO: Lage en config_load funksjon
const address = "123.123.123.123"
const N_BUTTONS = 3



func main (){

	

	//pollRate_ms := 25

	// Spørsmål til studass: er det greit å heller definere elevator på package level, slipper dermed å passe elevator pointer til alle funksjonene som skal endre på den??
	elev := elevator_uninitialized(address, N_FLOORS)
	elevio.Init(address, N_FLOORS)
	
	chDriverButtonRequests 		:= make(chan elevio.ButtonEvent)
	chDriverFloorSensor 		:= make(chan int)
	chDriverObstruction 		:= make(chan bool)
	chDriverStopButton 			:= make(chan bool)
	chDoorTimeout 				:= make(chan bool)

	go elevio.PollFloorSensor(chDriverFloorSensor)
	go elevio.PollButtons(chDriverButtonRequests)
	go elevio.PollObstructionSwitch(chDriverObstruction)
	go elevio.PollStopButton(chDriverStopButton)
	go doortimer(chDoorTimeout)

	if elevio.GetFloor() == -1{
		fsm_onInitBetweenFloors(&elev)
	}

	for {
		select{
		case newButtonRequest := <- chDriverButtonRequests:
			fmt.Println("Button pressed")
			fsm_onNewButtonRequest(&elev, newButtonRequest)
			// TODO: fikse at requesten sendes til master

		case floorArrivedAt := <- chDriverFloorSensor:
			fmt.Printf("Arrived at floor %d\n", floorArrivedAt)
			fsm_onFloorArrival(&elev, floorArrivedAt)

		case timedout := <- chDoorTimeout:
			fmt.Println("Door timeout")
			fsm_onDoorTimeout(&elev,timedout)
			//timer

		case  <- chDriverObstruction:
			fsm_onObstruction(&elev)
		}

	}
}

