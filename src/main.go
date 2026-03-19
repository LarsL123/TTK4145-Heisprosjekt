package main

import (
	"elevatorproject/src/config"
	"elevatorproject/src/donaldtrump"
	"elevatorproject/src/reelection"
	"elevatorproject/src/types"
	"fmt"
	"os"

	//"os/exec"
	"time"
)

func runSystem(id string, masterAliveCh chan struct{}, slaveAliveCh chan struct{}) { // Was called main before
	config.Load()

	isMasterCh := make(chan bool, 2)
	backupMasterCh := make(chan bool, 2)

	//Buffered for 12 which is N_FLOORS*3 which should be worst case
	forwardOrdersFromBackup := make(chan types.OrderEnvelope, 12) // backup-> master
	forwardOrdersFromDead := make(chan types.OrderEnvelope, 12)   // master killing itself because it detects another master -> slave

	master := donaldtrump.NewMaster(id, isMasterCh, forwardOrdersFromDead, forwardOrdersFromBackup, masterAliveCh)
	go master.Start()
	go reelection.ReelectionFSM(id, isMasterCh, backupMasterCh)
	go donaldtrump.RunBackup(backupMasterCh, forwardOrdersFromBackup)
	go donaldtrump.RunSlaveBrain(id, forwardOrdersFromDead, slaveAliveCh)
}

func main() {
	id := os.Args[1]

	masterAliveCh := make(chan struct{}, 1)
	slaveAliveCh := make(chan struct{}, 1)

	runSystem(id, masterAliveCh, slaveAliveCh)

	timeout := 5 * time.Second
	startupGrace := 10 * time.Second
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

// func RunWithWatchdog(id string) {
// 	masterAliveCh := make(chan struct{}, 1)
// 	slaveAliveCh := make(chan struct{}, 1)

// 	isMasterCh := make(chan bool, 1)
// 	forwardOrders := make(chan types.Order)

// 	master := donaldtrump.NewMaster(id, isMasterCh, forwardOrders)
// 	go master.Start(masterAliveCh)
// 	go reelection.ReelectionFSM(id, isMasterCh)
// 	go donaldtrump.RunSlaveBrain(id, forwardOrders, slaveAliveCh)

// 	watchdogTimeout := config.Cfg.WatchdogTimeout
// 	masterTimer := time.NewTimer(watchdogTimeout)
// 	slaveTimer := time.NewTimer(watchdogTimeout)

// 	for {
// 		select {
// 		case <-masterAliveCh:
// 			masterTimer.Reset(watchdogTimeout)
// 		case <-slaveAliveCh:
// 			slaveTimer.Reset(watchdogTimeout)
// 		case <-masterTimer.C:
// 			fmt.Println("Master goroutine dead - exiting and restarting with process pair")
// 			os.Exit(1)
// 		case <-slaveTimer.C:
// 			fmt.Println("Slave goroutine dead - exiting and restarting with process pair")
// 			os.Exit(1)
// 		}
// 	}
// }
