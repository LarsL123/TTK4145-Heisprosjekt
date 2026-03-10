package network

import (
	"Network-go/network/bcast"
	"Network-go/network/localip"
	"elevatorproject/src/config"
	"fmt"
	"sort"
	"time"
)

type Role int

const (
	Slave Role = 0
	Backup Role = 1
	Master Role = 2
)

type SlaveUpdate struct {
	Slaves []string
	New   string
	Lost  []string
}

type Heartbeat struct{
	ID string
	Role Role //Slave/Master//TODO Make enum.  ---> Brage satt opp enum for Role. Sjekk om det gir mening.
	IP string //Not in use. Hope we dont need it. 
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


	heartbeat := Heartbeat{id, Master, ip}
	enable := true

	for {
		select {
			case enable = <-isMaster:
			case <-time.After(config.Cfg.HeartbeatInterval):
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

			case <-time.After(config.Cfg.HeartbeatInterval):
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
			if time.Since(v) > config.Cfg.HeartbeatTimeout {
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


