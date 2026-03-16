package types

import "time"

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
	Floor     int
	Direction string
	Timestamp time.Time
}
type Assignments struct {
	Data map[string][4][2]bool
}

type cabOrder struct {
	Floor int
}

type RawMasterData struct {
    HallRequests [4][2]bool
    States       map[string]ElevatorState
}