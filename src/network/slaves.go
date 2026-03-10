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

// var masterip string  //Dont know hoe to do this properly.

func StartSlave(id string) *OrderSender{
	go ReplyToHeartbeat(id)

	sendOrdersCh := make(chan OrdersAndStateUpdate)
	ackCh := make(chan OrdersAndStateAck)

	go bcast.Transmitter(config.Cfg.MasterListenPort, sendOrdersCh)
	go bcast.Receiver(config.Cfg.SlaveListenPort, ackCh)

	orderSender := &OrderSender{
		SendCh: sendOrdersCh,
		AckIn: ackCh,
		AckResults: make(chan AckResult, 10), // buffered
	}

	return orderSender
}

func ReplyToHeartbeat(id string){
	receive := make(chan Heartbeat)
	go bcast.Receiver(config.Cfg.HeartbeatPort, receive)

	send := make(chan Heartbeat)
	go bcast.Transmitter(config.Cfg.SlaveHeartbeatReplyPort, send)

	// ip, err := localip.LocalIP()
	// if err != nil{
	// 	fmt.Println("Error when getting LocalIP")
	// 	fmt.Println("Aborting")
	// 	return
	// }
	
	fmt.Println("Reciving...")
	for {
		// beat := <-receive
		<-receive
		// fmt.Printf("Received heartbeat from id %s that is a %s with ip %s \n", beat.ID, beat.Role, beat.IP)

		// if beat.IP != masterip{
		// 	masterip = beat.IP
		// }

		reply := Heartbeat{id, "slave", ""}
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