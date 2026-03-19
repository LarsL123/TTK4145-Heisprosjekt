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

func Reelection(selfID string, toggelMaster chan bool, toggleBackup chan bool) {
	role := Slave

	// Setting this a little bit higher so that startup goes fine
	masterTimer := time.NewTicker(config.Cfg.NewMasterTimeoutTime * 3)
	backupTimer := time.NewTicker(config.Cfg.NewBackupTimeoutTime * 3)

	heartbeatTicker := time.NewTicker(config.Cfg.HeartbeatInterval)

	sendHeartbeatCh := make(chan Heartbeat)
	receiveHeartbeatCh := make(chan Heartbeat, 32)

	go bcast.Transmitter(config.Cfg.HeartbeatPort, sendHeartbeatCh)
	go bcast.Receiver(config.Cfg.HeartbeatPort, receiveHeartbeatCh)

	startRole := func(r Role) {
		role = r

		switch r {
		case Master:
			fmt.Println("I am master:", selfID)
			toggelMaster <- true
			toggleBackup <- true

		case Backup:
			fmt.Println("I am backup:", selfID)
			toggelMaster <- false
			toggleBackup <- true

		case Slave:
			fmt.Println("I am slave:", selfID)
			toggelMaster <- false
			toggleBackup <- false
		}
	}

	role = Slave
	fmt.Println("I am slave:", selfID)

	for {
		select {

		case hb := <-receiveHeartbeatCh:

			if hb.ID == selfID {
				continue
			}

			resetTimers(masterTimer, backupTimer)
			handleDuplicate(hb, role, selfID, startRole)

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

func resetTimers(masterTimer *time.Ticker, backupTimer *time.Ticker) {
	masterTimer.Reset(config.Cfg.NewMasterTimeoutTime)
	backupTimer.Reset(config.Cfg.NewBackupTimeoutTime)
}

func handleDuplicate(hb Heartbeat, role Role, selfID string, startRole func(Role)) {
	switch hb.Role {
	case Master:
		if role == Master {
			if hb.ID > selfID {
				fmt.Println("Higher ID master detected → stepping down")
				startRole(Slave)
			}
		}

	case Backup:
		if role == Backup {
			if hb.ID > selfID {
				fmt.Println("Higher ID backup detected → stepping down")
				startRole(Slave)
			}
		}
	}
}
