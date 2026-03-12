package elevatormanager

// Denne modulen skal være inspirert av utdelt single elevator control kode

import (
	"fmt"

	"elevatorproject/src/elevio"
)

const N_FLOORS = 4

// TODO: Lage en config_load funksjon
const address = "0.0.0.0:15657"
const N_BUTTONS = 3
const DOOR_OPEN_DURATION = 3 // [seconds]

func main(sendOrderCh chan elevio.ButtonEvent, receiveAssignmentsCh chan elevio.ButtonEvent) {

	//pollRate_ms := 25

	// Spørsmål til studass: er det greit å heller definere elevator på package level, slipper dermed å passe elevator pointer til alle funksjonene som skal endre på den??

	elevio.Init(address, N_FLOORS)

	driverButtonRequestsCh := make(chan elevio.ButtonEvent)
	driverFloorSensorCh := make(chan int)
	driverObstructionCh := make(chan bool)
	driverStopButtonCh := make(chan bool)
	//doorTimeoutCh 				:= make(chan int)

	go elevio.PollFloorSensor(driverFloorSensorCh)
	go elevio.PollButtons(driverButtonRequestsCh)
	go elevio.PollObstructionSwitch(driverObstructionCh)
	go elevio.PollStopButton(driverStopButtonCh)

	doortimer_init()

	if elevio.GetFloor() == -1 {
		fsm_onInitBetweenFloors()
	}

	for {
		select {
		case newAssignment := <- receiveAssignmentsCh:
			fmt.Println("Received Assignment")
			fsm_onNewAssignment(newAssignment)

		case newButtonRequest := <-driverButtonRequestsCh:
			fmt.Println("Button pressed")
			fsm_onNewButtonRequest(newButtonRequest,sendOrderCh)
			// TODO: fikse at requesten sendes til master

		case floorArrivedAt := <-driverFloorSensorCh:
			fmt.Printf("Arrived at floor %d\n", floorArrivedAt)
			fsm_onFloorArrival(floorArrivedAt)

		case <-doortimer.C:
			fmt.Println("Door timeout")
			fsm_onDoorTimeout()

		case obstruction := <-driverObstructionCh:
			fsm_onObstruction(obstruction)
		}
	}
}
