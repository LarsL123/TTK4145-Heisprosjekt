package elevatormanager

import (
	"time"
)

var doortimer *time.Timer

func doortimer_init() {
	doortimer = time.NewTimer(time.Millisecond)
	doortimer.Stop()
}

func doortimer_start() {
	doortimer.Stop()
	doortimer.Reset(time.Duration(elevator.doorOpenDuration) * time.Second)
}
