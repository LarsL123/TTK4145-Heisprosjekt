package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
    HeartbeatPort int `json:"heartbeatPort"`
	SlaveReplyPort int `json:"slaveReplyPort"`
}

var Cfg Config

var defaultValues = Config{
    HeartbeatPort: 15647,
	SlaveReplyPort: 15648,
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