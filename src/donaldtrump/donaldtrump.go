package donaldtrump

import (
	"Network-go/network/bcast"
	"elevatorproject/src/config"
	"elevatorproject/src/network"
	"fmt"
	"time"
)

/*

 Im gonne create the greates eleavtor youll ever see. Its gonne go to 1000 floors, yeah it can go to the moon.
 It will be the greatest and fastest elevator in the history of elevators.
 And we gonne make china pay for it.

*/

func RunMasterBrain(id string){
		isMaster := make(chan bool)
		slaveUpdate := network.StartMaster(id, isMaster)

		receiveOrdersAndStateUpdateCh := make(chan network.OrdersAndStateUpdate)
		orderAndStateAckCh := make(chan network.OrdersAndStateAck)
		

		sendAssignmentCh := make(chan network.AssignmentsAndOrders)
		assignmentAckCh := make(chan network.AssignementsAndOrdersAck)


		go bcast.Receiver(config.Cfg.MasterListenPort,receiveOrdersAndStateUpdateCh, assignmentAckCh)
		go bcast.Transmitter(config.Cfg.SlaveListenPort, orderAndStateAckCh, sendAssignmentCh)


		assignmentSender := &network.GenericSender[network.AssignmentsAndOrders, network.AssignementsAndOrdersAck]{
			SendCh: sendAssignmentCh,
			AckIn: assignmentAckCh,
			AckResults: make(chan network.AckResult, 10), // buffered OBS-OBS!! DO i need this??
		}

		msg := network.AssignmentsAndOrders{
					SourceId: id,
					UpdateNr: 1,
					OrdersAndState: "Ice will come to your home",
				}
		assignmentSender.UpdateAsyncGeneric(msg)

		count := 1

		for {
			select{
			case p := <-slaveUpdate:
				fmt.Printf("Slave update:\n")
				fmt.Printf("  Slaves:    %q\n", p.Slaves)
				fmt.Printf("  New:      %q\n", p.New)
				fmt.Printf("  Lost:     %q\n", p.Lost)
			case data := <- receiveOrdersAndStateUpdateCh: //Constant ack
				fmt.Printf("Received from slave: %s \n", data.OrdersAndState)
				fmt.Printf("count %d \n", count)
				count++

				if count > 3 {
					orderAndStateAckCh <- network.OrdersAndStateAck{UpdateNr: data.UpdateNr}
				}

			case <-time.After(6*time.Second): //New assignment to be distrebuted. 
				fmt.Println("Sending new assignments. ")

				msg := network.AssignmentsAndOrders{
					SourceId: id,
					UpdateNr: 1,
					OrdersAndState: "Ice will come to your home",
				}
				
				assignmentSender.UpdateAsyncGeneric(msg)
			}
		}
}
