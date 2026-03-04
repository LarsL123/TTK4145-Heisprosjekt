// Denne modulen kjører bare på master.
// Den skal holde styr på hvilken heis som skal ta hvilken ordre.
// Den skal få inn nye ordre fra orderManager, regne ut hvem som skal ta hvilken ordre og sende det videre til de andre heisene
// Master må også ha kontroll over alle slavene sine, hvordan skal dette implementeres?
// Master må sende ut heartbeats, men må slavesa gjøre det og??
//  - Nei, slaven svarer iAmSlave