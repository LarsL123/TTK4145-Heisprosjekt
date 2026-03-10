package reelection

/*

This module is present in all elevators. The main loop is set to trigger

*/

import (

	"Network-go/network/bcast"
	"elevatorproject/src/config"
	"elevatorproject/src/network"
	"time"

	// "golang.org/x/text/cases"
	// "golang.org/x/sys/windows/registry"

)

// config values defined here only for testing, will be removed
const heartBeatInterval = 1000 * time.Millisecond //Change to 15ms
const timeout = 2000 * time.Millisecond //Change to 500ms

// 
func ChooseRandomSlave(roleRegistry map[string]network.Role) {

	var idRandom string

	// Pick random id from registry
	for id, _ := range roleRegistry {
		idRandom = id
		if (roleRegistry[idRandom] == network.Slave) {
			break
		}
	}

	// Assign the random slave the role of backup
	roleRegistry[idRandom] = network.Backup

}

func SetAllToSlave(roleRegistry map[string]network.Role) {

	for k, _ := range roleRegistry {
		
		roleRegistry[k] = network.Slave

	}
}

func DetectMasterConflict(roleRegistry map[string]network.Role, conflictDetectedCh chan struct{}) {

	count := 0

	for _, role := range roleRegistry {
		if role == network.Master {
			count++
		}
	}

	if count > 1 {

		select {

		case conflictDetectedCh <- struct{}{}:
		default: // Important to not hault the system
		
		}
	}

}

func ReelectMaster(roleRegistry map[string]network.Role) {

	// Receive heartbeats
	heartbeatCh := make(chan network.Heartbeat)
	bcast.Receiver(config.Cfg.HeartbeatPort, heartbeatCh)
	conflictDetectedCh := make(chan struct{}, 1) // buffered of size 1
	noMasterCh := make(chan struct{}, 1) // buffered of size 1

	for {

		select {
		
		case <- noMasterCh:
			// TODO: self.Role = master
		case <- conflictDetectedCh:
			SetAllToSlave(roleRegistry)


		}
	}
}

func ReelectBackup(roleRegistry map[string]network.Role) {

	conflictDetectedCh := make(chan struct{}, 1) // buffered of size 1
	heartbeatCh := make(chan network.Heartbeat)
	bcast.Receiver(config.Cfg.HeartbeatPort, heartbeatCh)
	watchdog := time.NewTimer(timeout)

	for {

		select {
		
		case heartbeat := <- heartbeatCh:
		// Heartbeat is received
			roleRegistry[heartbeat.ID] = heartbeat.Role
			DetectMasterConflict(roleRegistry, conflictDetectedCh)

			if (!watchdog.Stop()) {
				<- watchdog.C
			}

			watchdog.Reset(timeout)
		case <- watchdog.C:
		// Triggered by watchdog timeout
			ChooseRandomSlave()
			
		}
	}
}

func SetupReelection(roleCh chan network.Role) {

	roleRegistry := map[string]network.Role{}
	
	for role := range roleCh {

		switch role {

		case network.Master:
			go ReelectBackup(roleRegistry)
		case network.Backup:
			go ReelectMaster(roleRegistry)
		case network.Slave:
			break
		
		}
	}
}