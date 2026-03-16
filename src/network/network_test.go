package network

import (
	"Network-go/network/bcast"
	"elevatorproject/src/config"
	"fmt"
	"testing"
)

func TestSender(t *testing.T){
	config.Load()

	sendOrdersCh := make(chan OrdersAndStateUpdate)
	ackCh := make(chan OrdersAndStateAck)

	go bcast.Transmitter(config.Cfg.MasterListenPort, sendOrdersCh)
	go bcast.Receiver(config.Cfg.MasterListenPort, ackCh)

	orderSender := &GenericSender[OrdersAndStateUpdate, OrdersAndStateAck]{
		SendCh: sendOrdersCh,
		AckIn: ackCh,
		AckResults: make(chan AckResult, 10), // buffered
	}

	msg := OrdersAndStateUpdate{"1", 2, "Hei ehi"}

	fmt.Println(msg.OrdersAndState)

	orderSender.SendAsyncWithAck(msg)

	res:=  <- orderSender.AckResults


	fmt.Println(res.UpdateNr)
}