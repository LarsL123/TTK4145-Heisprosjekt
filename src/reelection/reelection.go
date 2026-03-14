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
	"math/rand"
	"strconv"
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
	go SplitHeartbeats(selfId,heartbeatCh, masterHearbeatCh, backupHeartbeatCh)

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
			network.StartMaster(selfId, ctx)

		case network.Backup:
			fmt.Println("I am backup:  ", selfId)
			go DetectBackupConflict(selfId, backupHeartbeatCh, roleCh, ctx, cancel)
			go BackupSelfElectMaster(selfId, masterHearbeatCh, roleCh, ctx, cancel)
			network.StartBackup(selfId, ctx)
		
		case network.Slave:
			fmt.Println("I am slave:  ", selfId)
			go SlaveSelfElectBackup(selfId, backupHeartbeatCh, roleCh)
		}
	}
}


// Split heartbeats into separte channels, only containing their IDs
func SplitHeartbeats(selfID string, heartbeatCh chan network.Heartbeat, masterHeartbeatCh chan string, backupHeartbeatCh chan string) {

	for {
		heartbeat := <- heartbeatCh
		
		if (heartbeat.Role == network.Master) {
			if(heartbeat.ID != selfID){
				fmt.Println("Mster heartbeat: ", heartbeat.ID, heartbeat.Role)
				masterHeartbeatCh <- heartbeat.ID
			}
		}

		if (heartbeat.Role == network.Backup) {
			if(heartbeat.ID != selfID){
				fmt.Println("Backup heartbeat: ", heartbeat.ID, heartbeat.Role)
				backupHeartbeatCh <- heartbeat.ID
			}
		}
	}
}

// If more masters exists, set self to slave
func DetectMasterConflict(selfID string, masterHeartbeatCh chan string, roleCh chan network.Role, cancel context.CancelFunc) {
	// prevMasterID := ""

	//Daniel sjef hva skal denne sleep tiden være? Er litt cursed men funker som faen. 
	id, _ := strconv.Atoi(selfID)
	r := rand.New(rand.NewSource(time.Now().UnixNano() + int64(id)))
	sleep := time.Duration(r.Int63n(int64(time.Second*2))) //TODO change to better number. 


	for {
		// masterHeartbeat := <- masterHeartbeatCh
		<- masterHeartbeatCh

		// if (prevMasterID == "") {
		// 	prevMasterID = masterHeartbeat
		// }

		// Conflict detected -> Suicide
		// if (masterHeartbeat != selfID) {
			fmt.Println("DOUBLE MASTER DETECTED!!! -> Suicide ")
			// prevMasterID = ""
			cancel()
			time.Sleep(sleep)
			roleCh <- network.Slave
			return
		// }
		// prevMasterID = masterHeartbeat
	}
}

// If more backups exists, set self to slave
func DetectBackupConflict(selfID string, backupHeartbeatCh chan string, roleCh chan network.Role,  ctx context.Context, cancel context.CancelFunc) {
	// prevBackupID := ""
	
	//Daniel sjef hva skal denne sleep tiden være? Er litt cursed men funker som faen. 
	id, _ := strconv.Atoi(selfID)
	r := rand.New(rand.NewSource(time.Now().UnixNano() + int64(id)))
	sleep := time.Duration(r.Int63n(int64(time.Second*2))) //TODO change to better number. 

	for {
		select{
			case  <- backupHeartbeatCh: //backupHeartbeat := <- backupHeartbeatCh:
				// if (prevBackupID == ""){
				// 	prevBackupID = backupHeartbeat
				// }
				// Conflict detected -> Suicide and kill goroutine: BackupSelfElectMaster()
				// if (backupHeartbeat != prevBackupID) {
					fmt.Println("BACKUP CONFLICT", selfID)
					// prevBackupID = ""
					cancel()
					time.Sleep(sleep)
					roleCh <- network.Slave
					return
				// }

			case <- ctx.Done():
				return
		}
	}
}


// No master exists? Set yourself to master
func BackupSelfElectMaster(selfID string, masterHeartbeatCh chan string, roleCh chan network.Role,  ctx context.Context, cancel context.CancelFunc) {
	for {
		select {
		// Heartbeat received -> reset timer
		case <- masterHeartbeatCh:
			continue
			
		// Timeout -> suicide
		case <- time.After(time.Second * 2):
			fmt.Println("Jeg blir master nå wallah")
			cancel()
			roleCh <- network.Master
			return 	

		case <- ctx.Done():
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
				fmt.Println("Denne er viktig sjef hvorfor?")
				continue
			
			// Timeout -> suicide
			case <- time.After(config.Cfg.HeartbeatTimeout):
				fmt.Println("Gidder ikke slave. ")
				roleCh <- network.Backup
				return
		}
	}


}