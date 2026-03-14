package ordermanager

import (
	"fmt"
	"testing"
)

func TestCostFunction(t *testing.T) {

	OrdersCh := make(chan HRAInput)
    AssignmentsCh := make(chan map[string][][2]bool)

    // Set test input
	input := HRAInput{
        HallRequests: [4][2]bool{{false, false}, {true, false}, {false, false}, {false, true}},
        States: map[string]HRAElevState{
            "one": {
                Behavior:       "moving",
                Floor:          3,
                Direction:      "down",
                CabRequests:    [4]bool{false, false, false, true},
            },
            "two": {
                Behavior:       "idle",
                Floor:          0,
                Direction:      "stop",
                CabRequests:    [4]bool{false, false, false, false},
            },
        },
    }

    go ManageOrders(OrdersCh, AssignmentsCh)

    // Pass testing-input-data
    OrdersCh <- input

    for {

        // Receive optimal assignment
        output := <- AssignmentsCh

        // Print results
        fmt.Printf("output: \n")
        for k, v := range output {
            fmt.Printf("%6v :  %+v\n", k, v)
        }
    }
}