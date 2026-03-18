package donaldtrump

import (
	"elevatorproject/src/config"
	"testing"
)


func TestDonaldTrump(t *testing.T) {
	config.Load()
	// RunMasterBrain("1")	
}

func TestJDVance(t *testing.T) {
	config.Load()
	// RunSlaveBrain("2")
}

