package main

import (
	"elevatorproject/src/config"
	"elevatorproject/src/donaldtrump"
	"elevatorproject/src/reelection"
	"elevatorproject/src/types"
	"fmt"
	"os"
	"os/exec"
	"time"
)

func main_oldaswell() {

	if os.Getenv("IS_CHILD") == "1" {
		id := os.Getenv("ELEVATOR_ID")
		if id == "" {
			panic("ELEVATOR_ID not set")
		}
		RunWithWatchdog(id)
		return
	}

	if len(os.Args) < 2 {
		panic("You should run: ./elevator <id>")
	}

	id := os.Args[1]

	for {
		cmd := exec.Command(os.Args[0], os.Args[1:]...)
		cmd.Env = append(os.Environ(),
			"IS_CHILD=1",
			"ELEVATOR_ID="+id,
		)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		fmt.Println("Starting child process...")
		cmd.Run()

		fmt.Println("Child died - restarting in 200ms")
		time.Sleep(200 * time.Millisecond)

	}
}

func runSystem(id string,masterAliveCh chan struct{}, slaveAliveCh chan struct{}) { // Was called main before
	config.Load()

	// Our id can be anything. Here we pass it on the command line, using
	//  `go run main.go -id=our_id`
	//var id string
	// flag.StringVar(&id, "id", "", "id of this elevator")
	// flag.Parse()

	isMasterCh := make(chan bool, 1)
	forwardOrders := make(chan types.Order)

	master := donaldtrump.NewMaster(id, isMasterCh, forwardOrders)
	go master.Start(masterAliveCh)
	go reelection.ReelectionFSM(id, isMasterCh)
	//go donaldtrump.RunBackup(isMaster, forwardOrders) //TODO: fix this, has to be commented out as long as backup doesn't work
	go donaldtrump.RunSlaveBrain(id, forwardOrders,slaveAliveCh)
}

func main() {
	id := os.Args[1]

	masterAliveCh := make(chan struct{}, 1)
	slaveAliveCh := make(chan struct{}, 1)

	runSystem(id,masterAliveCh,slaveAliveCh)

	timeout := 5 * time.Second
	startupGrace := 10*time.Second
	masterAliveTimer := time.NewTimer(startupGrace)
	slaveAliveTimer := time.NewTimer(startupGrace)
	for {
		select {
		case <-masterAliveCh:
			masterAliveTimer.Reset(timeout)
		case <-slaveAliveCh:
			slaveAliveTimer.Reset(timeout)
		case <-masterAliveTimer.C:
			fmt.Println("Master dead - kms")
			os.Exit(1)
		case <-slaveAliveTimer.C:
			fmt.Println("Slave dead - kms")
			os.Exit(1)
		}

	}
}

func RunWithWatchdog(id string) {
	masterAliveCh := make(chan struct{}, 1)
	slaveAliveCh := make(chan struct{}, 1)

	isMasterCh := make(chan bool, 1)
	forwardOrders := make(chan types.Order)

	master := donaldtrump.NewMaster(id, isMasterCh, forwardOrders)
	go master.Start(masterAliveCh)
	go reelection.ReelectionFSM(id, isMasterCh)
	go donaldtrump.RunSlaveBrain(id, forwardOrders, slaveAliveCh)

	watchdogTimeout := config.Cfg.WatchdogTimeout
	masterTimer := time.NewTimer(watchdogTimeout)
	slaveTimer := time.NewTimer(watchdogTimeout)

	for {
		select {
		case <-masterAliveCh:
			masterTimer.Reset(watchdogTimeout)
		case <-slaveAliveCh:
			slaveTimer.Reset(watchdogTimeout)
		case <-masterTimer.C:
			fmt.Println("Master goroutine dead - exiting and restarting with process pair")
			os.Exit(1)
		case <-slaveTimer.C:
			fmt.Println("Slave goroutine dead - exiting and restarting with process pair")
			os.Exit(1)
		}
	}
}
