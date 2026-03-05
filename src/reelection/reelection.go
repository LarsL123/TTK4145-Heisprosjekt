package reelection
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
	ListenUDP: Heartbeat{ID, ...}

 Output:
	???
 Purpose:
	Make sure one and only one master exists at all times
 */

 /*

int backupCount; 

function Reelect(backup) {
	if backupCount == 0
		setMaster(self.ID)
	else
		masterID <- pick random ID
		setMaster(masterID)
}

function DeelectAll {
	for ID in elevators:
		setSlave(ID)
}


select:
	case receivedheartbeat := <- heartbeats
		reset -> watchdog

	case expiredWatchdog := <- watchdog
		Reelect()

	if more than one master
		DeelectAll()
		Reelect()

 */