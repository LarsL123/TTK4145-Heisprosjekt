package reelection

import (
	"fmt"
	"time"

	"Network-go/network/bcast"
	"elevatorproject/src/config"
)

type Role int

const (
	Slave  Role = 0
	Backup Role = 1
	Master Role = 2
)

type Heartbeat struct {
	ID   string
	Role Role 
}

func ReelectionFSM(selfID string, isMasterCh chan bool) {

	role := Slave

	masterTimer := time.NewTicker(config.Cfg.NewMasterTimeoutTime)
	backupTimer := time.NewTicker(config.Cfg.NewBackupTimeoutTime)

	// Setting up heartbeat
	heartbeatTicker := time.NewTicker(config.Cfg.HeartbeatInterval)
	sendHeartbeatCh := make(chan Heartbeat)
	go bcast.Transmitter(config.Cfg.HeartbeatPort, sendHeartbeatCh)
	heartbeatCh := make(chan Heartbeat, 32)
	go bcast.Receiver(config.Cfg.HeartbeatPort, heartbeatCh)

	startRole := func(r Role) {
		isMasterCh <- false //Er dette feil sjef??

		role = r

		switch r {
		case Master:
			fmt.Println("I am master:", selfID)
			isMasterCh <- true

		case Backup:
			fmt.Println("I am backup:", selfID)

		case Slave:
			fmt.Println("I am slave:", selfID)
		}
	}

	startRole(Slave)

	for {
		select {

		case hb := <-heartbeatCh:

			if hb.ID == selfID {
				continue
			}

			switch hb.Role {
			case Master:
				masterTimer.Reset(config.Cfg.NewMasterTimeoutTime)

				if role == Master {
					if hb.ID > selfID {
						fmt.Println("Higher ID master detected → stepping down")
						startRole(Slave)
					}
				}

			case Backup:
				backupTimer.Reset(config.Cfg.NewBackupTimeoutTime)

				if role == Backup {
					if hb.ID > selfID {
						fmt.Println("Higher ID backup detected → stepping down")
						startRole(Slave)
					}
				}
			}

		case <-masterTimer.C:
			if role == Backup {
				fmt.Println("Master dead → becoming master")
				startRole(Master)
			}

		case <-backupTimer.C:
			if role == Slave {
				fmt.Println("No backup → becoming backup")
				startRole(Backup)
			}

		case <-heartbeatTicker.C:
			if role != Slave {
				sendHeartbeatCh <- Heartbeat{ID: selfID, Role: role}
			}

		}
	}
}
