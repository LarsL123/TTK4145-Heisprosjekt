package network

import (
	"context"
	"sync"
)

type NetMessage interface {
    GetUpdateNr() int
}

//Maye this should just be an envolope for ack logic??
type OrdersAndStateUpdate struct {
	SourceId string
	UpdateNr int
	OrdersAndState string //Custom type from daniea (mae7tro)
}

type OrdersAndStateAck struct {
	UpdateNr int
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

type GenericSender [A NetMessage] struct {
    SendCh     chan<- A
    AckIn      <-chan A
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