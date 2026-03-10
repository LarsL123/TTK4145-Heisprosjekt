package config

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

type Config struct {
    HeartbeatPort int `json:"heartbeatPort"`
	SlaveHeartbeatReplyPort int `json:"slaveHeartbeatReplyPort"`
	SlaveListenPort int `json:"slaveListenPort"`
	MasterListenPort int `json:"masterListenPort"`
	AckRetryRateMs time.Duration `json:"ackRetryRateMs"`
}

var Cfg Config

var defaultValues = Config{
    HeartbeatPort: 15647,
	SlaveHeartbeatReplyPort: 15648,
	SlaveListenPort: 15649,
	MasterListenPort: 15650,
	AckRetryRateMs: 500,
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