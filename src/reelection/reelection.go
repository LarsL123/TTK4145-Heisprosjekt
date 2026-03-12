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
	"elevatorproject/src/config"
	"elevatorproject/src/network"
	"time"
	"context"
	// "golang.org/x/text/cases"
	// "golang.org/x/sys/windows/registry"
)

// A FSM running goroutines based on current role
func InitReelection(selfId string) {

	heartbeatCh := make(chan network.Heartbeat, 1)
	masterHearbeatCh := make(chan string) // Buffered? Check
	backupHeartbeatCh := make(chan string) // Buffered? Check
	roleCh := make(chan network.Role, 1) // Should only be one role at a time
	currentRoleCh := make(chan network.Role)

	// Receives and splits heartbeats based on role into separate channels
	go bcast.Receiver(config.Cfg.HeartbeatPort, heartbeatCh)
	go SplitHeartbeats(heartbeatCh, masterHearbeatCh, backupHeartbeatCh)

	// IMPORTANT: Tramsmitter logic, should it be elsewhere?
	var currentRole network.Role = network.Slave
	go bcast.Transmitter(config.Cfg.HeartbeatPort, currentRoleCh) // IMPORTANT: MAKE SURE IT IS ONLY SPAWNED ONCE PER AGENT

	// To allow goroutine killing
	var cancel context.CancelFunc

	for role := range roleCh {

		if cancel != nil {
            cancel()
        }

		// To allow goroutine killing
		ctx, newCancel := context.WithCancel(context.Background())
		cancel = newCancel

		switch role {

		case network.Master:
			currentRoleCh <-currentRole
			<-roleCh
			go DetectMasterConflict(selfId, masterHearbeatCh, roleCh)

		case network.Backup:
			<-roleCh
			go DetectBackupConflict(selfId, backupHeartbeatCh, roleCh, ctx, cancel)
			go BackupSelfElectMaster(selfId, masterHearbeatCh, roleCh, ctx)
		
		case network.Slave:
			<-roleCh
			go SlaveSelfElectBackup(selfId, masterHearbeatCh, roleCh)
		}
	}
}

// Split heartbeats into separte channels, only containing their IDs
func SplitHeartbeats(heartbeatCh chan network.Heartbeat, masterHeartbeatCh chan string, backupHeartbeatCh chan string) {

	for {
		heartbeat := <- heartbeatCh
		
		if (heartbeat.Role == network.Master) {
			masterHeartbeatCh <- heartbeat.ID
		}
		if (heartbeat.Role == network.Backup) {
			backupHeartbeatCh <- heartbeat.ID
		}
	}
}

// If more masters exists, set self to slave
func DetectMasterConflict(selfID string, masterHeartbeatCh chan string, roleCh chan network.Role) {

prevMasterID := ""

for {

	masterHeartbeat := <- masterHeartbeatCh
	
	// Conflict detected -> Suicide
	if (masterHeartbeat != prevMasterID) {
		roleCh <- network.Slave
		return
	}
	// Stash for comparison
	if (masterHeartbeat == prevMasterID || masterHeartbeat == "") {
		prevMasterID = masterHeartbeat
		}
	}
}

// If more backups exists, set self to slave
func DetectBackupConflict(selfID string, backupHeartbeatCh chan string, roleCh chan network.Role,  ctx context.Context, cancel context.CancelFunc) {

	prevBackupID := ""

	for {

		backupHeartbeat := <- backupHeartbeatCh
		
		// Conflict detected -> Suicide and kill goroutine: BackupSelfElectMaster()
		if (backupHeartbeat != prevBackupID) {
			roleCh <- network.Slave
			cancel()
			return
		}
		// Stash for comparison
		if (backupHeartbeat == prevBackupID || backupHeartbeat == "") { // nil for second expression
			prevBackupID = backupHeartbeat
		}
	}
}


// No master exists? Set yourself to master
func BackupSelfElectMaster(selfID string, masterHeartbeatCh chan string, roleCh chan network.Role,  ctx context.Context) {

	watchdog := time.NewTimer(config.Cfg.HeartbeatTimeout)

	for {
		select {
		// DetectBackupConflict murders this goroutine
		case <- ctx.Done():
			return
		// Heartbeat received -> reset timer
		case <- masterHeartbeatCh:
			watchdog.Reset(config.Cfg.HeartbeatTimeout)
			
		// Timeout -> suicide
		case <- watchdog.C:
			roleCh <- network.Master
			return 	
		}
	}	
}

// No backup exists? Set yourself to backup
func SlaveSelfElectBackup(selfID string, masterHeartBeatCh chan string, roleCh chan network.Role) {

	watchdog := time.NewTimer(config.Cfg.HeartbeatTimeout)

	for {
		select {
			// Heartbeat received -> reset timer
			case <- masterHeartBeatCh:
			watchdog.Reset(config.Cfg.HeartbeatTimeout)
			
			// Timeout -> suicide
			case <- watchdog.C:
			roleCh <- network.Backup
			return
		
		}
	}
}