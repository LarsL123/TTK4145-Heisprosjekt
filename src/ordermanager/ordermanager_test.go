package ordermanager

/*



*/

import (
    "testing"
    "fmt"
)

func TestCostFunction(t *testing.T) {

	ChOrders := make(chan HRAInput)
    ChAssignments := make(chan HRAInput)

	input := HRAInput{
        HallRequests: [][2]bool{{false, false}, {true, false}, {false, false}, {false, true}},
        States: map[string]HRAElevState{
            "one": {
                Behavior:       "moving",
                Floor:          3,
                Direction:      "up",
                CabRequests:    []bool{false, false, false, true},
            },
            "two": {
                Behavior:       "idle",
                Floor:          1,
                Direction:      "stop",
                CabRequests:    []bool{false, false, false, false},
            },
        },
    }

    go main(ChOrders, ChAssignments)

    // Pass testing-input-data
    ChOrders <- input

    for {
        


        output := <- ChAssignments
        
        
        // Print results for testing
        fmt.Printf("output: \n")
        for k, v := range *output {
            fmt.Printf("%6v :  %+v\n", k, v)
        }
    }
}