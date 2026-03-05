package main

import (
	"elevatorproject/src/network"
	"flag"
)

//TODO: Denne fila burde bare initialisere modulene våre,
// og kanalene de snakker sammen på

func main(){
	// Our id can be anything. Here we pass it on the command line, using
	//  `go run main.go -id=our_id`
	var id string
	flag.StringVar(&id, "id", "", "id of this peer")
	flag.Parse()

	
	done := make(chan bool)	
	isMaster := make(chan bool)
	if (id == "1"){
		go network.MasterActions(15647, id, isMaster)
	}else{
		go network.SlaveActions(id, done)
	}
	

	<-done
}