package network

import (
	"Network-go/network/bcast"
	"elevatorproject/src/config"
	"fmt"
	"sort"
	"time"
)

const heartBeatInterval = 1000 * time.Millisecond //Change to 15ms
const timeout = 2000 * time.Millisecond //Change to 500ms

type PeerUpdate struct {
	Peers []string
	New   string
	Lost  []string
}

type Heartbeat struct{
	ID string
	Role string //Slave/Master//TODO Make enum.
}

func SendHeartbeats(id string, isMaster <-chan bool) {
	//Burde isMaster vekk?
	send := make(chan Heartbeat)
	go bcast.Transmitter(config.Cfg.HeartbeatPort, send)

	heartbeat := Heartbeat{id, "master"}
	enable := true

	for {
		select {
			case enable = <-isMaster:
			case <-time.After(heartBeatInterval):
		}

		if enable {
			send <- heartbeat
		}
	}
}

func TrackSlaves( peerUpdateCh chan<- PeerUpdate){
	var p PeerUpdate
	lastSeen := make(map[string]time.Time)

	recive := make(chan Heartbeat)
	go bcast.Receiver(config.Cfg.SlaveReplyPort, recive)

	for {
		var id string
		
		select {
			case acc := <- recive:
				if acc.ID == "" {
					fmt.Println("Got invalid packet")
					continue
				}

				id = acc.ID

			case <-time.After(heartBeatInterval):
		}

		updated := false

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
