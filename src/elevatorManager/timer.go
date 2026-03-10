package elevatormanager
// Idk en timer i guess
import (
	"time"
	//"fmt"
)

// Mulig det finnes noen skjulte bugs her, må høre med studass. Men egt ganske sikker på at det burde funke.

var doortimer *time.Timer

func doortimer_init(){
	doortimer = time.NewTimer(time.Millisecond)
	doortimer.Stop()
}

func doortimer_start(){
	doortimer.Stop()
	doortimer.Reset(time.Duration(elevator.doorOpenDuration)*time.Second)
}
