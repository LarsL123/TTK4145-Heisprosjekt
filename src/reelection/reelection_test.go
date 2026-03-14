package reelection

import (
	"elevatorproject/src/config"
	"testing"
)


func TestMaster(t *testing.T){
	config.Load()

	// isMaster := make(chan bool)


	//Master
	// go network.StartMaster("1", isMaster)
	go InitReelection("1")

	select{}

	// <- time.After(2 *time.Second)

	// go network.StartSlave("2")
	// go InitReelection("2")
}

func TestDoubleMaster(t *testing.T){
	config.Load()

	// isMaster := make(chan bool)


	//Master
	// go network.StartMaster("1", isMaster)
	go InitReelectionMasterTest("2")

	select{}

	// <- time.After(2 *time.Second)

	// go network.StartSlave("2")
	// go InitReelection("2")
}

// func TestBackup(){

// }