package donaldtrump

import (
	"Network-go/network/bcast"
	"elevatorproject/src/config"
	"elevatorproject/src/network"
	"fmt"
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
		go bcast.Receiver(config.Cfg.MasterListenPort,receiveOrdersAndStateUpdateCh)

		orderAndStateAckCh := make(chan network.OrdersAndStateAck)
		go bcast.Transmitter(config.Cfg.SlaveListenPort, orderAndStateAckCh)

		count := 1

		for {
			select{
			case p := <-slaveUpdate:
				fmt.Printf("Slave update:\n")
				fmt.Printf("  Slaves:    %q\n", p.Slaves)
				fmt.Printf("  New:      %q\n", p.New)
				fmt.Printf("  Lost:     %q\n", p.Lost)
			case data := <- receiveOrdersAndStateUpdateCh:
				fmt.Printf("Received from slave: %s \n", data.OrdersAndState)
				fmt.Printf("count %d \n", count)
				count++

				if count > 3 {
					orderAndStateAckCh <- network.OrdersAndStateAck{UpdateNr: data.UpdateNr}
				}
			}
		}
}
