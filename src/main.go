package main

import (
	"elevatorproject/src/config"
	"elevatorproject/src/network"
	"flag"
	"fmt"
)

//TODO: Denne fila burde bare initialisere modulene våre,
// og kanalene de snakker sammen på

func main(){
	config.Load()
	//VIKTIG: Hvis man trenger å endre i denne fila, feks for teste andre deler av systemt så comment ut det som er her. 
	
	// Our id can be anything. Here we pass it on the command line, using
	//  `go run main.go -id=our_id`
	var id string
	flag.StringVar(&id, "id", "", "id of this elevator")
	flag.Parse()

	isMaster := make(chan bool)
	peerUpdate := make(chan network.PeerUpdate)


	if (id == "1"){
		go network.SendHeartbeats(id, isMaster)
		go network.TrackSlaves(peerUpdate)

		for {
			p := <-peerUpdate
			fmt.Printf("Slave update:\n")
			fmt.Printf("  Slaves:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)
		}
	}else{
		go network.ReplyToHeartbeat(id)
	}
	
	select {}
}