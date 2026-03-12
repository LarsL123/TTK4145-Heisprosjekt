package reelection

/*

This module is present in all elevators.

go SetupReelection()

The idea behind the case at the bottom of the file is that:
	if role == MASTER
		Ensure a backup is alive, if not find a new
	if role == BACKUP
		Ensure a master is alive, if not elect yoursellf
	if role == SLAVE
		Do nothing

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
var heartBeatInterval = config.Cfg.HeartbeatInterval //Change to 15ms
var timeout = config.Cfg.HeartbeatTimeout          //Change to 500ms

// Ensures roles-map does not store multiple elevators as backup
func ClearAllBackups(roles map[string]network.Role) {
	
	for k := range roles {
		if roles[k] == network.Backup {
		roles[k] = network.Slave
		}
	}	

}

// Elects a random slave to be backup
func ElectSlaveToBackup(roles map[string]network.Role) {

	var idRandom string

	// Looping through a map leads to a random id being chosen
	for id, _ := range roles {
		idRandom = id
		if roles[idRandom] == network.Slave {
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

// Sends a signal to a channel if more than one
// master or backup is detected in the roleRegistry
func DetectConflict(roles map[string]network.Role, masterConflictDetectedCh chan struct{}, backupConflictDetectedCh chan struct{}) {

	// IMPORTANT: conflictDetectedCh must be buffered to work,
	// e.g. conflictCh := make(chan struct{}, 1)

	masterCount := 0
	backupCount := 0

	for _, role := range roles {
		if role == network.Master {
			masterCount++
		}
		if role == network.Backup {
			backupCount++
		}
	}

	if masterCount > 1 {

		select {

		case masterConflictDetectedCh <- struct{}{}:
		default: // Important to not hault the system

		}
	}

	if backupCount > 1 {

		select {

		case backupConflictDetectedCh <- struct{}{}:
		default: // Important to not hault the system

		}
	}

}

// Signals if no master exists
func DetectNoMaster(roles map[string]network.Role, noMasterCh chan struct{}, heartbeatCh chan network.Heartbeat) {
	
	watchdog := time.NewTimer(timeout)
	
	for {

		select {

		case heartbeat := <-heartbeatCh:
		// Heartbeat is received
			
			if !watchdog.Stop() {
				<-watchdog.C
			}

			if heartbeat.Role == network.Master {
				watchdog.Reset(timeout)
			}

		case <-watchdog.C:
		// Watchdog timeout -> Alarm!
			noMasterCh <- struct{}{}

		}
	}

	
}



// A goroutine only given to the backup.
// If master dies or conflict is detected, set itself to master
func ReelectMaster(roles map[string]network.Role, noMasterCh chan struct{}, conflictDetectedCh chan struct{}, selfId string) {

	// Receive heartbeats
	//heartbeatCh := make(chan network.Heartbeat)
	//bcast.Receiver(config.Cfg.HeartbeatPort, heartbeatCh)

	// The encapsulated 
	/*
	conflictDetectedCh := make(chan struct{}, 1) // buffered of size 1
	noMasterCh := make(chan struct{}, 1)         // buffered of size 1
	*/



	for {

		select {

		case <-noMasterCh:
			if roles[selfId] == network.Backup { // If-statement could lead to bugs...
				roles[selfId] = network.Master
			}
		case <-conflictDetectedCh:
			SetAllToSlave(roles)
			roles[selfId] = network.Master

		}
	}
}

// A goroutine only given to the master.
// If no backup exist, elect a new backup

// Added input: heartbeatCh only for testing. REMOVE input argument!
func ReelectBackup(roles map[string]network.Role, heartbeatCh chan network.Heartbeat) {

	masterConflictDetectedCh := make(chan struct{}, 1) // buffered of size 1
	backupConflictDetectedCh := make(chan struct{}, 1) // buffered of size 1
	// heartbeatCh := make(chan network.Heartbeat)
	// go bcast.Receiver(config.Cfg.HeartbeatPort, heartbeatCh)
	watchdog := time.NewTimer(timeout)

	for {

		select {

		case heartbeat := <-heartbeatCh:
		// Heartbeat is received
			// println("Heartbeat received")
			roles[heartbeat.ID] = heartbeat.Role
			DetectConflict(roles, masterConflictDetectedCh, backupConflictDetectedCh)
			if !watchdog.Stop() {
				<-watchdog.C
			}

			watchdog.Reset(timeout)

		case <-watchdog.C:
		// Watchdog timeout
			// FOR TEST
			// println("timeout")
			// FOR TEST

			ClearAllBackups(roles)
			ElectSlaveToBackup(roles)
			watchdog.Reset(timeout)
		}
	}
}

