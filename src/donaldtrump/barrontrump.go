package donaldtrump

import (
	"Network-go/network/bcast"
	"elevatorproject/src/config"
	"elevatorproject/src/types"
)

func RunBackup(isMasterCh chan bool, forwardOrders chan types.OrderEnvelope) {
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
				pushOrdersToMaster(savedData, forwardOrders)
			}
		}
	}
}

func pushOrdersToMaster(data types.BackupData, forwardOrders chan<- types.OrderEnvelope) {
	// Sending all hallorders (Should probably be made into a function) started making pushOrdersToNewMaster, double check if it's correct before switching
	for floor := range N_FLOORS {
		for btn := range 2 {
			if data.HallRequests[floor][btn] {
				forwardOrders <- types.OrderEnvelope{ //Might be an idea to not send every order, but to send everything together when you have looped through here
					Order: types.Order{Floor: floor, Type: types.OrderType(btn)}, // Important to not use any of the other values in the orderenvelope as they are zero-values
				}
			}
		}
	}

	// Sending all cabOrders
	for id, cabOrder := range data.CabRequests {
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
