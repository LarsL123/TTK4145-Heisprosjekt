package network



/* Denne modulen skal:
1. Ha mulighet til å sende meldinger til de andre heisene
	- Hvis heisen er en slave trenger den bare å kunne:
		1. Gi beskjed til master når den har en ny order
		2. Svare på iAmAlive meldingene til master.
		3. Acknowledge at det har kommet en ny order.
	- Hvis heisen er en master:
		1. Sende iAmAlive meldinger
		2. Gi beskjed til heiser om at det har kommet en ny order.
		3. Svare tilbake til slavene som sier ifra om at det har kommet en ny order.

TCP eller udp?: Lars researcher forskjeller i bruk.

	- 

*/