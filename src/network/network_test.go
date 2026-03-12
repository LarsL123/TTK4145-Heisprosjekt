package network

import (
	"Network-go/network/bcast"
	"elevatorproject/src/config"
)

func TestGenereicSender(){
	sendOrdersCh := make(chan OrdersAndStateUpdate)
	ackCh := make(chan OrdersAndStateAck)

	go bcast.Transmitter(config.Cfg.MasterListenPort, sendOrdersCh)
	go bcast.Receiver(config.Cfg.SlaveListenPort, ackCh)

	orderSender := &GenericSender[OrdersAndStateUpdate, OrdersAndStateAck]{
		SendCh: sendOrdersCh,
		AckIn: OrdersAndStateAck,
		AckResults: make(chan AckResult, 10), // buffered
	}

	// orderSender.UpdateAsyncGeneric()

}