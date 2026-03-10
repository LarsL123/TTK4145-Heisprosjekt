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
			orderSender.UpdateAsync(msg, 3*time.Second)

			done := false

			for {
				select {
				case <-time.After(100 *time.Millisecond):
					if !done {
						count++
						msg := network.OrdersAndStateUpdate{
							SourceId: id,
							UpdateNr: count,
							OrdersAndState: "Moren din er mann",
						}

						orderSender.UpdateAsync(msg, 3*time.Second)
						done = true
					}
					
	

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