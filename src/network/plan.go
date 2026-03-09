package network

//Divide into slave network/master network.

/*
In master network we do: (BIG PICTURE)
	1: Poll slaves when needed.
	2: Send AssignmentsAndOrder
	3: Recive datachange, send to backup, acc to sender.  (brodcast)
	4: recive order done. (brodcast) //HELP: do we need acc?)

	Sync to backup and gucci.


In slave Network we do:
	1: func SendNewOrdersAndStateToMaster(sendNewOrderChannel chan bool) (ack and resend maybe)
    2: func ReciveAssignemnetsAndOrderFromMaster(reciveOrderChannel chan bool)
	3: Send datachange and take acc (bordcast)
	4: Send orderDone (bordcast)

In Backup Network:
	Recive sync, acc.

WE HAR ENDRA DET. mASTER POLLER OG SVALVE ER STATELESS OG BARE SVARER.

*/

