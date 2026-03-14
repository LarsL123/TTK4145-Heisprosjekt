package reelection

import (
	"elevatorproject/src/config"
	"testing"
)


func TestMasterCollison(t *testing.T){
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

func Test2(t *testing.T){
	config.Load()

	// isMaster := make(chan bool)


	//Master
	// go network.StartMaster("1", isMaster)
	go InitReelection("2")

	select{}

	// <- time.After(2 *time.Second)

	// go network.StartSlave("2")
	// go InitReelection("2")
}

func Test3(t *testing.T){
	config.Load()

	// isMaster := make(chan bool)


	//Master
	// go network.StartMaster("1", isMaster)
	go InitReelection("3")

	select{}

	// <- time.After(2 *time.Second)

	// go network.StartSlave("2")
	// go InitReelection("2")
}

func Test4(t *testing.T){
	config.Load()

	// isMaster := make(chan bool)


	//Master
	// go network.StartMaster("1", isMaster)
	go InitReelection("4")

	select{}

	// <- time.After(2 *time.Second)

	// go network.StartSlave("2")
	// go InitReelection("2")
}

// func TestBackup(){

// }