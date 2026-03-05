package network

import (
	"Network-go/network/bcast"
	"Network-go/network/conn"
	"fmt"
	"net"
	"time"
)

//The goal of this module is to see how many slaves the master have.
// This is done by sending heartbeats,
// listen to reponses snd then manage a list of the active slaves.

//Input: Only has to be started/stopped when promoted/revoked to/from master.
//Out: Can be pulled to return the current list of slaves.

type HelloMsg struct {
	Message string
	Iter    int
}

const heartBeatInterval = 15 * time.Millisecond

func MasterActions(port int, id string, isMaster <-chan bool) {

	conn := conn.DialBroadcastUDP(port)
	addr, _ := net.ResolveUDPAddr("udp4", fmt.Sprintf("255.255.255.255:%d", port))

	enable := true
	for {
		select {
		case enable = <-isMaster:
		case <-time.After(heartBeatInterval):
		}
		if enable {
			conn.WriteTo([]byte(id), addr)
		}
	}
	
}

func SlaveActions(id string, done chan bool){
	
	reciveBrodcast := make(chan HelloMsg)

	go bcast.Receiver(16569, reciveBrodcast)
	for {
		a := <-reciveBrodcast
		fmt.Printf("Received: %#v\n", a)
	}
	done <- true
}