package reelection

import (
	// "golang.org/x/text/cases"
	"time"
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
 Input:
	channel: heartbeats

 Output:
	channel: heartbeats

 Purpose:
	Make sure one and only one master exists at all times
*/

/*

LES:
Det meste under er bare brainstorm for å få ned tanker


select:
	case receivedheartbeat := <- heartbeats
		reset -> watchdog

	case expiredWatchdog := <- watchdog
		Reelect()

	if more than one master
		DeelectAll()
		Reelect()

*/


func setSlave(id string) {
	// TODO: implement
	// Depends on heartbeat struct to implement
}

func setMaster() {
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

func main() {

	heartbeats := make(chan string) // Should be struct, of type ID, heartbeats should be initialized in network or peers
	masterCount := make(chan int)

	var maxtime = 2 * time.Second // High value just for testing. Lower for operations
	watchdog := time.NewTimer(maxtime)


	for {

		select {

		case <- heartbeats:
		// Reset timer when heartbeat is received
			if !watchdog.Stop() {
				<- watchdog.C
			}
			watchdog.Reset(maxtime)

		case <- watchdog.C:
		// Timer has passed, value appears in channel: watchdog.C
			reelect()
			watchdog.Reset(maxtime)

		case master := <- masterCount: 
		// If a master is detected
			if master > 1 {
				deelectAll()
				reelect()
			}
		}
	}
}