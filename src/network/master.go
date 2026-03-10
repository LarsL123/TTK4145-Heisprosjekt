package network

import (
	"Network-go/network/bcast"
	"Network-go/network/localip"
	"elevatorproject/src/config"
	"fmt"
	"sort"
	"time"
)

// Fra Brage til Lars: burde vi sette alle kontanter i config i stedet? Gir iaf ferre includes...
const heartBeatInterval = 1000 * time.Millisecond //Change to 15ms
const timeout = 2000 * time.Millisecond //Change to 500ms

type SlaveUpdate struct {
	Slaves []string
	New   string
	Lost  []string
}

type Heartbeat struct{
	ID string
	Role string //Slave/Master//TODO Make enum.
	IP string
}

func StartMaster(id string, isMaster chan bool) chan SlaveUpdate{
		go SendHeartbeats(id, isMaster)

		heartBeatCh := make(chan Heartbeat)
		go bcast.Receiver(config.Cfg.SlaveHeartbeatReplyPort, heartBeatCh)

		slaveUpdate := make(chan SlaveUpdate)
		go TrackSlaves(heartBeatCh, slaveUpdate)

		return slaveUpdate
}

func SendHeartbeats(id string, isMaster <-chan bool) {
	//Burde isMaster vekk?
	sendCh := make(chan Heartbeat)
	go bcast.Transmitter(config.Cfg.HeartbeatPort, sendCh)

	ip, err := localip.LocalIP()

	if err != nil{
		fmt.Println("Failed to get IP. What do we do?")
	}


	heartbeat := Heartbeat{id, "master", ip}
	enable := true

	for {
		select {
			case enable = <-isMaster:
			case <-time.After(heartBeatInterval):
		}

		if enable {
			sendCh <- heartbeat
		}
	}
}

func TrackSlaves(heartBeatCh <-chan Heartbeat, slaveUpdateCh chan<- SlaveUpdate){
	lastSeen := make(map[string]time.Time)

	for {
		slaveID := ""
		
		select {
			case acc := <- heartBeatCh:
				if acc.ID == "" {
					fmt.Println("Got invalid packet")
					continue
				}

				slaveID = acc.ID

			case <-time.After(heartBeatInterval):
		}

	
		p := SlaveUpdate{}
		updated := false

		// Adding new connection
		p.New = ""
		if slaveID != "" {
			if _, idExists := lastSeen[slaveID]; !idExists {
				p.New = slaveID
				updated = true
			}

			lastSeen[slaveID] = time.Now()
		}

		// Removing dead connection
		p.Lost = make([]string, 0)
		for k, v := range lastSeen {
			if time.Since(v) > timeout {
				updated = true
				p.Lost = append(p.Lost, k)
				delete(lastSeen, k)
			}
		}

		// Sending update
		if updated {
			p.Slaves = make([]string, 0, len(lastSeen))

			for k := range lastSeen {
				p.Slaves = append(p.Slaves, k)
			}

			// fmt.Println(slaveIP)

			sort.Strings(p.Slaves)
			sort.Strings(p.Lost)

			slaveUpdateCh <- p
		}
	}
}


