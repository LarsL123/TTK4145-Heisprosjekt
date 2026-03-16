package network

import (
	"Network-go/network/bcast"
	"elevatorproject/src/config"
	//"fmt"
	"testing"
	"elevatorproject/src/types"
	"time"
)

// func TestSender(t *testing.T){
// 	config.Load()

// 	sendOrdersCh := make(chan OrdersAndStateUpdate)
// 	ackCh := make(chan OrdersAndStateAck)

// 	go bcast.Transmitter(config.Cfg.MasterListenPort, sendOrdersCh)
// 	go bcast.Receiver(config.Cfg.MasterListenPort, ackCh)

// 	orderSender := &GenericSender[OrdersAndStateUpdate, OrdersAndStateAck]{
// 		SendCh: sendOrdersCh,
// 		AckIn: ackCh,
// 		AckResults: make(chan AckResult, 10), // buffered
// 	}

// 	msg := OrdersAndStateUpdate{"1", 2, "Hei ehi"}

// 	fmt.Println(msg.OrdersAndState)

// 	orderSender.SendAsyncWithAck(msg)

// 	res:=  <- orderSender.AckResults


// 	fmt.Println(res.UpdateNr)
// }

func TestNewSender(t *testing.T){
	config.Load()
	sendOrdersToMasterCh := make(chan types.HallOrder)
	ackCh := make(chan types.HallOrderAck)
	sendOrdersFromElevatorch := make(chan types.HallOrder)

	go bcast.Transmitter(config.Cfg.MasterListenPort,sendOrdersToMasterCh)
	go bcast.Receiver(config.Cfg.MasterListenPort,ackCh)


	go sendOrdersWithAck(sendOrdersFromElevatorch, sendOrdersToMasterCh,ackCh)
	go masterAckingOrders(sendOrdersToMasterCh,ackCh)
	sendOrdersFromElevatorch <- types.HallOrder{
		UpdateNr: 1,
		Floor: 1,
		Direction: 1,
		Timestamp: time.Now(),
	}
	time.Sleep(5*time.Second)
		sendOrdersFromElevatorch <- types.HallOrder{
		UpdateNr: 2,
		Floor: 2,
		Direction: 1,
		Timestamp: time.Now(),
	}
	time.Sleep(3*time.Second)
}

func masterAckingOrders(receiveOrdersCh chan types.HallOrder, AckOrderCh chan types.HallOrderAck){
	for{
		receivedOrder := <- receiveOrdersCh
		AckOrderCh <- types.HallOrderAck{
			UpdateNr: receivedOrder.UpdateNr,
		}
	}
}