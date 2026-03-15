package reelection

import (
	"elevatorproject/src/config"
	"testing"
)


func TestClaimCrown(t *testing.T){

	config.Load()
	
	go ClaimCrown("1")

	select{}
}

func Test2(t *testing.T){

	config.Load()

	go InitReelection("2")

	select{}

}

func Test3(t *testing.T){

	config.Load()

	go InitReelection("3")
	
	select{}

}

func Test4(t *testing.T){

	config.Load()

	go InitReelection("4")

	select{}

}

func Test5(t *testing.T){

	config.Load()

	go InitReelection("5")

	select{}

}