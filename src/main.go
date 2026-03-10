package main

import (
	"elevatorproject/src/config"
	"elevatorproject/src/donaldtrump"
	"flag"
)

//TODO: Denne fila burde bare initialisere modulene våre,
// og kanalene de snakker sammen på

func main(){
	config.Load()
	
	// Our id can be anything. Here we pass it on the command line, using
	//  `go run main.go -id=our_id`
	var id string
	flag.StringVar(&id, "id", "", "id of this elevator")
	flag.Parse()

	
	if (id == "1"){
		donaldtrump.RunMasterBrain(id)
	}else{
		donaldtrump.RunSlaveBrain(id)
	}
}