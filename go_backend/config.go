package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	ListenPort     int               `json:"listen_port"`
	Debug          bool              `json:"debug"`
	EventEndpoint  string            `json:"eventEndpoint"`
	HealthEndpoint string            `json:"healthEndpoint"`
	AuthToken      string            `json:"AuthToken"`
	Devices        map[string]string `json:"devices"`
}

var AppConfig Config

func LoadConfig(path string) error {
	file, err := os.Open(path)
	if err != nil {
		Error(fmt.Sprintf("failed to open config file: %v", err))
		os.Exit(1)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&AppConfig); err != nil {
		Error(fmt.Sprintf("failed to decode config file: %v", err))
		os.Exit(1)
	}
	return nil
}
