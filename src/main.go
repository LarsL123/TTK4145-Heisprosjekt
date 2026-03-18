package main

import (
	"elevatorproject/src/config"
	"elevatorproject/src/donaldtrump"
	"elevatorproject/src/reelection"
	"flag"
)

// TODO: (slik jeg (Daniel) ser det nå (natt til mandag 16.03))
// 1. Få til kommunikasjonen mellom slaves og master
// 2. Etter det må vi få til initialiseringen, noe jeg tror vi har glemt litt
// 3. Dukker sikkert opp noe mer vi har glemt
// 4. Rename Donald Trump og James David Vance

func main() {
	config.Load()

	// Our id can be anything. Here we pass it on the command line, using
	//  `go run main.go -id=our_id`
	var id string
	flag.StringVar(&id, "id", "", "id of this elevator")
	flag.Parse()

	isMaster := make(chan bool)
	// dataFromBackup := make(chan types.BackupData)

	master := donaldtrump.NewMaster(id, isMaster)
	go master.Start()

	go reelection.ReelectionFSM(id, isMaster)
	go donaldtrump.RunSlaveBrain(id)

	select {}
}
