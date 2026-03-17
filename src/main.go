package main

import (
	"elevatorproject/src/config"
	"elevatorproject/src/donaldtrump"
	"flag"
	"fmt"
)

// TODO: (slik jeg (Daniel) ser det nå (natt til mandag 16.03))
// 1. Få til kommunikasjonen mellom slaves og master
// 2. Etter det må vi få til initialiseringen, noe jeg tror vi har glemt litt
// 3. Dukker sikkert opp noe mer vi har glemt
// 4. Rename Donald Trump og James David Vance

func main(){
	config.Load()
	
	// Our id can be anything. Here we pass it on the command line, using
	//  `go run main.go -id=our_id`
	var id string
	var port string
	flag.StringVar(&id, "id", "", "id of this elevator")
	flag.StringVar(&port, "port", "", "id of this elevator")
	flag.Parse()

	if port != ""{
		config.Cfg.ElevatorPort = port
	}
	fmt.Println("Running on port: ", config.Cfg.ElevatorPort)
	

	
	if (id == "1"){
		donaldtrump.RunMasterBrain(id)
	}else{
		donaldtrump.RunSlaveBrain(id, port)
	}
}