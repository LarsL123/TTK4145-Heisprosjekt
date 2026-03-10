package network

//Divide into slave network/master network.

/*
In master network we do: (BIG PICTURE)
	1: Recive new orders and state (have to ack) Done

	2: Send AssignmentsAndOrder to slave (would like ack)
	3: Send updated order list slave (for backup and turning off lights)
	4: recive order done.



In slave Network we do:
	1: func SendNewOrdersAndStateToMaster(sendNewOrderChannel chan bool) (ack and resend) done (missing resend)

    2: func ReciveAssignemnetsFromMaster(reciveOrderChannel chan bool)
	3: Recive updated order list.
	4: Send orderDone

In Backup Network:
	Recive sync, acc.

WE HAR ENDRA DET. mASTER POLLER OG SLAVE  BARE SVARER.

*/

