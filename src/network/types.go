package network

import (
	"context"
	"sync"
)

type NetMessage interface {
    GetUpdateNr() int
}

type NetAck interface {
    GetUpdateNr() int
}

//Maye this should just be an envolope for ack logic??
type OrdersAndStateUpdate struct {
	SourceId string
	UpdateNr int
	OrdersAndState string //Custom type from daniea (mae7tro)
}

func (s *OrdersAndStateUpdate) GetUpdateNr() int {
    return s.UpdateNr  
}

type OrdersAndStateAck struct {
	UpdateNr int
    Err error
}

func (s *OrdersAndStateAck) GetUpdateNr() int{
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

type AckResult struct {
    UpdateNr int
    Err      error
}

type GenericSender [A NetMessage, B NetAck] struct {
    SendCh     chan<- A
    AckIn      <-chan B
    AckResults chan AckResult

    cancelLast   context.CancelFunc // cancel previous pending send
    mu           sync.Mutex 
    lastUpdateNr int // Must use to prevent confict between different slaves
}


type AssignmentSender struct {
    SendCh     chan<- AssignmentsAndOrders
    AckIn      <-chan AssignementsAndOrdersAck
    AckResults chan AckResult

    cancelLast   context.CancelFunc // cancel previous pending send
    mu           sync.Mutex 
    lastUpdateNr int // Must use to prevent confict between different slaves
}

type OrderSender struct {
    SendCh     chan<- OrdersAndStateUpdate
    AckIn      <-chan OrdersAndStateAck
    AckResults chan AckResult

    cancelLast   context.CancelFunc // cancel previous pending send
    mu           sync.Mutex
    lastUpdateNr int // Must use to prevent confict between different slaves
}