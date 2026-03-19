package main

import (
	"elevatorproject/src/config"
	"elevatorproject/src/controllers"
	"elevatorproject/src/processpair"
	"elevatorproject/src/reelection"
	"elevatorproject/src/types"
	"os"
)

func main() {
	id := os.Args[1]
	config.Load()

	masterAliveCh := make(chan struct{}, 1)
	slaveAliveCh := make(chan struct{}, 1)

	toggleMaster := make(chan bool, 2)
	toggleBackup := make(chan bool, 2)

	forwardOrdersFromBackup := make(chan types.OrderEnvelope, 12)  // backup -> master
	ordersFromMasterCollison := make(chan types.OrderEnvelope, 12) // master killing itself because it detects another master -> slave

	master := controllers.NewMaster(id, toggleMaster)
	go master.Start(ordersFromMasterCollison, forwardOrdersFromBackup, masterAliveCh)

	slave := controllers.NewSlave(id)
	go slave.Start(ordersFromMasterCollison, slaveAliveCh)

	go controllers.RunBackup(toggleBackup, forwardOrdersFromBackup)

	go reelection.Reelection(id, toggleMaster, toggleBackup)

	processpair.KillDeadProcess(masterAliveCh, slaveAliveCh)
}
