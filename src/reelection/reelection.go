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