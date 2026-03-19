package donaldtrump

import (
	"Network-go/network/bcast"
	"elevatorproject/src/config"
	"elevatorproject/src/types"
)

func RunBackup(id string, isMasterCh chan bool, forwardOrders chan types.OrderEnvelope) {
	_ = id // Should probably be fixed by removing id as paramterer to RunBackup
	receiveBackupDataCh := make(chan types.BackupData, 10)
	sendBackupAckCh := make(chan types.BackupDataAck, 10)

	go bcast.Receiver(config.Cfg.BackupSendPort, receiveBackupDataCh)
	go bcast.Transmitter(config.Cfg.BackupReceivePort, sendBackupAckCh)

	var savedData types.BackupData

	for {
		select {
		case data := <-receiveBackupDataCh:
			savedData = data
			sendBackupAckCh <- types.BackupDataAck{UpdateNr: data.UpdateNr}
		case isMaster := <-isMasterCh:
			if isMaster {

				// Sending all hallorders (Should probably be made into a function) started making pushOrdersToNewMaster, double check if it's correct before switching
				for floor := range N_FLOORS {
					for btn := range 2 {
						if savedData.HallRequests[floor][btn] {
							forwardOrders <- types.OrderEnvelope{ //Might be an idea to not send every order, but to send everything together when you have looped through here
								Order: types.Order{Floor: floor, Type: types.OrderType(btn)}, // Important to not use any of the other values in the orderenvelope as they are zero-values
							}
						}
					}
				}

				// Sending all cabOrders
				for id, cabOrder := range savedData.CabRequests {
					for floor := range N_FLOORS {
						if cabOrder[floor] {
							forwardOrders <- types.OrderEnvelope{
								ElevatorID: id,
								Order: types.Order{
									Floor: floor,
									Type:  types.Cab, // Important to not use any of the other values in the orderenvelope as they are zero-values
								},
							}
						}
					}
				}
			}
		}
	}
}

//Started implementing this function: instead of the shit above here ^

// func pushOrdersToNewMaster(forwardOrdersCh chan types.OrderEnvelope, data types.BackupData) {
// 	for floor := range N_FLOORS {
// 		for button := range 2 {
// 			if data.HallRequests[floor][button] {
// 				forwardOrdersCh <- types.OrderEnvelope{Order: types.Order{Floor: floor, Type: types.OrderType(button)}}
// 			}
// 		}
// 	}

// 	for elevID, cabReq := range data.CabRequests {
// 		for floor := range N_FLOORS {
// 			if cabReq[floor] {
// 				forwardOrdersCh <- types.OrderEnvelope{
// 					ElevatorID: elevID,
// 					Order:      types.Order{Floor: floor, Type: types.Cab},
// 				}
// 			}
// 		}
// 	}
// }

//This is dead code
// func RunBackupLarsSin(isMasterCh <-chan bool, forwardOrdersCh chan types.Order) {

// 	var data types.BackupData

// 	reciveData := make(chan types.BackupData)
// 	go bcast.Receiver(config.Cfg.BackupPort, reciveData)

// 	for {
// 		select {
// 		case isMaster := <-isMasterCh:
// 			if isMaster {
// 				pushOrdersToNewMaster(forwardOrdersCh, data)
// 			}
// 		}

// 	}
// }
