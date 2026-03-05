package network

//Alot of andoutcode from "peers" but modifies to master/slave architecture.

import (
	"Network-go/network/conn"
	"fmt"
	"net"
	"sort"
	"time"
)

// The goal of this module is to see how many slaves the master have.
// This is done by sending heartbeats,
// listen to reponses snd then manage a list of the active slaves.

//Input: Only has to be started/stopped when promoted/revoked to/from master.
//Out: Can be pulled to return the current list of slaves.

//Ting som skal være i reeelection og ikke her:
// Check for double master.
// Check for no master.

type HelloMsg struct {
	Message string
	Iter    int
}

type PeerUpdate struct {
	Peers []string
	New   string
	Lost  []string
}

const heartBeatInterval = 15 * time.Millisecond
const timeout = 500 * time.Millisecond

func MasterActions(port int, id string, isMaster <-chan bool, peerUpdate chan PeerUpdate) {

	conn := conn.DialBroadcastUDP(port)
	addr, _ := net.ResolveUDPAddr("udp4", fmt.Sprintf("255.255.255.255:%d", port))

	go SendHeartbeats(conn, addr, isMaster)
	//go ReceveFromSlaves(conn, peerUpdate)


	
}

func ReceveFromSlaves(conn net.PacketConn, peerUpdateCh chan PeerUpdate){
	var buf [1024]byte
	var p PeerUpdate
	lastSeen := make(map[string]time.Time)

	for {
		updated := false

		conn.SetReadDeadline(time.Now().Add(heartBeatInterval))
		n, _, _ := conn.ReadFrom(buf[0:])

		id := string(buf[:n])

		// Adding new connection
		p.New = ""
		if id != "" {
			if _, idExists := lastSeen[id]; !idExists {
				p.New = id
				updated = true
			}

			lastSeen[id] = time.Now()
		}

		// Removing dead connection
		p.Lost = make([]string, 0)
		for k, v := range lastSeen {
			if time.Now().Sub(v) > timeout {
				updated = true
				p.Lost = append(p.Lost, k)
				delete(lastSeen, k)
			}
		}

		// Sending update
		if updated {
			p.Peers = make([]string, 0, len(lastSeen))

			for k, _ := range lastSeen {
				p.Peers = append(p.Peers, k)
			}

			sort.Strings(p.Peers)
			sort.Strings(p.Lost)
			peerUpdateCh <- p
		}
	}
}

func SendHeartbeats(conn net.PacketConn, addr *net.UDPAddr, isMaster <-chan bool){
	enable := true
	for {
		select {
		case enable = <-isMaster:
		case <-time.After(heartBeatInterval):
		}
		if enable {
			conn.WriteTo([]byte("master"), addr)
		}
	}

}

func SlaveActions(id string, conn net.PacketConn){

	var buf [1024]byte
	
	for {
		n, _, _ := conn.ReadFrom(buf[0:])

		id := string(buf[:n])

		if id != ""{
			fmt.Println("Recived from id: ", id)
			// conn.WriteTo([]byte(id + ":ack"), addr)
		}
	}
	

	
	// done <- true

}