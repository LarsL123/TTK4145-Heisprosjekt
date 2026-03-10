package elevatormanager
// Idk en timer i guess
import (
	"time"
	"fmt"
)

var doorTimer *time.Timer

func doortimer_start(){
	doorTimer = time.NewTimer(DOOR_OPEN_DURATION*time.Second)
}



func doortimer(receiver chan bool){
	var timer1 time.Timer	
	for {
		select{ 
		case runTimer := <- receiver:
			switch runTimer{
			case true:
				fmt.Println("created timer")
				timer1.Reset(DOOR_OPEN_DURATION*time.Second)
			case false:
				timer1.Stop()
			}
		case <-timer1.C:
			receiver <- true	
		}
	}
}