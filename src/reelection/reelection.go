package reelection

// ---
// Use the goroutine below to setup the FSM in all elevators.
// From there, all reelection logic is handled internally inside the FSM.
//
// go InitReelection(selfID)
//
// The FSM is fully agnostic to other elevators states, and state changes are
// fully dependent on if heartbeats from either master or backup arrives in time.
//
// The idea is to always ensure one master and one backup is alive.
// ---

import (
	"Network-go/network/bcast"
	"context"
	"elevatorproject/src/config"
	"elevatorproject/src/network"
	"fmt"
	"time"
	// "golang.org/x/text/cases"
	// "golang.org/x/sys/windows/registry"
)

// A FSM running goroutines based on current role
func InitReelection(selfId string) {

	heartbeatCh := make(chan network.Heartbeat)
	masterHearbeatCh := make(chan string) // Buffered? Check
	backupHeartbeatCh := make(chan string) // Buffered? Check
	roleCh := make(chan network.Role, 1) // Should only be one role at a time
	currentRoleCh := make(chan network.Role)

	// Receives and splits heartbeats based on role into separate channels
	go bcast.Receiver(config.Cfg.HeartbeatPort, heartbeatCh)
	go SplitHeartbeats(heartbeatCh, masterHearbeatCh, backupHeartbeatCh)

	// IMPORTANT: Tramsmitter logic, should it be elsewhere?
	// var currentRole network.Role = network.Slave
	go bcast.Transmitter(config.Cfg.HeartbeatPort, currentRoleCh) // IMPORTANT: MAKE SURE IT IS ONLY SPAWNED ONCE PER AGENT

	// To allow goroutine killing
	var cancel context.CancelFunc

	roleCh <- network.Slave

	for role := range roleCh {
		if cancel != nil {
            cancel()
        }

		// To allow goroutine killing
		ctx, newCancel := context.WithCancel(context.Background())
		cancel = newCancel

		switch role {

		case network.Master:
			fmt.Println("I am master:  ", selfId)
			// currentRoleCh <-currentRole
			go DetectMasterConflict(selfId, masterHearbeatCh, roleCh, cancel)
			go network.StartMaster(selfId, ctx)

		case network.Backup:
			fmt.Println("I am backup:  ", selfId)
			go DetectBackupConflict(selfId, backupHeartbeatCh, roleCh, ctx)
			go BackupSelfElectMaster(selfId, masterHearbeatCh, roleCh, cancel)
		
		case network.Slave:
			fmt.Println("I am slave:  ", selfId)
			go SlaveSelfElectBackup(selfId, backupHeartbeatCh, roleCh)
		}
	}
}

func InitReelectionMasterTest(selfId string) {

	heartbeatCh := make(chan network.Heartbeat)
	masterHearbeatCh := make(chan string) // Buffered? Check
	backupHeartbeatCh := make(chan string) // Buffered? Check
	roleCh := make(chan network.Role, 1) // Should only be one role at a time
	currentRoleCh := make(chan network.Role)

	// Receives and splits heartbeats based on role into separate channels
	go bcast.Receiver(config.Cfg.HeartbeatPort, heartbeatCh)
	go SplitHeartbeats(heartbeatCh, masterHearbeatCh, backupHeartbeatCh)

	// IMPORTANT: Tramsmitter logic, should it be elsewhere?
	// var currentRole network.Role = network.Slave
	go bcast.Transmitter(config.Cfg.HeartbeatPort, currentRoleCh) // IMPORTANT: MAKE SURE IT IS ONLY SPAWNED ONCE PER AGENT

	// To allow goroutine killing
	var cancel context.CancelFunc

	roleCh <- network.Master

	for role := range roleCh {
		if cancel != nil {
            cancel()
        }

		// To allow goroutine killing
		ctx, newCancel := context.WithCancel(context.Background())
		cancel = newCancel

		switch role {

		case network.Master:
			fmt.Println("I am master:  ", selfId)
			// currentRoleCh <-currentRole
			go DetectMasterConflict(selfId, masterHearbeatCh, roleCh, cancel)
			go network.StartMaster(selfId, ctx)

		case network.Backup:
			fmt.Println("I am backup:  ", selfId)
			go DetectBackupConflict(selfId, backupHeartbeatCh, roleCh, ctx)
			go BackupSelfElectMaster(selfId, masterHearbeatCh, roleCh, cancel)
		
		case network.Slave:
			fmt.Println("I am slave:  ", selfId)
			go SlaveSelfElectBackup(selfId, backupHeartbeatCh, roleCh)
		}
	}
}

// Split heartbeats into separte channels, only containing their IDs
func SplitHeartbeats(heartbeatCh chan network.Heartbeat, masterHeartbeatCh chan string, backupHeartbeatCh chan string) {

	for {
		heartbeat := <- heartbeatCh

		fmt.Println("Walla hjertet banker", heartbeat.ID)
		
		if (heartbeat.Role == network.Master) {
			masterHeartbeatCh <- heartbeat.ID
		}
		if (heartbeat.Role == network.Backup) {
			backupHeartbeatCh <- heartbeat.ID
		}
	}
}

// If more masters exists, set self to slave
func DetectMasterConflict(selfID string, masterHeartbeatCh chan string, roleCh chan network.Role, cancel context.CancelFunc) {
	prevMasterID := ""

	for {
		masterHeartbeat := <- masterHeartbeatCh

		if (prevMasterID == "") {
			prevMasterID = masterHeartbeat
		}

		// Conflict detected -> Suicide
		if (masterHeartbeat != prevMasterID) {
			fmt.Println("DOUBLE MASTER DETECTED!!! -> Suicide ")
			prevMasterID = ""
			cancel()
			roleCh <- network.Slave
			return
		}

		prevMasterID = masterHeartbeat
	}
}

// If more backups exists, set self to slave
func DetectBackupConflict(selfID string, backupHeartbeatCh chan string, roleCh chan network.Role,  ctx context.Context) {
	prevBackupID := ""

	for {
		select{
			case backupHeartbeat := <- backupHeartbeatCh:
				// Conflict detected -> Suicide and kill goroutine: BackupSelfElectMaster()
				if (backupHeartbeat != prevBackupID) {
					roleCh <- network.Slave
					return
				}

				// Stash for comparison
				if (backupHeartbeat == prevBackupID || backupHeartbeat == "") { // nil for second expression
					prevBackupID = backupHeartbeat
				}

			case <- ctx.Done():
				fmt.Println("Backup conflic function stopped", selfID)
				return
		}
	}
}


// No master exists? Set yourself to master
func BackupSelfElectMaster(selfID string, masterHeartbeatCh chan string, roleCh chan network.Role,  cancel context.CancelFunc) {
	for {
		select {
		// Heartbeat received -> reset timer
		case <- masterHeartbeatCh:
			continue
			
		// Timeout -> suicide
		case <- time.After(time.Second * 2):
			cancel()
			roleCh <- network.Master
			return 	
		}
	}	
}

// No backup exists? Set yourself to backup
func SlaveSelfElectBackup(selfID string, backupHeartbeatCh chan string, roleCh chan network.Role) {

	for {
		select {
			// Heartbeat received -> reset timer
			case <- backupHeartbeatCh:
				continue
			
			// Timeout -> suicide
			case <- time.After(config.Cfg.HeartbeatTimeout):
				fmt.Println("Gidder ikke slave. ")
				roleCh <- network.Backup
				return
		}
	}


}