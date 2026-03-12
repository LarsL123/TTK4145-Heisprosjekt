package donaldtrump

import (
	"elevatorproject/src/network"
	"fmt"
	"time"
)

func RunSlaveBrain(id string){
			orderSender := network.StartSlave(id)
			count := 1

			msg := network.OrdersAndStateUpdate{
						SourceId: id,
						UpdateNr: count,
						OrdersAndState: "Moren din er mann",
					}
			orderSender.UpdateAsyncGeneric(msg)

			for {
				select {
				case <-time.After(1000 *time.Millisecond): //Simulates updating orders. 
				
						count++
						msg := network.OrdersAndStateUpdate{
							SourceId: id,
							UpdateNr: count,
							OrdersAndState: "Moren din er mann",
						}

						orderSender.UpdateAsyncGeneric(msg)
						
				case res := <-orderSender.AckResults:

					//Check res.UpdateNr to check that the ACk is from the latest ctate change and is not old.
					if res.Err != nil {
						fmt.Println("Failed:", res.Err)
					} else {
						fmt.Println("ACK received for", res.UpdateNr)
					}
				}
			}
}