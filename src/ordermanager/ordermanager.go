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
	"elevatorproject/src/config"
	"elevatorproject/src/types"
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"
	"time"
)

// Struct members must be public in order to be accessible by json.Marshal/.Unmarshal
// This means they must start with a capital letter, so we need to use field renaming struct tags to make them camelCase
const N_FLOORS = types.N_FLOORS

type HRAElevState struct {
	Behavior    string         `json:"behaviour"`
	Floor       int            `json:"floor"`
	Direction   string         `json:"direction"`
	CabRequests [N_FLOORS]bool `json:"cabRequests"`
}

type HRAInput struct {
	HallRequests       [N_FLOORS][2]bool       `json:"hallRequests"`
	States             map[string]HRAElevState `json:"states"`
	SuspendedElevators map[string]types.SuspendedType
}

func ToHRAInput(hallRequests [N_FLOORS][2]bool, cabRequests map[string][N_FLOORS]bool, elevatorStates map[string]types.ElevatorState, suspendedElevators map[string]types.SuspendedType) HRAInput {

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
func ManageOrders(OrdersCh chan HRAInput, AssignmentsCh chan map[string][N_FLOORS][2]bool) {

	resendticker := time.NewTicker(config.Cfg.ResendAssignmentTime)
	var cachedAssignment map[string][N_FLOORS][2]bool

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
		select {
		case input := <-OrdersCh:
			fmt.Println("Calculating orders")
			

			for id, state := range input.SuspendedElevators {
				if state.IsSuspended {
					delete(input.States, id)
					fmt.Printf("Deleted elevatorstate: %s\n", id)
				}
			}

			if len(input.States) == 0 {
				fmt.Printf("No input states provided to ordermanager, skipping.\n")
				AssignmentsCh <- make(map[string][N_FLOORS][2]bool)
				continue
			}

			jsonBytes, err := json.Marshal(input)
			if err != nil {
				fmt.Println("json.Marshal error: ", err)
				return
			}

			ret, err := exec.Command("./src/ordermanager/"+hraExecutable, "-i", string(jsonBytes)).CombinedOutput()
			if err != nil {
				fmt.Println("exec.Command error: ", err)
				fmt.Println(string(ret))
				return
			}

			output := new(map[string][N_FLOORS][2]bool)

			err = json.Unmarshal(ret, &output)
			if err != nil {
				fmt.Println("json.Unmarshal error: ", err)
				return
			}

			resendticker.Reset(config.Cfg.ResendAssignmentTime)
			cachedAssignment = *output
			AssignmentsCh <- *output
		case <-resendticker.C:
			AssignmentsCh <- cachedAssignment
		}
	}
}
