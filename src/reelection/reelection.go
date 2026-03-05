package reelection

import "golang.org/x/text/cases"

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

/*

TODO: Implement system for watchdog

struct watchdog{
	time
}

func runWatchdog() {
	for {
	i++

	if i > 1000
	}
}
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

	*/
}

func main() {

	// Initialize all channels, done here for testing
	heartbeats := make(chan string) // Should be struct, of type ID
	watchdog := make(chan int)
	masterCount := make(chan int) // Should only be one

	for {
		select {
		// Check if master is alive
		case receivedheartbeat := <- heartbeats:
			// Reset watchdog


		// No master pings
		case <- watchdog:
			reelect()


		// More than one master
		case master := <- masterCount:
			if master > 1 {
				deelectAll()
				reelect()
			}
			// 
		}
	}


	// ListenUDP: Heartbeat{ID, ...}
	


}