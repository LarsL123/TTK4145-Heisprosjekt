package reelection

import (
	"fmt"
	"time"

	"Network-go/network/bcast"
	"elevatorproject/src/config"
	"elevatorproject/src/network"
)

// func InitReelection2(selfID string) {

// 	go ReelectionFSM(selfID)
// }

func ReelectionFSM(selfID string, isMasterCh chan bool) {

	role := network.Slave

	masterTimer := time.NewTicker(config.Cfg.NewMasterTimeoutTime) // TODO: Add to config
	backupTimer := time.NewTicker(config.Cfg.NewBackupTimeoutTime)

	// Setting up heartbeat
	heartbeatTicker := time.NewTicker(config.Cfg.HeartbeatInterval)
	sendHeartbeatCh := make(chan network.Heartbeat)
	go bcast.Transmitter(config.Cfg.HeartbeatPort, sendHeartbeatCh)
	heartbeatCh := make(chan network.Heartbeat, 32)
	go bcast.Receiver(config.Cfg.HeartbeatPort, heartbeatCh)

	startRole := func(r network.Role) {
		isMasterCh <- false //Er dette feil sjef??

		role = r

		switch r {
		case network.Master:
			fmt.Println("I am master:", selfID)
			isMasterCh <- true

		case network.Backup:
			fmt.Println("I am backup:", selfID)

		case network.Slave:
			fmt.Println("I am slave:", selfID)
		}
	}

	startRole(network.Slave)

	for {
		select {

		case hb := <-heartbeatCh:

			if hb.ID == selfID {
				continue
			}

			switch hb.Role {
			case network.Master:
				masterTimer.Reset(config.Cfg.NewMasterTimeoutTime)

				if role == network.Master {
					if hb.ID > selfID {
						fmt.Println("Higher ID master detected → stepping down")
						startRole(network.Slave)
					}
				}

			case network.Backup:
				backupTimer.Reset(config.Cfg.NewBackupTimeoutTime)

				if role == network.Backup {
					if hb.ID > selfID {
						fmt.Println("Higher ID backup detected → stepping down")
						startRole(network.Slave)
					}
				}
			}

		case <-masterTimer.C:
			if role == network.Backup {
				fmt.Println("Master dead → becoming master")
				startRole(network.Master)
			}

		case <-backupTimer.C:
			if role == network.Slave {
				fmt.Println("No backup → becoming backup")
				startRole(network.Backup)
			}

		case <-heartbeatTicker.C:
			if role != network.Slave {
				sendHeartbeatCh <- network.Heartbeat{ID: selfID, Role: role}
			}

		}
	}
}
