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

	//Defining listening ports
	masterListenPort := 15555
	elevatorListenPort := 15556
	
	// ------Elevator Side------
	elevatorTxCh := make(chan types.HallOrder)
	elevatorRxAckCh := make(chan types.HallOrderAck)
	elevatorInputCh := make(chan types.HallOrder)

	go bcast.Transmitter(masterListenPort,elevatorTxCh)
	go bcast.Receiver(elevatorListenPort,elevatorRxAckCh)

	go 	SendOrdersWithAck(elevatorInputCh,elevatorTxCh,elevatorRxAckCh)


	// -------Master side------
	masterRxCh := make(chan types.HallOrder)
	masterTxAckCh := make(chan types.HallOrderAck)

	go bcast.Receiver(masterListenPort,masterRxCh)
	go bcast.Transmitter(elevatorListenPort,masterTxAckCh)

	go masterAckingOrders(masterRxCh,masterTxAckCh)

	elevatorInputCh <- types.HallOrder{
		UpdateNr: 1,
		Floor: 1,
		Direction: 1,
		Timestamp: time.Now(),
	}
	time.Sleep(5*time.Second)
		elevatorInputCh <- types.HallOrder{
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