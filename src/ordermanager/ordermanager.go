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
    "time"
	// "golang.org/x/text/cases"
)

// Struct members must be public in order to be accessible by json.Marshal/.Unmarshal
// This means they must start with a capital letter, so we need to use field renaming struct tags to make them camelCase

type HRAElevState struct {
    Behavior    string      `json:"behaviour"`
    Floor       int         `json:"floor"` 
    Direction   string      `json:"direction"`
    CabRequests [4]bool      `json:"cabRequests"`
}

type HRAInput struct {
    HallRequests    [4][2]bool                   `json:"hallRequests"`
    States          map[string]HRAElevState     `json:"states"`
}

// Combines hallrequests and elevatorstates
func ToHRAInput(hallRequests [4][2]bool, elevatorStates map[string] types.ElevatorState) HRAInput{

    inputStates := make(map[string]HRAElevState)

    for id, elevatorState := range elevatorStates{
        inputStates[id] = HRAElevState{
            Behavior: elevatorState.Behaviour,
            Floor:elevatorState.Floor,
            Direction: elevatorState.Direction,
            CabRequests: elevatorState.CabRequests,
        }
    }

    return HRAInput{
        HallRequests: hallRequests,
        States: inputStates,
    }
}

/*
// Fetches the states of a single elevator
func HRAInputToState(input HRAInput, id string) HRAElevState{

    return input.States[id]
    
}
*/

func orderTimeout(id string, state types.ElevatorState, deadCh chan<- string) {
    orderTimer := time.After(10 * time.Second)

    var blockTimer <-chan time.Time
    if state.Obstructed {
        blockTimer = time.After(5 * time.Second)
    }

    select {
    case <-orderTimer:
        fmt.Printf("Elevator %s did not complete assignment in time\n", id)
        deadCh <- id

    case <-blockTimer: // nil if not obstructed, so never fires
        fmt.Printf("Elevator %s obstructed too long\n", id)
        deadCh <- id
    }
}


// Calculates optimal assignments based on orders
func ManageOrders(inputCh chan types.RawMasterData, assignmentsCh chan map[string][4][2]bool, deadCh chan<- string){

    hraExecutable := ""
    switch runtime.GOOS {
        case "linux":   hraExecutable  = "hall_request_assigner"
        case "windows": hraExecutable  = "hall_request_assigner.exe"
        default:        panic("OS not supported")
    }

    for {

        // Orders are received
        rawOrder := <- inputCh
        // HER SKAL order INNEHOLDE OBSTRUCTION, order må være av typen types.ElevatorState

        // Filter obstructed elevators before HRA
        filtered := make(map[string]types.ElevatorState)
        for id, state := range rawOrder.States {
            if !state.Obstructed {
                filtered[id] = state
            }
        }

        hraOrder := ToHRAInput(rawOrder.HallRequests, filtered)

        // Parse using ToHRAInput(hallRequests, elevatorStates)
        // Important for the rest of the program to work
        //

        fmt.Println(hraOrder)

        // JSON -> String
        jsonBytes, err := json.Marshal(hraOrder)
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

        assignment := new(map[string][4][2]bool)

        // Update output map with executable data, String -> JSON
        err = json.Unmarshal(ret, &assignment)
        if err != nil {
            fmt.Println("json.Unmarshal error: ", err)
            return
        }

        // Start timers only for elevators that got assignments
        for id, orders := range *assignment {
            if hasActiveOrder(orders) {
                go orderTimeout(id, rawOrder.States[id], deadCh)
            }
        }

        // Pass optimal assignments
        assignmentsCh <- *assignment;
    }
}

func hasActiveOrder(orders [4][2]bool) bool {
    for _, floor := range orders {
        for _, req := range floor {
            if req {
                return true
            }
        }
    }
    return false
}


/*

Brage note (16.03.2026 20:30):

Head is not working, tomorrow I need to setup logic such that:

- Each elevator that gets an assignment from the cost function gets a timer on their head to complete it. If the timer is not passed, then suspend the elevator
- Each elevator that gets obstructed for a set amount of time gets suspended. If obstruction stops, then remove suspension

*/