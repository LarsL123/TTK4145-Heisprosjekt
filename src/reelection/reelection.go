package reelection

/*

This module is present in all elevators.

A central part of the module is the roles map, which maps all ids
and roles to each other, on the format below. The idN is just for 
explanation and is not how the ids are set up in the finished system.

roles[string]network.Role {

	"id1": SLAVE,
	"id2": SLAVE,
	"id3": MASTER,
	"id4": SLAVE,
	"id5": BACKUP,
	...

}

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



// Elects a random slave to be backup
func ElectRandomSlave(roles map[string]network.Role) {

	var idRandom string

	// Looping through a map leads to a random id being chosen
	for id, _ := range roles {
		idRandom = id
		if (roles[idRandom] == network.Slave) {
			break
		}
	}

	// Assign the random slave the role of backup
	roles[idRandom] = network.Backup

}

// Ensures all elevators are slaves
func SetAllToSlave(roles map[string]network.Role) {

	for k := range roles {
		
		roles[k] = network.Slave

	}
}

// Sends a signal to the conflictDetectedCh if more than one
// master is detected in the roleRegistry
func DetectMasterConflict(roles map[string]network.Role, conflictDetectedCh chan struct{}) {

	count := 0

	for _, role := range roles {
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

// A goroutine only given to the backup.
// If master dies or conflict is detected, set itself to master
func ReelectMaster(roles map[string]network.Role, selfId string) {

	// Receive heartbeats
	heartbeatCh := make(chan network.Heartbeat)
	bcast.Receiver(config.Cfg.HeartbeatPort, heartbeatCh)

	conflictDetectedCh := make(chan struct{}, 1) // buffered of size 1
	noMasterCh := make(chan struct{}, 1) // buffered of size 1

	for {

		select {
		
		case <- noMasterCh:
			if roles[selfId] == network.Backup { // If-statement could lead to bugs...
				roles[selfId] = network.Master
			}
		case <- conflictDetectedCh:
			SetAllToSlave(roles)
			roles[selfId] = network.Master

		}
	}
}

// A goroutine only given to the master.
// If no backup exist, elect a new backup
func ReelectBackup(roles map[string]network.Role) {

	conflictDetectedCh := make(chan struct{}, 1) // buffered of size 1
	heartbeatCh := make(chan network.Heartbeat)
	bcast.Receiver(config.Cfg.HeartbeatPort, heartbeatCh)
	watchdog := time.NewTimer(timeout)

	for {

		select {
		
		case heartbeat := <- heartbeatCh:
		// Heartbeat is received
			roles[heartbeat.ID] = heartbeat.Role
			DetectMasterConflict(roles, conflictDetectedCh)

			if (!watchdog.Stop()) {
				<- watchdog.C
			}

			watchdog.Reset(timeout)

		case <- watchdog.C:
		// Triggered by watchdog timeout
			ElectRandomSlave(roles)
			
		}
	}
}

// The main goroutine of the module, run by all elevators.
// Starts a goroutine for reelection logic considering its role.
func SetupReelection(roleCh chan network.Role, selfId string) {

	roles := map[string]network.Role{}

	
	for role := range roleCh {

		switch role {

		case network.Master:
			go ReelectBackup(roles)
		case network.Backup:
			go ReelectMaster(roles, selfId)
		case network.Slave:
			// Slave does not have reelection responsibilities:
			// Do nothing
		}
	}
}