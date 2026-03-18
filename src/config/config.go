package config

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

type Config struct {
	HeartbeatPort      int           `json:"heartbeatPort"`
	HeartbeatReplyPort int           `json:"slaveHeartbeatReplyPort"`
	HeartbeatInterval  time.Duration `json:"heartbeatInterval"`
	HeartbeatTimeout   time.Duration `json:"heartbeatTimeout"`

	SlaveListenPort  int `json:"slaveListenPort"`
	MasterListenPort int `json:"masterListenPort"`
	BackupPort       int `json:"backupPort"`

	NewBackupTimeoutTime time.Duration `json:"newBackupTimeoutTime"`
	NewMasterTimeoutTime time.Duration `json:"newMasterTimeoutTime"`

	AckRetryRate time.Duration `json:"ackRetryRateMs"`
	AckTimeout   time.Duration `json:"ackTimeout"`

	ElevatorUpdateRate time.Duration `json:"elevatorUpdateRate"`

	MaxOrderSuspendTime    time.Duration `json:"maxOrderSuspendTime"`
	MaxElevatorSuspendTime time.Duration `json:"maxElevatorSuspendTime"`

	ResendAssignmentTime time.Duration `json:"resendAssignmentTime"`
	// N_FLOORS int `json:"nFloors"`
	// N_BUTTONS int `json:"nButtons"` //TODO: Er dette forksjellige i elevatorManager og orderManager? Isåfall hva gjør man?
}

var Cfg Config

var defaultValues = Config{
	HeartbeatPort:      15647,
	HeartbeatReplyPort: 15648,
	HeartbeatInterval:  50 * time.Millisecond,  //Change to 15ms
	HeartbeatTimeout:   500 * time.Millisecond, //Change to 500ms

	SlaveListenPort:  15649,
	MasterListenPort: 15650,
	BackupPort:       15651,

	AckRetryRate: 200 * time.Millisecond,
	AckTimeout:   4 * time.Second,

	ElevatorUpdateRate: 200 * time.Millisecond,

	MaxOrderSuspendTime:    9 * time.Second,
	MaxElevatorSuspendTime: 5 * time.Second,

	NewBackupTimeoutTime: 500 * time.Millisecond, //Needs to be bigger than heartbeatinterval
	NewMasterTimeoutTime: 500 * time.Millisecond,

	ResendAssignmentTime: 500 * time.Millisecond,
}

// Load returns the config, falling back to defaults
func Load() {
	Cfg = defaultValues

	data, readErr := os.ReadFile("config.json")
	if readErr != nil {
		log.Println("No config file found, using defaults")
		return
	}

	parseErr := json.Unmarshal(data, &Cfg)
	if parseErr != nil {
		log.Println("Invalid config file, using defaults:", parseErr)
	}
}
