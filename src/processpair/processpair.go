package processpair

import (
	"fmt"
	"os"
	"time"
)

func KillDeadProcess(masterAliveCh chan struct{}, slaveAliveCh chan struct{}) {
	timeout := 5 * time.Second
	startupGrace := 10 * time.Second
	masterAliveTimer := time.NewTimer(startupGrace)
	slaveAliveTimer := time.NewTimer(startupGrace)
	for {
		select {
		case <-masterAliveCh:
			masterAliveTimer.Reset(timeout)

		case <-slaveAliveCh:
			slaveAliveTimer.Reset(timeout)

		case <-masterAliveTimer.C:
			fmt.Println("Master dead - kms")
			os.Exit(1)

		case <-slaveAliveTimer.C:
			fmt.Println("Slave dead - kms")
			os.Exit(1)
		}
	}
}
