package reelection

import (
	// "golang.org/x/text/cases"
	"Network-go/network/bcast"
	"elevatorproject/src/config"
	"elevatorproject/src/network"
	"time"

	"golang.org/x/sys/windows/registry"
)

// Spørsmål til studass, hvordan burde man ordne reelection?
// Burde alle slavene kunne være backup, eller skal man ha en spesifikk backup?
// I så fall, hvis alle er backup, så trenger man vel
// at alle acknowledger at en order er mottatt fra en annen heis, hvis ikke kan man jo miste orders

// Siste problemet: Hva gjør man om en master skal integreres inn i et system der det er en master og en slave fra før?

/* Denne modulen skal:
1. Vite om master er død
2. Dersom master er død, gi beskjed om at vi må velge en ny master.
3. Bestemme hvem som blir ny master?
*/

/*

This module is practically an event based finite state machine.

 Input:
	channel: heartbeats

 Output:
	channel: heartbeats

 Purpose:
	Make sure one and only one master exists at all times.
	Internally, we need to handle the cases:
		1 master, continue
		< 1 master, elect exactly one master
		> 1 master, Remove all masters, elect exactly one master
*/


func setSlave(hb network.Heartbeat) {
	// TODO: implement
	hb.Role = 
}

func setMaster(hb network.Heartbeat) {
	// TODO: implement
	// Depends on heartbeat struct to implement
}

func reelect() {
	// Idea: Sleep for a short time. If no master is alive, elect yourself

	/*

	randomNumber = randint()
	sleep(random)

	if !masterIsAlive {
		setMaster(self.ID)
		masterIsAlive = True
	}
	*/
}

func deelectAll() {
	/*

	for id in IDs {
		setSlave(id)
	}

	masterIsAlive = False

	*/
}

func detectMasterConflict(registry map[hb.ID]hb.Role, conflictChannel chan) {

	count := 0

	for _, role := range registry {
		if role == MASTER { // Possible bug. Should it be "MASTER"?
			count++
		}
	}

	if count > 1 {
		select {
		case conflictChannel <- {}{}:
		default: // Important to not hault the system
		}
	}

}

func RunReelectionFSM() {

	conflictDetectedCh := make(chan struct{}, 1) // buffered of size 1
	registry := map[hb.ID]hb.Role{}

	// If the bcast use below is bad, for now it is written for testing
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

		// Timer has passed, value appears in channel: watchdog.

		case <- conflictDetectedCh:
			deelectAll()
			reelect()
		}
	}
}