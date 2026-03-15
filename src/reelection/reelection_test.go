package reelection

import (
	"elevatorproject/src/config"
	"testing"
)

// func TestClaimCrown(t *testing.T){

// 	config.Load()

// 	go ClaimCrown("1")

// 	select{}
// }

// func TestKillMaster(t *testing.T){
// 	config.Load()
// 	ctx, _ := context.WithCancel(context.Background())
// 	go network.SendHeartbeats("1", network.Master, ctx)

// 	select{}
// }

// func Test2(t *testing.T){

// 	config.Load()

// 	go InitReelection("2")

// 	select{}

// }

// func Test3(t *testing.T){

// 	config.Load()

// 	go InitReelection("3")

// 	select{}

// }

// func Test4(t *testing.T){

// 	config.Load()

// 	go InitReelection("4")

// 	select{}

// }

func TestNy2(t *testing.T){

	config.Load()

	go InitReelection2("2")

	select{}

}

func TestNy3(t *testing.T){

	config.Load()

	go InitReelection2("3")
	
	select{}

}

func TestNy4(t *testing.T){

	config.Load()

	go InitReelection2("4")

	select{}

}