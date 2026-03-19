package controllers

import (
	"Network-go/network/bcast"
	"elevatorproject/src/config"
	"elevatorproject/src/types"
	"fmt"
)

func RunBackup(isMasterCh chan bool, forwardOrders chan types.OrderEnvelope) {
	receiveBackupDataCh := make(chan types.BackupData, 10)
	sendBackupAckCh := make(chan types.BackupDataAck, 10)

	go bcast.Receiver(config.Cfg.BackupSendPort, receiveBackupDataCh)
	go bcast.Transmitter(config.Cfg.BackupReceivePort, sendBackupAckCh)

	var savedData types.BackupData

	for {
		select {
		case savedData = <-receiveBackupDataCh:
			sendBackupAckCh <- types.BackupDataAck{UpdateNr: savedData.UpdateNr} 
		case isMaster := <-isMasterCh:
			if isMaster {
				pushOrdersToMaster(savedData, forwardOrders)
			}
		}
	}
}

func pushOrdersToMaster(data types.BackupData, forwardOrders chan<- types.OrderEnvelope) {
	for floor := range N_FLOORS {
		for btn := range 2 {
			if data.HallRequests[floor][btn] {
				fmt.Println("Pushing HALL_ORDER from backup", floor, types.OrderType(btn).ToString())
				forwardOrders <- types.OrderEnvelope{ 
					Order: types.Order{Floor: floor, Type: types.OrderType(btn)}, // Important to not use any of the other values in the orderenvelope as they are zero-values
				}
			}
		}
	}

	// Sending all cabOrders
	for id, cabOrder := range data.CabRequests {
		for floor := range N_FLOORS {
			if cabOrder[floor] {
				fmt.Println("Pushing CAB_ORDER to master from backup floor ", floor)
				forwardOrders <- types.OrderEnvelope{
					ElevatorID: id,
					Order: types.Order{Floor: floor,Type:  types.Cab}, // Important to not use any of the other values in the orderenvelope as they are zero-values
				}
			}
		}
	}
}
