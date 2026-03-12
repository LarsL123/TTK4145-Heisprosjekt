package network

import (
	"context"
	"sync"
)

type NetMessage interface {
    GetUpdateNr() int
}

type AckResult struct {
    UpdateNr int
    Err      error
}

type GenericSender [SenderType NetMessage, ReciverType NetMessage] struct {
    SendCh     chan<- SenderType
    AckIn      <-chan ReciverType
    AckResults chan AckResult

    cancelLast   context.CancelFunc // cancel previous pending send
    mu           sync.Mutex 
    lastUpdateNr int // Must use to prevent confict between different slaves
}

//Maye this should just be an envolope for ack logic??
type OrdersAndStateUpdate struct {
	SourceId string
	UpdateNr int
	OrdersAndState string //Custom type from daniea (mae7tro)
}

func (s OrdersAndStateUpdate) GetUpdateNr() int {
    return s.UpdateNr  
}

type OrdersAndStateAck struct {
	UpdateNr int
    Err error
}

func (s OrdersAndStateAck) GetUpdateNr() int{
    return s.UpdateNr
}



type AssignmentsAndOrders struct {
    SourceId string
	UpdateNr int
	OrdersAndState string //Custom type from Brage Drage
}

type AssignementsAndOrdersAck struct {
	UpdateNr int
}

func (s AssignmentsAndOrders) GetUpdateNr() int {
    return s.UpdateNr  
}


func (s AssignementsAndOrdersAck) GetUpdateNr() int{
    return s.UpdateNr
}




// type AssignmentSender struct {
//     SendCh     chan<- AssignmentsAndOrders
//     AckIn      <-chan AssignementsAndOrdersAck
//     AckResults chan AckResult

//     cancelLast   context.CancelFunc // cancel previous pending send
//     mu           sync.Mutex 
//     lastUpdateNr int // Must use to prevent confict between different slaves
// }

// type OrderSender struct {
//     SendCh     chan<- OrdersAndStateUpdate
//     AckIn      <-chan OrdersAndStateAck
//     AckResults chan AckResult

//     cancelLast   context.CancelFunc // cancel previous pending send
//     mu           sync.Mutex
//     lastUpdateNr int // Must use to prevent confict between different slaves
// }