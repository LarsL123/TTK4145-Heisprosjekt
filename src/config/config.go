package config

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

type Config struct {
    HeartbeatPort int `json:"heartbeatPort"`
	HeartbeatReplyPort int `json:"slaveHeartbeatReplyPort"`
	HeartbeatInterval time.Duration `json:"heartbeatInterval"`
	HeartbeatTimeout time.Duration `json:"heartbeatTimeout"`

	SlaveListenPort int `json:"slaveListenPort"`
	MasterListenPort int `json:"masterListenPort"`

	AckRetryRate time.Duration `json:"ackRetryRateMs"`
	AckTimeout time.Duration `json:"ackTimeout"`

	ElevatorUpdateRate time.Duration `json:"elevatorUpdateRate"`

	// N_FLOORS int `json:"nFloors"` 
	// N_BUTTONS int `json:"nButtons"` //TODO: Er dette forksjellige i elevatorManager og orderManager? Isåfall hva gjør man?
}

var Cfg Config

var defaultValues = Config{
    HeartbeatPort: 15647,
	HeartbeatReplyPort: 15648,
	HeartbeatInterval: 1000 * time.Millisecond, //Change to 15ms
	HeartbeatTimeout: 2000 * time.Millisecond, //Change to 500ms

	SlaveListenPort: 15649,
	MasterListenPort: 15650,
	
	AckRetryRate: 500*time.Millisecond,
	AckTimeout: 3*time.Second,

	ElevatorUpdateRate: 2*time.Second,
}

// Load returns the config, falling back to defaults
func Load() {
    Cfg = defaultValues

	data, readErr := os.ReadFile("config.json")
	if readErr != nil {
        log.Println("No config file found, using defaults")
        return
    }

	parseErr := json.Unmarshal(data, &Cfg);
    if  parseErr != nil {
        log.Println("Invalid config file, using defaults:", parseErr)
    }
}