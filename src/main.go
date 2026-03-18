package main

import (
	"elevatorproject/src/config"
	"elevatorproject/src/donaldtrump"
	"elevatorproject/src/reelection"
	"elevatorproject/src/types"
	"flag"
)

func main() {
	config.Load()

	// Our id can be anything. Here we pass it on the command line, using
	//  `go run main.go -id=our_id`
	var id string
	flag.StringVar(&id, "id", "", "id of this elevator")
	flag.Parse()

	isMaster := make(chan bool)
	forwardOrders := make(chan types.Order)

	master := donaldtrump.NewMaster(id, isMaster, forwardOrders)
	go master.Start()
	go reelection.ReelectionFSM(id, isMaster)
	go donaldtrump.RunBackup(isMaster, forwardOrders) //TODO
	go donaldtrump.RunSlaveBrain(id, forwardOrders)

	select {}
}
