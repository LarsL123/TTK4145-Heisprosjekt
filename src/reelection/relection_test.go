package reelection

/*

To run all tests at once in terminal:
go test

To run specific test in terminal:
go test -run Test...()

*/

import (
//	"fmt"
//	"net"
//	"reelection"
	"testing"
	"elevatorproject/src/network"

)

// Works as intended!
func TestElectRandomSlave(t *testing.T) {

	exampleRoles := map[string]network.Role{

		"id1": network.Slave,
		"id2": network.Slave,
		"id3": network.Slave,
		"id4": network.Slave,
		"id5": network.Slave,
		"id6": network.Slave,
		"id7": network.Slave,
	
	}

	ElectRandomSlave(exampleRoles)

	for id, role := range exampleRoles {
		println(id, ": ", role, "\n")
	}

}


// Works as intended!
func TestSetAllToSlave(t* testing.T) {

	exampleRoles := map[string]network.Role{

		"id1": network.Slave,
		"id2": network.Slave,
		"id3": network.Slave,
		"id4": network.Backup,
		"id5": network.Slave,
		"id6": network.Master,
		"id7": network.Slave,
	
	}

	SetAllToSlave(exampleRoles)

	for id, role := range exampleRoles {
		println(id, ": ", role, "\n")
	}	

}


func TestDetectMasterConflict(t* testing.T) {

	conflictCh := make(chan struct{})
	exampleRoles := map[string]network.Role{

		"id1": network.Slave,
		"id2": network.Slave,
		"id3": network.Slave,
		"id4": network.Backup,
		"id5": network.Slave,
		"id6": network.Master,
		"id7": network.Slave,
	
	}

	go DetectMasterConflict(exampleRoles, conflictCh)

	// Set two ids to master

	for {

		<- conflictCh
		println("Conflict detected!")

	}

}

func TestReelectMaster(t* testing.T) {

	selfId := "id4"
	// conflictChannel

	// To test case 1, update to no masters

	exampleRoles := map[string]network.Role{

		"id1": network.Slave,
		"id2": network.Slave,
		"id3": network.Slave,
		"id4": network.Master,
		"id5": network.Master,
		"id6": network.Slave,
		"id7": network.Slave,
	
	}

	// case 2:
	// go DetectMasterConflict(exampleRoles, )

	ReelectMaster(exampleRoles, selfId)
	
	for id, role := range exampleRoles {
		println(id, ": ", role, "\n")
	}	

}

// Difficult to test!
func TestReelectBackup(t* testing.T) {

	/*

	exampleRoles := map[string]network.Role{

		"id1": network.Slave,
		"id2": network.Slave,
		"id3": network.Slave,
		"id4": network.Backup,
		"id5": network.Slave,
		"id6": network.Master,
		"id7": network.Slave,
	
	}

	go ReelectBackup(exampleRoles)

	bcast.Transmitter(config.Cfg.HeartbeatPort, heartbeatCh)

	*/

}

// This function includes all other logic from the module. 
// Could be considered a full module test.
func TestSetupReelection(t *testing.T) {

	// exampleRoleCh := make(chan )

	// Difficult to setup a reasonable interface to network
	// Should we send full registry or only changes...?

	/*
	exampleRoles := map[string]network.Role{
		
		"id1": network.Slave,
		"id2": network.Slave,
		"id3": network.Slave,
		"id4": network.Backup,
		"id5": network.Slave,
		"id6": network.Master,
		"id7": network.Slave,
	
	}



	go SetupReelection()
	*/
}