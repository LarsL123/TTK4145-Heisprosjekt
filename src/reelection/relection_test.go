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
	"Network-go/network/bcast"
	"elevatorproject/src/config"
	"elevatorproject/src/network"
	"testing"
	"time"
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

// Works as intended!
func TestDetectMasterConflict(t* testing.T) {

	conflictCh := make(chan struct{}, 1)
	exampleRoles := map[string]network.Role{

		"id1": network.Slave,
		"id2": network.Slave,
		"id3": network.Backup,
		"id4": network.Slave,
		"id5": network.Slave,
		"id6": network.Master,
		"id7": network.Slave,
	
	}

	// View roles
	for id, role := range exampleRoles {
		println(id, ": ", role, "\n")
	}

	// First iterate: Only one master present
	DetectMasterConflict(exampleRoles, conflictCh)

	time.Sleep(3 * time.Second)

	// Set two ids to master
	exampleRoles["id1"] = network.Master

	// View updated roles
	for id, role := range exampleRoles {
		println(id, ": ", role, "\n")
	}		

	time.Sleep(3 * time.Second)
	
	// Second iterate: More than one master present
	DetectMasterConflict(exampleRoles, conflictCh)

	// If conflict is detected, alarm
	select {

	case <- conflictCh:
		println("Conflict detected!")
	
	default:
		println("No conflict detected!")
	}

}

// Depends on an instance of ReelectBackup() to be tested
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
	// 

	ReelectMaster(exampleRoles, selfId)
	
	for id, role := range exampleRoles {
		println(id, ": ", role, "\n")
	}	

}

// For setting up a UDP transmitter, required for testing
// TestReelectBackup, TestReelectMaster and TestSetupReelection
// To use, run in separate terminal

// DOES NOT WORK, try at lab
func TestTransmitter(t *testing.T) {
	
	heartbeatCh := make(chan network.Heartbeat, 1)
	go bcast.Transmitter(config.Cfg.HeartbeatPort, heartbeatCh)

	println("bcast started!")

	time.Sleep(5 * time.Second)

	heartbeatCh <- network.Heartbeat {
		
		ID: "id2",
		Role: network.Master,
		IP: "23456", // Random for pass to work
		
	}

	println("heartbeat sent!")

	time.Sleep(30 * time.Second)

}


// Difficult to test!
func TestReelectBackup(t* testing.T) {

	exampleRoles := map[string]network.Role{

		"id1": network.Slave,
		"id2": network.Slave,
		"id3": network.Slave,
		"id4": network.Backup,
		"id5": network.Slave,
		"id6": network.Master,
		"id7": network.Slave,
	
	}

	println("Initial roles")
	for id, role := range exampleRoles {
		println(id, ": ", role, "\n")
	}	

	heartbeatCh := make(chan network.Heartbeat, 1)

	go ReelectBackup(exampleRoles, heartbeatCh)

	heartbeatCh <- network.Heartbeat {
		
		ID: "id2",
		Role: network.Backup,
		IP: "23456", // Random for pass to work
		
	}

	// println("heartbeat sent!")

	// Print out status after heartbeat
	println("Roles after heartbeats is sent before first timeout,\nshould be the same")
	for id, role := range exampleRoles {
		println(id, ": ", role, "\n")
	}	

	time.Sleep(time.Second)

	heartbeatCh <- network.Heartbeat {
		
		ID: "id2",
		Role: network.Backup,
		IP: "23456", // Random for pass to work
		
	}

	// println("heartbeat sent!")

	// Print out status after heartbeat, should print out the same
	println("Timeout triggers new backup")
	for id, role := range exampleRoles {
		println(id, ": ", role, "\n")
	}	

	time.Sleep(3 * time.Second)

	// Print out status after new timeout, should print different backup
	println("Timeout triggers new backup")
	for id, role := range exampleRoles {
		println(id, ": ", role, "\n")
	}	
	

	time.Sleep(30 * time.Second)

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