package ordermanager

/*

The order manager module recomputes the optimal assignments based on orders and states. Whenever new data enters the system, all hall requests gets reassigned.
This new data could be a new request, an updated state from some elevator, or an update on who is alive on the network. This redistribution means
that a request is not necessarily assigned to the same elevator for the duration of its lifetime, but can instead be re-assigned to a new elevator,
for example if a new idle elevator connects to the network, or the previously assigned elevator gets a lot of cab requests.

IMPORTANT:  The module is supposed to trigger when
                - A new order is received
                - New info about the peers in the network is received
                - Behaviour is updated for an elevator (Meaning e.g. idle -> moving)

                For the module to work, it needs an input as the struct HRAInput, meaning
                both a map of all elevator states (HRAElevState) and all hall requests, since
                the elevator states contain cab requests.

Input:
	channel: Orders

Output:
	channel: Assignments

*/

import (
	"elevatorproject/src/types"
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"
	// "golang.org/x/text/cases"
)

// Struct members must be public in order to be accessible by json.Marshal/.Unmarshal
// This means they must start with a capital letter, so we need to use field renaming struct tags to make them camelCase

type HRAElevState struct {
	Behavior    string  `json:"behaviour"`
	Floor       int     `json:"floor"`
	Direction   string  `json:"direction"`
	CabRequests [4]bool `json:"cabRequests"`
}

type HRAInput struct {
	HallRequests       [4][2]bool              `json:"hallRequests"`
	States             map[string]HRAElevState `json:"states"`
	SuspendedElevators map[string]types.SuspendedType
}

func ToHRAInput(hallRequests [4][2]bool, cabRequests map[string][4]bool, elevatorStates map[string]types.ElevatorState, suspendedElevators map[string]types.SuspendedType) HRAInput {

	inputStates := make(map[string]HRAElevState)

	for id, elevatorState := range elevatorStates {
		inputStates[id] = HRAElevState{
			Behavior:    elevatorState.Behaviour,
			Floor:       elevatorState.Floor,
			Direction:   elevatorState.Direction,
			CabRequests: cabRequests[id],
		}
	}

	return HRAInput{
		HallRequests:       hallRequests,
		States:             inputStates,
		SuspendedElevators: suspendedElevators,
	}
}

// Calculates optimal assignments based on orders
func ManageOrders(OrdersCh chan HRAInput, AssignmentsCh chan map[string][4][2]bool) {

	hraExecutable := ""
	switch runtime.GOOS {
	case "linux":
		hraExecutable = "hall_request_assigner"
	case "windows":
		hraExecutable = "hall_request_assigner.exe"
	default:
		panic("OS not supported")
	}

	for {

		// Order is received on input channel
		input := <-OrdersCh
		fmt.Println("Calculating orders")

		for id, suspended := range input.SuspendedElevators {
			if suspended.IsSuspended {

				delete(input.States, id)
				fmt.Printf("Deleted elevatorstate: %s\n", id)
			}
		}

		if len(input.States) == 0 {
			fmt.Printf("No input states provided to ordermanager, skipping.\n")
			continue
		}

		// JSON -> String
		jsonBytes, err := json.Marshal(input)
		if err != nil {
			fmt.Println("json.Marshal error: ", err)
			return
		}

		// Run cost function executable
		ret, err := exec.Command("./src/ordermanager/"+hraExecutable, "-i", string(jsonBytes)).CombinedOutput()
		if err != nil {
			fmt.Println("exec.Command error: ", err)
			fmt.Println(string(ret))
			return
		}

		output := new(map[string][4][2]bool)

		// Update output map with executable data, String -> JSON
		err = json.Unmarshal(ret, &output)
		if err != nil {
			fmt.Println("json.Unmarshal error: ", err)
			return
		}

		// Pass optimal assignments to output channel
		AssignmentsCh <- *output
	}
}
