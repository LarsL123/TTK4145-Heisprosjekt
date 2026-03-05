package ordermanager
// Denne modulen kjører bare på master.
// Den skal holde styr på hvilken heis som skal ta hvilken ordre.
// Den skal få inn nye ordre fra orderManager, regne ut hvem som skal ta hvilken ordre og sende det videre til de andre heisene
// Master må også ha kontroll over alle slavene sine, hvordan skal dette implementeres?
// Master må sende ut heartbeats, men må slavesa gjøre det og??
//  - Nei, slaven svarer iAmSlave


/*
Input:
	channel: Orders

Output:
	channel: Assignments

Purpose: 
    Computes the optimal assignments based on orders
*/

// The code below is a modified version of example.go from project_resources


import "os/exec"
import "fmt"
import "encoding/json"
import "runtime"

// Struct members must be public in order to be accessible by json.Marshal/.Unmarshal
// This means they must start with a capital letter, so we need to use field renaming struct tags to make them camelCase

type HRAElevState struct {
    Behavior    string      `json:"behaviour"`
    Floor       int         `json:"floor"` 
    Direction   string      `json:"direction"`
    CabRequests []bool      `json:"cabRequests"`
}

type HRAInput struct {
    HallRequests    [][2]bool                   `json:"hallRequests"`
    States          map[string]HRAElevState     `json:"states"`
}



func main(){

    hraExecutable := ""
    switch runtime.GOOS {
        case "linux":   hraExecutable  = "hall_request_assigner"
        case "windows": hraExecutable  = "hall_request_assigner.exe"
        default:        panic("OS not supported")
    }

    input := HRAInput{
        HallRequests: [][2]bool{{false, false}, {true, false}, {false, false}, {false, true}},
        States: map[string]HRAElevState{
            "one": HRAElevState{
                Behavior:       "moving",
                Floor:          2,
                Direction:      "up",
                CabRequests:    []bool{false, false, false, true},
            },
            "two": HRAElevState{
                Behavior:       "idle",
                Floor:          0,
                Direction:      "stop",
                CabRequests:    []bool{false, false, false, false},
            },
        },
    }

    jsonBytes, err := json.Marshal(input)
    if err != nil {
        fmt.Println("json.Marshal error: ", err)
        return
    }
    
    ret, err := exec.Command("./"+hraExecutable, "-i", string(jsonBytes)).CombinedOutput()
    if err != nil {
        fmt.Println("exec.Command error: ", err)
        fmt.Println(string(ret))
        return
    }
    
    output := new(map[string][][2]bool)
    err = json.Unmarshal(ret, &output)
    if err != nil {
        fmt.Println("json.Unmarshal error: ", err)
        return
    }
        
    fmt.Printf("output: \n")
    for k, v := range *output {
        fmt.Printf("%6v :  %+v\n", k, v)
    }

    // Pass output data to channel for module: orderDistributor to take
    Assignments := make(chan map[string][][2]bool) 
    Assignments <- *output;
}