package types

import (
	"time"
)

type ElevatorState struct {
	ID          string
	Floor       int
	Direction   string
	Behaviour   string
	CabRequests [4]bool
	CreatedAt   time.Time
	Obstructed  bool
}

type HallOrder struct {
	UpdateNr  int
	Floor     int //Refactor to order type
	Direction int
	CreatedAt time.Time
}

type HallOrderAck struct {
	UpdateNr int
}

type FinishedHallAssignments struct {
	UpdateNr  int
	Orders    []Order
	CreatedAt time.Time
}

type FinishedHallAssignmentsAck struct {
	UpdateNr int
}

type LivingMessage interface {
	GetCreationTime() time.Time
}

func (r HallOrder) GetCreationTime() time.Time {
	return r.CreatedAt
}

func (r FinishedHallAssignments) GetCreationTime() time.Time {
	return r.CreatedAt
}

type Assignements struct {
	Data map[string][4][2]bool
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
