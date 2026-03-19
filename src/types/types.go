package types

import (
	"time"
)

const N_FLOORS = 4

type ElevatorState struct {
	ID         string
	Floor      int
	Direction  string
	Behaviour  string
	CreatedAt  time.Time
	Obstructed bool
	Suspended  SuspendedType
}

type SuspendedType struct {
	IsSuspended bool
	TimeStamp   time.Time
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
	Assignments map[string][N_FLOORS][3]bool
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

func (orderType OrderType) ToString() string {
	switch orderType {
	case HallDown:
		return "HallDown"
	case HallUp:
		return "Hallup"
	case Cab:
		return "Cab"
	default:
		return "Invalid OrderType"
	}
}

type Order struct {
	Floor int
	Type  OrderType
}

type BackupData struct {
	UpdateNr     int
	HallRequests [N_FLOORS][2]bool
	CabRequests  map[string][N_FLOORS]bool
}

type BackupDataAck struct {
	UpdateNr int
}

type AssignedToAtTime struct {
	ElevatorId string
	TimeStamp  time.Time
}
