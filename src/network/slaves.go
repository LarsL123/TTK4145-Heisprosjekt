package network

import (
	"Network-go/network/bcast"
	"elevatorproject/src/config"
	"fmt"
)

//Alot of hadoutcode from "peers" but modifies to master/slave architecture.

// The goal of this module is to see how many slaves the master have.
// This is done by sending heartbeats,
// listen to reponses snd then manage a list of the active slaves.

//Input: Only has to be started/stopped when promoted/revoked to/from master.
//Out: Can be pulled to return the current list of slaves.

//Ting som skal være i reeelection og ikke her:
// Check for double master.
// Check for no master.

func ReplyToHeartbeat(id string){ //Add port to Config file. 
	recive := make(chan Heartbeat)
	go bcast.Receiver(config.Cfg.HeartbeatPort, recive)

	send := make(chan Heartbeat)
	go bcast.Transmitter(config.Cfg.SlaveReplyPort, send)
	

	fmt.Println("Reciving...")
	for {
		beat := <-recive
		fmt.Printf("Received heartbeat from id %s that is a %s\n", beat.ID, beat.Role)

		reply := Heartbeat{id, "slave"}
		send <-reply
	}
}







// func SlaveActions(id string, conn net.PacketConn){

// 	conn := conn.DialBroadcastUDP(port)

// 	var buf [1024]byte
	
// 	for {
// 		n, _, _ := conn.ReadFrom(buf[0:])

// 		id := string(buf[:n])

// 		if id != ""{
// 			fmt.Println("Recived from id: ", id)
// 			// conn.WriteTo([]byte(id + ":ack"), addr)
// 		}
// 	}

// 	// done <- true

// }