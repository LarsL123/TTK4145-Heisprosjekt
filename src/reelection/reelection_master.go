package reelection

/*

This module is only available to the master.
Its main task is to ensure 

*/

import (
	// "golang.org/x/text/cases"
	"Network-go/network/bcast"
	"elevatorproject/src/config"
	"elevatorproject/src/network"
	"time"

	// "golang.org/x/sys/windows/registry"
)

/*

Pseudocode:



*/


func ReelectBackup() {
	receive := make(chan network.Heartbeat)
	go bcast.Receiver(config.Cfg.HeartbeatPort, receive) // Bad to receive in reelect? Better to use channels?

	watchdog := time.NewTimer(timeout)

	var electionRunning bool = false

	for {

		select {

		case heartbeat <- receive:
		// Reset timer when heartbeat is received

			// Will fire to conflictDetectedCh if more than one masters
			registry[heartbeat.id] = heartbeat.Role
			detectMasterConflict(registry, conflictDetectedCh)

			if (!watchdog.Stop()) {
				<- watchdog.C
			}
			watchdog.Reset(timeout)

		case <- watchdog.C:
			reelect()
			watchdog.Reset(timeout)
		}
	}
}