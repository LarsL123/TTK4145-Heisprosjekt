package reelection

import (
	"context"
	"fmt"
	"time"

	"Network-go/network/bcast"
	"elevatorproject/src/config"
	"elevatorproject/src/network"
)


func InitReelection2(selfID string) {
    heartbeatCh := make(chan network.Heartbeat, 32)

    go bcast.Receiver(config.Cfg.HeartbeatPort, heartbeatCh)

    go ReelectionFSM(selfID, heartbeatCh)
}

func ReelectionFSM(selfID string, heartbeatCh chan network.Heartbeat) {

    role := network.Slave

    masterTimer := time.NewTimer(2 * time.Second)
    backupTimer := time.NewTimer(2 * time.Second)

    var cancel context.CancelFunc
    var ctx context.Context

    startRole := func(r network.Role) {
        if cancel != nil {
            cancel()
        }

        ctx, cancel = context.WithCancel(context.Background())
        role = r

        switch r {
        case network.Master:
            fmt.Println("I am master:", selfID)
            go network.StartMaster(selfID, ctx)

        case network.Backup:
            fmt.Println("I am backup:", selfID)
            go network.StartBackup(selfID, ctx)

        case network.Slave:
            fmt.Println("I am slave:", selfID)
        }
    }

    startRole(network.Slave)

    for {
        select {

        case hb := <-heartbeatCh:

            if hb.ID == selfID {
                continue
            }

            switch hb.Role {
				case network.Master:
					masterTimer.Reset(2 * time.Second)

					if role == network.Master {
						if hb.ID > selfID {
							fmt.Println("Higher ID master detected → stepping down")
							startRole(network.Slave)
						}
					}

				case network.Backup:
					backupTimer.Reset(2 * time.Second)
					
					if role == network.Backup {			
						if hb.ID > selfID {
							fmt.Println("Higher ID backup detected → stepping down")
							startRole(network.Slave)
						}
					}

			}

        case <-masterTimer.C:
            if role == network.Backup {
                fmt.Println("Master dead → becoming master")
                startRole(network.Master)
            }

            masterTimer.Reset(2 * time.Second)

        case <-backupTimer.C:
            if role == network.Slave {
                fmt.Println("No backup → becoming backup")
                startRole(network.Backup)
            }

            backupTimer.Reset(2 * time.Second)
        }
    }
}