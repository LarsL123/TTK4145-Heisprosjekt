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

func ReelectionFSM(selfID string, isMasterCh chan bool, backupMasterCh chan bool) {

	role := Slave

	// Setting this a little bit higher so that startup goes fine
	masterTimer := time.NewTicker(config.Cfg.NewMasterTimeoutTime * 3)
	backupTimer := time.NewTicker(config.Cfg.NewBackupTimeoutTime * 3)

	// Setting up heartbeat
	heartbeatTicker := time.NewTicker(config.Cfg.HeartbeatInterval)
	sendHeartbeatCh := make(chan Heartbeat)
	go bcast.Transmitter(config.Cfg.HeartbeatPort, sendHeartbeatCh)
	heartbeatCh := make(chan Heartbeat, 32)
	go bcast.Receiver(config.Cfg.HeartbeatPort, heartbeatCh)

	startRole := func(r Role) {
		role = r //Er dette feil sjef?? Ja dette er en bug, kan ende opp med å blocke channels....

		switch r {
		case Master:
			fmt.Println("I am master:", selfID)
			isMasterCh <- true
			backupMasterCh <- true

		case Backup:
			fmt.Println("I am backup:", selfID)
			isMasterCh <- false
			backupMasterCh <- true

		case Slave:
			fmt.Println("I am slave:", selfID)
			isMasterCh <- false
			backupMasterCh <- false
		}
	}

	role = Slave
	fmt.Println("I am slave:", selfID)

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
						backupTimer.Reset(config.Cfg.NewBackupTimeoutTime * 3)
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
