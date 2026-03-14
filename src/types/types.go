package types

import "time"

type ElevatorState struct {
	ID               string
	Floor            int  
	Direction        string
	Behaviour        string
	CabRequests      [4]bool
	CreatedAt        time.Time
}