// The main goroutine of the module, run by all elevators.
// Starts a goroutine for reelection logic considering its role.
/*
func SetupReelection(roleCh chan network.Role, selfId string) {

	roles := map[string]network.Role{}
	heartbeatCh := make(chan network.Heartbeat, 1)
	go bcast.Receiver(config.Cfg.HeartbeatPort, heartbeatCh)

	// TODO: Split heartbeats into two channels

	// REMOV
	conflictDetectedCh := make(chan struct{}, 1) // buffered of size 1
	noMasterCh := make(chan struct{}, 1)         // buffered of size 1
	// REMOVE

	for role := range roleCh {

		switch role {

		case network.Master:
			go ReelectBackup(roles, heartbeatCh) // REMOVE heartbeatCh after
		case network.Backup:
			go DetectNoMaster(roles, noMasterCh, heartbeatCh)
			go ReelectMaster(roles, noMasterCh, conflictDetectedCh, selfId)
		case network.Slave:
			// Slave does not have reelection responsibilities:
			// Do nothing
		}
	}
}
*/

/////////////////////////////////// NEW CODE ////////////////////////////////////////////
// WILL REMOVE ALL ABOVE WHEN NEW CODE IS SETUP

// Split heartbeat into separte channels, only containing their IDs
func SplitHeartbeats(heartbeatCh chan network.Heartbeat, masterHeartbeatCh chan string, backupHeartbeatCh chan string) {

	go bcast.Receiver(config.Cfg.HeartbeatPort, heartbeatCh)

	for {

		heartbeat := <- heartbeatCh
		
		if (heartbeat.Role == network.Master) {

			masterHeartbeatCh <- heartbeat.ID

		}

		if (heartbeat.Role == network.Backup) {

			backupHeartbeatCh <- heartbeat.ID

		}
	}


}

func SetupReelection(roleCh chan network.Role, selfId string) {

	heartbeatCh := make(chan network.Heartbeat, 1)
	masterHearbeatCh := make(chan string) // Buffered?
	backupHeartbeatCh := make(chan string) // Buffered?

	// Receive and split heartbeats
	go SplitHeartbeats(heartbeatCh, masterHearbeatCh, backupHeartbeatCh)

	for role := range roleCh {

		switch role {

		case network.Master:
			go DetectMasterConflict(selfId, masterHearbeatCh, roleCh)

		case network.Backup:
			go DetectBackupConflict(selfId, backupHeartbeatCh, roleCh)
			go BackupSelfElectMaster(selfId)
		
		case network.Slave:
			go SlaveSelfElectBackup(selfId, masterHearbeatCh)
		}
	}
}


func DetectMasterConflict(selfID string, masterHeartbeatCh chan string, roleCh chan network.Role) {

prevMasterID := ""

for {

	masterHeartbeat := <- masterHeartbeatCh
		
	if (masterHeartbeat != prevMasterID) {
	// Conflict detected
		
		// Suicide
		roleCh <- network.Slave
		return
	
	}
		
	if (masterHeartbeat == prevMasterID || masterHeartbeat == "") {
	// Stash for comparison
	
		prevMasterID = masterHeartbeat
		
		}
	}
}



func DetectBackupConflict(selfID string, backupHeartbeatCh chan string, roleCh chan network.Role) {

	prevBackupID := ""

	for {
		backupHeartbeat := <- backupHeartbeatCh
	
		prevBackupID = backupHeartbeat
		
		if (backupHeartbeat != prevBackupID) {
		// Conflict detected
		
			// Suicide
			roleCh <- network.Slave
			// TODO: Kill BackupSelfElectMaster()
			return
		}
		
		if (backupHeartbeat == prevBackupID || backupHeartbeat == "") { // nil for second expression
		// Stash for comparison
	
			prevBackupID = backupHeartbeat
		
		}
	}
}



func BackupSelfElectMaster(selfID string, masterHeartbeatCh chan string, roleCh chan network.Role) {

// bcast.Transmit(backupHeartbeatCh) TODO SETUP TRANSMIT
watchdog := time.NewTimer(timeout)

for {

	select {
	
		case <- masterHeartbeatCh:
		watchdog.Reset(timeout)
		
		case <- watchdog.C:
		
		// Suicide
		roleCh <- network.Master
		return 
	
	}

}	

}

func SlaveSelfElectBackup(selfID string, masterHeartBeatCh chan string) {

watchdog := time.NewTimer(timeout)

for {

	select {
	
		case <- masterHeartBeatCh:
		watchdog.Reset(timeout)
		
		case <- watchdog.C:
		// TODO: ELECT SELFID to MASTER
		watchdog.Reset(timeout)
	
	}

}

}