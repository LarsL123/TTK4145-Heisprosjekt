package donaldtrump

import (
	"Network-go/network/bcast"
	"elevatorproject/src/config"
	"elevatorproject/src/types"
)

 func RunBackup(isMasterCh <-chan bool, forwardOrdersCh chan types.Order){

 	var data types.BackupData

	reciveData := make(chan types.BackupData)
	go bcast.Receiver(config.Cfg.BackupPort, reciveData)


	for {
		select {
		case isMaster := <- isMasterCh:
			if(isMaster){
				pushOrdersToNewMaster(forwardOrdersCh, data)
			}
		}

	}
}

func  pushOrdersToNewMaster(forwardOrdersCh chan types.Order, data types.BackupData){
	for floor := range N_FLOORS{
		for button := range 2 {
			if(data.HallRequests[floor][button]){
				forwardOrdersCh <- types.Order{Floor: floor, Type: types.OrderType(button)}
			}
		}
	}

	for _,cabs := range data.CabRequests{
		for floor := range N_FLOORS{
			if(cabs[floor]){
				forwardOrdersCh <- types.Order{Floor: floor, Type: types.Cab}
			}
		}

	}
}