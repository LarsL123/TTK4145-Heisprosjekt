package types

import (
	"time"
)

type ElevatorState struct {
	ID         string
	Floor      int
	Direction  string
	Behaviour  string
	CreatedAt  time.Time
	Obstructed bool
}

type OrderEnvelope struct {
	ElevatorID string
	UpdateNr   int
	Order      Order
	CreatedAt  time.Time
}

type OrderAck struct {
	UpdateNr int
}

type FinishedHallAssignments struct {
	ElevatorID string
	UpdateNr   int
	Orders     []Order
	CreatedAt  time.Time
}

type FinishedHallAssignmentsAck struct {
	UpdateNr int
}

type LivingMessage interface {
	GetCreationTime() time.Time
}

func (r OrderEnvelope) GetCreationTime() time.Time {
	return r.CreatedAt
}

func (r FinishedHallAssignments) GetCreationTime() time.Time {
	return r.CreatedAt
}

type Assignments struct {
	Assignments map[string][4][3]bool
}

type CabOrder struct {
	Floor int
}

type OrderType int

const (
	HallUp   OrderType = 0
	HallDown OrderType = 1
	Cab      OrderType = 2
)

type Order struct {
	Floor int
	Type  OrderType
}
